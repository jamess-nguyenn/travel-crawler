package nuitee

import (
	"time"
	"travel-crawler/internal/accommodation"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/model"
)

func mapHotel(item Hotel, now time.Time) accommodation.Accommodation {
	return accommodation.Accommodation{
		Model: model.Model{
			Metadata:  helper.ToJSONObject(item),
			CreatedAt: now,
			CreatedBy: constant.SystemUserID,
			UpdatedAt: now,
			UpdatedBy: constant.SystemUserID,
		},
		ProviderCode:     constant.ProviderCodeNuitee,
		ExternalID:       item.ID,
		ExternalParentID: item.PrimaryHotelID,
		Name:             item.Name,
		Description:      item.HotelDescription,
		PropertyType:     helper.ParseString(item.HotelTypeID),
		Chain:            item.Chain,
		Currency:         item.Currency,
		Country:          item.Country,
		City:             item.City,
		Latitude:         item.Latitude,
		Longitude:        item.Longitude,
		Address:          item.Address,
		Zip:              item.Zip,
		MainPhoto:        item.MainPhoto,
		Thumbnail:        item.Thumbnail,
		Stars:            item.Stars,
		Rating:           item.Rating,
		ReviewCount:      item.ReviewCount,
	}
}

func mapHotels(items []Hotel) []accommodation.Accommodation {
	now := helper.NowUTC()

	return helper.MapSlice(items, func(item Hotel) accommodation.Accommodation {
		return mapHotel(item, now)
	})
}
