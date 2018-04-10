package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type City struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Borders []int64 `json:"borders,omitempty"`
}

type Cities []City

var mainDB *sql.DB

func initDatabase(databaseName string) {
	db, errOpenDB := sql.Open("sqlite3", databaseName)
	checkErr(errOpenDB)
	mainDB = db
	mainDB.Exec("create table if not exists cities (id integer, name text)")
	mainDB.Exec("create table if not exists borders (from integer, to integer)")
}

func main() {

	initDatabase("cities.db")

	router := mux.NewRouter()
	router.HandleFunc("/city/{id}", GetCity).Methods("GET")
	//router.HandleFunc("/city/{id}", UpdateCity).Methods("PUT")
	//router.HandleFunc("/city/{id}", RemoveCity).Methods("DELETE")
	router.HandleFunc("/city", CreateCity).Methods("POST")
	router.HandleFunc("/cities", ListCities).Methods("GET")
	//router.HandleFunc("/cities", RemoveCities).Methods("DELETE")
	//router.HandleFunc("/city/{from_id}/travel/{to_id}", GetPath).Methods("GET")
	fmt.Println("Serving at port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func ListCities(w http.ResponseWriter, r *http.Request) {
	rows, err := mainDB.Query("SELECT * FROM cities")
	checkErr(err)
	var cities Cities
	for rows.Next() {
		var city City
		_ = rows.Scan(&city.ID, &city.Name)
		cities = append(cities, city)
	}
	json.NewEncoder(w).Encode(cities)
}

func GetCity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, _ := mainDB.Prepare("SELECT * FROM cities where id = ?")
	rows, errQuery := stmt.Query(params["id"])
	checkErr(errQuery)
	var city City
	for rows.Next() {
		_ = rows.Scan(&city.ID, &city.Name)
	}
	json.NewEncoder(w).Encode(city)
}

func CreateCity(w http.ResponseWriter, r *http.Request) {
	var city City
	_ = json.NewDecoder(r.Body).Decode(&city)

	stmt, _ := mainDB.Prepare("INSERT INTO cities(name) values (?)")
	result, errExec := stmt.Exec(city.Name)
	checkErr(errExec)

	newID, errLast := result.LastInsertId()
	checkErr(errLast)
	city.ID = newID

	for _, value := range city.Borders {
		cityIdStr := strconv.FormatInt(city.ID, 10)
		valueStr := strconv.FormatInt(value, 10)
		borders, _ := mainDB.Prepare("INSERT INTO borders(from, to) values(" + cityIdStr + ",?)")
		result, err := borders.Exec(valueStr)
		checkErr(err)
		fmt.Println(result)
	}

	json.NewEncoder(w).Encode(city)
}

//
//func insert(w http.ResponseWriter, r *http.Request) {
//	name := r.FormValue("name")
//	var todo Todo
//	todo.Name = name
//	stmt, err := mainDB.Prepare("INSERT INTO todos(name) values (?)")
//	checkErr(err)
//	result, errExec := stmt.Exec(todo.Name)
//	checkErr(errExec)
//	newID, errLast := result.LastInsertId()
//	checkErr(errLast)
//	todo.ID = newID
//	jsonB, errMarshal := json.Marshal(todo)
//	checkErr(errMarshal)
//	fmt.Fprintf(w, "%s", string(jsonB))
//}
//
//func updateByID(w http.ResponseWriter, r *http.Request) {
//	name := r.FormValue("name")
//	id := r.URL.Query().Get(":id")
//	var todo Todo
//	ID, _ := strconv.ParseInt(id, 10, 0)
//	todo.ID = ID
//	todo.Name = name
//	stmt, err := mainDB.Prepare("UPDATE todos SET name = ? WHERE id = ?")
//	checkErr(err)
//	result, errExec := stmt.Exec(todo.Name, todo.ID)
//	checkErr(errExec)
//	rowAffected, errLast := result.RowsAffected()
//	checkErr(errLast)
//	if rowAffected > 0 {
//		jsonB, errMarshal := json.Marshal(todo)
//		checkErr(errMarshal)
//		fmt.Fprintf(w, "%s", string(jsonB))
//	} else {
//		fmt.Fprintf(w, "{row_affected=%d}", rowAffected)
//	}
//
//}
//
//func deleteByID(w http.ResponseWriter, r *http.Request) {
//	id := r.URL.Query().Get(":id")
//	stmt, err := mainDB.Prepare("DELETE FROM todos WHERE id = ?")
//	checkErr(err)
//	result, errExec := stmt.Exec(id)
//	checkErr(errExec)
//	rowAffected, errRow := result.RowsAffected()
//	checkErr(errRow)
//	fmt.Fprintf(w, "{row_affected=%d}", rowAffected)
//}

func checkErr(err error) {
	if err != nil {
		fmt.Println("error found")
		panic(err)
	}
}
