package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// A global variable that is incremented everytime a book is added.
// Used for providing a unique ID to each book
var count int

// Book Struct
type Book struct {
	ID     int     `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author Struct
type Author struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var books []Book

// Give all the books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(books)
}

// Give a book with some ID
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	for _, book := range books {
		if strconv.Itoa(book.ID) == vars["id"] {
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	json.NewEncoder(w).Encode(&Book{})
}

// Adds a new Book
func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	book.ID = count
	count++

	books = append(books, book)

	json.NewEncoder(w).Encode(book)
}

// Updates a book with some ID
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var tempBook Book
	for index, book := range books {
		if strconv.Itoa(book.ID) == vars["id"] {
			_ = json.NewDecoder(r.Body).Decode(&tempBook)
			tempBook.ID = index
			books[index] = tempBook
			json.NewEncoder(w).Encode(books[index])
			return
		}
	}

	json.NewEncoder(w).Encode(&Book{})
}

// Deletes the book with some ID
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var tempBook Book
	for index, book := range books {
		if strconv.Itoa(book.ID) == vars["id"] {
			tempBook = books[index]
			books = append(books[:index], books[index+1:]...)
			_ = json.NewEncoder(w).Encode(tempBook)
			return
		}
	}

	json.NewEncoder(w).Encode(&Book{})
}

func main() {
	// Initialize count with 0
	count = 0

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	// Initializing with some dummy data
	books = append(books, Book{ID: count, Isbn: "123456", Title: "Book 1", Author: &Author{FirstName: "Rehan", LastName: "Javed"}})
	count++
	books = append(books, Book{ID: count, Isbn: "456280", Title: "Book 2", Author: &Author{FirstName: "Haider", LastName: "Ali"}})
	count++

	// Handling the endpoints
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", addBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// Running the Server
	log.Fatal(http.ListenAndServe(":8080", router))
}
