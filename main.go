package main

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"math/rand"
	"net/http"
)

var db *leveldb.DB

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	url := r.PostFormValue("url")

	// TODO: check for collisions
	key := RandStringBytes(8)

	// TODO: check if url starts with http://
	err = db.Put([]byte(key), []byte(url), nil)

	if err != nil {
		panic(err)
	}

	fullUrl := "http:/" + r.URL.String() + r.Host + "/" + key + "\n"
	w.Write([]byte(fullUrl))
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// 1. get key from url
	// TODO: check if key is in valid format
	key := r.URL.Path[1:]
	// 2. look up key in database to get url
	newUrl, err := db.Get([]byte(key), nil)
	if err != nil {
		log.Fatal(err)
	}
	// 3. send redirect response to user with url
	http.Redirect(w, r, string(newUrl), http.StatusSeeOther)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL)

	if r.Method == "POST" {
		handlePost(w, r)
	} else if r.Method == "GET" {
		handleRedirect(w, r)
	}
}

func main() {
	var err error
	db, err = leveldb.OpenFile("urls", nil)

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", root)
	http.ListenAndServe(":8080", nil)

	defer db.Close()
}
