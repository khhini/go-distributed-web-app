package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khhini/go-distributed-web-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RecipeHandler godoc
type RecipeHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

// NewRecipesHandler godoc
func NewRecipesHandler(ctx context.Context, collection *mongo.Collection) *RecipeHandler {
	return &RecipeHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// NewRecipeHandler godoc
// @Summary      Add new recipe
// @Description  Add new recipe
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        recipe body models.Recipe false "recipe object"
// @Success      200  {object}  models.Recipe
// @Failure		 400  {string}  StatusBadRequest
// @Failure		 500  {string}  StatusInternalServerError
// @Router       /recipes [post]
func (handler *RecipeHandler) NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	result, err := handler.collection.InsertOne(handler.ctx, recipe)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while inserting a new recipe",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("New recipe added with id %s", result.InsertedID),
	})

}

// ListRecipesHandler godoc
// @Summary      List recipes
// @Description  get all recipes
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Recipe
// @Failure		 500  {string}  StatusInternalServerError
// @Router       /recipes [get]
func (handler *RecipeHandler) ListRecipesHandler(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer cur.Close(handler.ctx)

	recipes := make([]models.Recipe, 0)
	for cur.Next(handler.ctx) {
		var recipe models.Recipe
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
// @Param        recipe body models.Recipe false "recipe object"
// @Success      200  {object}  models.Recipe
// @Failure		 400  {string}  StatusBadRequest
// @Failure		 404  {string}  StatusNotFound
// @Failure		 500  {string}  StatusInternalServerError
// @Router       /recipes/{id} [put]
func (handler *RecipeHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	objectID, _ := primitive.ObjectIDFromHex(id)
	result, err := handler.collection.UpdateOne(handler.ctx, bson.M{
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
func (handler *RecipeHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(id)
	result, err := handler.collection.DeleteOne(handler.ctx, bson.M{
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
// @Success      200  {array}  models.Recipe
// @Failure		 500  {string}  StatusInternalServerError
// @Router       /recipes/search [get]
func (handler *RecipeHandler) SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	cur, err := handler.collection.Find(handler.ctx, bson.D{{"tags", bson.D{{"$all", bson.A{tag}}}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer cur.Close(handler.ctx)

	recipes := make([]models.Recipe, 0)
	for cur.Next(handler.ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, recipes)
}
