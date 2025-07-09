package lambda

import (
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"abodemine/domains/arc"
	"abodemine/lib/errors"
	"abodemine/lib/val"
	"abodemine/projects/packer/conf"
)

type Domain interface {
	HandleSecureDownloadLambdaEvent(r *arc.Request, in *HandleSecureDownloadLambdaEventInput) (*HandleSecureDownloadLambdaEventOutput, error)
}

type domain struct {
	config *conf.Config
}

type NewDomainInput struct {
	Config *conf.Config
}

func NewDomain(in *NewDomainInput) Domain {
	return &domain{
		config: in.Config,
	}
}

type HandleSecureDownloadLambdaEventInput struct {
	Event *events.APIGatewayV2HTTPRequest
}

type HandleSecureDownloadLambdaEventOutput struct {
	PresignedURL string
}

func (dom *domain) HandleSecureDownloadLambdaEvent(r *arc.Request, in *HandleSecureDownloadLambdaEventInput) (*HandleSecureDownloadLambdaEventOutput, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "8db16a44-2b71-4379-84fd-3821467cb4c4",
			Code:   errors.Code_INTERNAL,
			Detail: "Nil input.",
		}
	}

	switch {
	case dom.config == nil:
		return nil, &errors.Object{
			Id:     "8b4582f5-3c3c-49fb-805a-d8c8d5116c61",
			Code:   errors.Code_INTERNAL,
			Detail: "Nil config.",
		}
	case dom.config.File == nil:
		return nil, &errors.Object{
			Id:     "1eee0dc7-0ffd-4fe4-871e-bbfeeb6fd60d",
			Code:   errors.Code_INTERNAL,
			Detail: "Nil config file.",
		}
	case dom.config.File.Lambdas == nil:
		return nil, &errors.Object{
			Id:     "592fbeee-e093-4cb7-85c9-9bbd20a38987",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing Lambdas configuration.",
		}
	case dom.config.File.Lambdas.SecureDownload == nil:
		return nil, &errors.Object{
			Id:     "93cfaca7-1b12-47a9-a907-b0387c33ab96",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing Lambdas.SecureDownload configuration.",
		}
	case len(dom.config.File.Lambdas.SecureDownload.DynamodbTables) == 0:
		return nil, &errors.Object{
			Id:     "c0f2a1b4-3d5e-4b8f-9c7d-6a0e2f3b5c7d",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing Lambdas.SecureDownload.DynamodbTables configuration.",
		}
	case len(dom.config.File.Lambdas.SecureDownload.S3Buckets) == 0:
		return nil, &errors.Object{
			Id:     "a16895b3-585b-49cf-90b0-3cd2cefac071",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing Lambdas.SecureDownload.S3Buckets configuration.",
		}
	}

	secureDownloadTable, ok := dom.config.File.Lambdas.SecureDownload.DynamodbTables["secure-download"]
	if !ok {
		return nil, &errors.Object{
			Id:     "b418998f-0e50-4bc6-ad6e-fa6affd1a908",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing secure-download table configuration.",
		}
	}

	if secureDownloadTable.Name == "" {
		return nil, &errors.Object{
			Id:     "64df81e2-d60d-4555-a688-799cbe54198f",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing secure-download table name.",
		}
	}

	s3Bucket, ok := dom.config.File.Lambdas.SecureDownload.S3Buckets["secure-download"]
	if !ok {
		return nil, &errors.Object{
			Id:     "3664e880-a5fe-4cad-9e29-14425c502757",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing secure-download bucket configuration.",
		}
	}

	if s3Bucket.Name == "" {
		return nil, &errors.Object{
			Id:     "840023bf-d709-4811-894b-db5e0f8572ea",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing secure-download bucket name.",
		}
	}

	event := in.Event

	if event == nil {
		return nil, &errors.Object{
			Id:     "63612935-b2b3-4ba5-bab3-06b60bb31da9",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing event.",
		}
	}

	authHeader := val.Coalesce(event.Headers["Authorization"], event.Headers["authorization"])
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if token == "" {
		return nil, &errors.Object{
			Id:     "83bf18c0-6037-407e-b9ca-d98b537f1a17",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing auth token.",
		}
	}

	s3ObjectName := val.Coalesce(event.Headers["X-AbodeMine-S3-Object"], event.Headers["x-abodemine-s3-object"])

	if s3ObjectName == "" {
		return nil, &errors.Object{
			Id:     "cf351613-15ca-4743-bafa-c3e1553fe91e",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing S3 object name.",
		}
	}

	dynamodbClient := dynamodb.NewFromConfig(dom.config.AWS.Get("default"))

	getItemOut, err := dynamodbClient.GetItem(r.Context(), &dynamodb.GetItemInput{
		TableName: &secureDownloadTable.Name,
		Key: map[string]types.AttributeValue{
			"token": &types.AttributeValueMemberS{Value: token},
		},
	})
	if err != nil {
		return nil, &errors.Object{
			Id:     "5ecc482a-cebe-4d3e-bec9-a766a4ed761b",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to get DynamoDB item.",
			Cause:  err.Error(),
		}
	}

	if getItemOut.Item == nil {
		return nil, &errors.Object{
			Id:   "ec67047c-c40e-460e-afe0-787767dd2fce",
			Code: errors.Code_UNAUTHENTICATED,
		}
	}

	s3Client := s3.NewFromConfig(dom.config.AWS.Get("default"))
	presignClient := s3.NewPresignClient(s3Client)

	presignResult, err := presignClient.PresignGetObject(
		r.Context(),
		&s3.GetObjectInput{
			Bucket: &s3Bucket.Name,
			Key:    &s3ObjectName,
		},
		func(opts *s3.PresignOptions) {
			opts.Expires = 60 * time.Second
		})
	if err != nil {
		return nil, &errors.Object{
			Id:     "88a6437d-9e26-4c11-b54e-812c324f9b9e",
			Code:   errors.Code_INTERNAL,
			Detail: "Failed to presign S3 get object.",
			Cause:  err.Error(),
		}
	}

	out := &HandleSecureDownloadLambdaEventOutput{
		PresignedURL: presignResult.URL,
	}

	return out, nil
}
