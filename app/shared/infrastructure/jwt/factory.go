package jwt

import (
	"fmt"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewJWTServiceFromConfig, configuration.NewConf)
}

// NewJWTServiceFromConfig crea un nuevo servicio JWT usando la configuracióngo
func NewJWTServiceFromConfig(conf configuration.Conf) (*JWTService, error) {
	// Verificar que las claves RSA estén configuradas
	if conf.JWT_PRIVATE_KEY == "" || conf.JWT_PUBLIC_KEY == "" {
		fmt.Println("JWT_PRIVATE_KEY and JWT_PUBLIC_KEY must be configured to use RSA")
		return nil, nil
	}

	return NewJWTServiceWithRSA(conf.JWT_PRIVATE_KEY, conf.JWT_PUBLIC_KEY, conf.JWT_ISSUER)
}
