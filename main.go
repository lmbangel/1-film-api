package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"	
	"log"
	"encoding/json"
)

type Movie struct {
	Title 	string `json:"title"`
	Genre 	string `json:"genre"`
	ID 		string `json:"id"` 
}

var movies []Movie

func getMovies(writer http.ResponseWriter, request *http.Request){

	if len(movies) == 0 {		
		var new_movie Movie
		error := json.Unmarshal([]byte(`{"Title": "Creed 1", "Genre": "Drama", "ID": "1"}`), &new_movie)
		if error != nil{
			log.Fatalln("An error has occured: ", error)
			return
		}
		movies = append(movies, new_movie)
	}
	error := json.NewEncoder(writer).Encode(&movies)
	if error != nil{
		log.Fatalln("An error occured encoding response: ", error)
		return
	}
}

func getMovie(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	params := mux.Vars(request)
	id := params["id"] 
	for _, movie := range movies {
		if movie.ID == id {
			json.NewEncoder(writer).Encode(movie)
			return
		}
	}
}

func createMovie(writer http.ResponseWriter, request *http.Request){
	
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	var new_movie Movie
	error := json.NewDecoder(request.Body).Decode(&new_movie)
	if error != nil{
		log.Fatalln("An error occured, movie not created:  ", error)
		return 
	}
	movies = append(movies, new_movie);
	error = json.NewEncoder(writer).Encode(&movies)
	if error != nil{
		log.Fatalln("an error occured: ", error)
		return 
	}

}

/* Main ................................................*/

func main (){

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/movies", createMovie).Methods("POST")
	router.HandleFunc("/api/v1/movies", getMovies).Methods("GET")
	router.HandleFunc("/api/v1/movies/{id}", getMovie).Methods("GET")

	fmt.Println("Starting server at port 8080: ")
	error := http.ListenAndServe(":8080", router)
	if error != nil {
		log.Fatalln("An error occured listening to server: ", error);
		return
	}
}