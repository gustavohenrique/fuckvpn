arm:
	@GOOS=darwin GOARCH=arm64 go build -o fuckvpn-darwin main.go

linux:
	@GOOS=linux GOARCH=amd64 go build -o fuckvpn main.go
	@upx fuckvpn

windows:
	@GOOS=windows GOARCH=amd64 go build -o fuckvpn.exe main.go
