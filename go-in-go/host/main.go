package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bytecodealliance/wasmtime-go"
)

var vm *wasmtime.Instance
var store *wasmtime.Store

// go run ./go-in-go/host/
func main() {
	bs := readFile("go.wasm")
	vm, store = createWasmVM(bs)

	// 执行wasm中的导出方法
	fn := vm.GetExport(store, "multiply").Func()
	res, err := fn.Call(store, 2, 3)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("finish:", res.(int32))

	// 执行wasm的入口方法
	start := vm.GetFunc(store, "_start")
	if start != nil {
		start.Call(store)
	}
}

// 从文件中读取到byte数组
func readFile(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return bs
}

// 从字节数组创建wasm
func createWasmVM(code []byte) (*wasmtime.Instance, *wasmtime.Store) {
	cfg := wasmtime.NewConfig()
	cfg.SetInterruptable(true)
	engine := wasmtime.NewEngineWithConfig(cfg)
	module, err := wasmtime.NewModule(engine, code)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil
	}
	store := wasmtime.NewStore(engine)
	_, err = store.InterruptHandle()
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil
	}

	// 解决当wasm代码中有println等标准输出时，报错`thread '<unnamed>' panicked at 'called `Option::unwrap()` on a `None` value', crates\c-api\src\linker.rs:84:80`
	wasi := wasmtime.NewWasiConfig()
	wasi.InheritStdout()
	store.SetWasi(wasi)

	linker := wasmtime.NewLinker(engine)
	importHostFuncs(linker, store)
	err = linker.DefineWasi() // 解决报错`unknown import: `wasi_snapshot_preview1::fd_write` has not been defined`
	if err != nil {
		fmt.Println(err.Error())
	}

	inst, err := linker.Instantiate(store, module)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
		return nil, nil
	}

	return inst, store
}

// 向wasm中注入func
func importHostFuncs(linker *wasmtime.Linker, store *wasmtime.Store) {
	linker.DefineFunc(store, "xmh", "add", add)
}

func add(x, y int32) int32 {
	return x + y
}
