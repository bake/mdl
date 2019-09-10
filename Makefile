LDFLAGS=-ldflags "-X main.version=`git rev-parse HEAD`"

default:
	GOOS=darwin GOARCH=amd64 go build -o mdl-darwin-amd64 ${LDFLAGS}
	GOOS=linux GOARCH=amd64 go build -o mdl-linux-amd64 ${LDFLAGS}
