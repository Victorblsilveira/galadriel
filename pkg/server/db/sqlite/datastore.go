package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/HewlettPackard/galadriel/pkg/common/entity"
	"github.com/HewlettPackard/galadriel/pkg/server/db"
	"github.com/HewlettPackard/galadriel/pkg/server/db/criteria"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// Datastore is a SQL database accessor that provides convenient methods
// to perform CRUD operations for Galadriel entities.
// It implements the Datastore interface.
type Datastore struct {
	db      *sql.DB
	querier Querier
}

// NewDatastore creates a new instance of a Datastore object that connects to an SQLite database
// parsing the connString.
// The connString should be a file path to the SQLite database file.
func NewDatastore(connString string) (*Datastore, error) {
	openDB, err := sql.Open("sqlite3", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// enable foreign key constraint enforcement
	_, err = openDB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign key constraint enforcement: %w", err)
	}

	// validates if the schema in the DB matches the schema supported by the app, and runs the migrations if needed
	if err = validateAndMigrateSchema(openDB); err != nil {
		return nil, fmt.Errorf("failed to validate or migrate schema: %w", err)
	}

	return &Datastore{
		db:      openDB,
		querier: New(openDB),
	}, nil
}

func (d *Datastore) Close() error {
	return d.db.Close()
}

// CreateOrUpdateTrustDomain creates or updates the given TrustDomain in the underlying db, based on
// whether the given entity has an ID, in which case, it is updated.
func (d *Datastore) CreateOrUpdateTrustDomain(ctx context.Context, req *entity.TrustDomain) (*entity.TrustDomain, error) {
	if req.Name.String() == "" {
		return nil, errors.New("trustDomain trust domain is missing")
	}

	var trustDomain *TrustDomain
	var err error
	if req.ID.Valid {
		trustDomain, err = d.updateTrustDomain(ctx, req)
	} else {
		trustDomain, err = d.createTrustDomain(ctx, req)
	}
	if err != nil {
		return nil, err
	}

	response, err := trustDomain.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed converting trustDomain model to entity: %w", err)
	}

	return response, nil
}

func (d *Datastore) DeleteTrustDomain(ctx context.Context, trustDomainID uuid.UUID) error {
	if err := d.querier.DeleteTrustDomain(ctx, trustDomainID.String()); err != nil {
		return fmt.Errorf("failed deleting trust domain with ID=%q: %w", trustDomainID, err)
	}

	return nil
}

func (d *Datastore) ListTrustDomains(ctx context.Context, criteria *criteria.ListTrustDomainCriteria) ([]*entity.TrustDomain, error) {
	rows, err := db.ExecuteListTrustDomainQuery(ctx, d.db, criteria)
	if err != nil {
		return nil, fmt.Errorf("failed getting trust domain list: %w", err)
	}
	defer rows.Close()

	var domains []TrustDomain
	for rows.Next() {
		var t TrustDomain
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		domains = append(domains, t)
	}

	return trustDomainToEntity(domains)
}

func (d *Datastore) FindTrustDomainByID(ctx context.Context, trustDomainID uuid.UUID) (*entity.TrustDomain, error) {
	m, err := d.querier.FindTrustDomainByID(ctx, trustDomainID.String())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed looking up trust domain for ID=%q: %w", trustDomainID, err)
	}

	r, err := m.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed converting model trust domain to entity: %w", err)
	}

	return r, nil
}

func (d *Datastore) FindTrustDomainByName(ctx context.Context, name spiffeid.TrustDomain) (*entity.TrustDomain, error) {
	trustDomain, err := d.querier.FindTrustDomainByName(ctx, name.String())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed looking up trust domain for Trust Domain=%q: %w", name, err)
	}

	r, err := trustDomain.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed converting model trust domain to entity: %w", err)
	}

	return r, nil
}

