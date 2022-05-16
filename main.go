package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khhini/go-distributed-web-app/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Recipe godoc
type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

var recipes []Recipe
var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection

func init() {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Recipes API"
	docs.SwaggerInfo.Description = "This is a sample server Recipe server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
}

// IndexHandler godoc
// @Summary      Index endpoint
// @Description  Index endpoint
// @Tags         /
// @Accept       json
// @Produce      json
// @Success      200  {string}  StatusOK
// @Router       / [get]
func IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"ping": "ping",
	})
}

// HealthzHandler godoc
// @Summary      Health Check endpoint
// @Description  Health Check endpoint
// @Tags         /
// @Accept       json
// @Produce      json
// @Success      200  {string}  StatusOK
// @Router       / [get]
func HealthzHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"ping": "ping",
	})
}

// NewRecipeHandler godoc
// @Summary      Add new recipe
// @Description  Add new recipe
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        recipe body Recipe false "recipe object"
// @Success      200  {object}  Recipe
// @Failure		 400  {string}  StatusBadRequest
// @Failure		 500  {string}  StatusInternalServerError
// @Router       /recipes [post]
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	result, err := collection.InsertOne(ctx, recipe)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while inserting a new recipe",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("New recipe added with id %s", result.InsertedID),
		"recipeID": result.InsertedID,
	})

}

// ListRecipesHandler godoc
// @Summary      List recipes
// @Description  get all recipes
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Success      200  {array}  Recipe
// @Failure		 500  {string}  StatusInternalServerError
// @Router       /recipes [get]
func ListRecipesHandler(c *gin.Context) {
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer cur.Close(ctx)

	recipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, recipes)
}

// UpdateRecipeHandler godoc
// @Summary      Update recipe
// @Description  Update recipe
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param		 id path string false "recipe id"
// @Param        recipe body Recipe false "recipe object"
// @Success      200  {object}  Recipe
// @Failure		 400  {string}  StatusBadRequest
// @Failure		 404  {string}  StatusNotFound
// @Failure		 500  {string}  StatusInternalServerError
// @Router       /recipes/{id} [put]
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	objectID, _ := primitive.ObjectIDFromHex(id)
	result, err := collection.UpdateOne(ctx, bson.M{
		"_id": objectID,
	}, bson.D{{"$set", bson.D{
		{"name", recipe.Name},
		{"instructions", recipe.Instructions},
		{"ingredients", recipe.Ingredients},
		{"tags", recipe.Tags},
	}}})

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while inserting a new recipe",
		})
		return
	}
	if result.ModifiedCount > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Recipe has been updated",
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Recipe not found",
		})
	}
}

// DeleteRecipeHandler godoc
// @Summary      Delete recipe
// @Description  Delete recipe
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param		 id path string false "recipe id"
// @Success      200  {string}  StatusOK
// @Failure      404  {string}  StatusNotFound
// @Failure      500  {string}  StatusInternalServerError
// @Router       /recipes/{id} [delete]
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(id)
	result, err := collection.DeleteOne(ctx, bson.M{
		"_id": objectID,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while inserting a new recipe",
		})
		return
	}
	if result.DeletedCount > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Recipe has been deleted",
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Recipe Not Found",
		})
	}

}

// SearchRecipesHandler godoc
// @Summary      List recipes by tag
// @Description  get all recipes by tag
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        tag    query     string  false  "recipe search by tag"
// @Success      200  {array}  Recipe
// @Failure		 500  {string}  StatusInternalServerError
// @Router       /recipes/search [get]
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	cur, err := collection.Find(ctx, bson.D{{"tags", bson.D{{"$all", bson.A{tag}}}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer cur.Close(ctx)

	recipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, recipes)
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)
	router.GET("/", IndexHandler)
	router.GET("/healthz", IndexHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	SetupRouter().Run()
}
