package decoder_test

import (
	"context"

	v1 "github.com/syncromatics/proto-schema-registry/pkg/proto/schema/registry/v1"
	"google.golang.org/grpc"
)

type clientMock struct {
	GetSchemaFunction func(request *v1.GetSchemaRequest) (*v1.GetSchemaResponse, error)
}

func (m *clientMock) GetSchema(ctx context.Context, request *v1.GetSchemaRequest, o ...grpc.CallOption) (*v1.GetSchemaResponse, error) {
	return m.GetSchemaFunction(request)
}

func (m *clientMock) Ping(ctx context.Context, request *v1.PingRequest, o ...grpc.CallOption) (*v1.PingResponse, error) {
	panic("should never call ping")
}

func (m *clientMock) RegisterSchema(ctx context.Context, request *v1.RegisterSchemaRequest, o ...grpc.CallOption) (*v1.RegisterSchemaResponse, error) {
	panic("should never call register schema")
}
