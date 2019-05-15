package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TcpServerException struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RevTime      string             `json:"revTime,omitempty" bson:"revTime,omitempty"`
	Pin          string             `json:"pin,omitempty" bson:"pin,omitempty"`
	ProtocolType string             `json:"protocolType,omitempty" bson:"protocolType,omitempty"`
	Bytes        []byte             `json:"bytes,omitempty" bson:"bytes,omitempty"`
	ErrorReason  string             `json:"errorReason,omitempty" bson:"errorReason,omitempty"`
}

const PROTOCOL_TYPE_SMU string = "SMU"
