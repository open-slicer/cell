package main

import (
	"context"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/rs/xid"
)

type channel struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Owner  string  `json:"owner"`
	Parent *string `json:"parent,omitempty"`
}

func (c *channel) insert() error {
	if *c.Parent == "" {
		_, err := pg.Exec(
			context.Background(),
			"INSERT INTO channels (id, name, owner) VALUES ($1, $2, $3)",
			c.ID, c.Name, c.Owner,
		)
		return err
	}

	_, err := pg.Exec(
		context.Background(),
		"INSERT INTO channels (id, name, owner, parent) VALUES ($1, $2, $3, $4)",
		c.ID, c.Name, c.Owner, c.Parent,
	)
	return err
}

func (c *channel) parentExists() (bool, error) {
	var exists bool
	err := pg.QueryRow(
		context.Background(), "SELECT EXISTS(SELECT 1 FROM channels WHERE id = $1)", c.Parent,
	).Scan(&exists)
	return exists, err
}

func (c *channel) get() error {
	return pg.QueryRow(
		context.Background(), "SELECT id, name, owner, parent FROM channels WHERE id = $1", c.ID,
	).Scan(&c.ID, &c.Name, &c.Owner, &c.Parent)
}

type member struct {
	ID          string `json:"id"`
	User        string `json:"user"`
	Channel     string `json:"channel"`
	Permissions int64  `json:"permissions"`
}

func (m *member) insert() error {
	_, err := pg.Exec(
		context.Background(),
		"INSERT INTO members (id, \"user\", channel) VALUES ($1, $2, $3)",
		m.ID, m.User, m.Channel,
	)
	return err
}

func (m *member) exists() (bool, error) {
	var exists bool
	err := pg.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM members WHERE \"user\" = $1 AND channel = $2)",
		m.User, m.Channel,
	).Scan(&exists)
	return exists, err
}

func (m *member) get() error {
	return pg.QueryRow(
		context.Background(), "SELECT permissions FROM members WHERE \"user\" = $1 AND channel = $2)",
		m.User, m.Channel,
	).Scan(&m.Permissions)
}

func handleMembersGET(c *gin.Context) {
	member := member{
		User:    c.Param("member"),
		Channel: c.Param("channel"),
	}
	if err := member.get(); err != nil {
		if err != pgx.ErrNoRows {
			internalError(err).send(c)
			return
		}
		response{
			Code:    errorNotFound,
			Message: "Member doesn't exist",
			HTTP:    http.StatusNotFound,
		}.send(c)
		return
	}

	response{
		Code:    http.StatusOK,
		Message: "Member found",
		Data:    member,
	}.send(c)
}

type channelInsertion struct {
	Name   string `json:"name" binding:"required,gte=1,lte=32"`
	Parent string `json:"parent" binding:"lte=20"`
}

func handleChannelsPOST(c *gin.Context) {
	req := channelInsertion{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response{
			Code:    errorBindFailed,
			Message: "Failed to bind JSON",
			HTTP:    http.StatusBadRequest,
			Data:    err.Error(),
		}.send(c)
		return
	}

	claims := jwt.ExtractClaims(c)
	ch := channel{
		Name:   req.Name,
		Owner:  claims[identityKey].(string),
		Parent: &req.Parent,
	}
	if req.Parent != "" {
		exists, err := ch.parentExists()
		if err != nil {
			internalError(err).send(c)
			return
		}
		if !exists {
			response{
				Code:    errorParentNotExists,
				Message: "A parent channel with the given ID doesn't exist",
				HTTP:    http.StatusConflict,
			}.send(c)
			return
		}
	}

	ch.ID = xid.New().String()
	if err := ch.insert(); err != nil {
		internalError(err).send(c)
		return
	}

	m := member{
		ID:      xid.New().String(),
		User:    ch.Owner,
		Channel: ch.ID,
	}
	if err := m.insert(); err != nil {
		internalError(err).send(c)
		return
	}

	response{
		Code:    http.StatusCreated,
		Message: "Channel created and member created for owner",
		Data:    ch,
	}.send(c)
}

func handleChannelsGET(c *gin.Context) {
	ch := channel{
		ID: c.Param("channel"),
	}

	if err := ch.get(); err != nil {
		if err != pgx.ErrNoRows {
			internalError(err).send(c)
			return
		}
		response{
			Code:    errorNotFound,
			Message: "Channel doesn't exist",
			HTTP:    http.StatusNotFound,
		}.send(c)
		return
	}

	response{
		Code:    http.StatusOK,
		Message: "Channel found",
		Data:    ch,
	}.send(c)
}

