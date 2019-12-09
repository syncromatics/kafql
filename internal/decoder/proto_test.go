package decoder_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/syncromatics/kafql/internal/decoder"

	v1 "github.com/syncromatics/proto-schema-registry/pkg/proto/schema/registry/v1"
)

var (
	simpleSchema = `syntax = "proto3";
package gen;

message record {
	string name = 2;
}
`
)

func TestProtoDecoder(t *testing.T) {
	client := &clientMock{}

	decoder := decoder.NewProto(client)

	client.GetSchemaFunction = func(request *v1.GetSchemaRequest) (*v1.GetSchemaResponse, error) {
		assert.Equal(t, uint32(43), request.Id)
		schema, err := gzipWrite([]byte(simpleSchema))
		if err != nil {
			t.Fatal(err)
		}

		return &v1.GetSchemaResponse{
			Exists: true,
			Schema: schema,
		}, nil
	}

	s, err := decoder.Decode(context.Background(), 43, []byte{0x12, 0x02, 0x68, 0x69})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "{\"name\":\"hi\"}", s)
}

func gzipWrite(data []byte) ([]byte, error) {
	w := bytes.NewBuffer([]byte{})
	gw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	i, err := gw.Write(data)
	if err != nil {
		return nil, err
	}
	if i != len(data) {
		return nil, errors.Errorf("didnt write everything")
	}
	err = gw.Flush()
	if err != nil {
		return nil, err
	}
	gw.Close()
	return w.Bytes(), err
}
