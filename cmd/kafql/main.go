package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"

	"github.com/syncromatics/kafql/internal/decoder"
	"github.com/syncromatics/kafql/internal/kafka"
	"github.com/syncromatics/kafql/internal/server"

	v1 "github.com/syncromatics/proto-schema-registry/pkg/proto/schema/registry/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	settings, err := getSettingsFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(ctx)

	con, err := grpc.Dial(settings.ProtoRegistry, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := v1.NewRegistryAPIClient(con)
	proto := decoder.NewProto(client)

	conf := sarama.NewConfig()
	conf.Version = sarama.MaxVersion

	admin, err := sarama.NewClusterAdmin([]string{settings.KafkaBroker}, conf)
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := sarama.NewConsumer([]string{settings.KafkaBroker}, conf)
	if err != nil {
		log.Fatal(err)
	}

	service := kafka.NewService(admin, consumer, proto)

	server := server.NewServer(service)
	group.Go(server.Start(ctx, "/", settings.Port))

	eventChan := make(chan os.Signal)
	signal.Notify(eventChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("server started...")

	select {
	case <-eventChan:
	case <-ctx.Done():
	}

	fmt.Println("server stopping...")

	cancel()

	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
}
