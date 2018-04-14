package service

import (
	"errors"

	"../db"
)

type Borders []int64

type City struct {
	ID      int64   `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Borders Borders `json:"borders,omitempty"`
}

type Cities struct {
	Cities []City `json:"cities"`
}

type Path struct {
	Path Borders `json:"path,omitempty"`
}

func intInSlice(a int64, list Borders) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getCityBorders(city *City, idsToExclude Borders) error {
	var borders Borders
	stmt, _ := db.DB.Prepare("SELECT `to` FROM borders WHERE `from` = ?")
	rows, err := stmt.Query(city.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var border int64
		_ = rows.Scan(&border)
		if !intInSlice(border, idsToExclude) {
			borders = append(borders, border)
		}
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
	defer rows.Close()
	for rows.Next() {
		var city City
		_ = rows.Scan(&city.ID, &city.Name)

		err := getCityBorders(&city, Borders{})
		if err != nil {
			return cities, err
		}

		cities.Cities = append(cities.Cities, city)
	}
	return cities, nil
}

// Get a `City` object from database using `id`
func GetCity(cityId int64) (City, error) {
	stmt, _ := db.DB.Prepare("SELECT * FROM cities WHERE `id` = ?")
	rows, err := stmt.Query(cityId)
	var city City
	if err != nil {
		return city, err
	}
	defer rows.Close()
	for rows.Next() {
		_ = rows.Scan(&city.ID, &city.Name)

		err := getCityBorders(&city, Borders{})
		if err != nil {
			return city, err
		}
		return city, nil
	}
	return city, errors.New("City not found")
}

// Insert the `Borders` of a `City` on database in the two ways
func insertCityBorders(city *City) error {
	var borders Borders
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

func validateCity(city *City) error {
	if city.Name == "" {
		return errors.New("Empty city name")
	}
	return nil
}

// Create a `City` object on database
func CreateCity(city *City) error {
	err := validateCity(city)
	if err != nil {
		return err
	}
	stmt, _ := db.DB.Prepare("INSERT INTO cities(name) values (?)")
	result, err := stmt.Exec(city.Name)
	defer stmt.Close()
	if err != nil {
		return err
	}
	city.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return insertCityBorders(city)
}

// Remove all `Borders` of a `City` from the database
func removeCityBorders(cityID int64) error {
	stmt, _ := db.DB.Prepare("DELETE FROM borders WHERE `from` = ? or `to` = ?")
	_, err := stmt.Exec(cityID, cityID)
	defer stmt.Close()
	return err
}

// Update a `City` object on database
func UpdateCity(city *City) error {
	err := validateCity(city)
	if err != nil {
		return err
	}
	stmt, _ := db.DB.Prepare("UPDATE cities SET name = ? WHERE id = ?")
	_, err = stmt.Exec(city.Name, city.ID)
	defer stmt.Close()
	if err != nil {
		return err
	}
	err = removeCityBorders(city.ID)
	if err != nil {
		return err
	}
	err = insertCityBorders(city)
	return err
}

// Remove a specific `City` given the id
func RemoveCity(cityID int64) error {
	stmt, _ := db.DB.Prepare("DELETE FROM cities WHERE id = ?")
	_, err := stmt.Exec(cityID)
	defer stmt.Close()
	if err != nil {
		return err
	}
	return removeCityBorders(cityID)
}

// Remove all city and border data from database
func RemoveCities() error {
	defer db.DB.Close()
	_, err := db.DB.Exec("DELETE FROM cities")
	if err != nil {
		return err
	}
	_, errBorders := db.DB.Exec("DELETE FROM borders")
	return errBorders
}

// Find recursively a `Path` to a city and return it when found
func findPath(path Path, fromCityId int64, toCityId int64) (Path, error) {
	path.Path = append(path.Path, fromCityId)
	fromCity := City{ID: fromCityId}
	err := getCityBorders(&fromCity, path.Path)
	if err != nil {
		return Path{}, errors.New("Could not fetch city borders")
	}

	if intInSlice(toCityId, fromCity.Borders) {
		path.Path = append(path.Path, toCityId)
		return path, nil
	}

	for _, border := range fromCity.Borders {
		if !intInSlice(border, path.Path) {
			findPath, _ := findPath(path, border, toCityId)
			if len(findPath.Path) > 0 {
				return findPath, nil
			}
		}
	}

	lenPath := len(path.Path)
	if lenPath > 0 && path.Path[lenPath-1] != toCityId {
		return Path{}, errors.New("Path not found")
	}
	return path, nil
}

// Return a valid Path from a City to another
func GetPath(fromId int64, toId int64) (Path, error) {
	return findPath(Path{}, fromId, toId)
}
