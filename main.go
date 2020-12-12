package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var db *leveldb.DB

func handlePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	url := r.PostFormValue("url")
	if !strings.HasPrefix(url, "http") {
		w.Write([]byte("Error: Invalid URL\n"))
		return
	}

	h := sha256.New()
	h.Write([]byte(url))
	hash := hex.EncodeToString(h.Sum(nil))
	// take first 4 bytes of hash
	key := hash[:8]

	err = db.Put([]byte(key), []byte(url), nil)
	if err != nil {
		panic(err)
	}

	fullUrl := "http:/" + r.URL.String() + r.Host + "/" + key + "\n"
	w.Write([]byte(fullUrl))
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// 1. get key from url
	key := r.URL.Path[1:]
	// 2. check if key is in valid format
	matched, err := regexp.MatchString(`^([a-f0-9]{8})$`, key)
	if err != nil {
		log.Fatal(err)
	}
	if !matched {
		http.NotFound(w, r)
		return
	}
	// 3. look up key in database to get url
	newUrl, err := db.Get([]byte(key), nil)
	if err != nil {
		log.Fatal(err)
	}
	// 4. send redirect response to user with url
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
