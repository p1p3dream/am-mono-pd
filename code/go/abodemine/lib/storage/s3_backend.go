package storage

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"

	"abodemine/lib/errors"
	"abodemine/lib/val"
)

type S3Backend struct {
	AWS aws.Config `json:"-"`

	Bucket string `json:"bucket,omitempty"`
}

func (b *S3Backend) Path() string {
	return "/"
}

func (b *S3Backend) Type() BackendType {
	return BackendTypeS3
}

func (b *S3Backend) PathSeparator() string {
	return "/"
}

func (b *S3Backend) PathJoin(elem ...string) string {
	return path.Join(elem...)
}

func (b *S3Backend) Put(ctx context.Context, obj *Object) error {
	return nil
}

func (b *S3Backend) Get(ctx context.Context, obj *Object) (ObjectReader, error) {
	s3Client := s3.NewFromConfig(b.AWS)
	targetPath := path.Join(obj.Dir, obj.Name)

	if len(targetPath) > 0 && targetPath[0] == '/' {
		targetPath = targetPath[1:]
	}

	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(b.Bucket),
		Key:    aws.String(targetPath),
	}

	result, err := s3Client.GetObject(ctx, getObjectInput)
	if err != nil {
		return nil, &errors.Object{
			Id:     "222d21ad-c873-41c3-a4e1-d0a10e39acd6",
			Code:   errors.Code_NOT_FOUND,
			Detail: "Failed to get object.",
			Cause:  err.Error(),
		}
	}

	return &S3Reader{
		body: result.Body,
		size: obj.Size,
	}, nil
}

func (b *S3Backend) Delete(ctx context.Context, obj *Object) error {
	return nil
}

func (b *S3Backend) List(ctx context.Context, prefix string, options *ListOptions) ([]*Object, error) {
	s3Client := s3.NewFromConfig(b.AWS)
	targetPath := path.Clean(prefix)

	if targetPath == "." {
		targetPath = "/"
	}

	targetPrefix := val.Ternary(
		targetPath == "/",
		targetPath,
		targetPath+"/",
	)

	if len(targetPrefix) > 0 && targetPrefix[0] == '/' {
		targetPrefix = targetPrefix[1:]
	}

	var continuationToken *string

	listInput := &s3.ListObjectsV2Input{
		Bucket: val.PtrRef(b.Bucket),
		// Ensure that prefix is considered a directory.
		Prefix:            val.PtrRef(targetPrefix),
		ContinuationToken: continuationToken,
	}

	if !options.Recursive {
		listInput.Delimiter = val.PtrRef("/")
	}

	objects := []*Object{}

	// Paginate through results
	paginator := s3.NewListObjectsV2Paginator(s3Client, listInput)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, &errors.Object{
				Id:     "20fba024-7496-46a7-93c5-ad14e84d598a",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed paginate next page.",
				Cause:  err.Error(),
			}
		}

		if options.DirsOnly || !options.FilesOnly {
			for _, entry := range page.CommonPrefixes {
				obj := &Object{
					Dir: targetPath,
					Name: strings.TrimSuffix(
						strings.TrimPrefix(
							val.PtrDeref(entry.Prefix),
							targetPrefix,
						),
						"/",
					),
					isDirectory: true,
				}

				objects = append(objects, obj)
			}
		}

		if !options.DirsOnly || options.FilesOnly {
			for _, entry := range page.Contents {
				obj := &Object{
					Dir:  targetPath,
					Name: strings.TrimPrefix(val.PtrDeref(entry.Key), targetPrefix),
				}

				if options.WithSize {
					obj.Size = val.PtrDeref(entry.Size)
				}

				objects = append(objects, obj)
			}
		}

		listInput.ContinuationToken = page.NextContinuationToken
	}

	return objects, nil
}

type S3Reader struct {
	body        io.ReadCloser
	readerAt    io.ReaderAt
	closed      bool
	mu          sync.Mutex
	size        int64
	tmpFilePath string
}

func (r *S3Reader) Read(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.body == nil {
		return 0, &errors.Object{
			Id:     "5be8c03e-785d-4e20-8264-ac676dd9825a",
			Code:   errors.Code_UNKNOWN,
			Detail: "Body is nil.",
		}
	}

	if r.closed {
		return 0, &errors.Object{
			Id:     "30b53d34-0523-4933-9705-6edcb3aececc",
			Code:   errors.Code_UNKNOWN,
			Detail: "Body is closed.",
		}
	}

	n, err := r.body.Read(p)
	if err != nil {
		return 0, &errors.Object{
			Id:     "2ef3d6ad-6a30-4356-a793-7fe6bee468db",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to read body.",
			Cause:  err.Error(),
		}
	}

	return n, nil
}

