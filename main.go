package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// Order структура заказа
type Order struct {
	ID        int       `json:"id"`
	Item      string    `json:"item"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	orders []Order
	nextID = 1
	mu     sync.Mutex
)

func main() {
	// 1. Раздаем статические файлы из папки static
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// 2. API для работы с заказами
	http.HandleFunc("/orders", orderHandler)

	// 3. Настройка порта для Render.com
	// Render сам передает номер порта через переменную окружения PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Если запускаем локально и порта в настройках нет
	}

	fmt.Printf("🚀 Сервер запущен! Переходи по адресу: http://localhost:%s\n", port)

	// Запуск сервера
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	case http.MethodPost:
		var newOrder Order
		err := json.NewDecoder(r.Body).Decode(&newOrder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		newOrder.ID = nextID
		newOrder.CreatedAt = time.Now()
		nextID++
		orders = append(orders, newOrder)
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newOrder)
	}
}
