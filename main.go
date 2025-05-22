package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gregdel/pushover"
	"github.com/sethvargo/go-envconfig"
)

// resdiary doesn't use RFC3339 time so we need to customise
type CustomTime time.Time

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Trim quotes from the string
	s := string(b)
	s = s[1 : len(s)-1] // Remove the surrounding quotes

	// Parse the time in the expected format
	parsedTime, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}

	*ct = CustomTime(parsedTime)
	return nil
}

// Convert back to time.Time when needed
func (ct CustomTime) ToTime() time.Time {
	return time.Time(ct)
}

func (ct CustomTime) ToString() string {
	return time.Time(ct).Format("2006-01-02T15:04:05")
}

type ApiResponse struct {
	TimeSlots []struct {
		TimeSlot                      CustomTime `json:"TimeSlot"`
		IsLeaveTimeRequired           bool       `json:"IsLeaveTimeRequired"`
		LeaveTime                     string     `json:"LeaveTime"`
		ServiceID                     int        `json:"ServiceId"`
		HasStandardAvailability       bool       `json:"HasStandardAvailability"`
		AvailablePromotions           []any      `json:"AvailablePromotions"`
		StandardAvailabilityFeeAmount float64    `json:"StandardAvailabilityFeeAmount"`
	} `json:"TimeSlots"`
	Promotions                               []any `json:"Promotions"`
	StandardAvailabilityMayRequireCreditCard bool  `json:"StandardAvailabilityMayRequireCreditCard"`
}

type Options struct {
	DisablePushover                  bool   `env:"DISABLE_PUSHOVER,default=false"`
	PushoverApiKey                   string `env:"PUSHOVER_API_KEY,required"`
	PushoverRecipient                string `env:"PUSHOVER_RECIPIENT,required"`
	ReservationDate                  string `env:"RESERVATION_DATE,default=2024-12-21"`
	RestaurantNames                  string `env:"RESTAURANT_NAMES,default=ChesilRectory"`
	RestaurantCovers                 string `env:"RESTAURANT_COVERS,default=2"`
	ReservationIgnoreThresholdHour   int    `env:"RESERVATION_IGNORE_THRESHOLD_HOUR,default=21"` // Ignore slots that are after 21hours
	ReservationIgnoreThresholdMinute int    `env:"RESERVATION_IGNORE_THRESHOLD_HOUR,default=0"`  // Ignore slots that are after 21hours
}

func checkSlotIsValid(slot time.Time, cutoffHour int, cutoffMinute int) bool {
	cutoff := time.Date(
		slot.Year(), slot.Month(), slot.Day(),
		cutoffHour, cutoffMinute, 0, 0, slot.Location(),
	)

	return slot.Before(cutoff)
}

func main() {
	ctx := context.Background()
	var o Options
	if err := envconfig.Process(ctx, &o); err != nil {
		log.Fatal(err)
	}

	app := pushover.New(o.PushoverApiKey)
	recipient := pushover.NewRecipient(o.PushoverRecipient)
	restaurants := strings.Split(o.RestaurantNames, ",")

	for _, restaurantName := range restaurants {
		trimmedRestaurantName := strings.TrimSpace(restaurantName)
		log.Printf("Checking %s", trimmedRestaurantName)
		url := fmt.Sprintf("https://booking.resdiary.com/api/Restaurant/"+
			"%s/AvailabilitySearch?date=%s&covers=%s"+
			"&channelCode=ONLINE&areaId=0&availabilityType=Reservation", trimmedRestaurantName, o.ReservationDate, o.RestaurantCovers)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error: failed to GET url %s: %v", url, err)
			continue
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body for restaurant %s: %v", trimmedRestaurantName, err)
			continue
		}

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error: Received non-OK HTTP status: %s for restaurant %s", resp.Status, trimmedRestaurantName)
			continue
		}

		// Parse the JSON response into the struct
		var apiResponse ApiResponse
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			log.Printf("Error parsing JSON for restaurant %s: %v", trimmedRestaurantName, err)
			continue
		}

		if len(apiResponse.TimeSlots) > 0 {
			for _, slot := range apiResponse.TimeSlots {
				if checkSlotIsValid(slot.TimeSlot.ToTime(), o.ReservationIgnoreThresholdHour, o.ReservationIgnoreThresholdMinute) {
					log.Printf("found slot at %s: %s", trimmedRestaurantName, slot.TimeSlot.ToString())
					if !o.DisablePushover {
						message := pushover.NewMessage(fmt.Sprintf("found slot at %s: %s", trimmedRestaurantName, slot.TimeSlot.ToString()))
						_, err := app.SendMessage(message, recipient)
						if err != nil {
							fmt.Println(err)
						}
					}
				} else {
					log.Printf("Unacceptable timeslot found at %s: %s", trimmedRestaurantName, slot.TimeSlot.ToString())
				}
			}
		} else {
			log.Printf("no slots found at %s for %s", trimmedRestaurantName, o.ReservationDate)
		}
	}
}
