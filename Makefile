build:
	GOOS=windows GOARCH=amd64 go build -o bin/golox.exe main.go
	GOOS=linux GOARCH=amd64 go build -o bin/golox main.go