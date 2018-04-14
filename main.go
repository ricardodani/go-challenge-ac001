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
		log.Fatal(err)
		panic(err)
	}
	router := routes.GetRouter()
	fmt.Println("Serving at port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
