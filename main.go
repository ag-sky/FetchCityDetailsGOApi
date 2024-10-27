package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
)

type City struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Population int    `json:"population"`
    Area     float64 `json:"area"`
    Country  string  `json:"country"`
}

var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("mysql", "ag_sky:Radheybol@tcp(127.0.0.1:3306)/cityDetailDB")
    if err != nil {
        log.Fatal(err)
    }
}

func getCities(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT * FROM cityDetails")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    cities := []City{}
    for rows.Next() {
        var city City
        if err := rows.Scan(&city.ID, &city.Name, &city.Population, &city.Area, &city.Country); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        cities = append(cities, city)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cities)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/cities", getCities).Methods("GET")

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}
