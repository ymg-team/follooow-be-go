package handlers

import (
	"context"
	"encoding/json"
	"follooow-be/configs"
	"follooow-be/models"
	"follooow-be/responses"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var newsCollection *mongo.Collection = configs.GetCollection(configs.DB, "news")

// var validate = validator.New()

// handle of GET /news
func ListNews(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var news []models.NewsModel

	filterListData := bson.M{}

	// handling limit, by default 6
	var limit int64
	var page int64

	if c.QueryParam("limit") != "" {
		i, err := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
		}
		limit = i
	} else {
		limit = int64(6)
	}

	// handling page, by default 1
	if c.QueryParam("page") != "" {
		i, err := strconv.ParseInt(c.QueryParam("page"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
		}
		page = (i - 1) * limit
	} else {
		page = int64(0)
	}

	optsListData := options.Find().SetLimit(limit).SetSkip(page)

	// handling filter by language
	if c.QueryParam("lang") != "" {
		filterListData["lang"] = c.QueryParam("lang")
	}

	// handling filter by search keyword [DONE]
	if c.QueryParam("search") != "" {
		filterListData["title"] = bson.M{"$regex": c.QueryParam("search"), "$options": "i"}
	}

	// handling filter by influencer id keyword [DONE]
	if c.QueryParam("influencer_ids") != "" {
		idsArr := strings.Split(c.QueryParam("influencer_ids"), ",")
		filterListData["influencers"] = bson.M{"$in": idsArr}
	}

	// by default sortby last update [DONE]
	if c.QueryParam("order_by") == "created_on" {
		optsListData = optsListData.SetSort(bson.D{{"created_on", -1}})
	} else {
		optsListData = optsListData.SetSort(bson.D{{"updated_on", -1}})
	}

	// handling filter by tags [DONE]
	if c.QueryParam("tags") != "" {
		filterListData["tags"] = bson.M{"$in": strings.Split(c.QueryParam("tags"), ",")}
	}

	// get data from database
	results, err := newsCollection.Find(ctx, filterListData, optsListData)

	// reading data from db in an optimal way
	defer results.Close(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// get count data from database
	// see https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/count/#example
	count, err := newsCollection.CountDocuments(ctx, filterListData)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// normalize db results
	for results.Next(ctx) {
		var singleNews models.NewsModel
		var influencers []models.InfluencerSmallDataModel
		if err = results.Decode(&singleNews); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
		}

		// convert influencers to data
		if len(singleNews.Influencers) > 0 {
			// get all influencers on post
			// max result is 20
			optsListDataInfluencers := options.Find().SetLimit(20)

			// filter generator
			filterListDataInfluencers := bson.D{}

			idsArr := singleNews.Influencers
			var idsObjId []primitive.ObjectID

			// normalize ids
			for key := range idsArr {
				objId, _ := primitive.ObjectIDFromHex(idsArr[key])
				idsObjId = append(idsObjId, objId)
			}

			filterListDataInfluencers = bson.D{{"_id", bson.M{"$in": idsObjId}}}

			// get data from database
			resultsInfluencers, err := influencersCollection.Find(ctx, filterListDataInfluencers, optsListDataInfluencers)
			defer resultsInfluencers.Close(ctx)
			// normalize db results
			for resultsInfluencers.Next(ctx) {
				var singleInfluencer models.InfluencerSmallDataModel
				if err = resultsInfluencers.Decode(&singleInfluencer); err != nil {
					return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
				}

				influencers = append(influencers, singleInfluencer)
			}

			singleNews.Influencers = nil
			singleNews.InfluencersData = influencers
		}
		// end of get all influencers on post

		news = append(news, singleNews)
	}

	// check is no data available
	if len(news) < 1 {
		return c.JSON(http.StatusOK, responses.GlobalResponse{Status: http.StatusNoContent, Message: "News not available", Data: &echo.Map{"news": news, "total": count}})
	} else {
		return c.JSON(http.StatusOK, responses.GlobalResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"news": news, "total": count}})
	}
}

