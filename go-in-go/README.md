# 在go中运行go

## 依赖

安装[tinygo](https://tinygo.org/)

## usage

```shell
go mod tidy && go mod vendor
tinygo build -o go.wasm -target wasm ./go-in-go/guest/
go run ./go-in-go/host/
```