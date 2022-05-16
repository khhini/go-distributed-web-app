package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/khhini/go-distributed-web-app/docs"
	"github.com/khhini/go-distributed-web-app/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var err error
var client *mongo.Client

var recipesHandler *handlers.RecipeHandler
var authHandler *handlers.AuthHandler

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
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	userCollection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_SERVER"),
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping()
	log.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
	authHandler = handlers.NewAuthHandler(ctx, userCollection)
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

// SetupRouter godoc
func SetupRouter() *gin.Engine {
	router := gin.Default()
	store, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte(os.Getenv("JWT_SECRET")))
	router.Use(sessions.Sessions("recipes_api", store))

	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)

	router.POST("/signin", authHandler.SignInJWTHandler)
	router.POST("/signout", authHandler.SignOutHandler)
	router.POST("/refresh", authHandler.RefreshJWTHandler)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthJWTMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	}

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
