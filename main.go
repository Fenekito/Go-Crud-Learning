package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gofor-little/env"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movies []Movie

	rows, err := db.Query("Select * from movie")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description); err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}

	json.NewEncoder(w).Encode(movies)
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	var movie Movie
	if err := db.QueryRow("Select * FROM movie WHERE id = ?", id).Scan(&movie.ID, &movie.Title, &movie.Description); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(movie)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	_, err := db.Exec("UPDATE movie SET title = ?, description = ? WHERE id = ?", movie.Title, movie.Description, id)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(movie)
}

func DeleteMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	result, err := db.Exec("DELETE FROM movie WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(result)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	_, err := db.Exec("INSERT INTO movie(title, description) VALUES (?, ?)", movie.Title, movie.Description)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(movie)
}

var db *sql.DB

func main() {
	db = ConnectDB()
	defer db.Close()
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/movies", GetMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", GetMovie).Methods("GET")
	r.HandleFunc("/movies", CreateMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", UpdateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", DeleteMovies).Methods("DELETE")
	log.Println("Server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func ConnectDB() *sql.DB {
	if err := env.Load(".env"); err != nil {
		log.Fatal(err)
	}

	cfg := mysql.Config{
		User:                 env.Get("DBUSER", ""),
		Passwd:               env.Get("DBPASS", ""),
		Net:                  "tcp",
		Addr:                 env.Get("ADDR", ""),
		DBName:               env.Get("DATABASE", ""),
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Database Connected Successfully!")
	return db
}
