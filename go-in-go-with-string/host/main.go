package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bytecodealliance/wasmtime-go"
)

var vm *wasmtime.Instance
var store *wasmtime.Store

// go run ./go-in-go-with-string/host/
func main() {
	bs := readFile("str.wasm")
	vm, store = createWasmVM(bs)

	// 创建用于测试的数据并转为地址
	str := "test"
	addr := str2addr(str)
	// 调用guest的方法
	fn := vm.GetExport(store, "modifyByGuest").Func()
	res, err := fn.Call(store, addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 从返回地址中读取字符串
	addr = res.(int32)
	str = addr2str(addr)
	fmt.Println("finish:", str)
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
	linker.DefineFunc(store, "xmh", "modifyByHost", modifyByHost)
}

// 根据地址从wasm的内存中读取字符串
func addr2str(addr int32) string {
	// 获取wasm的内存
	mem := vm.GetExport(store, "memory").Memory().UnsafeData(store)
	// 从前4字节读取字符串长度
	size := int32(binary.LittleEndian.Uint32(mem[addr:]))
	// 创建byte数组存放字符串
	data := make([]byte, size-1) // size-1是因为最后为空白字符
	// 跳过前4字节读取数据
	copy(data, mem[addr+4:])
	return string(data)
}

// 将字符串写到wasm的内存中并返回地址
func str2addr(s string) int32 {
	// 获取wasm的内存
	mem := vm.GetExport(store, "memory").Memory().UnsafeData(store)

	// 向wasm申请一块内存空间
	ex := vm.GetExport(store, "wasm_alloc")
	if ex == nil {
		fmt.Println("handle error!")
		return 0
	}
	fn := ex.Func()
	if fn == nil {
		fmt.Println("handle error!")
		return 0
	}
	vaddr, e := fn.Call(store, len(s)+4+1) // +4是因为前4个字节放字符串长度；+1是因为最后为空白字符
	if e != nil {
		panic(e)
	}
	addr := vaddr.(int32)

	// 先写入字符串长度
	binary.LittleEndian.PutUint32(mem[addr:], uint32(len(s)+1))
	// 写入字符串数据
	copy(mem[addr+4:], s)
	// 写入空白字符到最后
	mem[addr+4+int32(len(s))] = 0

	return addr
}

// 在host中实现的方法
func modifyByHost(addr int32) int32 {
	// 读取字符串
	str := addr2str(addr)
	// 修改字符串
	str = "PrefixByHost-" + str
	// 将字符串转换成地址返回
	addr = str2addr(str)
	return addr
}
