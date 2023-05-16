package endpoints

import (
	"net"

	"github.com/HewlettPackard/galadriel/pkg/common/jwt"
	"github.com/HewlettPackard/galadriel/pkg/server/catalog"
	"github.com/sirupsen/logrus"
)

// Config represents the configuration of the Galadriel Server Endpoints
type Config struct {
	// TPCAddr is the address to bind the TCP listener to.
	TCPAddress *net.TCPAddr

	// LocalAddress is the local address to bind the listener to.
	LocalAddress net.Addr

	JWTIssuer jwt.Issuer

	JWTValidator jwt.Validator

	Catalog catalog.Catalog

	Logger logrus.FieldLogger
}
