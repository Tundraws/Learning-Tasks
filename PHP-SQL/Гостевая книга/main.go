package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var tmpl *template.Template

func main() {
	// Удаляем старый файл базы данных (для тестирования)
	os.Remove("guestbook.db")

	// Инициализация базы данных
	if err := InitDB(); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer CloseDB()

	// Загрузка шаблонов
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Ошибка загрузки шаблонов: %v", err)
	}

	// Настройка маршрутов
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := GetMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Messages []Message
	}{
		Messages: messages,
	}

	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	content := r.FormValue("content")

	if err := AddMessage(name, content); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
