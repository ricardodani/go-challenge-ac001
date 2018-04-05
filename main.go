package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type City struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Borders []int  `json:"borders"`
}

var cities []City

func GetCity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err == nil {
		for _, city := range cities {
			if city.ID == id {
				json.NewEncoder(w).Encode(city)
				break
			}
		}
	}
}

func UpdateCity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err == nil {
		for _, city := range cities {
			if city.ID == id {
				json.NewEncoder(w).Encode(city)
				break
			}
		}
	}
}

func RemoveCity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err == nil {
		for index, city := range cities {
			if city.ID == id {
				cities = append(cities[:index], cities[index+1:]...)
				break
			}
		}
	}
}

func CreateCity(w http.ResponseWriter, r *http.Request) {
	var city City
	_ = json.NewDecoder(r.Body).Decode(&city)
	cities = append(cities, city)
	json.NewEncoder(w).Encode(city)
}

func ListCities(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(cities)
}

func RemoveCities(w http.ResponseWriter, r *http.Request) {
	cities = append(cities[:0])
}

func GetPath(w http.ResponseWriter, r *http.Request) {}

func add_cities() {
	cities = append(cities, City{ID: 1, Name: "City 1", Borders: []int{3}})
	cities = append(cities, City{ID: 2, Name: "City 2", Borders: []int{4}})
	cities = append(cities, City{ID: 3, Name: "City 3", Borders: []int{1, 4}})
	cities = append(cities, City{ID: 4, Name: "City 4", Borders: []int{2, 3}})
}

func route_views() {
	router := mux.NewRouter()
	router.HandleFunc("/city/{id}", GetCity).Methods("GET")
	router.HandleFunc("/city/{id}", UpdateCity).Methods("PUT")
	router.HandleFunc("/city/{id}", RemoveCity).Methods("DELETE")
	router.HandleFunc("/city/{id}", CreateCity).Methods("POST")
	router.HandleFunc("/cities", ListCities).Methods("GET")
	router.HandleFunc("/cities", RemoveCities).Methods("DELETE")
	router.HandleFunc("/city/{from_id}/travel/{to_id}", GetPath).Methods("GET")
	fmt.Println("Serving at port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func main() {
	add_cities()
	route_views()
}
