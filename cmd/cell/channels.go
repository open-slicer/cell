package main

import (
	"context"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"net/http"
)

type channel struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Owner  string  `json:"owner"`
	Parent *string `json:"parent,omitempty"`
}

type channelInsertion struct {
	Name   string `json:"name" binding:"required,gte=1,lte=32"`
	Parent string `json:"parent" binding:"lte=20"`
}

func (req *channelInsertion) insert(requesterID string) response {
	origin := false
	if req.Parent != "" {
		var exists bool
		if err := pg.QueryRow(
			context.Background(), "SELECT EXISTS(SELECT 1 FROM channels WHERE id = $1)", req.Parent,
		).Scan(&exists); err != nil {
			return internalError(err)
		}
		if !exists {
			return response{
				Code:    errorParentNotExists,
				Message: "A parent channel with the given ID doesn't exist",
				HTTP:    http.StatusConflict,
			}
		}
	} else {
		origin = true
	}

	c := channel{
		ID:     idNode.Generate().String(),
		Name:   req.Name,
		Owner:  requesterID,
		Parent: &req.Parent,
	}

	var err error
	if origin {
		_, err = pg.Exec(
			context.Background(),
			"INSERT INTO channels (id, name, owner) VALUES ($1, $2, $3)",
			c.ID, c.Name, c.Owner,
		)
	} else {
		_, err = pg.Exec(
			context.Background(),
			"INSERT INTO channels (id, name, owner, parent) VALUES ($1, $2, $3, $4)",
			c.ID, c.Name, c.Owner, c.Parent,
		)
	}
	if err != nil {
		return internalError(err)
	}

	if _, err := pg.Exec(
		context.Background(),
		"INSERT INTO members (id, \"user\", channel) VALUES ($1, $2, $3)",
		idNode.Generate().String(), c.Owner, c.ID,
	); err != nil {
		return internalError(err)
	}

	return response{
		Code:    http.StatusCreated,
		Message: "Channel created and member created for owner",
		Data:    c,
	}
}

func handleChannelsPOST(c *gin.Context) {
	channel := channelInsertion{}
	if err := c.ShouldBindJSON(&channel); err != nil {
		response{
			Code:    errorBindFailed,
			Message: "Failed to bind JSON",
			HTTP:    http.StatusBadRequest,
			Data:    err.Error(),
		}.send(c)
		return
	}

	claims := jwt.ExtractClaims(c)
	channel.insert(claims[identityKey].(string)).send(c)
}

func (c *channel) get() response {
	var fChannel channel
	if err := pg.QueryRow(
		context.Background(), "SELECT id, name, owner, parent FROM channels WHERE id = $1", c.ID,
	).Scan(&fChannel.ID, &fChannel.Name, &fChannel.Owner, &fChannel.Parent); err != nil {
		if err != pgx.ErrNoRows {
			return internalError(err)
		}
		return response{
			Code:    errorNotFound,
			Message: "Channel doesn't exist",
			HTTP:    http.StatusNotFound,
		}
	}

	return response{
		Code:    http.StatusOK,
		Message: "Channel found",
		Data:    fChannel,
	}
}

func handleChannelsGET(c *gin.Context) {
	channel := channel{
		ID: c.Param("id"),
	}
	channel.get().send(c)
}

type invite struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
	Owner   string `json:"owner"`
}

type inviteInsertion struct {
	Name string `json:"name" binding:"required,gte=4,lte=32"`
}

func (req *inviteInsertion) insert(requesterID, channelID string) response {
	if !commonNameRegex.MatchString(req.Name) {
		return response{
			Code:    errorNotCommonName,
			Message: "Name didn't match the commonName regex",
			HTTP:    http.StatusBadRequest,
			Data:    commonNameRegex.String(),
		}
	}

	var exists bool
	if err := pg.QueryRow(
		context.Background(), "SELECT EXISTS(SELECT 1 FROM invites WHERE name = $1)", req.Name,
	).Scan(&exists); err != nil {
		return internalError(err)
	}
	if exists {
		return response{
			Code:    errorExists,
			Message: "An invite with the given name already exists",
			HTTP:    http.StatusConflict,
		}
	}

	if err := pg.QueryRow(
		context.Background(), "SELECT EXISTS(SELECT 1 FROM channels WHERE id = $1)", channelID,
	).Scan(&exists); err != nil {
		return internalError(err)
	}
	if !exists {
		return response{
			Code:    errorNotFound,
			Message: "A channel with the given ID doesn't exist",
			HTTP:    http.StatusNotFound,
		}
	}

	i := invite{
		Name:    req.Name,
		Channel: channelID,
		Owner:   requesterID,
	}
	if _, err := pg.Exec(
		context.Background(),
		"INSERT INTO invites (name, channel, owner) VALUES ($1, $2, $3)",
		i.Name, i.Channel, i.Owner,
	); err != nil {
		return internalError(err)
	}

	return response{
		Code:    http.StatusCreated,
		Message: "Invite created",
		Data:    i,
	}
}

func handleInvitesPOST(c *gin.Context) {
	invite := inviteInsertion{}
	if err := c.ShouldBindJSON(&invite); err != nil {
		response{
			Code:    errorBindFailed,
			Message: "Failed to bind JSON",
			HTTP:    http.StatusBadRequest,
			Data:    err.Error(),
		}.send(c)
		return
	}

	claims := jwt.ExtractClaims(c)
	invite.insert(claims[identityKey].(string), c.Param("id")).send(c)
}

func (i *invite) get() response {
	var fInvite invite
	if err := pg.QueryRow(
		context.Background(), "SELECT name, owner, channel FROM invites WHERE name = $1", i.Name,
	).Scan(&fInvite.Name, &fInvite.Owner, &fInvite.Channel); err != nil {
		if err != pgx.ErrNoRows {
			return internalError(err)
		}
		return response{
			Code:    errorNotFound,
			Message: "Invite doesn't exist",
			HTTP:    http.StatusNotFound,
		}
	}

	return response{
		Code:    http.StatusOK,
		Message: "Invite found",
		Data:    fInvite,
	}
}

func handleInvitesGET(c *gin.Context) {
	invite := invite{
		Name: c.Param("name"),
	}
	invite.get().send(c)
}
