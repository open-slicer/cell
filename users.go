package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nbutton23/zxcvbn-go"
	"github.com/rs/xid"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
)

var usernameRegex = regexp.MustCompile(fmt.Sprintf(
	"^(?=.{%d,32}$)(?![_.])(?!.*[_.]{2})[a-zA-Z0-9._]+(?<![_.])$",
	viper.GetInt("options.min_username_length"),
))

type user struct {
	ID           []byte `json:"id" bson:"_id"`
	Username     string `json:"username" bson:"username"`
	DisplayName  string `json:"display_name" bson:"display_name"`
	PublicKey    []byte `json:"public_key" bson:"public_key"`
	Password     string `json:"password,omitempty" bson:"-"`
	PasswordHash []byte `json:"-" bson:"password_hash"`
}

func (u *user) insert() response {
	if someZero(u.Username, u.Password) {
		return response{
			Code:    errorMissingField,
			Message: "`username` and `password` are required",
			HTTP:    http.StatusBadRequest,
			Data:    []string{"username", "password"},
		}
	}

	if len(u.Password) > 72 {
		return response{
			Code:    errorTooLarge,
			Message: "Password must be less than 72 characters",
			HTTP:    http.StatusBadRequest,
		}
	}
	if len(u.DisplayName) > 32 {
		return response{
			Code:    errorTooLarge,
			Message: "Username must be less than 32 characters",
			HTTP:    http.StatusBadRequest,
		}
	}
	if !usernameRegex.MatchString(u.Username) {
		return response{
			Code:    errorDidntMatch,
			Message: "Username didn't match the required regex",
			HTTP:    http.StatusBadRequest,
			Data:    usernameRegex.String(),
		}
	}

	passStrength := zxcvbn.PasswordStrength(u.Password, []string{u.Username, u.DisplayName})
	if passStrength.Score < 3 {
		return response{
			Code:    errorPasswordInsecure,
			Message: fmt.Sprintf("Password is not secure under zxcvbn (got %d, want >=3)", passStrength.Score),
			HTTP:    http.StatusBadRequest,
			Data:    passStrength,
		}
	}

	var fetchedUser user

	ctx, _ := context.WithTimeout(context.Background(), callTimeout)
	if err := db.users.FindOne(ctx, bson.M{
		"username": u.Username,
	}).Decode(&fetchedUser); err == nil {
		return response{
			Code:    errorExists,
			Message: fmt.Sprintf("A user with the username %s already exists", u.Username),
			HTTP:    http.StatusConflict,
			Data:    fetchedUser,
		}
	} else if err != mongo.ErrNoDocuments {
		// TODO: Capture this error with Sentry.
		return response{
			Code:    errorInternalError,
			Message: "Failed to find existing users",
			HTTP:    http.StatusInternalServerError,
		}
	}

	var err error
	u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Capture this error with Sentry.
		return response{
			Code:    errorInternalError,
			Message: "Failed to hash password",
			HTTP:    http.StatusInternalServerError,
		}
	}

	u.ID = xid.New().Bytes()
	if u.DisplayName == "" {
		u.DisplayName = u.Username
	} else {
		u.DisplayName = strings.TrimSpace(u.DisplayName)
	}

	if _, err := db.users.InsertOne(ctx, u); err != nil {
		return response{
			Code:    errorInternalError,
			Message: "Failed to create user",
			HTTP:    http.StatusInternalServerError,
		}
	}
	// Make sure to hide the password.
	u.Password = ""
	return response{
		Code:    http.StatusCreated,
		Message: "User created",
		Data:    u,
	}
}

func handleUsersPost(c *gin.Context) {
	user := user{}
	if err := c.BindJSON(&user); err != nil {
		response{
			Code:    errorBindFailed,
			Message: "Failed to bind JSON",
			HTTP:    http.StatusBadRequest,
		}.send(c)
		return
	}

	user.insert().send(c)
}