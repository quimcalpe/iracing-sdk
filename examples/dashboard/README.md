## Compile

```bash
go get github.com/gorilla/websocket
go build -o dashboard.exe main.go
```

## Run

Default host and port:
```bash
dashboard.exe
```

With custom host and port:
```bash
dashboard.exe -addr=192.168.1.100:8888
```
