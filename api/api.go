package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"../service"
)

func checkErr(err error, statusCode int, w http.ResponseWriter) bool {
	if err != nil {
		log.Println(err)
		w.WriteHeader(statusCode)
		return true
	}
	return false
}

func ListCities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cities, err := service.GetCities()
	if checkErr(err, http.StatusInternalServerError, w) {
		return
	}

	err = json.NewEncoder(w).Encode(cities)
	checkErr(err, http.StatusInternalServerError, w)
}

func GetCity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	cityID, err := strconv.ParseInt(params["id"], 10, 0)
	if checkErr(err, http.StatusBadRequest, w) {
		return
	}

	city, err := service.GetCity(cityID)
	if checkErr(err, http.StatusNotFound, w) {
		return
	}

	err = json.NewEncoder(w).Encode(city)
	checkErr(err, http.StatusInternalServerError, w)
}

func CreateCity(w http.ResponseWriter, r *http.Request) {
	city := service.NewCity()
	err := json.NewDecoder(r.Body).Decode(&city)
	if checkErr(err, http.StatusBadRequest, w) {
		return
	}

	err = service.CreateCity(&city)
	if checkErr(err, http.StatusInternalServerError, w) {
		return
	}

	err = service.InsertCityBorders(&city)
	if checkErr(err, http.StatusInternalServerError, w) {
		return
	}

	err = json.NewEncoder(w).Encode(city)
	if checkErr(err, http.StatusInternalServerError, w) {
		return
	}

	cityURL := fmt.Sprintf("/city/%d", city.ID)
	w.Header().Set("Location", cityURL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func UpdateCity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	city := service.NewCity()
	err := json.NewDecoder(r.Body).Decode(&city)
	if checkErr(err, http.StatusBadRequest, w) {
		return
	}

	params := mux.Vars(r)
	city.ID, err = strconv.ParseInt(params["id"], 10, 0)
	if checkErr(err, http.StatusBadRequest, w) {
		return
	}

	err = service.UpdateCity(&city)
	if checkErr(err, http.StatusInternalServerError, w) {
		return
	}

	err = service.RemoveCityBorders(&city)
	if checkErr(err, http.StatusInternalServerError, w) {
		return
	}

	err = service.InsertCityBorders(&city)
	if checkErr(err, http.StatusInternalServerError, w) {
		return
	}

	err = json.NewEncoder(w).Encode(city)
	checkErr(err, http.StatusInternalServerError, w)
}

func setJsonContentType(w *http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func RemoveCity(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	setJsonContentType(&w)

	params := mux.Vars(r)
	cityID, err := strconv.ParseInt(params["id"], 10, 0)
	if checkErr(err, http.StatusBadRequest, w) {
		return
	}

	err = service.RemoveCity(cityID)
	if checkErr(err, http.StatusInternalServerError, w) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func RemoveCities(w http.ResponseWriter, r *http.Request) {
	setJsonContentType(&w)

	err := service.RemoveCities()
	if checkErr(err, http.StatusInternalServerError, w) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetPath(w http.ResponseWriter, r *http.Request) {
	setJsonContentType(&w)

	params := mux.Vars(r)

	fromID, err := strconv.ParseInt(params["from_id"], 10, 0)
	if checkErr(err, http.StatusBadRequest, w) {
		return
	}

	toID, err := strconv.ParseInt(params["to_id"], 10, 0)
	if checkErr(err, http.StatusBadRequest, w) {
		return
	}

	path, err := service.GetPath(fromId, toId)
	if checkErr(err, http.StatusInternalServerError, w) {
		return
	}

	err = json.NewEncoder(w).Encode(path)
	checkErr(err, http.StatusInternalServerError, w)
}
