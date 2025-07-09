package user

import (
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"abodemine/domains/arc"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/extutils"
	"abodemine/lib/val"
)

type Repository interface {
	InsertUserRecord(r *arc.Request, in *InsertUserRecordInput) (*InsertUserRecordOutput, error)
	SelectUserRecord(r *arc.Request, in *SelectUserRecordInput) (*SelectUserRecordOutput, error)
}

type repository struct{}

type NewRepositoryInput struct{}

func NewRepository(in *NewRepositoryInput) *repository {
	return &repository{}
}

type InsertUserRecordInput struct {
	Record *User
}

type InsertUserRecordOutput struct {
	Record *User
}

func (rep *repository) InsertUserRecord(r *arc.Request, in *InsertUserRecordInput) (*InsertUserRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("users").
		Columns(
			"id",
			"created_at",
			"updated_at",
			"meta",
			"organization_id",
			"username",
			"email",
			"role_id",
			"external_id",
		).
		Values(
			in.Record.Id,
			in.Record.CreatedAt,
			in.Record.UpdatedAt,
			in.Record.Meta,
			in.Record.OrganizationId,
			in.Record.Username,
			in.Record.Email,
			func() any {
				if in.Record.RoleId != uuid.Nil {
					return in.Record.RoleId
				}

				if in.Record.RoleName != "" {
					return squirrel.Expr("(select id from roles where name = ?)", in.Record.RoleName)
				}

				return ""
			}(),
			in.Record.ExternalId,
		).
		Suffix(`
			RETURNING
				id,
				created_at,
				updated_at,
				meta,
				organization_id,
				username,
				email,
				role_id,
				(select name from roles where id = role_id) as role_name,
				external_id
		`)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "e8e49c05-e0c0-41d7-ad4f-bf6cd9f60f2e",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresSaas, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "581ee592-a0be-4599-b40d-3bcea8b7c791")
	}

	var (
		email      pgtype.Text
		roleId     pgtype.UUID
		roleName   pgtype.Text
		externalId pgtype.Text
	)

	record := new(User)

	if err := row.Scan(
		&record.Id,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.Meta,
		&record.OrganizationId,
		&record.Username,
		&email,
		&roleId,
		&roleName,
		&externalId,
	); err != nil {
		return nil, &errors.Object{
			Id:     "f1b828d1-2cc7-4ede-a785-1941ec7d7efa",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to insert row.",
			Cause:  err.Error(),
		}
	}

	if email.Valid {
		record.Email = email.String
	}

	if roleId.Valid {
		record.RoleId, err = val.UUIDFromBytes(roleId.Bytes[:])
		if err != nil {
			return nil, errors.Forward(err, "585dd6ed-b5c7-4895-95e4-d593c7262daf")
		}
	}

	if roleName.Valid {
		record.RoleName = roleName.String
	}

	if externalId.Valid {
		record.ExternalId = externalId.String
	}

	out := &InsertUserRecordOutput{
		Record: record,
	}

	return out, nil
}

type SelectUserRecordInput struct {
	OrganizationId uuid.UUID
	Id             uuid.UUID
	ExternalId     string
}

type SelectUserRecordOutput struct {
	Record *User
}

func (rep *repository) SelectUserRecord(r *arc.Request, in *SelectUserRecordInput) (*SelectUserRecordOutput, error) {
	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select(
			"users.id",
			"users.organization_id",
			"users.username",
			"roles.name",
		).
		From("users").
		Join("roles ON roles.id = users.role_id").
		Where("users.organization_id = ?", in.OrganizationId)

	if in.Id != uuid.Nil {
		builder = builder.Where("users.id = ?", in.Id)
	}

	if in.ExternalId != "" {
		builder = builder.Where("users.external_id = ?", in.ExternalId)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, &errors.Object{
			Id:     "cc1ff6a0-4cba-4ab1-b020-01d6a9788434",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to build SQL.",
			Cause:  err.Error(),
		}
	}

	row, err := extutils.PgxQueryRow(r, consts.ConfigKeyPostgresSaas, sql, args)
	if err != nil {
		return nil, errors.Forward(err, "52548080-87d0-42e2-bd53-8b5a8675e2db")
	}

	record := new(User)

	if err := row.Scan(
		&record.Id,
		&record.OrganizationId,
		&record.Username,
		&record.RoleName,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &errors.Object{
				Id:     "d815c177-862e-4fe3-a40b-679ac17a93d4",
				Code:   errors.Code_NOT_FOUND,
				Detail: "Record not found.",
			}
		}

		return nil, &errors.Object{
			Id:     "2e7b0167-5cc1-4566-b21b-ee855c8ad109",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to fetch row.",
			Cause:  err.Error(),
		}
	}

	out := &SelectUserRecordOutput{
		Record: record,
	}

	return out, nil
}
