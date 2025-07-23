package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type PageData struct {
	VisitCount  int
	CurrentTime string
}

func initDB() {
	var err error
	connStr := "user=postgres dbname=counter password=1111 host=localhost sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func updateCounter() (int, error) {
	// Увеличиваем счетчик на 1 и возвращаем новое значение
	var count int
	err := db.QueryRow(`
		UPDATE page_visits 
		SET visit_count = visit_count + 1 
		WHERE id = 1 
		RETURNING visit_count
	`).Scan(&count)
	return count, err
}
func getCounter() (int, error) {
	var count int
	err := db.QueryRow("SELECT visit_count FROM page_visits WHERE id = 1").Scan(&count)
	return count, err
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Логируем входящий запрос
	log.Printf("Request: %s %s", r.Method, r.URL.Path)

	// Обрабатываем только GET-запросы к корневому пути
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// Игнорируем запросы favicon.ico
	if strings.Contains(r.URL.Path, "favicon.ico") {
		http.NotFound(w, r)
		return
	}

	count, err := updateCounter()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	currentTime := time.Now().Format("15:04")

	data := PageData{
		VisitCount:  count,
		CurrentTime: currentTime,
	}

	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Template execute error:", err)
	}
}

func main() {
	initDB()
	defer db.Close()

	// Проверяем соединение с БД
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
