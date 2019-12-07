module github.com/syncromatics/kafql

go 1.13

require (
	github.com/99designs/gqlgen v0.10.2
	github.com/Shopify/sarama v1.24.1
	github.com/avast/retry-go v2.4.3+incompatible // indirect
	github.com/bsm/sarama-cluster v2.1.15+incompatible // indirect
	github.com/burdiyan/kafkautil v0.0.0-20190131162249-eaf83ed22d5b
	github.com/confluentinc/confluent-kafka-go v1.1.0
	github.com/emicklei/protobuf2map v0.0.0-20190903085639-467487a25d28
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/websocket v1.2.0
	github.com/lovoo/goka v0.1.3 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rs/cors v1.6.0
	github.com/samuel/go-zookeeper v0.0.0-20190923202752-2cc03de413da // indirect
	github.com/stretchr/testify v1.4.0
	github.com/syncromatics/proto-schema-registry v0.7.1
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/vektah/gqlparser v1.2.0
	github.com/wvanbergen/kazoo-go v0.0.0-20180202103751-f72d8611297a // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	google.golang.org/grpc v1.25.1
	gopkg.in/confluentinc/confluent-kafka-go.v1 v1.1.0
)

replace github.com/emicklei/protobuf2map => github.com/annymsMthd/protobuf2map v0.0.0-20191207033210-cf6e103222d7
