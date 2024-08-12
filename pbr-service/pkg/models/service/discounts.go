package service

import (
	"fmt"

	"github.com/Manas8803/Walmart_Project_Backend/pbr-service/pkg/models/db"
)

type Discount struct {
	BeaconID      string `json:"beacon_id"`
	DiscountOffer string `json:"discount_offer"`
}

func FetchDiscountByBeaconID(beaconID string) (*Discount, error) {
	dbDiscount, err := db.FetchDiscountByBeaconID(beaconID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch discount from database: %w", err)
	}

	if dbDiscount == nil {
		return nil, nil
	}

	discount := &Discount{
		BeaconID:      dbDiscount.BeaconID,
		DiscountOffer: dbDiscount.DiscountOffer,
	}

	return discount, nil
}
