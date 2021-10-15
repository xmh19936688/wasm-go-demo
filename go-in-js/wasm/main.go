package main

// tinygo build -o static/tiny.wasm -target wasm ./go-in-js/wasm/
func main() {
	println("add:", add(2, 3))
}

// 通过`go:wasm-module`指定module名称 通过`export`指定导出名称
//go:wasm-module xmh
//export add
func add(x, y int32) int32

//export multiply
func multiply(x, y int32) int32 {
	return x*y + add(x, y) // return 23
	//return x * y // return 15
}
