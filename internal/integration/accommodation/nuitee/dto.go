package nuitee

import "time"

type AccessibilityAttributes struct {
	Attributes                                 []any   `json:"attributes"`
	ShowerChair                                *bool   `json:"showerChair"`
	EntranceType                               *string `json:"entranceType"`
	PetFriendly                                string  `json:"petFriendly"`
	RampAngle                                  float64 `json:"rampAngle"`
	RampLength                                 float64 `json:"rampLength"`
	EntranceDoorWidth                          float64 `json:"entranceDoorWidth"`
	RoomMaxGuestsNumber                        int     `json:"roomMaxGuestsNumber"`
	DistanceFromTheElevatorToTheAccessibleRoom float64 `json:"distanceFromTheElevatorToTheAccessibleRoom"`
}

type Hotel struct {
	ID                      string                  `json:"id"`
	PrimaryHotelID          *string                 `json:"primaryHotelId"`
	Name                    string                  `json:"name"`
	HotelDescription        string                  `json:"hotelDescription"`
	HotelTypeID             int                     `json:"hotelTypeId"`
	Chain                   string                  `json:"chain"`
	Currency                string                  `json:"currency"`
	Country                 string                  `json:"country"`
	City                    string                  `json:"city"`
	Latitude                *float64                `json:"latitude"`
	Longitude               *float64                `json:"longitude"`
	Address                 string                  `json:"address"`
	Zip                     string                  `json:"zip"`
	MainPhoto               *string                 `json:"main_photo"`
	Thumbnail               *string                 `json:"thumbnail"`
	Stars                   float64                 `json:"stars"`
	Rating                  *float64                `json:"rating"`
	ReviewCount             uint32                  `json:"reviewCount"`
	FacilityIDs             []int                   `json:"facilityIds"`
	AccessibilityAttributes AccessibilityAttributes `json:"accessibilityAttributes"`
	DeletedAt               *time.Time              `json:"deletedAt"`
}

type HotelResponse struct {
	Data  []Hotel `json:"data"`
	Total int     `json:"total"`
}
