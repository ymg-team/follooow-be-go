package controllers

import (
	"context"
	"follooow-be/configs"
	"follooow-be/models"
	"follooow-be/responses"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var influencersCollection *mongo.Collection = configs.GetCollection(configs.DB, "influencers")
var validate = validator.New()

func ListInfluencers(c echo.Context) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var influencers []models.InfluencerModel

	filterListData := bson.D{}

	// handling limit, by default 6
	var limit int64
	if c.QueryParam("limit") != "" {
		i, err := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.InfluencerResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
		}
		limit = i
	} else {
		limit = int64(6)
	}

	optsListData := options.Find().SetLimit(limit)

	// handling filter by search keyword [DONE]
	if c.QueryParam("search") != "" {
		filterListData = bson.D{{"name", bson.M{"$regex": c.QueryParam("search"), "$options": "i"}}}
	}

	// handling filter by label [DONE]
	if c.QueryParam("label") != "" {
		filterListData = bson.D{{"label", bson.M{"$in": strings.Split(c.QueryParam("label"), ",")}}}
	}

	// handling filter by gender

	// by default sortby last update [DONE]
	optsListData = optsListData.SetSort(bson.D{{"updated_on", -1}})

	// get data from database
	results, err := influencersCollection.Find(ctx, filterListData, optsListData)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.InfluencerResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// get count data from database
	// see https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/count/#example
	count, err := influencersCollection.CountDocuments(ctx, filterListData)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.InfluencerResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// reading data from db in an optimal way
	defer results.Close(ctx)

	// normalize db results
	for results.Next(ctx) {
		var singleInfluencer models.InfluencerModel
		if err = results.Decode(&singleInfluencer); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.InfluencerResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
		}

		influencers = append(influencers, singleInfluencer)
	}

	return c.JSON(http.StatusOK, responses.InfluencerResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"influencers": influencers, "total": count}})
}
