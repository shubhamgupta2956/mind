build:
	@echo "[*] Building as ./target/mind"
	@test -d target || mkdir target
	@go build -o ./target/mind cmd/mind/main.go
	@echo "[+] Build complete!"

.PHONY: build
