package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type OriginalPacket struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Protocol    string             `json:"protocol,omitempty"`
	Ip          string             `json:"ip,omitempty"`
	RevTime     string             `json:"revTime,omitempty"`
	Bytes       []byte             `json:"bytes,omitempty"`
	BytesLength int                `json:"bytesLength,omitempty" bson:"bytesLength,omitempty"`
	Pin         string             `json:"pin,omitempty"`

	Source string `json:"source,omitempty"`
}
