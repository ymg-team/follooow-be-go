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

var influencersCollection *mongo.Collection = configs.GetCollection(configs.DB, "influencers")

// var validate = validator.New()

// handler of GET /influencers
func ListInfluencers(c echo.Context) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var influencers []models.InfluencerModel

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

	// handling filter by search keyword [DONE]
	if c.QueryParam("search") != "" {
		filterListData["name"] = bson.M{"$regex": c.QueryParam("search"), "$options": "i"}
	}

	// handling filter by label [DONE]
	if c.QueryParam("label") != "" {
		filterListData["label"] = bson.M{"$in": strings.Split(c.QueryParam("label"), ",")}
	}

	// handling filter by label [DONE]
	if c.QueryParam("gender") != "" {
		gender := strings.ToLower(c.QueryParam("gender"))
		if gender == "f" || gender == "m" {
			filterListData["gender"] = gender
		}
	}

	// handling filter by nationality [DONE]
	if c.QueryParam("nationality") != "" {
		nationality := strings.ToUpper(c.QueryParam("nationality"))
		filterListData["nationality"] = nationality
	}

	// handling filter by gender

	// by default sortby last update [DONE]
	if c.QueryParam("order_by") == "created_on" {
		optsListData = optsListData.SetSort(bson.D{{"created_on", -1}})
	} else {
		optsListData = optsListData.SetSort(bson.D{{"updated_on", -1}})
	}

	// get data from database
	results, err := influencersCollection.Find(ctx, filterListData, optsListData)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// get count data from database
	// see https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/count/#example
	count, err := influencersCollection.CountDocuments(ctx, filterListData)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// reading data from db in an optimal way
	defer results.Close(ctx)

	// normalize db results
	for results.Next(ctx) {
		var singleInfluencer models.InfluencerModel
		if err = results.Decode(&singleInfluencer); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
		}

		influencers = append(influencers, singleInfluencer)
	}

	// check is no data available
	if len(influencers) < 1 {
		return c.JSON(http.StatusOK, responses.GlobalResponse{Status: http.StatusNoContent, Message: "Influencers not available", Data: &echo.Map{"influencers": influencers, "total": count}})
	} else {
		return c.JSON(http.StatusOK, responses.GlobalResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"influencers": influencers, "total": count}})
	}
}

// handler of GET /influencers/:id
func DetailInfluencers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get influencer_id
	influencerId := c.Param("influencer_id")
	var influencer models.InfluencerModel

	objId, _ := primitive.ObjectIDFromHex(influencerId)

	err := influencersCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&influencer)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// update visits + 1
	_, err = influencersCollection.UpdateOne(ctx, bson.D{{"_id", objId}}, bson.D{{"$set", bson.D{{"visits", influencer.Visits + 1}}}})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	return c.JSON(http.StatusOK, responses.GlobalResponse{Status: http.StatusOK, Message: "OK", Data: &echo.Map{"influencer": influencer}})
}

// handler of GET /influencers/quick-find
func QuickFindInfluencers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var influencers []models.InfluencerSmallDataModel

	// max result is 20
	optsListData := options.Find().SetLimit(20)

	// filter generator
	filterListData := bson.D{}

	if c.QueryParam("ids") != "" {
		idsArr := strings.Split(c.QueryParam("ids"), ",")
		var idsObjId []primitive.ObjectID

		// normalize ids
		for key := range idsArr {
			objId, _ := primitive.ObjectIDFromHex(idsArr[key])

			idsObjId = append(idsObjId, objId)
		}

		filterListData = bson.D{{"_id", bson.M{"$in": idsObjId}}}
	}

	// get data from database
	results, err := influencersCollection.Find(ctx, filterListData, optsListData)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// normalize db results
	for results.Next(ctx) {
		var singleInfluencer models.InfluencerSmallDataModel
		if err = results.Decode(&singleInfluencer); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
		}

		influencers = append(influencers, singleInfluencer)
	}

	return c.JSON(http.StatusOK, responses.GlobalResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"influencers": influencers}})
}

// handler of POST /influencer
func AddInfluencer(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now().UnixNano() / int64(time.Millisecond)

	payload := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{Status: http.StatusBadRequest, Message: "Error parsing json", Data: nil})
	} else {
		new_data := bson.D{
			{"name", payload["name"]},
			{"bio", payload["bio"]},
			{"code", payload["code"]},
			{"avatar", payload["avatar"]},
			{"updated_on", now},
			{"created_on", now},
			{"nationality", payload["nationality"]},
			{"gender", payload["gender"]},
			{"socials", payload["socials"]},
			{"label", payload["label"]},
			{"visits", 1}}

		_, err := influencersCollection.InsertOne(ctx, new_data)

		if err != nil {
			return c.JSON(http.StatusBadRequest, responses.GlobalResponse{Status: http.StatusBadRequest, Message: "Error insert data", Data: nil})
		} else {
			return c.JSON(http.StatusCreated, responses.GlobalResponse{Status: http.StatusCreated, Message: "Success add influencer", Data: nil})
		}

	}
}

// handler of PUT /influencer/:id
func UpdateInfluencer(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get influencer_id
	influencerId := c.Param("influencer_id")
	var influencer models.InfluencerModel

	objId, _ := primitive.ObjectIDFromHex(influencerId)

	// check is data available in db
	err := influencersCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&influencer)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "Error", Data: nil})
	}

	// get payload
	payload := make(map[string]interface{})
	err = json.NewDecoder(c.Request().Body).Decode(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{Status: http.StatusBadRequest, Message: "Error parsing json", Data: nil})
	} else {
		// start update
		filter := bson.D{{"_id", objId}}

		new_data := bson.D{
			{"name", payload["name"]},
			{"bio", payload["bio"]},
			{"code", payload["code"]},
			{"avatar", payload["avatar"]},
			{"updated_on", time.Now().UnixNano() / int64(time.Millisecond)},
			{"nationality", payload["nationality"]},
			{"gender", payload["gender"]},
			{"socials", payload["socials"]},
			{"label", payload["label"]},
		}

		update := bson.D{{"$set", new_data}}

		_, err := influencersCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "Error update database", Data: nil})
		} else {
			return c.JSON(http.StatusOK, responses.GlobalResponse{Status: http.StatusOK, Message: "Success update influencer", Data: nil})
		}
	}
}
