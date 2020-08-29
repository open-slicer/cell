package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

type websocketServer struct{}

func (s websocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"locket"},
	})
	if err != nil {
		log.Error().Err(err).Msg("Couldn't accept ws connection")
		return
	}
	defer c.Close(websocket.StatusInternalError, "Internal server error")

	if c.Subprotocol() != "locket" {
		_ = c.Close(websocket.StatusPolicyViolation, "Clients are required to speak the locket protocol")
		return
	}

	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		err = echo(r.Context(), c, l)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			log.Error().Err(err).Msg("WebSocket handler failed")
			return
		}
	}
}

func echo(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}
	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}
	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}

	err = w.Close()
	return err
}
