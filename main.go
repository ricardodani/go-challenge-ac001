package main

import (
	"fmt"
	"log"
	"net/http"

	"./db"
	"./routes"
)

func main() {
	err := db.InitDatabase("cities.db")
	if err != nil {
		log.Fatal("Could not init database")
		panic(err)
	}
	router := routes.GetRouter()
	fmt.Println("Serving at port 3001")
	log.Fatal(http.ListenAndServe(":3001", router))
}
