package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var wg sync.WaitGroup

type WaterStatus struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func main() {
	e := echo.New()

	e.GET("/update", UpdateData)

	go makeRequests()

	e.Logger.Fatal(e.Start(":8888"))
}

func makeRequests() {
	defer wg.Done()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Make a request to your own endpoint
			resp, err := http.Get("http://localhost:8888/update")
			if err != nil {
				fmt.Println("Error making request:", err)
				continue
			}
			defer resp.Body.Close()

			fmt.Println()
		}
	}
}

func UpdateData(c echo.Context) error {
	log.Println("Execute update data water & wind", time.Now())
	water := WaterStatus{
		Water: rand.Intn(15) + 1,
		Wind:  rand.Intn(15) + 1,
	}

	// Load existing data
	existingData, err := readFromJSON("data.json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load data"})
	}

	// Append new data
	existingData = append(existingData, water)

	// Save data to JSON file
	err = saveToJSON("data.json", existingData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save data"})
	}

	// Unmarshal the JSON array
	var wotah []WaterStatus
	fileContent, err := ioutil.ReadFile("data.json")
	err = json.Unmarshal(fileContent, &wotah)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save data"})
	}

	// Get the last element in the array
	var statusWater string = ""
	var statusWind string = ""

	if len(wotah) > 0 {
		lastStatus := wotah[len(wotah)-1]
		fmt.Println("Last Status", lastStatus)
		switch {
		case lastStatus.Water < 5:
			statusWater = "Aman"
		case lastStatus.Water >= 6 && lastStatus.Water <= 8:
			statusWater = "Siaga"
		case lastStatus.Water > 8:
			statusWater = "Bahaya"
		}

		switch {
		case lastStatus.Wind < 6:
			statusWind = "Aman"
		case lastStatus.Wind >= 7 && lastStatus.Wind <= 15:
			statusWind = "Siaga"
		case lastStatus.Wind > 15:
			statusWind = "Bahaya"
		}

		fmt.Printf("Water: %#v", lastStatus.Water)
		fmt.Printf(" - Status: %s\n", statusWater)
		fmt.Printf("Wind: %#v", lastStatus.Wind)
		fmt.Printf(" - Status: %s\n", statusWind)
	} else {
		fmt.Println("No data in the JSON file.")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Data saved successfully "})
}

func saveToJSON(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func readFromJSON(filename string) ([]WaterStatus, error) {
	var data []WaterStatus

	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
