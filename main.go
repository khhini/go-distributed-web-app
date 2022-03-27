package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khhini/go-distributed-web-app/docs"
	"github.com/rs/xid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Recipe godoc
type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe

func init() {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Recipes API"
	docs.SwaggerInfo.Description = "This is a sample server Recipe server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}

// IndexHandler godoc
func IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"ping": "pong",
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
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)

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

	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}
	recipe.ID = id
	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)
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
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}
	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
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
	listOfRecipes := make([]Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
}

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)
	router.GET("/", IndexHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run()
}