func (d *Datastore) CreateOrUpdateBundle(ctx context.Context, req *entity.Bundle) (*entity.Bundle, error) {
	var bundle *Bundle
	var err error
	if req.ID.Valid {
		bundle, err = d.updateBundle(ctx, req)
	} else {
		bundle, err = d.createBundle(ctx, req)
	}
	if err != nil {
		return nil, err
	}

	response, err := bundle.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed converting trust domain model to entity: %w", err)
	}

	return response, nil
}

func (d *Datastore) FindBundleByID(ctx context.Context, bundleID uuid.UUID) (*entity.Bundle, error) {
	bundle, err := d.querier.FindBundleByID(ctx, bundleID.String())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed looking up bundle with ID=%q: %w", bundleID, err)
	}

	b, err := bundle.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed converting model bundle to entity: %w", err)
	}

	return b, nil
}

func (d *Datastore) FindBundleByTrustDomainID(ctx context.Context, trustDomainID uuid.UUID) (*entity.Bundle, error) {
	trustDomain, err := d.querier.FindBundleByTrustDomainID(ctx, trustDomainID.String())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed looking up bundle for ID=%q: %w", trustDomainID, err)
	}

	td, err := trustDomain.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed converting model bundle to entity: %w", err)
	}

	return td, nil
}

func (d *Datastore) ListBundles(ctx context.Context) ([]*entity.Bundle, error) {
	bundles, err := d.querier.ListBundles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting bundle list: %w", err)
	}

	result := make([]*entity.Bundle, len(bundles))
	for i, m := range bundles {
		r, err := m.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed converting model bundle to entity: %w", err)
		}
		result[i] = r
	}

	return result, nil
}

func (d *Datastore) DeleteBundle(ctx context.Context, bundleID uuid.UUID) error {
	if err := d.querier.DeleteBundle(ctx, bundleID.String()); err != nil {
		return fmt.Errorf("failed deleting bundle with ID=%q: %w", bundleID, err)
	}

	return nil
}

func (d *Datastore) CreateJoinToken(ctx context.Context, req *entity.JoinToken) (*entity.JoinToken, error) {
	id := uuid.New()
	params := CreateJoinTokenParams{
		ID:            id.String(),
		Token:         req.Token,
		ExpiresAt:     req.ExpiresAt,
		TrustDomainID: req.TrustDomainID.String(),
	}
	joinToken, err := d.querier.CreateJoinToken(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed creating join token: %w", err)
	}

	ent, err := joinToken.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed converting model join token to entity: %w", err)
	}

	return ent, nil
}

func (d *Datastore) FindJoinTokensByID(ctx context.Context, joinTokenID uuid.UUID) (*entity.JoinToken, error) {
	joinToken, err := d.querier.FindJoinTokenByID(ctx, joinTokenID.String())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed looking up join token with ID=%q: %w", joinTokenID, err)
	}

	ent, err := joinToken.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed converting model join token to entity: %w", err)
	}

	return ent, nil
}

func (d *Datastore) FindJoinTokensByTrustDomainID(ctx context.Context, trustDomainID uuid.UUID) ([]*entity.JoinToken, error) {
	tokens, err := d.querier.FindJoinTokensByTrustDomainID(ctx, trustDomainID.String())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed looking up join token for Name ID=%q: %w", trustDomainID, err)
	}

	result := make([]*entity.JoinToken, len(tokens))
	for i, t := range tokens {
		ent, err := t.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed converting model join token to entity: %w", err)
		}
		result[i] = ent
	}

	return result, nil
}

func (d *Datastore) ListJoinTokens(ctx context.Context) ([]*entity.JoinToken, error) {
	tokens, err := d.querier.ListJoinTokens(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed looking up join tokens: %w", err)
	}

	result := make([]*entity.JoinToken, len(tokens))
	for i, t := range tokens {
		ent, err := t.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed converting model join token to entity: %w", err)
		}
		result[i] = ent
	}

	return result, nil
}

