package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gregdel/pushover"
	"github.com/sethvargo/go-envconfig"
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

type Options struct {
	DisablePushover   bool   `env:"DISABLE_PUSHOVER,default=false"`
	PushoverApiKey    string `env:"PUSHOVER_API_KEY,required"`
	PushoverRecipient string `env:"PUSHOVER_RECIPIENT,required"`
	ReservationDate   string `env:"RESERVATION_DATE,default=2024-12-21"`
	RestaurantName    string `env:"RESTAURANT_NAME,default=ChesilRectory"`
	RestaurantCovers  string `env:"RESTAURANT_COVERS,default=2"`
}

func main() {
	ctx := context.Background()
	var o Options
	if err := envconfig.Process(ctx, &o); err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("https://booking.resdiary.com/api/Restaurant/"+
		"%s/AvailabilitySearch?date=%s&covers=%s"+
		"&channelCode=ONLINE&areaId=0&availabilityType=Reservation", o.RestaurantName, o.ReservationDate, o.RestaurantCovers)

	app := pushover.New(o.PushoverApiKey)
	recipient := pushover.NewRecipient(o.PushoverRecipient)

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
			if !o.DisablePushover {
				message := pushover.NewMessage(fmt.Sprintf("found slot: %s", slot.TimeSlot))
				_, err := app.SendMessage(message, recipient)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	} else {
		log.Printf("no slots found at %s for %s", o.RestaurantName, o.ReservationDate)
	}
}
