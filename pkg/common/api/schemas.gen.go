// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
)

// Defines values for ConsentStatus.
const (
	Approved ConsentStatus = "approved"
	Denied   ConsentStatus = "denied"
	Pending  ConsentStatus = "pending"
)

// ApiError defines model for ApiError.
type ApiError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// BundleDigest base64 encoded SHA-256 digest of the bundle
type BundleDigest = string

// Certificate X.509 certificate in PEM format
type Certificate = string

// ConsentStatus defines model for ConsentStatus.
type ConsentStatus string

// Date defines model for Date.
type Date = openapi_types.Date

// JWT defines model for JWT.
type JWT = string

// JoinToken defines model for JoinToken.
type JoinToken = UUID

// PageNumber The number of items to skip before starting to collect the result set.
type PageNumber = int

// PageSize The numbers of items to return.
type PageSize = int

// Relationship defines model for Relationship.
type Relationship struct {
	CreatedAt           time.Time        `json:"created_at"`
	Id                  UUID             `json:"id"`
	TrustDomainAConsent ConsentStatus    `json:"trust_domain_a_consent"`
	TrustDomainAId      UUID             `json:"trust_domain_a_id"`
	TrustDomainAName    *TrustDomainName `json:"trust_domain_a_name,omitempty"`
	TrustDomainBConsent ConsentStatus    `json:"trust_domain_b_consent"`
	TrustDomainBId      UUID             `json:"trust_domain_b_id"`
	TrustDomainBName    *TrustDomainName `json:"trust_domain_b_name,omitempty"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

// SPIFFEID defines model for SPIFFEID.
type SPIFFEID = string

// Signature base64 encoded signature of the bundle
type Signature = string

// TrustBundle SPIFFE Trust bundle in JSON format
type TrustBundle = string

// TrustDomain defines model for TrustDomain.
type TrustDomain struct {
	CreatedAt         time.Time       `json:"created_at"`
	Description       *string         `json:"description,omitempty"`
	HarvesterSpiffeId *SPIFFEID       `json:"harvester_spiffe_id,omitempty"`
	Id                UUID            `json:"id"`
	Name              TrustDomainName `json:"name"`

	// OnboardingBundle SPIFFE Trust bundle in JSON format
	OnboardingBundle *TrustBundle `json:"onboarding_bundle,omitempty"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

// TrustDomainName defines model for TrustDomainName.
type TrustDomainName = string

// UUID defines model for UUID.
type UUID = openapi_types.UUID

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
