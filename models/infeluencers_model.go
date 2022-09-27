package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InfluencerModel struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" validate:"required"`
	Name        string             `json:"name,omitempty" validate:"required"`
	Avatar      string             `json:"avatar,omitempty"`
	Bio         string             `json:"bio,omitempty"`
	UpdatedOn   int                `json:"updated_on,omitempty" bson:"updated_on,omitempty"`
	Nationality string             `json:"nationality,omitempty"`
	Gender      string             `json:"gender,omitempty"`
	Visits      int                `json:"visits,omitempty"`
	Socials     []influencerSocial `json:"socials,omitempty"`
	Label       []string           `json:"label,omitempty"`
}

type influencerSocial struct {
	Link string `json:"link,omitempty"`
	Type string `json:"type,omitempty"`
}