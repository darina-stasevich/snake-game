
1. clone repo

```bash
git clone https://github.com/darina-stasevich/snake-game.git
cd snake-game
```

2. go to project folder, then build docker image

```bash
docker build -t snake-game:latest .
```

3. start game

(for Linux using Wayland)
```bash
xhost +local:
docker run -it --rm --net=host -e DISPLAY=$DISPLAY -e WAYLAND_DISPLAY=$WAYLAND_DISPLAY -e XDG_RUNTIME_DIR=$XDG_RUNTIME_DIR -v $XDG_RUNTIME_DIR/$WAYLAND_DISPLAY:$XDG_RUNTIME_DIR/$WAYLAND_DISPLAY snake-game:latest
```
(for Linux using X11)
```bash
xhost +local:
docker run -it --rm -e DISPLAY=$DISPLAY -v /tmp/.X11-unix:/tmp/.X11-unix snake-game:latest
```