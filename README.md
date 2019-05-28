# Klog, Web-based Cloud Gaming Service
- **SEA**: [http://cloud.webgame2d.com](http://cloud.webgame2d.com) - **US West**: [http://usw.cloud.webgame2d.com](http://usw.cloud.webgame2d.com)
- **US East**: [http://use.cloud.webgame2d.com](http://use.cloud.webgame2d.com) - **Europe**: [http://eu.cloud.webgame2d.com](http://eu.cloud.webgame2d.com)
- [**Game Instruction**](document/instruction/)
For the best gaming experience, please select the closest region to you.   

**Video demo**: https://www.youtube.com/watch?v=koqWB1VKflo

---

Klog is an open source Cloud Gaming Service building on [WebRTC](https://github.com/pion) using browser as the main platform.  
  
Klog aims to bring the most convenient gaming experience to gamer. You can play any games on your browser directly, which is fully compatible on multi-platform like Desktop, Android, IOS. This flexibility enables modern online gaming experience to retro games starting with NES in this current release.  

Note: The current state of POGO are not optimized for production. The service will still experience lag under heavy traffic. You can try hosting your own service following the instruction in the next session.  

![screenshot](document/img/landing-page.gif)

## Feature
1. Cloud gaming: Game logic is handled in server and streamed to user.
2. Cross-platform compatibility: The game is run on webbrowser, the most universal builtin app. No console, plugin, external app or devices are needed. The device must support webRTC to perform streaming. Joystick is also supported.
3. Vertically scaled: Services are distributed. The closest server with highest free resource will be assigned to user. (In development)
4. Cloud storage: Game state is storing on online storage, so you can come back to continue playing in a game.

## Run on local by Docker

You try hosting the server yourself by running `./run_local_docker.sh`. It will spawn a docker environment and you can access the emulator on `localhost:8000`.  

## Development environment

Install Golang https://golang.org/doc/install . Because the project uses GoModule, so it requires Go1.11 version.

Install dependencies  

  * Install [libvpx](https://www.webmproject.org/code/) and [pkg-config](https://www.freedesktop.org/wiki/Software/pkg-config/)
```
# Ubuntu
apt-get install -y pkg-config libvpx-dev libopus-dev libopusfile-dev

# MacOS
brew install libvpx pkg-config libopus libopusfile

# Windows
... not tested yet ...
```

And run 
  * `./run_local.sh`
  * The scripts includes build the binary using Go module

## Documentation
- ✏ [Design Doc](document/designdoc/)  
- 💿 [Implementation Doc](document/implementation/)  

## Follow up

This project demos the simplest cloud game with NES. Integrating with other emulator like GBA, NDS will also be possible. I'm welcome for the contribution.

## Credits

* *Pion* Webrtc team for the incredible Golang Webrtc library and their supports https://github.com/pion/webrtc/.  
* *fogleman* for the awesome NES emulator https://github.com/fogleman/nes.  
* *poi5305* for the video encoding https://github.com/poi5305/go-yuv2webRTC.  
* And last but not least, my longtime friend Tri as the co-author. 

## Contributor

Nguyen Huu Thanh  
https://www.linkedin.com/in/huuthanhnguyen/  

Tri Dang Minh  
https://trich.im  

