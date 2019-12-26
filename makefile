migrate:
	go run ./cmd/chatserver-admin/main.go --db-disable-tls=1 migrate
useradd:
	go run ./cmd/chatserver-admin/main.go --db-disable-tls=1 useradd admin@example.com pass123
