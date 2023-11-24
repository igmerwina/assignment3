package service

import (
	"assignment3/model"
	"assignment3/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func UpdateData(c echo.Context) error {
	log.Println("--> Execute update data water & wind\n", time.Now())

	water := model.WaterStatus{
		Water: rand.Intn(15) + 1,
		Wind:  rand.Intn(15) + 1,
	}

	// Load existing data
	existingData, err := util.ReadFromJSON("data.json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load data"})
	}

	// Append new data
	existingData = append(existingData, water)

	// Save data to JSON file
	err = util.SaveToJSON("data.json", existingData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save data"})
	}

	// Unmarshal the JSON array
	var wotah []model.WaterStatus
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
