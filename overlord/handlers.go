package overlord

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/giongto35/cloud-game/cws"
	"github.com/giongto35/cloud-game/overlord/gamelist"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

const (
	gameboyIndex = "./static/gameboy2.html"
	debugIndex   = "./static/gameboy2.html"
	gamePath     = "games"
)

type Server struct {
	roomToServer map[string]string
	// workerClients are the map serverID to worker Client
	workerClients map[string]*WorkerClient
}

var upgrader = websocket.Upgrader{}
var errNotFound = errors.New("Not found")

func NewServer() *Server {
	return &Server{
		// Mapping serverID to client
		workerClients: map[string]*WorkerClient{},
		// Mapping roomID to server
		roomToServer: map[string]string{},
	}
}

// GetWeb returns web frontend
func (o *Server) GetWeb(w http.ResponseWriter, r *http.Request) {
	indexFN := gameboyIndex

	bs, err := ioutil.ReadFile(indexFN)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(bs)
}

// WSO handles all connections from a new worker to overlord
func (o *Server) WSO(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connected")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Overlord: [!] WS upgrade:", err)
		return
	}
	// Register new server
	serverID := uuid.Must(uuid.NewV4()).String()
	log.Println("Overlord: A new server connected to Overlord", serverID)

	// Register to workersClients map the client connection
	client := NewWorkerClient(c)
	o.workerClients[serverID] = client
	defer o.cleanConnection(client, serverID)

	// Sendback the ID to server
	client.Send(
		cws.WSPacket{
			ID:   "serverID",
			Data: serverID,
		},
		nil,
	)

	// registerRoom event from a server, when server created a new room.
	// RoomID is global so it is managed by overlord.
	client.Receive("registerRoom", func(resp cws.WSPacket) cws.WSPacket {
		log.Println("Overlord: Received registerRoom ", resp.Data, serverID)
		o.roomToServer[resp.Data] = serverID
		return cws.WSPacket{
			ID: "registerRoom",
		}
	})

	// getRoom returns the server ID based on requested roomID.
	client.Receive("getRoom", func(resp cws.WSPacket) cws.WSPacket {
		log.Println("Overlord: Received a getroom request")
		log.Println("Result: ", o.roomToServer[resp.Data])
		return cws.WSPacket{
			ID:   "getRoom",
			Data: o.roomToServer[resp.Data],
		}
	})

	client.Listen()
}

// WSO handles all connections from frontend to overlord
func (o *Server) WS(w http.ResponseWriter, r *http.Request) {
	log.Println("Browser connected to overlord")
	//TODO: Add it back
	defer func() {
		if r := recover(); r != nil {
			log.Println("Warn: Something wrong. Recovered in ", r)
		}
	}()

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("[!] WS upgrade:", err)
		return
	}
	defer c.Close()

	// Set up server
	// SessionID will be the unique per frontend connection
	sessionID := uuid.Must(uuid.NewV4()).String()
	serverID, err := o.findBestServer()
	if err != nil {
		log.Fatal(err)
	}

	client := NewBrowserClient(c)

	// Setup session
	wssession := &Session{
		ID:            sessionID,
		handler:       o,
		BrowserClient: client,
		WorkerClient:  o.workerClients[serverID],
		ServerID:      serverID,
	}
	// TODO:?
	//defer wssession.Close()
	log.Println("New client will conect to server", wssession.ServerID)

	wssession.RouteBrowser()

	wssession.BrowserClient.Send(cws.WSPacket{
		ID:   "gamelist",
		Data: gamelist.GetEncodedGameList(gamePath),
	}, nil)

	// If peerconnection is done (client.Done is signalled), we close peerconnection
	go func() {
		<-client.Done
		// Notify worker to clean session
		wssession.WorkerClient.Send(
			cws.WSPacket{
				ID:        "terminateSession",
				SessionID: sessionID,
			},
			nil,
		)

		//log.Println("Socket terminated, detach connection")
		//h.detachPeerConn(wssession.peerconnection)
	}()

	wssession.BrowserClient.Listen()
}

// findBestServer returns the best server for a session
func (o *Server) findBestServer() (string, error) {
	// TODO: Find best Server by latency, currently return by ping
	if len(o.workerClients) == 0 {
		return "", errors.New("No server found")
	}

	r := rand.Intn(len(o.workerClients))
	for k, _ := range o.workerClients {
		if r == 0 {
			return k, nil
		}
		r--
	}

	return "", errors.New("No server found")
}

func (o *Server) cleanConnection(client *WorkerClient, serverID string) {
	log.Println("Unregister server from overlord")
	// Remove serverID from servers
	delete(o.workerClients, serverID)
	// Clean all rooms connecting to that server
	for roomID, roomServer := range o.roomToServer {
		if roomServer == serverID {
			delete(o.roomToServer, roomID)
		}
	}

	client.Close()
}
