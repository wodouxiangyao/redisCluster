SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go env -w GOSUMDB=off
go env -w GOPROXY=https://goproxy.cn,direct
go build -o redis-cluster main/main.go