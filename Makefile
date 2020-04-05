default:
	go build cmd/msnap.go cmd/conf.go

test:
	go test ./lib/...
