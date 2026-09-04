package main

import (
	"fmt"
	"log"
	"net/http"

	"books-api/handlers"
	"books-api/middleware"
)

func main() {
	// Создаём собственный маршрутизатор.
	mux := http.NewServeMux()

	// Работа со списком книг:
	// GET  /books
	// POST /books
	mux.HandleFunc("/books", handlers.BooksHandler)

	// Работа с конкретной книгой:
	// GET    /books/{id}
	// PUT    /books/{id}
	// DELETE /books/{id}
	mux.HandleFunc("/books/", handlers.BookByIDHandler)

	mux.HandleFunc("/slow", handlers.HandleSlowBook)

	// Важно явно указать тип http.Handler.
	// Middleware возвращают http.Handler, а не *http.ServeMux.
	var handler http.Handler = mux

	// Получается цепочка:
	// Recovery -> Logging -> ServeMux -> конкретный handler.
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.RecoveryMiddleware(handler)

	fmt.Println("Сервер запущен на http://localhost:8080")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
