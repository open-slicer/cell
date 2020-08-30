package main

import "net/http"

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
	if req.Parent == "" {
		req.Parent = "origin"
	}

	c := channel{
		ID:     idNode.Generate().String(),
		Name:   req.Name,
		Owner:  requesterID,
		Parent: req.Parent,
	}
	return response{
		Code:    http.StatusCreated,
		Message: "Channel created",
		Data:    c,
	}
}
