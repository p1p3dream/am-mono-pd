package storage

import (
	"context"
	"os"
	"path/filepath"

	"abodemine/lib/errors"
)

type LocalBackend struct {
	FilesystemPath string `json:"path,omitempty"`
}

func (b *LocalBackend) Path() string {
	return b.FilesystemPath
}

func (b *LocalBackend) Type() BackendType {
	return BackendTypeLocal
}

func (b *LocalBackend) PathSeparator() string {
	return string(filepath.Separator)
}

func (b *LocalBackend) PathJoin(elem ...string) string {
	return filepath.Join(elem...)
}

func (b *LocalBackend) Put(ctx context.Context, obj *Object) error {
	return nil
}

func (b *LocalBackend) Get(ctx context.Context, obj *Object) (ObjectReader, error) {
	targetPath := filepath.Join(obj.Dir, obj.Name)

	file, err := os.Open(targetPath)
	if err != nil {
		return nil, &errors.Object{
			Id:     "230af5a5-e6df-4241-8894-b92a85527656",
			Code:   errors.Code_FAILED_PRECONDITION,
			Detail: "Failed to open file.",
			Cause:  err.Error(),
		}
	}

	return file, nil
}

func (b *LocalBackend) Delete(ctx context.Context, obj *Object) error {
	return nil
}

func (b *LocalBackend) List(ctx context.Context, prefix string, options *ListOptions) ([]*Object, error) {
	targetPath := filepath.Join(b.FilesystemPath, prefix)

	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return nil, &errors.Object{
			Id:     "ff1edd92-578c-49f4-a191-2cbe8b7f33b3",
			Code:   errors.Code_FAILED_PRECONDITION,
			Detail: "Failed to list directory.",
			Cause:  err.Error(),
		}
	}

	objects := []*Object{}

	for _, entry := range entries {
		if (entry.IsDir() && options.FilesOnly) ||
			(!entry.IsDir() && options.DirsOnly) {
			continue
		}

		obj := &Object{
			Dir:  targetPath,
			Name: entry.Name(),

			isDirectory: entry.IsDir(),
		}

		if options.WithSize {
			info, err := entry.Info()
			if err != nil {
				return nil, &errors.Object{
					Id:     "f0a22612-91f5-40a0-985c-f85858c01406",
					Code:   errors.Code_FAILED_PRECONDITION,
					Detail: "Failed to get info of entry.",
					Cause:  err.Error(),
				}
			}

			obj.Size = info.Size()
		}

		objects = append(objects, obj)
	}

	return objects, nil
}
