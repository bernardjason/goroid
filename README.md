# goroid, asteroid style game in Go using ebitengine

Asteroids style game where space is larger than the visible screen ut does wrap around. Can build to desktop as well as wasm

Left,Right and up arrow keys to move. Space to fire

It's my first program with Go btw...

See it deployed here https://bjason.org/en/go/golangasteroids/


To run
```
go run main.go
```

wasm build
```
./wasm.sh
cd html
python3 -m http.server
```
then check http://127.0.0.1:8000
