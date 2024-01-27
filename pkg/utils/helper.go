package utils

import (
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func LocalTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc)
}

func LocationTime() (*time.Location, error) {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("Error: handle sending failed: %s", err.Error())
		return nil, errors.New("error: failed to load location")
	}
	return loc, nil
}

func ConvertStringTimeToTime(t string) time.Time {
	layout := "2006-01-02 15:04:05.999 -0700 MST"
	result, err := time.Parse(layout, t)
	if err != nil {
		log.Printf("Error: Parse time failed: %s", err.Error())
	}
	return result
}

func ConvertToObjectId(id string) primitive.ObjectID {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return objectId
}