type invite struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
	Owner   string `json:"owner"`
}

func (i *invite) insert() error {
	_, err := pg.Exec(
		context.Background(),
		"INSERT INTO invites (name, channel, owner) VALUES ($1, $2, $3)",
		i.Name, i.Channel, i.Owner,
	)
	return err
}

func (i *invite) exists() (bool, error) {
	var exists bool
	err := pg.QueryRow(
		context.Background(), "SELECT EXISTS(SELECT 1 FROM invites WHERE name = $1)", i.Name,
	).Scan(&exists)
	return exists, err
}

func (i *invite) channelExists() (bool, error) {
	var exists bool
	err := pg.QueryRow(
		context.Background(), "SELECT EXISTS(SELECT 1 FROM channels WHERE id = $1)", i.Channel,
	).Scan(&exists)
	return exists, err
}

func (i *invite) get() error {
	return pg.QueryRow(
		context.Background(), "SELECT owner, channel FROM invites WHERE name = $1", i.Name,
	).Scan(&i.Owner, &i.Channel)
}

type inviteInsertion struct {
	Name string `json:"name" binding:"required,gte=4,lte=32"`
}

func handleInvitesPOST(c *gin.Context) {
	req := inviteInsertion{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response{
			Code:    errorBindFailed,
			Message: "Failed to bind JSON",
			HTTP:    http.StatusBadRequest,
			Data:    err.Error(),
		}.send(c)
		return
	}
	if !commonNameRegex.MatchString(req.Name) {
		response{
			Code:    errorNotCommonName,
			Message: "Name didn't match the commonName regex",
			HTTP:    http.StatusBadRequest,
			Data:    commonNameRegex.String(),
		}.send(c)
		return
	}

	claims := jwt.ExtractClaims(c)
	i := invite{
		Name:    req.Name,
		Channel: c.Param("channel"),
		Owner:   claims[identityKey].(string),
	}
	exists, err := i.exists()
	if err != nil {
		internalError(err).send(c)
		return
	}
	if exists {
		response{
			Code:    errorExists,
			Message: "An invite with the given name already exists",
			HTTP:    http.StatusConflict,
		}.send(c)
		return
	}

	channelExists, err := i.channelExists()
	if err != nil {
		internalError(err).send(c)
		return
	}
	if !channelExists {
		response{
			Code:    errorNotFound,
			Message: "A channel with the given ID doesn't exist",
			HTTP:    http.StatusNotFound,
		}.send(c)
		return
	}

	if err := i.insert(); err != nil {
		internalError(err).send(c)
		return
	}
	response{
		Code:    http.StatusCreated,
		Message: "Invite created",
		Data:    i,
	}.send(c)
}

func handleInvitesGET(c *gin.Context) {
	i := invite{
		Name: c.Param("invite"),
	}
	if err := i.get(); err != nil {
		if err != pgx.ErrNoRows {
			internalError(err).send(c)
			return
		}
		response{
			Code:    errorNotFound,
			Message: "Invite doesn't exist",
			HTTP:    http.StatusNotFound,
		}.send(c)
		return
	}

	response{
		Code:    http.StatusOK,
		Message: "Invite found",
		Data:    i,
	}.send(c)
}

func handleInvitesAcceptGET(c *gin.Context) {
	i := invite{
		Name: c.Param("invite"),
	}
	if err := i.get(); err != nil {
		if err != pgx.ErrNoRows {
			internalError(err).send(c)
			return
		}
		response{
			Code:    errorNotFound,
			Message: "Invite doesn't exist",
			HTTP:    http.StatusNotFound,
		}.send(c)
		return
	}

	claims := jwt.ExtractClaims(c)
	m := member{
		User:    claims[identityKey].(string),
		Channel: i.Channel,
	}

	exists, err := m.exists()
	if err != nil {
		internalError(err).send(c)
		return
	}
	if exists {
		response{
			Code:    errorExists,
			Message: "Invite already accepted",
			HTTP:    http.StatusConflict,
			Data:    m.Channel,
		}.send(c)
		return
	}

	m.ID = xid.New().String()
	if err := m.insert(); err != nil {
		internalError(err).send(c)
		return
	}
	response{
		Code:    http.StatusCreated,
		Message: "Invite accepted, member created",
	}.send(c)
}
