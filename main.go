package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type JSON struct {
	Results struct {
		Date    string `json:"date"`
		Sunrise string `json:"sunrise"`
		Sunset  string `json:"sunset"`
	} `json:"results"`
	Status string `json:"status"`
}

type Mode string

const (
	Color Mode = "color"
	White Mode = "white"
	Warm  Mode = "warm"
)

// TODO: Dodaj strukta dla tego co ma shelly by nie DDOS'ować żarówki

func AmIhome(w http.ResponseWriter, req *http.Request) {
	// Get data from my iPhone's request to know if I'm home
	// If no, then don;t turn on the light bulbs
}

func turnBulb(turn bool, brightness int, mode Mode) {
	resp, err := http.Get("http://192.168.100.50/color/0?turn=off")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp)
}

func checkSun() {

	var sd JSON

	for {

		if sd.Results.Date == "" && sd.Results.Sunrise == "" && sd.Results.Sunset == "" {

			resp, err := http.Get("https://api.sunrisesunset.io/json?lat=52.23547&lng=21.04191&time_format=24")
			if err != nil {
				log.Fatalln(err)
			}

			// fmt.Println(resp)

			// json.NewDecoder(resp.Body).Decode(&sd)
			body, err := io.ReadAll(resp.Body) // response body is []byte
			if err != nil {
				fmt.Printf("Error, nie udało się przechwicycić danych! %s\n", err)
			}

			fmt.Println(string(body))
			fmt.Println()

			if err := json.Unmarshal(body, &sd); err != nil { // Parse []byte to go struct pointer
				fmt.Println("Nie odmarszal JaySyna")
			}

			err = resp.Body.Close()
			if err != nil {
				fmt.Println("Nie udało się zamknąć ciała!")
			}

			fmt.Println("Sunrise:", sd)

			if sd.Status != "OK" {
				sd = JSON{}
			}

			time.Sleep(10 * time.Second)
		}

		now := time.Now()

		fmt.Println(now.Format("2006-01-02"))
		time.Sleep(10 * time.Second)

		if sd.Results.Date != now.Format("2006-01-02") {
			fmt.Println("Data się NIE zgadza!")
			sd = JSON{}
		} else {
			fmt.Println("Data się ZGADZA")
		}

		// Zachód

		sunset, err := time.Parse("15:04:05", sd.Results.Sunset)
		if err != nil {
			fmt.Println("Nie udało sie uprasować czasu")
		}

		// Wschód

		sunrise, err := time.Parse("15:04:05", sd.Results.Sunrise)
		if err != nil {
			fmt.Println("Nie udało sie uprasować czasu")
		}

		var (
			turn       bool
			brightness int
			mode       Mode
		)

		if now.Compare(sunrise) == -1 /* Jeżeli wcześniej od wschodu */ &&
			now.Compare(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 6, 0, 0, 0, time.Now().Location())) < 0 /* ale po 6 rano */ {
			fmt.Println("Włącz")
			turn = true
			brightness = 10
		}

		if now.Compare(sunrise) < 0 /* jeżeli później niż wschód */ && now.Compare(sunset) > 0 /* Ale przed zachodem */ {
			turn = false
		}

		if now.Compare(sunset) /* po zachodzie słońca */ <= 0 &&
			now.Compare(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 20, 0, 0, 0, time.Now().Location())) > 0 /* ale przed 20 */ {
			turn = true
		}

		// // Sen

		if now.Compare(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 20, 0, 0, 0, time.Now().Location())) < 0 {
			fmt.Println("Tryb cieniejszy z zółtym światłem")
			turn = true
			brightness = 50
			// resp, err := http.Get("192.168.100.50/color/0?brightness=50") JAK ZROBIĆ WARM MODE?
		}

		if now.Compare(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 22, 0, 0, 0, time.Now().Location())) < 0 {
			fmt.Println("Najciemniejszy")
			turn = true
			brightness = 1
		}

		if now.Compare(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 0, 0, 0, time.Now().Location())) < 0 {
			fmt.Println("Wyłącz")
			turn = false
			// resp, err := http.Get("http://192.168.100.50/color/0?turn=off")
		}

	}
}

func main() {

	go checkSun()

	http.HandleFunc("/zapal", turnOn)
	http.HandleFunc("/zgas", turnOff)

	fmt.Println(http.ListenAndServe(":8123", nil))
}
