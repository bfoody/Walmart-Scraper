client:
	go build -o bin/client ./services/client/cmd/client/client.go

hub:
	go build -o bin/hub ./services/hub/cmd/hub/hub.go
