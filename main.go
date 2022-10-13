package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movies struct {
	ID       int       `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	params := mux.Vars(r)
	param_id, _ := strconv.Atoi(params["id"])

	for index, movie := range movies {
		if movie.ID == param_id {
			json.NewEncoder(w).Encode(movies[index])
			return
		}
	}
	result := map[string]string{"status": fmt.Sprintf("no movie with id = %d", param_id)}
	json.NewEncoder(w).Encode(result)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	movie := Movies{}

	err := json.NewDecoder(r.Body).Decode(&movie)

	if err != nil {
		log.Fatal(err.Error())
	}

	movie.ID = rand.Intn(1000000)
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	updatedMovie := Movies{}
	err := json.NewDecoder(r.Body).Decode(&updatedMovie)

	if err != nil {
		log.Fatal(err.Error())
	}

	params := mux.Vars(r)
	param_id, _ := strconv.Atoi(params["id"])

	for index, movie := range movies {
		if param_id == movie.ID {
			movies[index] = updatedMovie
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
	result := map[string]string{"status": fmt.Sprintf("no movie with id = %d", param_id)}
	json.NewEncoder(w).Encode(result)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	param_id, _ := strconv.Atoi(params["id"])

	for index, movie := range movies {
		if param_id == movie.ID {
			movies = append(movies[:index], movies[index+1:]...)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
	result := map[string]string{"status": fmt.Sprintf("no movie with id = %d", param_id)}
	json.NewEncoder(w).Encode(result)
}

var movies []Movies

func main() {

	movies = append(movies, Movies{1, "2121", "movie satu", &Director{"Sandeep", "Sharma"}})
	movies = append(movies, Movies{2, "2212", "movie dua", &Director{"Jack", "Ma"}})

	r := mux.NewRouter()
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	port := "localhost:8000"
	log.Fatal(http.ListenAndServe(port, r))
}
