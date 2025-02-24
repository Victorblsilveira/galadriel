package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/galadriel/cmd/common/cli"
	httputil "github.com/HewlettPackard/galadriel/cmd/common/http"
	"github.com/HewlettPackard/galadriel/pkg/common/api"
	"github.com/HewlettPackard/galadriel/pkg/common/entity"
	"github.com/HewlettPackard/galadriel/pkg/server/api/admin"
	"github.com/google/uuid"
)

const (
	errorRequestFailed        = "failed to send request: %v"
	errUnmarshalRelationships = "failed to unmarshal relationships: %v"
	errUnmarshalTrustDomains  = "failed to unmarshal trust domain: %v"
	errUnmarshalJoinToken     = "failed to unmarshal join token: %v"
)

// GaladrielAPIClient represents an API client for the Galadriel Server API.
type GaladrielAPIClient interface {
	CreateTrustDomain(context.Context, api.TrustDomainName) (*entity.TrustDomain, error)
	GetTrustDomainByName(context.Context, api.TrustDomainName) (*entity.TrustDomain, error)
	ListTrustDomains(context.Context) ([]*entity.TrustDomain, error)
	DeleteTrustDomainByName(context.Context, api.TrustDomainName) error
	UpdateTrustDomainByName(context.Context, api.TrustDomainName, string) (*entity.TrustDomain, error)
	CreateRelationship(context.Context, *entity.Relationship) (*entity.Relationship, error)
	GetRelationshipByID(context.Context, uuid.UUID) (*entity.Relationship, error)
	GetRelationships(context.Context, api.ConsentStatus, api.TrustDomainName) (*entity.Relationship, error)
	GetJoinToken(context.Context, api.TrustDomainName, int32) (*entity.JoinToken, error)
}

type galadrielAdminClient struct {
	client *admin.Client
}

// NewGaladrielUDSClient creates a Galadriel API client that connects to the Galadriel Server
// using a Unix Domain Socket (UDS) specified by the socketPath parameter.
func NewGaladrielUDSClient(socketPath string, httpClient *http.Client) (GaladrielAPIClient, error) {
	baseURLOption := admin.WithBaseURL(cli.LocalhostURL)

	adminClient, err := admin.NewClient(socketPath, baseURLOption)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate the Admin Client: %v", err)
	}

	if httpClient == nil {
		httpClient = httputil.NewUDSHTTPClient(socketPath)
	}

	adminClient.Client = httpClient

	return &galadrielAdminClient{client: adminClient}, nil
}

func (g *galadrielAdminClient) GetTrustDomainByName(ctx context.Context, trustDomainName api.TrustDomainName) (*entity.TrustDomain, error) {
	res, err := g.client.GetTrustDomainByName(ctx, trustDomainName)
	if err != nil {
		return nil, fmt.Errorf(errorRequestFailed, err)
	}
	defer res.Body.Close()

	body, err := httputil.ReadResponse(res)
	if err != nil {
		return nil, err
	}

	trustDomain, err := unmarshalJSONToTrustDomain(body)
	if err != nil {
		return nil, err
	}

	return trustDomain, nil
}

func (g *galadrielAdminClient) ListTrustDomains(ctx context.Context) ([]*entity.TrustDomain, error) {
	res, err := g.client.ListTrustDomains(ctx)
	if err != nil {
		return nil, fmt.Errorf(errorRequestFailed, err)
	}
	defer res.Body.Close()

	body, err := httputil.ReadResponse(res)
	if err != nil {
		return nil, err
	}

	var trustDomains []*api.TrustDomain
	if err := json.Unmarshal(body, &trustDomains); err != nil {
		return nil, fmt.Errorf(errUnmarshalRelationships, err)
	}

	tds := make([]*entity.TrustDomain, 0, len(trustDomains))
	for i, td := range trustDomains {
		trustDomain, err := td.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed to convert trust domain %d: %v", i, err)
		}
		tds = append(tds, trustDomain)
	}

	return tds, nil
}

func (g *galadrielAdminClient) DeleteTrustDomainByName(ctx context.Context, trustDomainName api.TrustDomainName) error {
	res, err := g.client.DeleteTrustDomainByName(ctx, trustDomainName)
	if err != nil {
		return fmt.Errorf(errorRequestFailed, err)
	}
	defer res.Body.Close()

	_, err = httputil.ReadResponse(res)
	if err != nil {
		return err
	}

	return nil
}