func (d *Datastore) UpdateJoinToken(ctx context.Context, joinTokenID uuid.UUID, used bool) (*entity.JoinToken, error) {
	params := UpdateJoinTokenParams{
		ID:   joinTokenID.String(),
		Used: used,
	}

	jt, err := d.querier.UpdateJoinToken(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed updating join token with ID=%q, %w", joinTokenID, err)
	}

	ent, err := jt.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed converting model join token to entity: %w", err)
	}

	return ent, nil
}

func (d *Datastore) DeleteJoinToken(ctx context.Context, joinTokenID uuid.UUID) error {
	if err := d.querier.DeleteJoinToken(ctx, joinTokenID.String()); err != nil {
		return fmt.Errorf("failed deleting join token with ID=%q, %w", joinTokenID, err)
	}

	return nil
}

func (d *Datastore) FindJoinToken(ctx context.Context, token string) (*entity.JoinToken, error) {
	joinToken, err := d.querier.FindJoinToken(ctx, token)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed looking up join token: %w", err)
	}

	ent, err := joinToken.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed converting model join token to entity: %w", err)
	}

	return ent, nil
}

func (d *Datastore) CreateOrUpdateRelationship(ctx context.Context, req *entity.Relationship) (*entity.Relationship, error) {
	var relationship *Relationship
	var err error
	if req.ID.Valid {
		relationship, err = d.updateRelationship(ctx, req)
	} else {
		relationship, err = d.createRelationship(ctx, req)
	}
	if err != nil {
		return nil, err
	}

	response, err := relationship.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed converting relationship model to entity: %w", err)
	}

	return response, nil
}

func (d *Datastore) FindRelationshipByID(ctx context.Context, relationshipID uuid.UUID) (*entity.Relationship, error) {
	relationship, err := d.querier.FindRelationshipByID(ctx, relationshipID.String())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed looking up relationship for ID=%q: %w", relationshipID, err)
	}

	response, err := relationship.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed converting relationship model to entity: %w", err)
	}

	return response, nil
}

func (d *Datastore) FindRelationshipsByTrustDomainID(
	ctx context.Context,
	trustDomainID uuid.UUID,
) ([]*entity.Relationship, error) {

	params := FindRelationshipsByTrustDomainIDParams{
		TrustDomainAID: trustDomainID.String(),
		TrustDomainBID: trustDomainID.String(),
	}
	relationships, err := d.querier.FindRelationshipsByTrustDomainID(ctx, params)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed looking up relationships for TrustDomainID %q: %w", trustDomainID, err)
	}

	result := make([]*entity.Relationship, len(relationships))
	for i, m := range relationships {
		ent, err := m.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed converting relationship model to entity: %w", err)
		}
		result[i] = ent
	}

	return result, nil
}

