SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go env -w GOSUMDB=off
go env -w GOPROXY=https://goproxy.cn,direct
go build -ldflags "-w -s" -o redis-cluster main/main.go && upx.exe redis-cluster