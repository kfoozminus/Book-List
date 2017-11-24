package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/bmizerany/pat"
)

type Book struct {
	Id     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Author string `json:"author,omitempty"`
}

type Response struct {
	Success int    `json:"success"`
	Message string `json:"message,omitempty"`
	Book    []Book `json:"book,omitempty"`
}

var bookList []Book

//var mu sync.Mutex
var ind int

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the BookList RESTful API!")
}

func addBook(w http.ResponseWriter, r *http.Request) {

	ind++
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.Id = ind
	bookList = append(bookList, book)
	var _Book []Book
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Success: 1, Message: "Added Book Successfully!", Book: append(_Book, book)})
}

func showBooks(w http.ResponseWriter, r *http.Request) {
	//if bookList is empty
	if len(bookList) == 0 {
		json.NewEncoder(w).Encode(Response{Success: 1, Message: "No Book Added Yet"})
	} else {
		json.NewEncoder(w).Encode(Response{Success: 1, Message: "The Book List", Book: bookList})
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	var delBook Book
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		//not valid
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid ID"})
		return
	}
	for i, book := range bookList {
		if book.Id == id {
			delBook = book
			bookList = append(bookList[:i], bookList[i+1:]...)
			//json.NewEncoder(w).Encode(delBook)
			var _Book []Book
			json.NewEncoder(w).Encode(Response{Success: 1, Message: "Deleted Book Successfully!", Book: append(_Book, delBook)})
			return
		}
	}
	//not found
	//w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Success: 0, Message: "Book Not Found"})
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		//not valid
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid ID"})
		return
	}
	for i, book := range bookList {
		if book.Id == id {
			_ = json.NewDecoder(r.Body).Decode(&bookList[i])
			bookList[i].Id = id
			var _Book []Book
			json.NewEncoder(w).Encode(Response{Success: 1, Message: "Updated Book Info Successfully!", Book: append(_Book, bookList[i])})
			return
		}
	}
	//not found
	//w.WriteHeader(http.StatusNotFound)
	fmt.Println(Response{Success: 0, Message: "Book Not Found"})
	json.NewEncoder(w).Encode(Response{Success: 0, Message: "Book Not Found"})
}

func main() {

	m := pat.New()
	m.Get("/", http.HandlerFunc(homePage))
	m.Get("/book", http.HandlerFunc(showBooks))
	m.Post("/book/", http.HandlerFunc(addBook))
	m.Put("/book/:id", http.HandlerFunc(updateBook))
	m.Del("/book/:id", http.HandlerFunc(deleteBook))

	http.Handle("/", m)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
