package handlers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/khhini/go-distributed-web-app/models"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler godoc
type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

// Claims godoc
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// JWTOutput godoc
type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

// NewAuthHandler godoc
func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// AuthJWTMiddleware godoc
func (handler *AuthHandler) AuthJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := c.GetHeader("Authorization")
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if tkn == nil || !tkn.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()

	}
}

// SignInJWTHandler godoc
// @Summary      Sigin API with username and password
// @Description  Sigin API with username and password
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        user body models.User false "user object"
// @Success      200  {object}  JWTOutput
// @Failure		 400  {string}  StatusBadRequest
// @Failure		 404  {string}  StatusNotFound
// @Failure		 401  {string}  StatusUnauthorized
// @Failure		 500  {string}  StatusInternalServerError
// @Router       /singin [post]
func (handler *AuthHandler) SignInJWTHandler(c *gin.Context) {
	var user models.User
	var dbUser models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"username": user.Username,
	})
	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Username or Password",
		})
		return
	}

	cur.Decode(&dbUser)
	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Username or Password",
		})
		return
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}

	c.JSON(http.StatusOK, jwtOutput)
}

// RefreshHandler godoc
// @Summary      Refresh JWT Token
// @Description  Refresh JWT Token
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "jwt token"
// @Success      200  {object}  JWTOutput
// @Failure		 400  {string}  StatusBadRequest
// @Failure		 404  {string}  StatusNotFound
// @Failure		 401  {string}  StatusUnauthorized
// @Failure		 500  {string}  StatusInternalServerError
// @Router       /signin [post]
func (handler *AuthHandler) RefreshJWTHandler(c *gin.Context) {
	tokenValue := c.GetHeader("Authorization")
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	if tkn == nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token is not expired yet",
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwtOutput)
}

// AuthSessionMiddleware godoc
func (handler *AuthHandler) AuthSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionToken := session.Get("token")
		if sessionToken == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Not logged",
			})
			c.Abort()
		}
		c.Next()
	}
}

// SignInSessionHandler godoc
func (handler *AuthHandler) SignInSessionHandler(c *gin.Context) {
	var user models.User
	var dbUser models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"username": user.Username,
	})
	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Username or Password",
		})
		return
	}

	cur.Decode(&dbUser)
	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Username or Password",
		})
		return
	}

	sessionToken := xid.New().String()
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Set("token", sessionToken)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "User signed in"})
}

// SignOutHandler godoc
// @Summary      Sign Out user
// @Description  Sign Out user by removing cookies
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Success      200  {string}  StatusOK
// @Failure		 400  {string}  StatusBadRequest
// @Failure		 404  {string}  StatusNotFound
// @Failure		 401  {string}  StatusUnauthorized
// @Failure		 500  {string}  StatusInternalServerError
// @Router       /signout [post]
func (handler *AuthHandler) SignOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "Signed out...",
	})
}
