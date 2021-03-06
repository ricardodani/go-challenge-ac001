package routes

import (
	"github.com/gorilla/mux"

	"../api"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/city/{id}", api.GetCity).Methods("GET")
	router.HandleFunc("/city/{id}", api.UpdateCity).Methods("PUT")
	router.HandleFunc("/city/{id}", api.RemoveCity).Methods("DELETE")
	router.HandleFunc("/city", api.CreateCity).Methods("POST")
	router.HandleFunc("/cities", api.ListCities).Methods("GET")
	router.HandleFunc("/cities", api.RemoveCities).Methods("DELETE")
	router.HandleFunc("/city/{from_id}/travel/{to_id}", api.GetPath).Methods("GET")
	return router
}