func (r *S3Reader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}

	if r.body == nil {
		return &errors.Object{
			Id:     "83fa34fd-366a-4435-8ca6-557fa2021f86",
			Code:   errors.Code_UNKNOWN,
			Detail: "Body is nil.",
		}
	}

	defer func() {
		if r.readerAt == nil {
			return
		}

		if err := os.Remove(r.tmpFilePath); err != nil {
			log.Error().
				Err(&errors.Object{
					Id:     "cd6c46c3-4718-4c51-be09-77c5d989c38a",
					Code:   errors.Code_UNKNOWN,
					Detail: "Failed to remove temporary file.",
					Cause:  err.Error(),
				}).
				Send()
		}
	}()

	if err := r.body.Close(); err != nil {
		return &errors.Object{
			Id:     "d4dc8995-520d-4238-a4fc-8ba046227c4d",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to close body.",
			Cause:  err.Error(),
		}
	}

	r.closed = true

	return nil
}

// ReadAt caches the body into disk and reads from the cache.
// This is to ensure we are compliant with io.ReaderAt interface, which
// requires parallel reads to be safe. Moreover, we don't want to
// buffer the entire body in memory to avoid OOM.
//
// This method was added to support stdlib's zip.NewReader.
func (r *S3Reader) ReadAt(p []byte, off int64) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

READAT:
	if r.readerAt != nil {
		n, err := r.readerAt.ReadAt(p, off)
		if err != nil {
			return 0, &errors.Object{
				Id:     "122dc423-34ad-4abd-b1e9-79df874926d9",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to read body.",
				Cause:  err.Error(),
			}
		}

		return n, nil
	}

	// Create a temporary file to cache the body.

	randName, err := val.NewUUID4()
	if err != nil {
		return 0, errors.Forward(err, "9e3a6627-355c-4f1a-bd38-9df0c1bda444")
	}

	tmpDir := filepath.Join(os.TempDir(), "s3_reader")
	if err := os.MkdirAll(tmpDir, 0750); err != nil {
		return 0, &errors.Object{
			Id:     "5f055938-a3ae-4695-b3a6-61df739c4363",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create temporary directory.",
			Cause:  err.Error(),
		}
	}

	tmpFilePath := filepath.Join(tmpDir, randName.String())

	// We write the file first and close the s3 body to ensure we
	// don't leave any open file descriptors.
	// Moreover, we want to have a read-only object to ensure that
	// we don't accidentally modify the file during future operations.
	// This is a reader object, after all.

	tmpWOFile, err := os.OpenFile(tmpFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return 0, &errors.Object{
			Id:     "b0a7a9b5-7b4d-4f3f-8f0c-9c9c1f4a5b0d",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create temporary file.",
			Cause:  err.Error(),
		}
	}

	defer func() {
		if r.readerAt != nil {
			return
		}

		if err := tmpWOFile.Close(); err != nil {
			log.Error().
				Err(&errors.Object{
					Id:     "63be4367-455e-4c96-937e-bd927f7c07de",
					Code:   errors.Code_UNKNOWN,
					Detail: "Failed to close temporary file.",
					Cause:  err.Error(),
				}).
				Send()
		}

		if err := os.Remove(tmpFilePath); err != nil {
			log.Error().
				Err(&errors.Object{
					Id:     "dab723ad-b2fc-4654-9f7b-70a506f7cd54",
					Code:   errors.Code_UNKNOWN,
					Detail: "Failed to remove temporary file.",
					Cause:  err.Error(),
				}).
				Send()
		}
	}()

	if _, err := io.Copy(tmpWOFile, r.body); err != nil {
		return 0, &errors.Object{
			Id:     "6cce0c15-0677-4d44-a4d0-61cda447013f",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to copy body to temporary file.",
			Cause:  err.Error(),
		}
	}

	if err := r.body.Close(); err != nil {
		return 0, &errors.Object{
			Id:     "05a61e46-56c5-4829-97b9-b1c3006f5e3a",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to close body.",
			Cause:  err.Error(),
		}
	}

	tmpROFile, err := os.Open(tmpFilePath)
	if err != nil {
		return 0, &errors.Object{
			Id:     "35674f9d-8a96-4306-b0e5-98c159bd0278",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create temporary file.",
			Cause:  err.Error(),
		}
	}

	r.body = tmpROFile
	r.readerAt = tmpROFile
	r.tmpFilePath = tmpFilePath

	goto READAT
}