func (d *Datastore) ListRelationships(ctx context.Context, criteria *criteria.ListRelationshipsCriteria) ([]*entity.Relationship, error) {
	rows, err := db.ExecuteListRelationshipsQuery(ctx, d.db, criteria, db.SQLite)
	if err != nil {
		return nil, fmt.Errorf("failed looking up relationships: %w", err)
	}
	defer rows.Close()

	var relationships []Relationship
	for rows.Next() {
		var m Relationship
		if err := rows.Scan(&m.ID, &m.TrustDomainAID, &m.TrustDomainBID, &m.TrustDomainAConsent, &m.TrustDomainBConsent, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		relationships = append(relationships, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during row iteration: %w", err)
	}

	return relationshipsToEntity(relationships)
}

func (d *Datastore) DeleteRelationship(ctx context.Context, relationshipID uuid.UUID) error {
	if err := d.querier.DeleteRelationship(ctx, relationshipID.String()); err != nil {
		return fmt.Errorf("failed deleting relationship ID=%q: %w", relationshipID, err)
	}

	return nil
}

func (d *Datastore) createTrustDomain(ctx context.Context, req *entity.TrustDomain) (*TrustDomain, error) {
	id := uuid.New()
	params := CreateTrustDomainParams{
		ID:   id.String(),
		Name: req.Name.String(),
	}
	if req.Description != "" {
		params.Description = sql.NullString{
			String: req.Description,
			Valid:  true,
		}
	} else {
		params.Description = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	td, err := d.querier.CreateTrustDomain(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed creating new trust domain: %w", err)
	}
	return &td, nil
}

func (d *Datastore) updateTrustDomain(ctx context.Context, req *entity.TrustDomain) (*TrustDomain, error) {
	params := UpdateTrustDomainParams{
		ID: req.ID.UUID.String(),
	}

	if req.Description != "" {
		params.Description = sql.NullString{
			String: req.Description,
			Valid:  true,
		}
	}

	td, err := d.querier.UpdateTrustDomain(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed updating trust domain: %w", err)
	}
	return &td, nil
}

func (d *Datastore) createRelationship(ctx context.Context, req *entity.Relationship) (*Relationship, error) {
	id := uuid.New()

	if req.TrustDomainAConsent == "" {
		req.TrustDomainAConsent = entity.ConsentStatusPending
	}
	if req.TrustDomainBConsent == "" {
		req.TrustDomainBConsent = entity.ConsentStatusPending
	}
	if req.CreatedAt.IsZero() {
		req.CreatedAt = time.Now()
	}
	if req.UpdatedAt.IsZero() {
		req.UpdatedAt = time.Now()
	}
	params := CreateRelationshipParams{
		ID:                  id.String(),
		TrustDomainAID:      req.TrustDomainAID.String(),
		TrustDomainBID:      req.TrustDomainBID.String(),
		TrustDomainAConsent: string(req.TrustDomainAConsent),
		TrustDomainBConsent: string(req.TrustDomainBConsent),
		CreatedAt:           req.CreatedAt,
		UpdatedAt:           req.UpdatedAt,
	}

	relationship, err := d.querier.CreateRelationship(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed creating new relationship: %w", err)
	}

	return &relationship, nil
}

func (d *Datastore) updateRelationship(ctx context.Context, req *entity.Relationship) (*Relationship, error) {
	params := UpdateRelationshipParams{
		ID:                  req.ID.UUID.String(),
		TrustDomainAConsent: string(req.TrustDomainAConsent),
		TrustDomainBConsent: string(req.TrustDomainBConsent),
	}

	relationship, err := d.querier.UpdateRelationship(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed updating relationship: %w", err)
	}

	return &relationship, nil
}

func (d *Datastore) createBundle(ctx context.Context, req *entity.Bundle) (*Bundle, error) {
	id := uuid.New()
	params := CreateBundleParams{
		ID:                 id.String(),
		Data:               req.Data,
		Digest:             req.Digest,
		Signature:          req.Signature,
		SigningCertificate: req.SigningCertificate,
		TrustDomainID:      req.TrustDomainID.String(),
	}

	bundle, err := d.querier.CreateBundle(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed creating new bundle: %w", err)
	}

	return &bundle, nil
}

func (d *Datastore) updateBundle(ctx context.Context, req *entity.Bundle) (*Bundle, error) {
	params := UpdateBundleParams{
		ID:                 req.ID.UUID.String(),
		Data:               req.Data,
		Digest:             req.Digest,
		Signature:          req.Signature,
		SigningCertificate: req.SigningCertificate,
	}

	bundle, err := d.querier.UpdateBundle(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed updating bundle: %w", err)
	}

	return &bundle, nil
}

func relationshipsToEntity(models []Relationship) ([]*entity.Relationship, error) {
	result := make([]*entity.Relationship, len(models))

	for i, m := range models {
		ent, err := m.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed converting relationship model to entity: %w", err)
		}
		result[i] = ent
	}

	return result, nil
}

func trustDomainToEntity(models []TrustDomain) ([]*entity.TrustDomain, error) {
	result := make([]*entity.TrustDomain, len(models))

	for i, m := range models {
		r, err := m.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed converting model trust domain to entity: %v", err)
		}
		result[i] = r
	}

	return result, nil
}
