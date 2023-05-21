cp $(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js
GOOS=js GOARCH=wasm go build -o ./html/main.wasm .
