package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

type websocketServer struct {
	subscriberMessageBuffer int
	publishLimiter          *rate.Limiter
	subscribersMu           sync.Mutex
	subscribers             map[string][]*subscriber
}

func newWebsocketServer() *websocketServer {
	cs := &websocketServer{
		subscriberMessageBuffer: 16,
		subscribers:             make(map[string][]*subscriber),
		publishLimiter:          rate.NewLimiter(rate.Every(time.Second/2), 10),
	}
	return cs
}

type subscriber struct {
	id        string
	messages  chan []byte
	closeSlow func()
}

var currentlyListening = map[string]bool{}

func (srv *websocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.subscribeHandler(w, r)
}

func (srv *websocketServer) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Add error codes and a better way to handle this.
	splitToken := strings.Split(r.Header.Get("Authorization"), "Bearer")
	if len(splitToken) < 2 {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Authorization header should be formatted like 'prefix TOKEN'"))
		return
	}

	token, err := jwt.Parse(strings.TrimSpace(splitToken[1]), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Invalid token"))
		return
	}
	userID, ok := claims["id"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Malformed token"))
		return
	}

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Couldn't accept connection")
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	err = srv.subscribe(r.Context(), c, userID)
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("Couldn't subscribe")
		return
	}
}

func (srv *websocketServer) subscribe(ctx context.Context, c *websocket.Conn, userID string) error {
	ctx = c.CloseRead(ctx)

	if listening, ok := currentlyListening[userID]; !ok || !listening {
		currentlyListening[userID] = true
		go func() {
			pubsub := rdb.Subscribe(ctx, "locketuser:"+userID)
			_, err := pubsub.Receive(ctx)
			if err != nil {
				return
			}
			ch := pubsub.Channel()

			for msg := range ch {
				srv.publish(userID, []byte(msg.Payload))
			}
		}()
	}

	s := &subscriber{
		id:       xid.New().String(),
		messages: make(chan []byte, srv.subscriberMessageBuffer),
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "Connection couldn't keep up; too slow")
		},
	}
	srv.addSubscriber(userID, s)
	defer srv.deleteSubscriber(userID, s.id)

	for {
		select {
		case msg := <-s.messages:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (srv *websocketServer) publish(userID string, msg []byte) bool {
	srv.subscribersMu.Lock()
	defer srv.subscribersMu.Unlock()

	_ = srv.publishLimiter.Wait(context.Background())

	subs, ok := srv.subscribers[userID]
	if !ok || subs == nil {
		return false
	}

	for _, s := range subs {
		select {
		case s.messages <- msg:
			return true
		default:
			go s.closeSlow()
			return false
		}
	}
	return false
}

func (srv *websocketServer) addSubscriber(userID string, s *subscriber) {
	srv.subscribersMu.Lock()
	defer srv.subscribersMu.Unlock()

	subs := srv.subscribers[userID]
	srv.subscribers[userID] = append(subs, s)
}

func (srv *websocketServer) deleteSubscriber(userID string, subID string) {
	srv.subscribersMu.Lock()
	defer srv.subscribersMu.Unlock()

	subs := srv.subscribers[userID]
	for i, sub := range subs {
		if sub != nil && sub.id == subID {
			subs[i] = subs[len(subs)-1]
			subs[len(subs)-1] = nil
			srv.subscribers[userID] = subs[:len(subs)-1]
			break
		}
	}
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}
