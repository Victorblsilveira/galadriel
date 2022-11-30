// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: join_tokens.sql

package datastore

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
)

const createJoinToken = `-- name: CreateJoinToken :one
INSERT INTO join_tokens(token, expires_at, trust_domain_id)
VALUES ($1, $2, $3)
RETURNING id, trust_domain_id, token, used, expires_at, created_at, updated_at
`

type CreateJoinTokenParams struct {
	Token         string
	ExpiresAt     time.Time
	TrustDomainID pgtype.UUID
}

func (q *Queries) CreateJoinToken(ctx context.Context, arg CreateJoinTokenParams) (JoinToken, error) {
	row := q.queryRow(ctx, q.createJoinTokenStmt, createJoinToken, arg.Token, arg.ExpiresAt, arg.TrustDomainID)
	var i JoinToken
	err := row.Scan(
		&i.ID,
		&i.TrustDomainID,
		&i.Token,
		&i.Used,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteJoinToken = `-- name: DeleteJoinToken :exec
DELETE
FROM join_tokens
WHERE id = $1
`

func (q *Queries) DeleteJoinToken(ctx context.Context, id pgtype.UUID) error {
	_, err := q.exec(ctx, q.deleteJoinTokenStmt, deleteJoinToken, id)
	return err
}

const findJoinToken = `-- name: FindJoinToken :one
SELECT id, trust_domain_id, token, used, expires_at, created_at, updated_at
FROM join_tokens
WHERE token = $1
`

func (q *Queries) FindJoinToken(ctx context.Context, token string) (JoinToken, error) {
	row := q.queryRow(ctx, q.findJoinTokenStmt, findJoinToken, token)
	var i JoinToken
	err := row.Scan(
		&i.ID,
		&i.TrustDomainID,
		&i.Token,
		&i.Used,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findJoinTokenByID = `-- name: FindJoinTokenByID :one
SELECT id, trust_domain_id, token, used, expires_at, created_at, updated_at
FROM join_tokens
WHERE id = $1
`

func (q *Queries) FindJoinTokenByID(ctx context.Context, id pgtype.UUID) (JoinToken, error) {
	row := q.queryRow(ctx, q.findJoinTokenByIDStmt, findJoinTokenByID, id)
	var i JoinToken
	err := row.Scan(
		&i.ID,
		&i.TrustDomainID,
		&i.Token,
		&i.Used,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findJoinTokensByTrustDomainID = `-- name: FindJoinTokensByTrustDomainID :many
SELECT id, trust_domain_id, token, used, expires_at, created_at, updated_at
FROM join_tokens
WHERE trust_domain_id = $1
`

func (q *Queries) FindJoinTokensByTrustDomainID(ctx context.Context, trustDomainID pgtype.UUID) ([]JoinToken, error) {
	rows, err := q.query(ctx, q.findJoinTokensByTrustDomainIDStmt, findJoinTokensByTrustDomainID, trustDomainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []JoinToken
	for rows.Next() {
		var i JoinToken
		if err := rows.Scan(
			&i.ID,
			&i.TrustDomainID,
			&i.Token,
			&i.Used,
			&i.ExpiresAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listJoinTokens = `-- name: ListJoinTokens :many
SELECT id, trust_domain_id, token, used, expires_at, created_at, updated_at
FROM join_tokens
ORDER BY created_at DESC
`

func (q *Queries) ListJoinTokens(ctx context.Context) ([]JoinToken, error) {
	rows, err := q.query(ctx, q.listJoinTokensStmt, listJoinTokens)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []JoinToken
	for rows.Next() {
		var i JoinToken
		if err := rows.Scan(
			&i.ID,
			&i.TrustDomainID,
			&i.Token,
			&i.Used,
			&i.ExpiresAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateJoinToken = `-- name: UpdateJoinToken :one
UPDATE join_tokens
    SET used = $2,
        updated_at = now()
WHERE id = $1
RETURNING id, trust_domain_id, token, used, expires_at, created_at, updated_at
`

type UpdateJoinTokenParams struct {
	ID   pgtype.UUID
	Used bool
}

func (q *Queries) UpdateJoinToken(ctx context.Context, arg UpdateJoinTokenParams) (JoinToken, error) {
	row := q.queryRow(ctx, q.updateJoinTokenStmt, updateJoinToken, arg.ID, arg.Used)
	var i JoinToken
	err := row.Scan(
		&i.ID,
		&i.TrustDomainID,
		&i.Token,
		&i.Used,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
