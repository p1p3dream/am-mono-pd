package distsync

import (
	"context"
	"time"
)

type LockStatusCode int

const (
	LockStatusUnknown       LockStatusCode = 0
	LockStatusNotPresent    LockStatusCode = 1
	LockStatusAcquiredRead  LockStatusCode = 2
	LockStatusAcquiredWrite LockStatusCode = 3
	LockStatusReleased      LockStatusCode = 4
	LockStatusExpiredRead   LockStatusCode = 5
	LockStatusExpiredWrite  LockStatusCode = 6
)

func (c LockStatusCode) String() string {
	switch c {
	case LockStatusUnknown:
		return "UNKNOWN"
	case LockStatusNotPresent:
		return "NOT_PRESENT"
	case LockStatusAcquiredRead:
		return "ACQUIRED_READ"
	case LockStatusAcquiredWrite:
		return "ACQUIRED_WRITE"
	case LockStatusReleased:
		return "RELEASED"
	case LockStatusExpiredRead:
		return "EXPIRED_READ"
	case LockStatusExpiredWrite:
		return "EXPIRED_WRITE"
	default:
		return "INVALID"
	}
}

type LockType int

const (
	LockTypeRead  LockType = 0
	LockTypeWrite LockType = 1
)

type Lock struct {
	Id   string
	Type LockType
	Ttl  time.Duration

	// Do not wait for the lock to be released.
	NoPolling bool

	// Do not create the lock if it has readers.
	NoReaders bool
}

type LockStatus struct {
	Code            LockStatusCode
	ReaderCount     int64
	ReaderExpiresAt time.Time
	WriterExpiresAt time.Time
}

type Locker interface {
	Extend(ctx context.Context) error
	Lock(ctx context.Context, lock *Lock) error
	Status(ctx context.Context, id string) (*LockStatus, error)
	Unlock(ctx context.Context) error
}
