package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"	
	"log"
	"encoding/json"
)

type Movie struct {
	Title string `json:"title"`
	Genre string `json:"genre"`
}

var movies []Movie

func getMovies(writer http.ResponseWriter, request *http.Request){

	var new_movie Movie
	error := json.Unmarshal([]byte(`{"Title": "Creed 1", "Genre": "Drama"}`), &new_movie)
	if error != nil{
		log.Fatalln("An error has occured: ", error)
		return
	}

	error = json.NewEncoder(writer).Encode(&new_movie);
	if error != nil{
		log.Fatalln("An error occured encoding response: ", error)
		return
	}
}

func main (){

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/movies", getMovies).Methods("GET")

	fmt.Println("Starting server at port 8080: ")
	error := http.ListenAndServe(":8080", router)
	if error != nil{
		log.Fatalln("An error occured listening to server: ", error);
		return
	}
}