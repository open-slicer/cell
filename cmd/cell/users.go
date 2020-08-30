package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"net/http"
	"regexp"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nbutton23/zxcvbn-go"
	"github.com/rs/xid"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var usernameRegex = regexp.MustCompile("^[A-Za-z0-9]+(?:[ _-][A-Za-z0-9]+)*$")

type user struct {
	ID           string `json:"id" bson:"_id"`
	Username     string `json:"username" bson:"username"`
	DisplayName  string `json:"display_name" bson:"display_name"`
	PublicKey    []byte `json:"public_key" bson:"public_key"`
	PasswordHash []byte `json:"-" bson:"password_hash"`
}

type userInsertion struct {
	Username    string `json:"username" binding:"required,gte=1,lte=32"`
	DisplayName string `json:"display_name" binding:"lte=32"`
	PublicKey   []byte `json:"public_key" binding:"required"`
	Password    string `json:"password" binding:"required,gte=1,lte=72"`
}

type userLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (req *userInsertion) insert() response {
	if !usernameRegex.MatchString(req.Username) {
		return response{
			Code:    errorBindFailed,
			Message: "Username didn't match the required regex",
			HTTP:    http.StatusBadRequest,
			Data:    usernameRegex.String(),
		}
	}

	passStrength := zxcvbn.PasswordStrength(req.Password, []string{req.Username, req.DisplayName})
	if passStrength.Score < 3 {
		return response{
			Code:    errorPasswordInsecure,
			Message: fmt.Sprintf("Password is not secure under zxcvbn (got %d, want >=3)", passStrength.Score),
			HTTP:    http.StatusBadRequest,
			Data:    passStrength,
		}
	}

	var fetchedID string

	if err := pg.QueryRow(
		context.Background(), "SELECT id FROM users WHERE username = $1", req.Username,
	).Scan(&fetchedID); err == nil {
		return response{
			Code:    errorExists,
			Message: "A user with the given username already exists",
			HTTP:    http.StatusConflict,
			Data:    fetchedID,
		}
	} else if err != pgx.ErrNoRows {
		return internalError(err)
	}

	u := user{
		ID:        xid.New().String(),
		Username:  req.Username,
		PublicKey: req.PublicKey,
	}
	if u.DisplayName == "" {
		u.DisplayName = u.Username
	} else {
		u.DisplayName = strings.TrimSpace(u.DisplayName)
	}

	var err error
	u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return internalError(err)
	}

	if _, err := pg.Exec(
		context.Background(),
		"INSERT INTO users (id, username, display_name, password_hash, public_key) VALUES ($1, $2, $3, $4, $5)",
		u.ID, u.Username, u.DisplayName, u.PasswordHash, u.PublicKey,
	); err != nil {
		return internalError(err)
	}
	return response{
		Code:    http.StatusCreated,
		Message: "User created",
		Data:    u,
	}
}

func handleUsersPost(c *gin.Context) {
	user := userInsertion{}
	if err := c.ShouldBindJSON(&user); err != nil {
		response{
			Code:    errorBindFailed,
			Message: "Failed to bind JSON",
			HTTP:    http.StatusBadRequest,
			Data:    err.Error(),
		}.send(c)
		return
	}

	user.insert().send(c)
}

func (u *user) get() response {
	if err := mng.users.FindOne(context.Background(), bson.M{
		"_id": u.ID,
	}).Decode(u); err != nil {
		if err != mongo.ErrNoDocuments {
			return internalError(err)
		}
		return response{
			Code:    errorNotFound,
			Message: "User doesn't exist",
			HTTP:    http.StatusNotFound,
		}
	}

	return response{
		Code:    http.StatusOK,
		Message: "User found",
		Data:    u,
	}
}

func handleUsersGet(c *gin.Context) {
	user := user{
		ID: c.Param("id"),
	}
	user.get().send(c)
}

const identityKey = "id"

func getAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "slicer",
		Key:         []byte(viper.GetString("security.secret")),
		IdentityKey: identityKey,
		MaxRefresh:  dayDuration * 7,
		TokenLookup: "header: Authorization, query: token",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(user); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return user{
				ID: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var req userLogin
			if err := c.ShouldBindJSON(&req); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			var userDoc user

			if err := mng.users.FindOne(context.Background(), bson.M{
				"username": req.Username,
			}).Decode(&userDoc); err != nil {
				if err != mongo.ErrNoDocuments {
					captureException(err)
				}
				return nil, jwt.ErrFailedAuthentication
			}

			if err := bcrypt.CompareHashAndPassword(userDoc.PasswordHash, []byte(req.Password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return userDoc, nil
		},
	})
}
