package distsync

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"

	"abodemine/lib/errors"
	"abodemine/lib/val"
)

type DynamoDB struct {
	// General config.

	// The interval at which to poll for the lock.
	PollInterval time.Duration

	// DynamoDB config.

	Client    *dynamodb.Client
	TableName string

	mu   sync.Mutex
	lock *Lock
}

func (l *DynamoDB) Extend(ctx context.Context) error {
	lock := l.lock

	if lock == nil {
		return &errors.Object{
			Id:     "24d4acbf-9991-4de5-83f9-bbb528d52c24",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Lock is required.",
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	expiresAt := now.Add(lock.Ttl).Unix()

	var update expression.UpdateBuilder

	switch lock.Type {
	case LockTypeRead:
		update = expression.Set(
			expression.Name("reader_expires_at"),
			expression.Value(expiresAt),
		)
	case LockTypeWrite:
		update = expression.Set(
			expression.Name("writer_expires_at"),
			expression.Value(expiresAt),
		)
	default:
		return &errors.Object{
			Id:     "946b036f-72e0-4e58-a8d1-edeb2e32ace0",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Invalid lock type.",
			Meta: map[string]any{
				"lock_type": lock.Type,
			},
		}
	}

	expr, err := expression.NewBuilder().
		WithUpdate(update).
		Build()
	if err != nil {
		return &errors.Object{
			Id:     "31c862d4-8fe8-44fe-8317-1d9c73cf6150",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create expression.",
			Cause:  err.Error(),
			Meta: map[string]any{
				"lock_type": l.lock.Type,
			},
		}
	}

	_, err = l.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &l.TableName,
		Key: map[string]types.AttributeValue{
			"lock_id": &types.AttributeValueMemberS{Value: lock.Id},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return &errors.Object{
			Id:     "f9034cdd-2409-4098-bb19-0c392da5ae51",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to update lock item.",
			Cause:  err.Error(),
		}
	}

	return nil
}

func (l *DynamoDB) Lock(ctx context.Context, lock *Lock) error {
	if lock == nil {
		return &errors.Object{
			Id:     "650b6075-e7b6-4eb4-96a0-a72f0c2b897d",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Lock is required.",
		}
	}

	l.mu.Lock()

	if l.lock != nil {
		l.mu.Unlock()
		return &errors.Object{
			Id:     "77922016-66a7-422e-a8f4-7a98bf2e9692",
			Code:   errors.Code_FAILED_PRECONDITION,
			Detail: "Locker already in use.",
		}
	}

	l.lock = lock
	l.mu.Unlock()

	for {
		select {
		case <-ctx.Done():
			return &errors.Object{
				Id:     "f913924a-4d3b-489c-9641-c4f9df9069a1",
				Code:   errors.Code_CANCELED,
				Detail: "Context was canceled.",
				Cause:  ctx.Err().Error(),
			}
		default:
			acquired, err := l.tryAcquire(ctx)
			if err != nil {
				return errors.Forward(
					err,
					"1bc8cc4c-f537-4735-be68-1e86092a3948",
				)
			}

			if acquired {
				return nil
			}

			time.Sleep(l.PollInterval)
		}
	}
}

func (l *DynamoDB) tryAcquire(ctx context.Context) (bool, error) {
	lock := l.lock
	now := time.Now()
	expiresAt := now.Add(lock.Ttl).Unix()

	// We only care about writer status, since readers can be concurrent.
	conditions := expression.Or(
		expression.AttributeNotExists(expression.Name("lock_id")),
		expression.AttributeNotExists(expression.Name("writer_present")),
		expression.Equal(
			expression.Name("writer_present"),
			expression.Value(false),
		),
		expression.LessThan(
			expression.Name("writer_expires_at"),
			expression.Value(now.Unix()),
		),
	)

	if l.lock.NoReaders {
		conditions = expression.And(
			conditions,
			expression.Or(
				expression.AttributeNotExists(expression.Name("reader_count")),
				expression.Equal(
					expression.Name("reader_count"),
					expression.Value(0),
				),
			),
		)
	}

	var update expression.UpdateBuilder

	switch lock.Type {
	case LockTypeRead:
		update = expression.
			Set(expression.Name("reader_expires_at"), expression.Value(expiresAt)).
			Add(expression.Name("reader_count"), expression.Value(1))
	case LockTypeWrite:
		// We immediately set writer_present to true to prevent new
		// readers/writers from acquiring the lock.
		// Later we check if there are remaning readers and wait for them
		// before returning true (for the acquire op, i.e., this func).
		update = expression.
			Set(expression.Name("writer_expires_at"), expression.Value(expiresAt)).
			Set(expression.Name("writer_present"), expression.Value(true))
	default:
		return false, &errors.Object{
			Id:     "d1bc40c3-0a0d-4057-8b97-dba17dd8dd81",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Invalid lock type.",
			Meta: map[string]any{
				"lock_type": lock.Type,
			},
		}
	}

	expr, err := expression.NewBuilder().
		WithCondition(conditions).
		WithUpdate(update).
		Build()
	if err != nil {
		return false, &errors.Object{
			Id:     "8f65bdc6-3f80-4248-b4e7-ee65c886b95b",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create expression.",
			Cause:  err.Error(),
			Meta: map[string]any{
				"lock_type": l.lock.Type,
			},
		}
	}

	updateItemOut, err := l.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &l.TableName,
		Key: map[string]types.AttributeValue{
			"lock_id": &types.AttributeValueMemberS{Value: lock.Id},
		},
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueAllNew,
	})
	if err != nil {
		if ok := errors.As(
			err,
			val.PtrRef(new(types.ConditionalCheckFailedException)),
		); ok {
			if l.lock.NoPolling {
				// If we're not waiting for the lock to be released,
				// return the failed condition check as an error.
				return false, &errors.Object{
					Id:     "8f8ba2db-bcb0-4abb-ab79-4623360cff11",
					Code:   errors.Code_UNKNOWN,
					Detail: "Lock condition failed.",
					Cause:  err.Error(),
				}
			}

			return false, nil
		}

		return false, &errors.Object{
			Id:     "a91bd3bf-3413-4a79-96ae-0086ecf77966",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to update lock item.",
			Cause:  err.Error(),
		}
	}

	if lock.Type == LockTypeRead {
		return true, nil
	}

	readerCountStr, ok := updateItemOut.Attributes["reader_count"].(*types.AttributeValueMemberN)
	if !ok || readerCountStr.Value == "0" {
		// Key is new or no reader has acquired the lock.
		return true, nil
	}

WAIT_READERS:
	for {
		// Wait for reader_count to be 0 to write safely.

		log.Info().
			Str("reader_count", readerCountStr.Value).
			Msg("Waiting for readers to release lock.")

		time.Sleep(l.PollInterval)

		select {
		case <-ctx.Done():
			return false, &errors.Object{
				Id:     "19c57ac4-d323-4eec-a919-fe3b204fb5ef",
				Code:   errors.Code_CANCELED,
				Detail: "Context was canceled.",
				Cause:  ctx.Err().Error(),
			}
		default:
			getItemOut, err := l.Client.GetItem(ctx, &dynamodb.GetItemInput{
				TableName:      &l.TableName,
				ConsistentRead: val.PtrRef(true),
				Key: map[string]types.AttributeValue{
					"lock_id": &types.AttributeValueMemberS{Value: lock.Id},
				},
			})
			if err != nil {
				return false, &errors.Object{
					Id:     "10e7477b-3929-4f93-a8fa-7f36cc149e40",
					Code:   errors.Code_UNKNOWN,
					Detail: "Failed to get lock item.",
					Cause:  err.Error(),
				}
			}

			if getItemOut.Item == nil {
				// This is a critical error.
				return false, &errors.Object{
					Id:     "5c2057ab-236d-460b-8b62-764120c88578",
					Code:   errors.Code_DATA_LOSS,
					Detail: "Lock item not found.",
				}
			}

			readerCountStr, ok = getItemOut.Item["reader_count"].(*types.AttributeValueMemberN)
			if !ok {
				// Field should exists.
				return false, &errors.Object{
					Id:     "9588da0c-b448-4959-9491-a08202c9efe3",
					Code:   errors.Code_FAILED_PRECONDITION,
					Detail: "Failed to read reader_count.",
				}
			}

			readerCount, err := strconv.ParseInt(readerCountStr.Value, 10, 64)
			if err != nil {
				return false, &errors.Object{
					Id:     "75d134b5-8f44-4fc5-b35c-aba1850ad9ec",
					Code:   errors.Code_UNKNOWN,
					Detail: "Failed to parse reader_count.",
					Cause:  err.Error(),
				}
			}

			if readerCount == 0 {
				break WAIT_READERS
			}
		}
	}

	return true, nil
}

