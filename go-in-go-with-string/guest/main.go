package main

import (
	"encoding/binary"
	"unsafe"
)

// tinygo build -o str.wasm -target wasi ./go-in-go-with-string/guest/
func main() {

}

// 在host中实现的方法
//go:wasm-module xmh
//export modifyByHost
func modifyByHost(addr int32) int32

// 在guest中实现的方法
//export modifyByGuest
func modifyByGuest(addr int32) int32 {
	// 读取字符串
	str := addr2str(addr)
	// 在guest中修改字符串
	str += "-SuffixByGuest"
	// 将字符串转换成地址
	addr = str2addr(str)
	// 调用host再次修改字符串
	addr = modifyByHost(addr)
	return addr
}

// 给host调用申请内存
//export wasm_alloc
func wasm_alloc(size int32) int32 {
	buf := make([]byte, size)
	return int32(uintptr(unsafe.Pointer(&buf[0])))
}

// 从字符串转为包含数据的地址
func str2addr(str string) int32 {
	// 创建一个byte数组，前4个字节放字符串长度，后面放字符串，最后是空白字符
	s := make([]byte, len(str)+4+1)
	binary.LittleEndian.PutUint32(s[0:], uint32(len(str)+1))
	copy(s[4:], str)
	s[4+len(str)] = 0
	return int32(uintptr(unsafe.Pointer(&s[0])))
}

// 从包含字符串的地址读取字符串
func addr2str(addr int32) string {
	// 从前4个字节中读取长度
	ptr0 := (*byte)(unsafe.Pointer(uintptr(addr)))
	ptr1 := (*byte)(unsafe.Pointer(uintptr(addr + 1)))
	ptr2 := (*byte)(unsafe.Pointer(uintptr(addr + 2)))
	ptr3 := (*byte)(unsafe.Pointer(uintptr(addr + 3)))
	bs := []byte{*ptr0, *ptr1, *ptr2, *ptr3}
	size := int32(binary.LittleEndian.Uint32(bs[:]))

	// 读取字符串数据
	bs = make([]byte, size-1) // size-1是因为最后一个为空白字符
	for i := int32(0); i < size-1; i++ {
		bs[i] = *(*byte)(unsafe.Pointer(uintptr(addr + 4 + i)))
	}
	return string(bs)
}
