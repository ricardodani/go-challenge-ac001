package main

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/ricardodani/go-challenge-ac001/db"
	"gitlab.com/ricardodani/go-challenge-ac001/routes"
)

func main() {
	err := db.InitDatabase("cities.db")
	if err != nil {
		log.Fatal("Could not init database")
		panic(err)
	}
	router := routes.GetRouter()
	fmt.Println("Serving at port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
