package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/bmizerany/pat"
)

var storage = make(map[string][]string)
var mu sync.Mutex

type data struct {
	Author string
	Book   string
}

type updateData struct {
	Old string
	New string
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the booklist!")
}

func adder(w http.ResponseWriter, r *http.Request) {

	mu.Lock()

	if r.Method == "GET" {

		author := r.URL.Query().Get(":author")
		book := r.URL.Query().Get(":book")
		storage[author] = append(storage[author], book)

	} else {

		var value data
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&value)

		if err != nil {
			fmt.Println(err)
		} else {
			storage[value.Author] = append(storage[value.Author], value.Book)
		}

	}

	mu.Unlock()
}

func changeAuthorBooks(old, new string) {
	if old != new {
		for _, book := range storage[old] {
			storage[new] = append(storage[new], book)
		}
		delete(storage, old)
	}
}

func updateAuthor(w http.ResponseWriter, r *http.Request) {
	mu.Lock()

	if r.Method == "GET" {
		old := r.URL.Query().Get(":old")
		new := r.URL.Query().Get(":new")
		changeAuthorBooks(old, new)
	} else {

		var value updateData
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&value)

		if err != nil {
			fmt.Println(err)
		} else {
			old := value.Old
			new := value.New
			changeAuthorBooks(old, new)
		}
	}
	mu.Unlock()
}

func changeBooks(old, new string) {
	for _, books := range storage {
		for i, book := range books {
			if book == old {
				books[i] = new
				return
			}
		}
	}
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	mu.Lock()

	if r.Method == "GET" {
		old := r.URL.Query().Get(":old")
		new := r.URL.Query().Get(":new")
		changeBooks(old, new)
	} else {

		var value updateData
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&value)

		if err != nil {
			fmt.Println(err)
		} else {
			old := value.Old
			new := value.New
			changeBooks(old, new)
		}
	}
	mu.Unlock()
}

func show(w http.ResponseWriter, r *http.Request) {
	for author, books := range storage {
		fmt.Fprintln(w, author+" :")
		for _, book := range books {
			fmt.Fprintln(w, book)
		}
		fmt.Fprintln(w, "")
	}
}

func main() {

	m := pat.New()
	m.Get("/", http.HandlerFunc(home))

	m.Get("/add/:author/:book", http.HandlerFunc(adder))
	m.Get("/update-author/:old/:new", http.HandlerFunc(updateAuthor))
	m.Get("/update-book/:old/:new", http.HandlerFunc(updateBook))
	m.Get("/show", http.HandlerFunc(show))

	m.Post("/add/", http.HandlerFunc(adder))
	m.Post("/update-author/", http.HandlerFunc(updateAuthor))
	m.Post("/update-book/", http.HandlerFunc(updateBook))

	http.Handle("/", m)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
