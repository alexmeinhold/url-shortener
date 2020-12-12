# url-shortener
Url shortener in Go

## Setup
```bash
go get github.com/syndtr/goleveldb/leveldb
go build main.go
./main
```

## Usage
```bash
curl -d "url=http://stackoverflow.com" -X POST http://localhost:8080
# -> http://localhost:8080/81b30731
```
