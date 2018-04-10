package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"gitlab.com/ricardodani/go-challenge-ac001/db"
	"gitlab.com/ricardodani/go-challenge-ac001/types"
)

func ListCities(w http.ResponseWriter, r *http.Request) {
	cities, err := db.GetCities()
	checkErr(err)
	json.NewEncoder(w).Encode(cities)
}

func GetCity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cityID, _ := strconv.ParseInt(params["id"], 10, 0)
	city, err := db.GetCity(cityID)
	checkErr(err)
	json.NewEncoder(w).Encode(city)
}

func CreateCity(w http.ResponseWriter, r *http.Request) {
	var city types.City
	_ = json.NewDecoder(r.Body).Decode(&city)
	checkErr(db.CreateCity(&city))
	checkErr(db.InsertCityBorders(&city))
	json.NewEncoder(w).Encode(city)
	// TODO: return proper status code
	// TODO: use permanent redirect to GetCity fund, change location
}

func UpdateCity(w http.ResponseWriter, r *http.Request) {
	var city types.City
	_ = json.NewDecoder(r.Body).Decode(&city)

	params := mux.Vars(r)
	city.ID, _ = strconv.ParseInt(params["id"], 10, 0)

	checkErr(db.UpdateCity(&city))
	checkErr(db.RemoveCityBorders(&city))
	checkErr(db.InsertCityBorders(&city))

	json.NewEncoder(w).Encode(city)
	// TODO: check rows affected result.RowsAffected()
	// TODO: return proper status code
}

func RemoveCity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cityID, err := strconv.ParseInt(params["id"], 10, 0)
	checkErr(err)
	checkErr(db.RemoveCity(cityID))
	json.NewEncoder(w).Encode(nil)
	// TODO: check rows affected result.RowsAffected()
	// TODO: return proper status code
}

func RemoveCities(w http.ResponseWriter, r *http.Request) {
	checkErr(db.RemoveCities())
	json.NewEncoder(w).Encode(nil)
	// TODO: check rows affected result.RowsAffected()
	// TODO: return proper status code
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
