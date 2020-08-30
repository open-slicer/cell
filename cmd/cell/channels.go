package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"net/http"
)

type channel struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	Parent string `json:"parent"`
}

type channelInsertion struct {
	Name   string `json:"name" binding:"required,gte=1,lte=32"`
	Parent string `json:"parent" binding:"lte=20"`
}

func (req *channelInsertion) insert(requesterID string) response {
	if req.Parent != "" {
		if _, err := pg.Exec(
			context.Background(), "SELECT 1 FROM channels WHERE id = $1", req.Parent,
		); err != nil {
			if err != pgx.ErrNoRows {
				return internalError(err)
			}
			return response{
				Code:    errorParentNotExists,
				Message: "A parent channel with the provided ID doesn't exist",
				HTTP:    http.StatusBadRequest,
			}
		}
	} else {
		req.Parent = "origin"
	}

	c := channel{
		ID:     idNode.Generate().String(),
		Name:   req.Name,
		Owner:  requesterID,
		Parent: req.Parent,
	}

	if _, err := pg.Exec(
		context.Background(),
		"INSERT INTO channels (id, name, owner, parent) VALUES ($1, $2, $3, $4)",
		c.ID, c.Name, c.Owner, c.Parent,
	); err != nil {
		return internalError(err)
	}
	if _, err := pg.Exec(
		context.Background(), "INSERT INTO members (id, channel) VALUES ($1, $2)", c.Owner, c.ID,
	); err != nil {
		return internalError(err)
	}

	return response{
		Code:    http.StatusCreated,
		Message: "Channel created and member created for owner",
		Data:    c,
	}
}
