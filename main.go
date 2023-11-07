package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Structure
type Movie struct {
	Id       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
}

var movies []Movie

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{Id: "1", Isbn: "123456", Title: "Learning Go", Director: &Director{FirstName: "Carlos", LastName: "Salazar", DateOfBirth: "25-06-1990"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Length of movies array: %d\n", len(movies))
	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range movies {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	//movie.id = strconv.Itoa(rand.Intn(10000)) // CREATE A NEW ID FOR MOVIE WITH A RANDOM NUMBER
	movie.Id = createIdForMovie()
	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movie)
}

func createIdForMovie() string {
	lastId := 0

	for _, item := range movies {
		itemId, _ := strconv.Atoi(item.Id)
		if itemId > lastId {
			lastId = itemId
		}
	}

	return strconv.Itoa(lastId + 1)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.Id == params["id"] {
			movies = deleteFromSlice(movies, index)    // we delete the movie with the index
			var movie Movie                            // Then create a new movie with the same id
			_ = json.NewDecoder(r.Body).Decode(&movie) // decode the body of the request and assign it to the new movie
			movie.Id = item.Id                         // then assign the id of the movie to be updated
			movies = append(movies, movie)             // finally we append it to the movies slice
			json.NewEncoder(w).Encode(movie)           // and return the movie
			return
		}
	}
}

func deleteFromSlice(slice []Movie, index int) []Movie {
	return append(slice[:index], slice[index+1:]...) // Taking the slice of the array with movies[:index] and appending the rest of the array with movies[index+1:]
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.Id == params["id"] {
			movies = deleteFromSlice(movies, index)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}
