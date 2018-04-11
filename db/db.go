package db

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"../types"
)

var mainDB *sql.DB

func InitDatabase(databaseName string) error {
	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		return err
	}
	mainDB = db
	mainDB.Exec("create table if not exists cities (id integer, name text)")
	mainDB.Exec("create table if not exists borders (from integer, to integer)")
	return nil
}

func GetCities() (types.Cities, error) {
	rows, err := mainDB.Query("SELECT * FROM cities")
	if err != nil {
		return nil, err
	}
	var cities types.Cities
	for rows.Next() {
		var city types.City
		_ = rows.Scan(&city.ID, &city.Name)
		cities = append(cities, city)
	}
	return cities, nil
}

func GetCity(cityId int64) (types.City, error) {
	stmt, _ := mainDB.Prepare("SELECT * FROM cities where id = ?")
	rows, err := stmt.Query(cityId)
	var city types.City
	if err != nil {
		return city, err
	}
	for rows.Next() {
		_ = rows.Scan(&city.ID, &city.Name)
	}
	return city, nil
}

func CreateCity(city *types.City) error {
	stmt, _ := mainDB.Prepare("INSERT INTO cities(name) values (?)")
	result, err := stmt.Exec(city.Name)
	if err != nil {
		return err
	}
	city.ID, _ = result.LastInsertId()
	return nil
}

func UpdateCity(city *types.City) error {
	stmt, _ := mainDB.Prepare("UPDATE cities SET name = ? WHERE id = ?")
	_, err := stmt.Exec(city.Name, city.ID)
	return err
}

func InsertCityBorders(city *types.City) error {
	var borders []int64
	for _, value := range city.Borders {
		cityIdStr := strconv.FormatInt(city.ID, 10)
		valueStr := strconv.FormatInt(value, 10)
		stmt, _ := mainDB.Prepare("INSERT INTO borders(from, to) values(?, ?)")
		_, err := stmt.Exec(cityIdStr, valueStr)
		if err != nil {
			return err
		}
		borders = append(borders, value)
	}
	city.Borders = borders
	return nil
}

func RemoveCityBorders(city *types.City) error {
	stmt, _ := mainDB.Prepare("DELETE FROM borders WHERE from = ?")
	_, err := stmt.Exec(city.ID)
	if err != nil {
		return err
	}
	var borders []int64
	city.Borders = borders
	return nil
}

func RemoveCity(cityID int64) error {
	stmt, _ := mainDB.Prepare("DELETE FROM cities WHERE id = ?")
	_, err := stmt.Exec(cityID)
	return err
}

func RemoveCities() error {
	_, err := mainDB.Exec("DELETE FROM cities")
	return err
}
