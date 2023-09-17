package repositories

import (
	"context"
	"follooow-be/configs"
	"follooow-be/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var influencersCollections *mongo.Collection = configs.GetCollection(configs.DB, "influencers")

// function to get detail influencer by influencer_id
// auto increase visits + 1 if data found on DB
func GetDetailInfluencers(ctx context.Context, influencer_id string) (error, models.InfluencerModel) {

	var influencer models.InfluencerModel
	objId, _ := primitive.ObjectIDFromHex(influencer_id)

	err := influencersCollections.FindOne(ctx, bson.M{"_id": objId}).Decode(&influencer)

	if err == nil {
		// increase visits
		influencersCollections.UpdateOne(ctx, bson.D{{"_id", objId}}, bson.D{{"$set", bson.D{{"visits", influencer.Visits + 1}}}})
	}

	return err, influencer
}
