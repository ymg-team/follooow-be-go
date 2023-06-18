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

var influencersCollections *mongo.Collection = configs.GetCollection(configs.DB, "influencers")

func GetDetailInfluencers(influencer_id string) (error, models.InfluencerModel) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var influencer models.InfluencerModel
	objId, _ := primitive.ObjectIDFromHex(influencer_id)

	err := influencersCollections.FindOne(ctx, bson.M{"_id": objId}).Decode(&influencer)

	if err == nil {
		// increase visits
		influencersCollections.UpdateOne(ctx, bson.D{{"_id", objId}}, bson.D{{"$set", bson.D{{"visits", influencer.Visits + 1}}}})
	}

	return err, influencer
}