func (l *DynamoDB) Status(ctx context.Context, id string) (*LockStatus, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	getItemOut, err := l.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      &l.TableName,
		ConsistentRead: val.PtrRef(true),
		Key: map[string]types.AttributeValue{
			"lock_id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, &errors.Object{
			Id:     "26abbc60-f42e-4ff9-a4eb-122b33a9df50",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to get lock item.",
			Cause:  err.Error(),
		}
	}

	item := getItemOut.Item

	if item == nil {
		return &LockStatus{
			Code: LockStatusNotPresent,
		}, nil
	}

	now := time.Now()

	writerPresent, ok := item["writer_present"]
	if ok {
		writerPresentAttribute, ok := writerPresent.(*types.AttributeValueMemberBOOL)
		if !ok {
			return nil, &errors.Object{
				Id:     "c5f0b9c8-7e4a-4d3d-b8a1-6c2b2f0a5e3f",
				Code:   errors.Code_FAILED_PRECONDITION,
				Detail: "Failed to convert writer_present.",
			}
		}

		if writerPresentAttribute.Value {
			// If writer is present, there MUST be a writer_expires_at.
			writerExpiresAt, ok := item["writer_expires_at"]
			if !ok {
				return nil, &errors.Object{
					Id:     "7576f061-5828-4b28-8cf7-5db7f66d0c33",
					Code:   errors.Code_FAILED_PRECONDITION,
					Detail: "Failed to read writer_expires_at.",
				}
			}

			writerExpiresAtAttribute, ok := writerExpiresAt.(*types.AttributeValueMemberN)
			if !ok {
				return nil, &errors.Object{
					Id:     "940cd490-7982-42ce-a88d-1ea5c54801fb",
					Code:   errors.Code_FAILED_PRECONDITION,
					Detail: "Failed to convert writer_expires_at.",
				}
			}

			writerExpiresAtInt64, err := strconv.ParseInt(writerExpiresAtAttribute.Value, 10, 64)
			if err != nil {
				return nil, &errors.Object{
					Id:     "f8e3a5d2-1c0c-4f3b-8a9d-5c2b2f0a5e3f",
					Code:   errors.Code_UNKNOWN,
					Detail: "Failed to parse writer_expires_at.",
					Cause:  err.Error(),
				}
			}

			writerExpiresAtTime := time.Unix(writerExpiresAtInt64, 0)

			if writerExpiresAtTime.After(now) {
				return &LockStatus{
					Code:            LockStatusAcquiredWrite,
					WriterExpiresAt: writerExpiresAtTime,
				}, nil
			}

			return &LockStatus{
				Code: LockStatusExpiredWrite,
			}, nil
		}
	}

	readerCount, ok := item["reader_count"]
	if !ok {
		// Lock has only written. No reads.
		return &LockStatus{
			Code: LockStatusReleased,
		}, nil
	}

	readerCountAttribute, ok := readerCount.(*types.AttributeValueMemberN)
	if !ok {
		return nil, &errors.Object{
			Id:     "15a761b1-1326-4407-adc2-c2aad14e82cb",
			Code:   errors.Code_FAILED_PRECONDITION,
			Detail: "Failed to convert reader_count.",
		}
	}

	readerCountInt64, err := strconv.ParseInt(readerCountAttribute.Value, 10, 64)
	if err != nil {
		return nil, &errors.Object{
			Id:     "df581a02-301e-4df7-8882-d1e585f24a24",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to parse reader_count.",
			Cause:  err.Error(),
		}
	}

	if readerCountInt64 == 0 {
		return &LockStatus{
			Code: LockStatusReleased,
		}, nil
	}

	// If reader is present, there MUST be a reader_expires_at.
	readerExpiresAt, ok := item["reader_expires_at"]
	if !ok {
		return nil, &errors.Object{
			Id:     "53c9ae35-8bc3-4188-9a52-d81e5753e09a",
			Code:   errors.Code_FAILED_PRECONDITION,
			Detail: "Failed to read reader_expires_at.",
		}
	}

	readerExpiresAtAttribute, ok := readerExpiresAt.(*types.AttributeValueMemberN)
	if !ok {
		return nil, &errors.Object{
			Id:     "b8a56e2b-a15d-42fd-98d9-7ab6faa70378",
			Code:   errors.Code_FAILED_PRECONDITION,
			Detail: "Failed to convert reader_expires_at.",
		}
	}

	readerExpiresAtInt64, err := strconv.ParseInt(readerExpiresAtAttribute.Value, 10, 64)
	if err != nil {
		return nil, &errors.Object{
			Id:     "ca110768-3a53-426f-8bc6-ab7f9da717dc",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to parse reader_expires_at.",
			Cause:  err.Error(),
		}
	}

	readerExpiresAtTime := time.Unix(readerExpiresAtInt64, 0)

	if readerExpiresAtTime.After(now) {
		return &LockStatus{
			Code:            LockStatusAcquiredRead,
			ReaderCount:     readerCountInt64,
			ReaderExpiresAt: readerExpiresAtTime,
		}, nil
	}

	return &LockStatus{
		Code: LockStatusExpiredRead,
	}, nil
}

