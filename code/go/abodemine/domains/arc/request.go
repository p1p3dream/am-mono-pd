package arc

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"abodemine/lib/errors"
	"abodemine/lib/flags"
	"abodemine/lib/val"
)

type Request struct {
	id      uuid.UUID
	ctx     context.Context
	flags   *val.Cache[uint, struct{}]
	domain  Domain
	session ServerSession

	pgxTx *val.Cache[string, pgx.Tx]
}

type CloneRequestOption func(*Request)

func CloneRequestWithContext(ctx context.Context) CloneRequestOption {
	return func(r *Request) {
		r.ctx = ctx
	}
}

func CloneRequestWithPgxTx(k string, tx pgx.Tx) CloneRequestOption {
	return func(r *Request) {
		if r.pgxTx == nil {
			r.pgxTx = val.NewCache[string, pgx.Tx]()
		}

		r.pgxTx.Set(k, tx)
	}
}

func (r *Request) Clone(options ...CloneRequestOption) *Request {
	newRequest := &Request{
		id:      r.id,
		ctx:     r.ctx,
		domain:  r.domain,
		flags:   r.flags,
		session: r.session,
	}

	for _, option := range options {
		option(newRequest)
	}

	return newRequest
}

func (r *Request) Id() uuid.UUID {
	return r.id
}

func (r *Request) Context() context.Context {
	return r.ctx
}

func (r *Request) Dom() Domain {
	return r.domain
}

func (r *Request) HasFlag(k uint) bool {
	if r.flags == nil {
		return false
	}

	return r.flags.Has(k)
}

func (r *Request) Session() ServerSession {
	return r.session
}

func (r *Request) SetSession(s ServerSession) {
	if r.session != nil {
		// This is a CRITICAL error because the session should only be set once.
		// If you're seeing this error you should ask for advice.
		panic("Session already set.")
	}

	r.session = s
}

func (r *Request) SelectPgxTx(k string) (pgx.Tx, bool) {
	if r.pgxTx == nil {
		return nil, false
	}

	return r.pgxTx.Select(k)
}

func (r *Request) CasbinEnforce(key string, vals ...any) error {
	cb, err := r.Dom().SelectCasbin(key)
	if err != nil {
		return errors.Forward(err, "a791e0c3-f92d-4eba-9150-c4a647e7189b")
	}

	rvals := append([]any{r.Session().RoleName()}, vals...)

	authOk, err := cb.Enforce(rvals...)
	if err != nil {
		return &errors.Object{
			Id:     "c32b779e-fccb-49b9-bebc-92363b1a7f20",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to enforce auth.",
			Cause:  err.Error(),
		}
	}

	if !authOk {
		return &errors.Object{
			Id:     "d425b7c6-9768-4a8a-b6f6-f64cc79c4759",
			Code:   errors.Code_UNAUTHENTICATED,
			Detail: "User is not authorized to perform this action.",
			Meta: map[string]any{
				"key":  key,
				"vals": vals,
			},
		}
	}

	return nil
}

type CreateRequestInput struct {
	Id      uuid.UUID
	Context context.Context
	Flags   []string
}

type CreateRequestOutput struct {
	Request *Request
}

func (dom *domain) CreateRequest(in *CreateRequestInput) (*Request, error) {
	r := &Request{
		domain: dom,
	}

	if len(dom.flags) > 0 || len(in.Flags) > 0 {
		r.flags = val.NewCache[uint, struct{}]()

		for _, flag := range dom.flags {
			v, ok := flags.Select(flag)
			if !ok {
				return nil, &errors.Object{
					Id:     "400911a4-dc88-43bf-a92b-0ac1df493d3d",
					Code:   errors.Code_INVALID_ARGUMENT,
					Detail: "Unknown flag.",
					Meta: map[string]any{
						"flag": flag,
					},
				}
			}

			r.flags.Set(v, struct{}{})
		}

		for _, flag := range in.Flags {
			if len(flag) == 0 {
				return nil, &errors.Object{
					Id:     "3b008eaa-f350-4a65-bb16-ff02652d81ec",
					Code:   errors.Code_INVALID_ARGUMENT,
					Detail: "Empty flag.",
				}
			}

			var removeFlag bool

			if flag[0] == '-' {
				removeFlag = true
				flag = flag[1:]
			}

			v, ok := flags.Select(flag)
			if !ok {
				return nil, &errors.Object{
					Id:     "fe2a125d-f581-4d9b-9742-b0243250147a",
					Code:   errors.Code_INVALID_ARGUMENT,
					Detail: "Unknown flag.",
					Meta: map[string]any{
						"flag": flag,
					},
				}
			}

			if removeFlag {
				r.flags.Del(v)
			} else {
				r.flags.Set(v, struct{}{})
			}
		}
	}

	if in.Id != uuid.Nil {
		r.id = in.Id
	} else {
		v, err := val.NewUUID4()
		if err != nil {
			return nil, errors.Forward(err, "89e3919a-8172-466e-836c-bba7ec73e853")
		}

		r.id = v
	}

	if in.Context != nil {
		r.ctx = in.Context
	} else {
		r.ctx = context.TODO()
	}

	return r, nil
}
