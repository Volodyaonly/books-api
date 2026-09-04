package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// HandleSlowBook — пример долгой операции с контекстом.
func HandleSlowBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	fmt.Println("📚 Начали обработку запроса (долгая операция)...")

	select {
	case <-time.After(5 * time.Second): // имитация долгой работы
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("✅ Книга успешно обработана"))
		fmt.Println("✅ Завершено без таймаута")
	case <-ctx.Done(): // запрос прерван или истёк таймаут
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write([]byte("⏰ Таймаут: операция заняла слишком много времени"))
		fmt.Println("⛔ Работа остановлена:", ctx.Err())
		return
	}
}
