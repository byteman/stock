set GOOS=windows
set GOARCH=386
go build -ldflags "-s -w"  -v
upx tool.exe