func (l *DynamoDB) Unlock(ctx context.Context) error {
	lock := l.lock

	if lock == nil {
		return &errors.Object{
			Id:     "24e3ea13-cd51-481f-b370-177d045192c2",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Lock is required.",
		}
	}

	var update expression.UpdateBuilder

	switch lock.Type {
	case LockTypeRead:
		// We don't update reader_expires_at because there might be
		// other active concurrent readers.
		update = expression.
			Add(expression.Name("reader_count"), expression.Value(-1))
	case LockTypeWrite:
		update = expression.
			Set(expression.Name("writer_present"), expression.Value(false)).
			Set(expression.Name("writer_expires_at"), expression.Value(0))
	default:
		return &errors.Object{
			Id:     "f20e1698-0e6a-43fb-bdbd-38e57e16d834",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Invalid lock type.",
			Meta: map[string]any{
				"lock_type": lock.Type,
			},
		}
	}

	expr, err := expression.NewBuilder().
		WithUpdate(update).
		Build()
	if err != nil {
		return &errors.Object{
			Id:     "1a7d5186-ee92-43d8-90a9-cac3ca877c0d",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create expression.",
			Cause:  err.Error(),
			Meta: map[string]any{
				"lock_type": l.lock.Type,
			},
		}
	}

	_, err = l.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &l.TableName,
		Key: map[string]types.AttributeValue{
			"lock_id": &types.AttributeValueMemberS{Value: lock.Id},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return &errors.Object{
			Id:     "05b80090-d2bc-4eee-a100-a44a7297b1a8",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to update lock item.",
			Cause:  err.Error(),
		}
	}

	return nil
}
