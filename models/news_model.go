package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsModel struct {
	Id              primitive.ObjectID         `json:"id,omitempty" bson:"_id,omitempty" validate:"required"`
	Title           string                     `json:"title,omitempty" validate:"required"`
	Content         string                     `json:"content,omitempty" validate:"required"`
	Views           int                        `json:"views,omitempty" validate:"required"`
	CreatedOn       int                        `json:"created_on,omitempty" bson:"created_on,omitempty" validate:"required"`
	UpdatedOn       int                        `json:"updated_on,omitempty" bson:"updated_on,omitempty" validate:"required"`
	Tags            []string                   `json:"tags,omitempty" validate:"required"`
	Influencers     []string                   `json:"influencers,omitempty"`
	InfluencersData []InfluencerSmallDataModel `json:"influencers_data,omitempty"`
}
