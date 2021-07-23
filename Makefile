client:
	go build -o bin/client ./services/client/cmd/client/client.go

hub:
	go build -o bin/hub ./services/hub/cmd/hub/hub.go

apitest:
	go build -o bin/apitest ./services/client/cmd/apitest/apitest.go

migrate:
	docker run -v $(PWD)/services/hub/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://${SCR_DATABASE_USERNAME}:${SCR_DATABASE_PASSWORD}@${SCR_DATABASE_URL}:${SCR_DATABASE_PORT}/${SCR_DATABASE_NAME} up
