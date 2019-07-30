package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/syndtr/goleveldb/leveldb"
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

func post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	url := r.PostFormValue("url")

	// TODO: check for collisions
	key := RandStringBytes(8)

	// TODO: check if url starts with http://
	fmt.Printf("put %s into db\n", url)
	err = db.Put([]byte(key), []byte(url), nil)

	if err != nil {
		panic(err)
	}

	fullUrl := "http:/" + r.URL.String() + r.Host + "/" + key + "\n"
	w.Write([]byte(fullUrl))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	newUrl, err := db.Get([]byte(key), nil)

	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, string(newUrl), http.StatusSeeOther)
}

func main() {
	var err error
	db, err = leveldb.OpenFile("urls", nil)

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", post).Methods("POST")
	r.HandleFunc("/{id}", redirect)
	http.ListenAndServe(":8080", r)

	defer db.Close()
}
