migrate:
	go run ./cmd/chatserver-admin/main.go --db-disable-tls=1 migrate
