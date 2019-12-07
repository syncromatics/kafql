package decoder

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"

	pd "github.com/emicklei/protobuf2map"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	v1 "github.com/syncromatics/proto-schema-registry/pkg/proto/schema/registry/v1"
)

type Proto struct {
	client v1.RegistryAPIClient
}

func NewProto(registryClient v1.RegistryAPIClient) *Proto {
	return &Proto{registryClient}
}

func (p *Proto) Decode(ctx context.Context, id uint32, message []byte) (string, error) {
	schema, err := p.client.GetSchema(ctx, &v1.GetSchemaRequest{
		Id: id,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to get schema from client")
	}

	if !schema.Exists {
		return "", errors.Errorf("schema '%d' does not exist", id)
	}

	zr, err := gzip.NewReader(bytes.NewBuffer(schema.Schema))
	if err != nil {
		return "", errors.Wrap(err, "failed to unzip schema")
	}
	defer zr.Close()

	defs := pd.NewDefinitions()
	err = defs.ReadFrom("schema.proto", zr)
	if err != nil {
		return "", errors.Wrap(err, "failed to readfrom schema")
	}

	dec := pd.NewDecoder(defs, proto.NewBuffer(message))
	result, err := dec.Decode("gen", "record")
	if err != nil && err != pd.ErrEndOfMessage { // ErrEndOfMessage is okay and means the whole message was decoded...
		return "", errors.Wrap(err, "failed to decode")
	}

	r, err := json.Marshal(result)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal object")
	}

	return string(r), nil
}
