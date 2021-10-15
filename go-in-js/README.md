# 在浏览器中运行go

## 依赖

安装[tinygo](https://tinygo.org/)

## usage

```shell
tinygo build -o static/tiny.wasm -target wasm ./go-in-js/wasm/
go run ./go-in-js/http-server/main.go
```

打开浏览器访问`http://localhost:8080/`，然后打开`F12`或`Ctrl+Shift+I`，观察控制台
