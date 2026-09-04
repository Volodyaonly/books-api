package handlers

import (
	"log"
	"net/http"
	"sync"

	"books-api/models"
	"books-api/utils"
)

var (
	books = []models.Book{
		{
			ID:        1,
			Title:     "Мастер и Маргарита",
			Author:    "Булгаков",
			Available: true,
		},
		{
			ID:        2,
			Title:     "Гарри Поттер",
			Author:    "Роулинг",
			Available: false,
		},
	}

	nextID = 3
	mu     sync.RWMutex
)

func BooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBooks(w)

	case http.MethodPost:
		createBook(w, r)

	default:
		writeError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
	}
}

func BookByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBookByID(w, r)

	case http.MethodPut:
		updateBook(w, r)

	case http.MethodDelete:
		deleteBook(w, r)

	default:
		writeError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
	}
}

func getBooks(w http.ResponseWriter) {
	mu.RLock()

	// Создаём копию слайса, чтобы не держать блокировку
	// во время отправки HTTP-ответа.
	result := make([]models.Book, len(books))
	copy(result, books)

	mu.RUnlock()

	writeJSON(w, http.StatusOK, result)
}

func getBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIDFromPath(r.URL.Path, "/books/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	mu.RLock()

	for _, book := range books {
		if book.ID == id {
			mu.RUnlock()

			writeJSON(w, http.StatusOK, book)
			return
		}
	}

	mu.RUnlock()

	writeError(w, http.StatusNotFound, "Книга не найдена")
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book

	if err := utils.DecodeJSON(r, &newBook); err != nil {
		writeError(w, http.StatusBadRequest, "Невалидный JSON: "+err.Error())
		return
	}

	if newBook.Title == "" {
		writeError(w, http.StatusBadRequest, "Название книги обязательно")
		return
	}

	if newBook.Author == "" {
		writeError(w, http.StatusBadRequest, "Автор книги обязателен")
		return
	}

	mu.Lock()

	newBook.ID = nextID
	nextID++

	books = append(books, newBook)

	mu.Unlock()

	writeJSON(w, http.StatusCreated, newBook)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIDFromPath(r.URL.Path, "/books/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	var updatedBook models.Book

	if err := utils.DecodeJSON(r, &updatedBook); err != nil {
		writeError(w, http.StatusBadRequest, "Невалидный JSON: "+err.Error())
		return
	}

	if updatedBook.Title == "" {
		writeError(w, http.StatusBadRequest, "Название книги обязательно")
		return
	}

	if updatedBook.Author == "" {
		writeError(w, http.StatusBadRequest, "Автор книги обязателен")
		return
	}

	mu.Lock()

	for i, book := range books {
		if book.ID == id {
			books[i].Title = updatedBook.Title
			books[i].Author = updatedBook.Author
			books[i].Available = updatedBook.Available

			result := books[i]

			mu.Unlock()

			writeJSON(w, http.StatusOK, result)
			return
		}
	}

	mu.Unlock()

	writeError(w, http.StatusNotFound, "Книга не найдена")
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseIDFromPath(r.URL.Path, "/books/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	mu.Lock()

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)

			mu.Unlock()

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	mu.Unlock()

	writeError(w, http.StatusNotFound, "Книга не найдена")
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	if err := utils.WriteJSON(w, status, data); err != nil {
		log.Printf("ошибка отправки JSON: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	if err := utils.WriteError(w, status, message); err != nil {
		log.Printf("ошибка отправки ответа с ошибкой: %v", err)
	}
}
