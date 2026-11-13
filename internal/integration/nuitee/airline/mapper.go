package airline

import (
	"time"
	airlinedimension "travel-crawler/internal/airline"

	"travel-crawler/internal/constant"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/model"
)

func mapAirline(item Airline, now time.Time) airlinedimension.Airline {
	return airlinedimension.Airline{
		Model: model.Model{
			Metadata:  helper.ToJSONObject(item),
			CreatedAt: now,
			CreatedBy: constant.SystemUserID,
			UpdatedAt: now,
			UpdatedBy: constant.SystemUserID,
		},
		ProviderCode: constant.ProviderCodeNuitee,
		IATACode:     item.IATA,
		Name:         item.Name,
		Alliance:     item.Alliance,
		Active:       item.Active,
		Logo:         item.Logo,
	}
}

func mapAirlines(items []Airline) []airlinedimension.Airline {
	now := helper.NowUTC()

	return helper.MapSlice(items, func(item Airline) airlinedimension.Airline {
		return mapAirline(item, now)
	})
}
