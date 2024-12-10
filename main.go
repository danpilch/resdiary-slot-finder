package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ApiResponse struct {
	TimeSlots []struct {
		TimeSlot                      string  `json:"TimeSlot"`
		IsLeaveTimeRequired           bool    `json:"IsLeaveTimeRequired"`
		LeaveTime                     string  `json:"LeaveTime"`
		ServiceID                     int     `json:"ServiceId"`
		HasStandardAvailability       bool    `json:"HasStandardAvailability"`
		AvailablePromotions           []any   `json:"AvailablePromotions"`
		StandardAvailabilityFeeAmount float64 `json:"StandardAvailabilityFeeAmount"`
	} `json:"TimeSlots"`
	Promotions                               []any `json:"Promotions"`
	StandardAvailabilityMayRequireCreditCard bool  `json:"StandardAvailabilityMayRequireCreditCard"`
}

func main() {
	date := "2024-12-21"
	url := "https://booking.resdiary.com/api/Restaurant/ChesilRectory/AvailabilitySearch?date=" + date + "&covers=2&channelCode=ONLINE&areaId=0&availabilityType=Reservation"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("failed to GET url")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: Received non-OK HTTP status: %s", resp.Status)
	}

	// Parse the JSON response into the struct
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	if len(apiResponse.TimeSlots) > 0 {
		for _, slot := range apiResponse.TimeSlots {
			log.Printf("found slot: %s", slot.TimeSlot)
		}
	} else {
		fmt.Println("nothing found")
	}
}
