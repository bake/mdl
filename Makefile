LDFLAGS=-ldflags "-X main.version=`git describe`"

default:
	GOOS=darwin GOARCH=amd64 go build -o mdl-darwin-amd64 ${LDFLAGS}
	GOOS=linux GOARCH=amd64 go build -o mdl-linux-amd64 ${LDFLAGS}
	GOOS=windows GOARCH=amd64 go build -o mdl-windows-amd64.exe ${LDFLAGS}
