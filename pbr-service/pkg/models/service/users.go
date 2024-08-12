package service

import (
	"github.com/Manas8803/Walmart_Project_Backend/pbr-service/pkg/models/db"
)

type User struct {
	UserID        string   `json:"user_id"`
	Email         string   `json:"email"`
	ConnectionIDs []string `json:"connection_ids"`
}

func FetchUserByID(userID string) (*User, error) {
	dbUser, err := db.FetchUserByID(userID)
	if err != nil {
		return nil, err
	}

	if dbUser == nil {
		return nil, nil
	}

	user := &User{
		UserID:        dbUser.UserID,
		Email:         dbUser.Email,
		ConnectionIDs: dbUser.ConnectionIDs,
	}

	return user, nil
}
