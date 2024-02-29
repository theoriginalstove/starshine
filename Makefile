.PHONY: install
install:
	go build -o starshine-server ./cmd/server/main.go
	sudo cp ./starshine-server /usr/local/bin/starshine-server
	sudo systemctl start starshine-daemon
