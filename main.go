package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"	
	"log"
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/joho/godotenv"
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
		env_default_movie, exists := os.LookupEnv("DEFAULT_MOVIE")
		if exists == true {
			error := json.Unmarshal([]byte(env_default_movie), &new_movie)
			if error != nil{
				log.Fatalln("An error has occured: ", error)
				return
			}
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


func createJsonFile(writer http.ResponseWriter, request *http.Request){

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	file , _ := json.MarshalIndent(movies, "", " ")
	_ = ioutil.WriteFile("test_file.json", file, 0644)
}

func init(){
	error := godotenv.Load()
	if error != nil{
		log.Print(".env file not found.")
	}
}

/* Main ................................................*/

func main (){

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/movies", createMovie).Methods("POST")
	router.HandleFunc("/api/v1/movies", getMovies).Methods("GET")
	router.HandleFunc("/api/v1/movies/{id}", getMovie).Methods("GET")

	router.HandleFunc("/api/v1/write/json", createJsonFile).Methods("GET")

	router.HandleFunc("/api/v1/readfiile", func(writer http.ResponseWriter, request *http.Request){
		
		file,_ := ioutil.ReadFile("test_file.json")
		var movieList []Movie
		_ = json.Unmarshal([]byte(file), &movieList)
		_ = json.NewEncoder(writer).Encode(&movieList)

	}).Methods("GET")

	fmt.Println("Starting server at port 8080: ")
	env_default_movie, exists := os.LookupEnv("DEFAULT_MOVIE")
	if exists == true {
		fmt.Println(env_default_movie)
	}
	error := http.ListenAndServe(":8080", router)
	if error != nil {
		log.Fatalln("An error occured listening to server: ", error);
		return
	}
}