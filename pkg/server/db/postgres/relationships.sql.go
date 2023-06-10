// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: relationships.sql

package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
)

const createRelationship = `-- name: CreateRelationship :one
INSERT INTO relationships(trust_domain_a_id, trust_domain_b_id, trust_domain_a_consent, trust_domain_b_consent, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, trust_domain_a_id, trust_domain_b_id, trust_domain_a_consent, trust_domain_b_consent, created_at, updated_at
`

type CreateRelationshipParams struct {
	TrustDomainAID      pgtype.UUID
	TrustDomainBID      pgtype.UUID
	TrustDomainAConsent ConsentStatus
	TrustDomainBConsent ConsentStatus
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (q *Queries) CreateRelationship(ctx context.Context, arg CreateRelationshipParams) (Relationship, error) {
	row := q.queryRow(ctx, q.createRelationshipStmt, createRelationship,
		arg.TrustDomainAID,
		arg.TrustDomainBID,
		arg.TrustDomainAConsent,
		arg.TrustDomainBConsent,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Relationship
	err := row.Scan(
		&i.ID,
		&i.TrustDomainAID,
		&i.TrustDomainBID,
		&i.TrustDomainAConsent,
		&i.TrustDomainBConsent,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteRelationship = `-- name: DeleteRelationship :exec
DELETE
FROM relationships
WHERE id = $1
`

func (q *Queries) DeleteRelationship(ctx context.Context, id pgtype.UUID) error {
	_, err := q.exec(ctx, q.deleteRelationshipStmt, deleteRelationship, id)
	return err
}

const findRelationshipByID = `-- name: FindRelationshipByID :one
SELECT id, trust_domain_a_id, trust_domain_b_id, trust_domain_a_consent, trust_domain_b_consent, created_at, updated_at
FROM relationships
WHERE id = $1
`

func (q *Queries) FindRelationshipByID(ctx context.Context, id pgtype.UUID) (Relationship, error) {
	row := q.queryRow(ctx, q.findRelationshipByIDStmt, findRelationshipByID, id)
	var i Relationship
	err := row.Scan(
		&i.ID,
		&i.TrustDomainAID,
		&i.TrustDomainBID,
		&i.TrustDomainAConsent,
		&i.TrustDomainBConsent,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findRelationshipsByTrustDomainID = `-- name: FindRelationshipsByTrustDomainID :many
SELECT id, trust_domain_a_id, trust_domain_b_id, trust_domain_a_consent, trust_domain_b_consent, created_at, updated_at
FROM relationships
WHERE trust_domain_a_id = $1
   OR trust_domain_b_id = $1
`

func (q *Queries) FindRelationshipsByTrustDomainID(ctx context.Context, trustDomainAID pgtype.UUID) ([]Relationship, error) {
	rows, err := q.query(ctx, q.findRelationshipsByTrustDomainIDStmt, findRelationshipsByTrustDomainID, trustDomainAID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Relationship
	for rows.Next() {
		var i Relationship
		if err := rows.Scan(
			&i.ID,
			&i.TrustDomainAID,
			&i.TrustDomainBID,
			&i.TrustDomainAConsent,
			&i.TrustDomainBConsent,
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

const updateRelationship = `-- name: UpdateRelationship :one
UPDATE relationships
SET trust_domain_a_consent = $2,
    trust_domain_b_consent = $3,
    updated_at             = now()
WHERE id = $1
RETURNING id, trust_domain_a_id, trust_domain_b_id, trust_domain_a_consent, trust_domain_b_consent, created_at, updated_at
`

type UpdateRelationshipParams struct {
	ID                  pgtype.UUID
	TrustDomainAConsent ConsentStatus
	TrustDomainBConsent ConsentStatus
}

func (q *Queries) UpdateRelationship(ctx context.Context, arg UpdateRelationshipParams) (Relationship, error) {
	row := q.queryRow(ctx, q.updateRelationshipStmt, updateRelationship, arg.ID, arg.TrustDomainAConsent, arg.TrustDomainBConsent)
	var i Relationship
	err := row.Scan(
		&i.ID,
		&i.TrustDomainAID,
		&i.TrustDomainBID,
		&i.TrustDomainAConsent,
		&i.TrustDomainBConsent,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
