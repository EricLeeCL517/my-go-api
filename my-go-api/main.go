package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"text/template"
)

type BookResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Total  int    `json:"total"`
	Data   []Book `json:"data"`
}

type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Genre       string `json:"genre"`
	Description string `json:"description"`
	ISBN        string `json:"isbn"`
	Image       string `json:"image"`
	Published   string `json:"published"`
	Publisher   string `json:"publisher"`
}

var (
	originalBooks []Book
	books         []Book
	mu            sync.Mutex
)

func getBook() ([]Book, error) {
	url := "https://fakerapi.it/api/v1/books"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var BookResponse BookResponse
	err = json.NewDecoder(resp.Body).Decode(&BookResponse)
	if err != nil {
		return nil, err
	}

	return BookResponse.Data, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if len(originalBooks) == 0 {
		var err error
		originalBooks, err = getBook()
		if err != nil {
			http.Error(w, "Failed to fetch books: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if len(books) == 0 {
		books = make([]Book, len(originalBooks))
		copy(books, originalBooks)
	}

	tmpl := template.Must(template.ParseFiles("html/index.html"))

	err := tmpl.Execute(w, books)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func addBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newBook Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	newBook.ID = len(books) + 1
	books = append(books, newBook)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newBook)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var updatedBook Book
	err = json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, book := range books {
		if book.ID == id {
			updatedBook.ID = id
			books[i] = updatedBook
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(updatedBook)
			if err != nil {
				http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}

func resetBooksHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	books = make([]Book, len(originalBooks))
	copy(books, originalBooks)

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	fs := http.FileServer(http.Dir("html"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler)
	http.HandleFunc("/add", addBookHandler)
	http.HandleFunc("/update", updateBookHandler)
	http.HandleFunc("/delete", deleteBookHandler)
	http.HandleFunc("/reset", resetBooksHandler)

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
