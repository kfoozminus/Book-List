package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

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

type User struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Name          string `json:"name,omitempty"`
	LastSessionID string `json:"lastsessionid,omitempty"`
}

var bookList []Book
var userList = make(map[string]User)

//var mu sync.Mutex

func isAuthorized(r *http.Request) bool {
	user, pass, err := r.BasicAuth()
	if err == false {
		//fmt.Println(user, pass)
		if userList[user].Password == pass {
			return true
		}
	} else {
		cookie, err := r.Cookie("SessionID")
		if err != nil {
			return false
		}
		sessionID := cookie.Value
		creden := strings.Split(sessionID, ":")
		user := creden[0]
		sessionID = creden[1]

		expectedSessionID := userList[user].LastSessionID

		if expectedSessionID == sessionID {
			return true
		}
	}
	return false
}

var ind int

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the BookList RESTful API!")
}

func addBook(w http.ResponseWriter, r *http.Request) {

	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		ind++
		book.Id = ind
		bookList = append(bookList, book)
		var _Book []Book
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Response{Success: 1, Message: "Added Book Successfully!", Book: append(_Book, book)})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid/Inefficient information"})
	}
}

func showBooks(w http.ResponseWriter, r *http.Request) {

	if isAuthorized(r) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "No Authorization Provided"})
		//http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

	//if bookList is empty
	if len(bookList) == 0 {
		json.NewEncoder(w).Encode(Response{Success: 1, Message: "No Book Added Yet"})
	} else {
		json.NewEncoder(w).Encode(Response{Success: 1, Message: "The Book List", Book: bookList})
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	if isAuthorized(r) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "No Authorization Provided"})
		//http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

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
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Success: 0, Message: "Book Not Found"})
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	if isAuthorized(r) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "No Authorization Provided"})
		//http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

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
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Success: 0, Message: "Book Not Found"})
}

/*func logout(user string) {
	userList[user].LastSessionID = ""
}*/

func login(w http.ResponseWriter, r *http.Request) {
	if isAuthorized(r) == true {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Please logout to login again!"})
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {

		if val, ok := userList[user.Username]; ok {

			if user.Password == val.Password {

				cookieValue := val.Username + ":" + val.Username + strconv.Itoa(rand.Intn(100000000))
				expire := time.Now().AddDate(0, 0, 1)
				cookie := http.Cookie{Name: "SessionID", Value: cookieValue, Expires: expire, HttpOnly: true}
				http.SetCookie(w, &cookie)
				json.NewEncoder(w).Encode(Response{Success: 1, Message: "Login Successful"})

			} else {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(Response{Success: 0, Message: "Password doesn't match"})
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Response{Success: 0, Message: "User Not Found"})
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Response{Success: 0, Message: "Login Unsuccessgul"})
}

func register(w http.ResponseWriter, r *http.Request) {
	if isAuthorized(r) == true {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Please logout to login again!"})
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid Info"})
		return
	}

	if _, ok := userList[user.Username]; ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Username already exists"})
		return
	}

	userList[user.Username] = user
}

func main() {

	m := pat.New()
	m.Get("/", http.HandlerFunc(homePage))
	m.Get("/book", http.HandlerFunc(showBooks))
	m.Post("/book/", http.HandlerFunc(addBook))
	m.Put("/book/:id", http.HandlerFunc(updateBook))
	m.Del("/book/:id", http.HandlerFunc(deleteBook))

	m.Post("/book/login", http.HandlerFunc(login))
	m.Post("/book/register", http.HandlerFunc(register))

	http.Handle("/", m)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
