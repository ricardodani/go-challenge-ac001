package service

import (
	"errors"
	"fmt"

	"../db"
)

type City struct {
	ID      int64   `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Borders []int64 `json:"borders,omitempty"`
}

type Cities struct {
	Cities []City `json:"cities"`
}

type Path struct {
	Path []int64 `json:"path,omitempty"`
}

func NewCity() City {
	var city City
	return city
}

func getCityBorders(city *City) error {
	var borders []int64
	stmt, _ := db.DB.Prepare("SELECT `to` FROM borders WHERE `from` = ?")
	rows, err := stmt.Query(city.ID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var border int64
		_ = rows.Scan(&border)
		borders = append(borders, border)
	}
	city.Borders = borders
	return nil
}

// Return all Cities from database
func GetCities() (Cities, error) {
	var cities Cities
	rows, err := db.DB.Query("SELECT * FROM cities")
	if err != nil {
		return cities, err
	}
	for rows.Next() {
		var city City
		_ = rows.Scan(&city.ID, &city.Name)

		err := getCityBorders(&city)
		if err != nil {
			return cities, err
		}

		cities.Cities = append(cities.Cities, city)
	}
	return cities, nil
}

// Get a `City` object from database using `id`
func GetCity(cityId int64) (City, error) {
	stmt, _ := db.DB.Prepare("SELECT * FROM cities where id = ?")
	rows, err := stmt.Query(cityId)
	var city City
	if err != nil {
		return city, err
	}
	for rows.Next() {
		_ = rows.Scan(&city.ID, &city.Name)

		err := getCityBorders(&city)
		if err != nil {
			return city, err
		}
		return city, nil
	}
	return city, errors.New("City not found")
}

// Create a `City` object on database
func CreateCity(city *City) error {
	stmt, _ := db.DB.Prepare("INSERT INTO cities(name) values (?)")
	result, err := stmt.Exec(city.Name)
	if err != nil {
		return err
	}
	city.ID, _ = result.LastInsertId()
	return nil
}

// Update a `City` object on database
func UpdateCity(city *City) error {
	stmt, _ := db.DB.Prepare("UPDATE cities SET name = ? WHERE id = ?")
	_, err := stmt.Exec(city.Name, city.ID)
	return err
}

// Insert the `Borders` of a `City` on database in the two ways
func InsertCityBorders(city *City) error {
	var borders []int64
	for _, value := range city.Borders {
		if value == city.ID {
			return errors.New("Could not set the City itself as border")
		}
		_, errInsert := db.DB.Exec(
			"INSERT OR IGNORE INTO borders(`from`, `to`) VALUES ($1, $2)",
			city.ID,
			value,
		)
		if errInsert != nil {
			return errInsert
		}
		_, errInsertInverse := db.DB.Exec(
			"INSERT OR IGNORE INTO borders(`from`, `to`) VALUES ($1, $2)",
			value,
			city.ID,
		)
		if errInsertInverse != nil {
			return errInsertInverse
		}
		borders = append(borders, value)
	}
	city.Borders = borders
	return nil
}

// Remove all `Borders` of a `City` from the database
func RemoveCityBorders(city *City) error {
	stmt, _ := db.DB.Prepare("DELETE FROM borders WHERE `from` = ? or `to` = ?")
	_, err := stmt.Exec(city.ID, city.ID)
	if err != nil {
		return err
	}
	var borders []int64
	city.Borders = borders
	return nil
}

// Remove a specific `City` given the id
func RemoveCity(cityID int64) error {
	stmt, _ := db.DB.Prepare("DELETE FROM cities WHERE id = ?")
	_, err := stmt.Exec(cityID)
	return err
}

// Remove all city and border data from database
func RemoveCities() error {
	_, err := db.DB.Exec("DELETE FROM cities")
	if err != nil {
		return err
	}
	_, errBorders := db.DB.Exec("DELETE FROM borders")
	return errBorders
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Find recursively a valid path and update the given Path pointer
func findPath(path *Path, start int64, end int64) error {
	// Warning: in progress! not working yet
	stmt, _ := db.DB.Prepare("SELECT `to` FROM borders WHERE `from` = ?")
	rows, err := stmt.Query(start)
	if err != nil {
		return err
	}
	for rows.Next() {
		path.Path = append(path.Path, start)
		fmt.Println(path.Path)
		var border int64
		_ = rows.Scan(&border)
		if border == end {
			path.Path = append(path.Path, end)
			return nil
		}
		if contains(path.Path, border) == false {
			return findPath(path, border, end)
		}
	}
	return errors.New("No path found")
}

// Return a valid Path from a City to another
func GetPath(fromId int64, toId int64) (Path, error) {
	var path Path
	fmt.Println(path.Path, fromId, toId)
	err := findPath(&path, fromId, toId)
	if err != nil {
		return path, err
	}
	return path, nil
}
