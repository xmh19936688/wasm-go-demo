# 在go中运行go

## 依赖

安装[tinygo](https://tinygo.org/)

## usage

```shell
go mod tidy && go mod vendor
tinygo build -o str.wasm -target wasi ./go-in-go-with-string/guest/
go run ./go-in-go-with-string/host/
```