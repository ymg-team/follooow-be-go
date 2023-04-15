package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InfluencerModel struct {
	Id          primitive.ObjectID      `json:"id,omitempty" bson:"_id,omitempty" validate:"required"`
	Name        string                  `json:"name,omitempty" validate:"required"`
	Avatar      string                  `json:"avatar,omitempty"`
	Bio         string                  `json:"bio,omitempty"`
	UpdatedOn   int                     `json:"updated_on,omitempty" bson:"updated_on,omitempty"`
	Nationality string                  `json:"nationality,omitempty"`
	Gender      string                  `json:"gender,omitempty"`
	Visits      int                     `json:"visits,omitempty"`
	Socials     []InfluencerSocial      `json:"socials,omitempty"`
	Label       []string                `json:"label,omitempty"`
	Views       int                     `json:"views,omitempty"`
	Code        string                  `json:"code,omitempty"`
	BestMoments []InfluencerBestMoments `json:"best_moments,omitempty" bson:"best_moments,omitempty"`
}

type InsertInfluencerModel struct {
	Name        string                  `json:"name,omitempty" validate:"required"`
	Avatar      string                  `json:"avatar,omitempty"`
	Bio         string                  `json:"bio,omitempty"`
	UpdatedOn   int                     `json:"updated_on,omitempty" bson:"updated_on,omitempty"`
	Nationality string                  `json:"nationality,omitempty"`
	Gender      string                  `json:"gender,omitempty"`
	Visits      int                     `json:"visits,omitempty"`
	Socials     []InfluencerSocial      `json:"socials,omitempty"`
	Label       []string                `json:"label,omitempty"`
	Views       int                     `json:"views,omitempty"`
	Code        string                  `json:"code,omitempty"`
	BestMoments []InfluencerBestMoments `json:"best_moments,omitempty" bson:"best_moments,omitempty"`
}

type InfluencerSmallDataModel struct {
	Id     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" validate:"required"`
	Name   string             `json:"name,omitempty" validate:"required"`
	Avatar string             `json:"avatar,omitempty"`
}

type InfluencerSocial struct {
	Link  string `json:"link,omitempty"`
	Type  string `json:"type,omitempty"`
	Title string `json:"title,omitempty"`
}

type InfluencerBestMoments struct {
	Image      string                     `json:"image,omitempty"`
	Text       string                     `json:"text,omitempty"`
	Year       string                     `json:"year,omitempty"`
	Background string                     `json:"background,omitempty"`
	Style      InfluencerBestMomentsStyle `json:"style,omitempty"`
}

type InfluencerBestMomentsStyle struct {
	Margin string `json:"margin,omitempty"`
}
