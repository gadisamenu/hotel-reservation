package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId    primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	RoomId    primitive.ObjectID `bson:"roomId,omitempty" json:"roomId,omitempty"`
	FromDate  time.Time          `bson:"fromDate" json:"fromDate"`
	ToDate    time.Time          `bson:"toDate" json:"toDate"`
	NumPerson int                `bson:"numPerson" json:"numPerson"`
	Canceled  bool               `bson:"canceled" json:"canceled"`
}