func (g *galadrielAdminClient) UpdateTrustDomainByName(ctx context.Context, trustDomainName api.TrustDomainName, description string) (*entity.TrustDomain, error) {
	payload := api.TrustDomain{Name: trustDomainName, Description: &description}
	res, err := g.client.PutTrustDomainByName(ctx, trustDomainName, payload)
	if err != nil {
		return nil, fmt.Errorf(errorRequestFailed, err)
	}
	defer res.Body.Close()

	body, err := httputil.ReadResponse(res)
	if err != nil {
		return nil, err
	}

	trustDomain, err := unmarshalJSONToTrustDomain(body)
	if err != nil {
		return nil, err
	}

	return trustDomain, nil
}

func (g *galadrielAdminClient) CreateTrustDomain(ctx context.Context, trustDomainName api.TrustDomainName) (*entity.TrustDomain, error) {
	payload := admin.PutTrustDomainJSONRequestBody{Name: trustDomainName}

	res, err := g.client.PutTrustDomain(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf(errorRequestFailed, err)
	}
	defer res.Body.Close()

	body, err := httputil.ReadResponse(res)
	if err != nil {
		return nil, err
	}

	trustDomain, err := unmarshalJSONToTrustDomain(body)
	if err != nil {
		return nil, err
	}

	return trustDomain, nil
}

func (g *galadrielAdminClient) CreateRelationship(ctx context.Context, rel *entity.Relationship) (*entity.Relationship, error) {
	payload := admin.PutRelationshipJSONRequestBody{TrustDomainAName: rel.TrustDomainAName.String(), TrustDomainBName: rel.TrustDomainBName.String()}
	res, err := g.client.PutRelationship(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf(errorRequestFailed, err)
	}
	defer res.Body.Close()

	body, err := httputil.ReadResponse(res)
	if err != nil {
		return nil, err
	}

	relationship, err := unmarshalJSONToRelationship(body)
	if err != nil {
		return nil, err
	}

	return relationship, nil
}

func (g *galadrielAdminClient) GetRelationships(ctx context.Context, consentStatus api.ConsentStatus, trustDomainName api.TrustDomainName) (*entity.Relationship, error) {
	payload := &admin.GetRelationshipsParams{ConsentStatus: &consentStatus, TrustDomainName: &trustDomainName}

	res, err := g.client.GetRelationships(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf(errorRequestFailed, err)
	}
	defer res.Body.Close()

	body, err := httputil.ReadResponse(res)
	if err != nil {
		return nil, err
	}

	relationship, err := unmarshalJSONToRelationship(body)
	if err != nil {
		return nil, err
	}

	return relationship, nil
}

func (g *galadrielAdminClient) GetRelationshipByID(ctx context.Context, relID uuid.UUID) (*entity.Relationship, error) {
	res, err := g.client.GetRelationshipByID(ctx, relID)
	if err != nil {
		return nil, fmt.Errorf(errorRequestFailed, err)
	}
	defer res.Body.Close()

	body, err := httputil.ReadResponse(res)
	if err != nil {
		return nil, err
	}

	relationship, err := unmarshalJSONToRelationship(body)
	if err != nil {
		return nil, err
	}

	return relationship, nil
}

func (g *galadrielAdminClient) GetJoinToken(ctx context.Context, trustDomainName api.TrustDomainName, ttl int32) (*entity.JoinToken, error) {
	params := &admin.GetJoinTokenParams{Ttl: ttl}
	res, err := g.client.GetJoinToken(ctx, trustDomainName, params)
	if err != nil {
		return nil, fmt.Errorf(errorRequestFailed, err)
	}
	defer res.Body.Close()

	body, err := httputil.ReadResponse(res)
	if err != nil {
		return nil, err
	}

	var joinToken *entity.JoinToken
	if err = json.Unmarshal(body, &joinToken); err != nil {
		return nil, fmt.Errorf(errUnmarshalJoinToken, err)
	}

	return joinToken, nil
}

func unmarshalJSONToTrustDomain(body []byte) (*entity.TrustDomain, error) {
	var trustDomain *entity.TrustDomain
	if err := json.Unmarshal(body, &trustDomain); err != nil {
		return nil, fmt.Errorf(errUnmarshalTrustDomains, err)
	}

	return trustDomain, nil
}

func unmarshalJSONToRelationship(body []byte) (*entity.Relationship, error) {
	var relationship *entity.Relationship
	if err := json.Unmarshal(body, &relationship); err != nil {
		return nil, fmt.Errorf(errUnmarshalRelationships, err)
	}

	return relationship, nil
}
