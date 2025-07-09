package errors

import (
	"strings"

	"github.com/rs/zerolog/log"
)

type Chain struct {
	Objects []*Object `json:"objects,omitempty"`
}

func NewChain(err ...*Object) *Chain {
	return &Chain{
		Objects: err,
	}
}

func (c *Chain) Add(err ...*Object) {
	c.Objects = append(c.Objects, err...)
}

func (c *Chain) Empty() bool {
	return c == nil || len(c.Objects) == 0
}

func (c *Chain) Copy() []*Object {
	return c.Objects[:]
}

func (c *Chain) First() *Object {
	if len(c.Objects) == 0 {
		return nil
	}

	return c.Objects[0]
}

func (c *Chain) Last() *Object {
	if len(c.Objects) == 0 {
		// Ensure that we always return an error.
		// If you find yourself here, you probably have a chain bug.
		return &Object{
			Id:   "04ba6bbb-150c-4ab4-968a-ac8cfcc8603a",
			Code: Code_INTERNAL,
		}
	}

	return c.Objects[len(c.Objects)-1]
}

func (c *Chain) Error() string {
	objects := c.Objects
	s := new(strings.Builder)

	for i := len(objects) - 1; i >= 0; i-- {
		o := objects[i]

		s.WriteString(o.Error())

		if i > 0 {
			s.WriteString(" => ")
		}
	}

	return s.String()
}

func AsChain(err error) *Chain {
	switch v := err.(type) {
	case *Object:
		return NewChain(v)
	case *Chain:
		return v
	default:
		return NewChain(&Object{
			Id:    "f65d986a-7a01-4442-8f3c-d68df0793b9b",
			Code:  Code_UNKNOWN,
			Cause: err.Error(),
		})
	}
}

// Forward should be used when we want to return an existing error while ensuring traceability.
func Forward(err error, id string) *Chain {
	var chain *Chain

	switch v := err.(type) {
	case *Object:
		chain = NewChain(v)
	case *Chain:
		chain = v
	default:
		chain = NewChain(&Object{
			Id:    "7b024d1e-fe83-4999-ba04-16ce8dd62657",
			Code:  Code_UNKNOWN,
			Cause: err.Error(),
		})
	}

	chain.Add(&Object{
		Id: id,
	})

	return chain
}

// Wrap adds an error to the chain.
func Wrap(err error, wrapper *Object) *Chain {
	var chain *Chain

	switch v := err.(type) {
	case *Object:
		chain = NewChain(v)
	case *Chain:
		chain = v
	default:
		chain = NewChain(&Object{
			Id:    "3a7894b8-a8c1-427a-b5b3-36565328c5ec",
			Code:  Code_UNKNOWN,
			Cause: err.Error(),
		})
	}

	chain.Add(wrapper)

	return chain
}

func Sanitize(e error, debug bool) *Chain {
	if e == nil {
		return NewChain(Internal("9aa08800-1d37-4dfe-b0eb-12494cdc938f"))
	}

	var chain *Chain

	switch v := e.(type) {
	case *Object:
		chain = NewChain(v)
	case *Chain:
		if v.Empty() {
			return NewChain(Internal("66909f9e-93a6-49be-aca3-e8d46838654f"))
		}

		chain = v
	default:
		chain = NewChain(&Object{
			Id:    "dfe782d7-74a4-425f-8740-eb64197a133a",
			Code:  Code_UNKNOWN,
			Cause: e.Error(),
		})
	}

	var firstError *Object
	var outObjects []*Object

	objects := chain.Objects

	for i := len(objects) - 1; i >= 0; i-- {
		object := objects[i]

		if object == nil {
			// To-Do: full stack trace log.

			log.Warn().
				Str("id", "38a23afc-167c-46db-9ef4-4a4794d931a7").
				Msg("Nil error object found in chain.")

			continue
		}

		isInternalError := object.Detail == "" ||
			object.Code == Code_UNKNOWN ||
			object.Code == Code_INTERNAL

		// Ensure we don't change the original error for post-processing.
		objectCopy := *object

		if !isInternalError && firstError == nil {
			firstError = &objectCopy
		}

		if !debug {
			objectCopy.Cause = ""

			if isInternalError {
				continue
			}
		}

		outObjects = append(outObjects, &objectCopy)
	}

	if len(outObjects) == 0 {
		return NewChain(Internal("6abda17f-f7f8-461e-929e-019379fa9e24"))
	}

	if firstError == nil {
		// We should always return a readable (sanitized) error.
		// If we can't find one, we should return a generic error
		// and log this issue so we can improve the chain with
		// more information.

		log.Warn().
			Str("id", "1c0386a8-cba2-4764-ae9f-96d97ac90351").
			Err(e).
			Msg("No first error found in chain.")

		firstError = Internal("25867a84-7cdb-44c5-ac0e-5203220da30c")
	}

	if !debug {
		return &Chain{Objects: []*Object{firstError}}
	}

	// Ensure the client always receive the sanitized error first, then
	// the full chain if enabled.
	return &Chain{Objects: append([]*Object{firstError}, outObjects...)}
}

func Unauthenticated(err error, id string) *Chain {
	return Wrap(err, &Object{
		Id:     id,
		Code:   Code_UNAUTHENTICATED,
		Detail: "Unauthenticated.",
	})
}
