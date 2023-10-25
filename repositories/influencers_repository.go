package repositories

import (
	"context"
	"follooow-be/configs"
	"follooow-be/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var InfluencersCollections *mongo.Collection = configs.GetCollection(configs.DB, "influencers")

// function to get detail influencer by influencer_id
// auto increase visits + 1 if data found on DB
func GetDetailInfluencers(ctx context.Context, influencer_id string) (error, models.InfluencerModel) {

	var influencer models.InfluencerModel
	objId, _ := primitive.ObjectIDFromHex(influencer_id)

	err := InfluencersCollections.FindOne(ctx, bson.M{"_id": objId}).Decode(&influencer)

	if err == nil {
		// get count data
		// filter generaor
		filterListData := bson.M{}
		var influencerIds []string
		influencerIds = append(influencerIds, influencer_id)
		filterListData["influencers"] = bson.M{"$in": influencerIds}

		countNews, _ := NewsCollections.CountDocuments(ctx, filterListData)
		countGallery, _ := GalleryCollections.CountDocuments(ctx, filterListData)

		influencer.Stats.TotalNews = int(countNews)
		influencer.Stats.TotalGallery = int(countGallery)

		// increase visits
		InfluencersCollections.UpdateOne(ctx, bson.D{{"_id", objId}}, bson.D{{"$set", bson.D{{"visits", influencer.Visits + 1}}}})
	}

	return err, influencer
}

// function to update influencers last update
// trigger when create/update post, gallery relate to influencers
func InfluencersUpdateOnToNow(ctx context.Context, influencersIds []string) error {
	var objectIds []primitive.ObjectID

	now := time.Now().UnixNano() / int64(time.Millisecond)

	// convert influencer id to objectId
	for key := range influencersIds {
		objId, _ := primitive.ObjectIDFromHex(influencersIds[key])

		objectIds = append(objectIds, objId)
	}

	// start update db
	_, err := InfluencersCollections.UpdateMany(ctx, bson.D{{"_id", bson.M{"$in": objectIds}}}, bson.D{{"$set", bson.D{{"updated_on", now}}}})

	return err
}
