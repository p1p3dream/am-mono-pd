package storage

import (
	"context"
	"io"
)

type ListOptions struct {
	DirsOnly  bool
	FilesOnly bool
	WithSize  bool
	Recursive bool
}

type BackendType uint

const (
	BackendTypeLocal BackendType = iota
	BackendTypeS3
)

type Backend interface {
	Path() string
	Type() BackendType
	PathSeparator() string
	PathJoin(elem ...string) string

	Put(ctx context.Context, obj *Object) error
	Get(ctx context.Context, obj *Object) (ObjectReader, error)
	Delete(ctx context.Context, obj *Object) error
	List(ctx context.Context, prefix string, options *ListOptions) ([]*Object, error)
}

// The object's path should be the combination of Dir and Name.
type Object struct {
	// Dir is the path to the directory where the object is stored.
	Dir string `json:"dir,omitempty"`

	// Name is the name of the object.
	Name string `json:"name,omitempty"`

	// The size of the object in bytes.
	Size int64 `json:"size,omitempty"`

	isDirectory bool
}

type NewObjectInput struct {
	Dir         string
	Name        string
	Size        int64
	IsDirectory bool
}

// NewObject creates a new object with the given input.
// The current use case is to create objects for testing purposes,
// and uses outside of this scope are not recommended.
func NewObject(in *NewObjectInput) *Object {
	return &Object{
		Dir:         in.Dir,
		Name:        in.Name,
		Size:        in.Size,
		isDirectory: in.IsDirectory,
	}
}

func (o *Object) IsDirectory() bool {
	return o.isDirectory
}

type ObjectReader interface {
	io.Reader
	io.Closer
	io.ReaderAt
}
