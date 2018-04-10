setup:
	@go get github.com/mattn/go-sqlite3
	@echo "Building..." && go build

run: setup
	@./go-challenge-ac001