// handle of GET /news/:id
func DetailNews(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get influencer_id
	newsId := c.Param("news_id")
	var news models.NewsModel
	var influencers []models.InfluencerSmallDataModel

	objId, _ := primitive.ObjectIDFromHex(newsId)

	filterListData := bson.M{}

	filterListData["_id"] = objId
	// filterListData["lang"] = c.QueryParam("lang")

	// handling filter by language
	if c.QueryParam("lang") != "" {
		filterListData["lang"] = c.QueryParam("lang")
	}

	err := newsCollection.FindOne(ctx, filterListData).Decode(&news)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// update views + 1
	_, err = newsCollection.UpdateOne(ctx, bson.D{{"_id", objId}}, bson.D{{"$set", bson.D{{"views", news.Views + 1}}}})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// get influencers data

	// convert influencers to data
	if len(news.Influencers) > 0 {
		// get all influencers on post
		// max result is 20
		optsListDataInfluencers := options.Find().SetLimit(20)

		// filter generator
		filterListDataInfluencers := bson.D{}

		idsArr := news.Influencers
		var idsObjId []primitive.ObjectID

		// normalize ids
		for key := range idsArr {
			objId, _ := primitive.ObjectIDFromHex(idsArr[key])
			idsObjId = append(idsObjId, objId)
		}

		// filter generator
		filterListDataInfluencers = bson.D{{"_id", bson.M{"$in": idsObjId}}}

		// get data from database
		resultsInfluencers, err := influencersCollection.Find(ctx, filterListDataInfluencers, optsListDataInfluencers)
		defer resultsInfluencers.Close(ctx)
		// normalize db results
		for resultsInfluencers.Next(ctx) {
			var singleInfluencer models.InfluencerSmallDataModel
			if err = resultsInfluencers.Decode(&singleInfluencer); err != nil {
				return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
			}

			influencers = append(influencers, singleInfluencer)
		}

		news.Influencers = nil
		news.InfluencersData = influencers
	}
	// end of get all influencers on post

	return c.JSON(http.StatusOK, responses.GlobalResponse{Status: http.StatusOK, Message: "OK", Data: &echo.Map{"news": news}})
}

// handle of POST /news
func CreateNews(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// var payload models.PayloadNews
	var payload models.PayloadNews
	err := json.NewDecoder(c.Request().Body).Decode(&payload)
	now := time.Now().UnixNano() / int64(time.Millisecond)

	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{Status: http.StatusBadRequest, Message: "Error parsing json", Data: nil})
	} else {
		new_data := bson.D{
			{"title", payload.Title},
			{"views", 1},
			{"updated_on", now},
			{"created_on", now},
			{"thumbnail", payload.Thumbnail},
			{"content", payload.Content},
			{"tags", payload.Tags},
			{"influencers", payload.Influencers},
			{"lang", payload.Lang},
		}

		_, err := newsCollection.InsertOne(ctx, new_data)

		if err != nil {
			return c.JSON(http.StatusBadRequest, responses.GlobalResponse{Status: http.StatusBadRequest, Message: "Error insert data", Data: nil})
		} else {

			var idsObjId []primitive.ObjectID

			// normalize ids
			influencers := payload.Influencers
			for key := range influencers {
				objId, _ := primitive.ObjectIDFromHex(influencers[key])

				idsObjId = append(idsObjId, objId)
			}

			_, err = influencersCollection.UpdateMany(ctx, bson.D{{"_id", bson.M{"$in": idsObjId}}}, bson.D{{"$set", bson.D{{"updated_on", now}}}})

			return c.JSON(http.StatusCreated, responses.GlobalResponse{Status: http.StatusCreated, Message: "Success create news", Data: nil})
		}
	}
}

// handle of POST /news/:id
func UpdateNews(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get news id
	newsId := c.Param("news_id")
	var news models.NewsModel

	objId, _ := primitive.ObjectIDFromHex(newsId)

	// check is data available in db
	errFind := newsCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&news)

	if errFind != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "Error", Data: nil})
	}

	// var payload models.PayloadNews
	var payload models.PayloadNews
	err := json.NewDecoder(c.Request().Body).Decode(&payload)
	now := time.Now().UnixNano() / int64(time.Millisecond)

	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{Status: http.StatusBadRequest, Message: "Error parsing json", Data: nil})
	} else {
		new_data := bson.D{
			{"title", payload.Title},
			{"updated_on", now},
			{"thumbnail", payload.Thumbnail},
			{"content", payload.Content},
			{"tags", payload.Tags},
			{"influencers", payload.Influencers},
			{"lang", payload.Lang},
		}

		filter := bson.D{{"_id", objId}}
		update := bson.D{{"$set", new_data}}

		_, err := newsCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return c.JSON(http.StatusBadRequest, responses.GlobalResponse{Status: http.StatusBadRequest, Message: "Error update news", Data: nil})
		} else {

			var idsObjId []primitive.ObjectID

			// normalize ids
			influencers := payload.Influencers
			for key := range influencers {
				objId, _ := primitive.ObjectIDFromHex(influencers[key])

				idsObjId = append(idsObjId, objId)
			}

			_, err = influencersCollection.UpdateMany(ctx, bson.D{{"_id", bson.M{"$in": idsObjId}}}, bson.D{{"$set", bson.D{{"updated_on", now}}}})

			return c.JSON(http.StatusCreated, responses.GlobalResponse{Status: http.StatusCreated, Message: "Success update news", Data: nil})
		}
	}
}
