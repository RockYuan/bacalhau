package s3

import (
	_ "embed"
	"errors"
	"fmt"

	"github.com/ipfs/go-cid"
	ipldcodec "github.com/ipld/go-ipld-prime/codec/dagjson"
	dslschema "github.com/ipld/go-ipld-prime/schema/dsl"

	"github.com/bacalhau-project/bacalhau/pkg/model/spec"
	"github.com/bacalhau-project/bacalhau/pkg/model/spec/storage"
)

//go:embed spec.ipldsch
var schema []byte

func load() *storage.Schema {
	dsl, err := dslschema.ParseBytes(schema)
	if err != nil {
		panic(err)
	}
	return (*storage.Schema)(dsl)
}

var (
	Schema              *storage.Schema = load()
	StorageType         cid.Cid         = Schema.Cid()
	defaultModelEncoder                 = ipldcodec.Encode
	defaultModelDecoder                 = ipldcodec.Decode
	EncodingError                       = errors.New("encoding S3StorageSpec to spec.Storage")
	DecodingError                       = errors.New("decoding spec.Storage to S3StorageSpec")
)

type S3StorageSpec struct {
	Bucket         string
	Key            string
	ChecksumSHA256 string
	VersionID      string
	Endpoint       string
	Region         string
}

func (e *S3StorageSpec) AsSpec(name, mount string) (spec.Storage, error) {
	s, err := storage.Encode(name, mount, e, defaultModelEncoder, Schema)
	if err != nil {
		return spec.Storage{}, errors.Join(EncodingError, err)
	}
	return s, nil
}

func Decode(spec spec.Storage) (*S3StorageSpec, error) {
	if spec.Schema != Schema.Cid() {
		return nil, fmt.Errorf("unexpected spec schema %s: %w", spec, DecodingError)
	}
	out, err := storage.Decode[S3StorageSpec](spec, defaultModelDecoder)
	if err != nil {
		return nil, errors.Join(DecodingError, err)
	}
	return out, nil
}
