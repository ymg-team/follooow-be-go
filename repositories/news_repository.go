package repositories

import (
	"context"
	"follooow-be/configs"
	"follooow-be/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// struct of GetDetailNews() params
type DetailNewsParams struct {
	NewsId string
	Lang   string
}

var newsCollections *mongo.Collection = configs.GetCollection(configs.DB, "news")

// function to to get detail by news_id
// auto increase visits + 1 if data found on DB
func GetDetailNews(ctx context.Context, params DetailNewsParams) (error, models.NewsModel) {
	var news models.NewsModel
	var influencers []models.InfluencerSmallDataModel

	objId, _ := primitive.ObjectIDFromHex(params.NewsId)
	filterListData := bson.M{
		"_id":  objId,
		"lang": params.Lang,
	}

	err := newsCollections.FindOne(ctx, filterListData).Decode(&news)

	if err == nil {
		// increase visits
		newsCollections.UpdateOne(ctx, bson.D{{"_id", objId}}, bson.D{{"$set", bson.D{{"views", news.Views + 1}}}})

		// get list related influencers, max results is 20
		var idsObjId []primitive.ObjectID

		// variable to save all influencer_id
		idsArr := news.Influencers

		// normalize ids
		for key := range idsArr {
			objId, _ := primitive.ObjectIDFromHex(idsArr[key])
			idsObjId = append(idsObjId, objId)
		}
		optsListDataInfluencers := options.Find().SetLimit(20)

		// filter generator
		filterLastDataInfluencers := bson.D{{"_id", bson.M{"$in": idsObjId}}}

		// get influencer data from database
		resultsInfluencers, _ := influencersCollections.Find(ctx, filterLastDataInfluencers, optsListDataInfluencers)
		defer resultsInfluencers.Close(ctx)

		// normalize db results
		for resultsInfluencers.Next(ctx) {
			var singleInfluencer models.InfluencerSmallDataModel
			if err = resultsInfluencers.Decode(&singleInfluencer); err != nil {
				return err, news
			}

			influencers = append(influencers, singleInfluencer)
		}

		news.Influencers = nil
		news.InfluencersData = influencers

	}

	return err, news

}
