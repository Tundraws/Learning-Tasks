package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Инициализация кэша городов
	initCityCache()

	// Настройка маршрутов
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/cities", citiesHandler)
	http.HandleFunc("/api/calculate", calculateHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка загрузки шаблона: %v", err), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка рендеринга шаблона: %v", err), http.StatusInternalServerError)
	}
}

func citiesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cities, err := getCities()
	if err != nil {
		log.Printf("Ошибка получения городов: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Не удалось получить список городов",
			"details": err.Error(),
		})
		return
	}

	if err := json.NewEncoder(w).Encode(cities); err != nil {
		log.Printf("Ошибка кодирования городов: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Ошибка формирования ответа",
		})
	}
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получаем параметры из URL
	city := r.URL.Query().Get("city")
	weight := r.URL.Query().Get("weight")

	// Вычисляем стоимость доставки
	result, err := calculateDelivery(city, weight)
	if err != nil {
		log.Printf("Ошибка расчета доставки: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Не удалось рассчитать стоимость доставки",
			"details": err.Error(),
		})
		return
	}

	// Возвращаем успешный результат
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Ошибка кодирования результата: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Ошибка формирования ответа",
		})
	}
}
