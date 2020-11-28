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
# then navigate to returned link to be redirected to url
```
