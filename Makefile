client:
	go build -o bin/client ./services/client/cmd/client/client.go

hub:
	go build -o bin/hub ./services/hub/cmd/hub/hub.go

apitest:
	go build -o bin/apitest ./services/client/cmd/apitest/apitest.go
