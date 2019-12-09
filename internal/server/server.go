package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/syncromatics/kafql/internal/models"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/syncromatics/kafql/internal/kafka"
)

type server struct {
	client *kafka.Service
}

// Run runs the kafql server
type Run struct {
	schema graphql.ExecutableSchema
}

func (r *Run) Start(ctx context.Context, route string, port int) func() error {
	mux := http.NewServeMux()
	mux.Handle(
		route,
		handler.GraphQL(r.schema,
			handler.WebsocketKeepAliveDuration(10*time.Second),
			handler.WebsocketUpgrader(websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}),
		),
	)
	mux.Handle("/playground", handler.Playground("GraphQL", route))
	handler := cors.AllowAll().Handler(mux)

	srv := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%d", port),
	}
	srv.SetKeepAlivesEnabled(true)

	cancel := make(chan error)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			select {
			case cancel <- errors.Wrap(err, "failed to serve http"):
			default:
			}
		}
	}()

	return func() error {
		select {
		case <-ctx.Done():
			srv.Close()
			return nil
		case msg := <-cancel:
			return msg
		}
	}
}

// NewServer creates a new kafql server
func NewServer(client *kafka.Service) *Run {
	s := &server{client}
	schema := NewExecutableSchema(Config{
		Resolvers: s,
	})
	return &Run{schema}
}

func (s *server) Query() QueryResolver {
	return s
}

func (s *server) Subscription() SubscriptionResolver {
	return s
}

func (s *server) GetTopics(ctx context.Context) ([]*models.Topic, error) {
	topics, err := s.client.GetTopics()
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (s *server) RecordAdded(ctx context.Context, topic string, key string) (<-chan *models.Record, error) {
	sub, err := s.client.Subscribe(ctx, topic, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe to topic")
	}

	return sub, nil
}
