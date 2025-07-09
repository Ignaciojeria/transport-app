package encryption

import (
	"fmt"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewClientCredentialsEncryptionFromConfig, configuration.NewConf)
}

// NewClientCredentialsEncryptionFromConfig crea un nuevo servicio de encriptación usando la configuración
func NewClientCredentialsEncryptionFromConfig(conf configuration.Conf) (*ClientCredentialsEncryption, error) {
	// Verificar que la clave de encriptación esté configurada
	if conf.CLIENT_CREDENTIALS_ENCRYPTION_KEY == "" {
		return nil, fmt.Errorf("CLIENT_CREDENTIALS_ENCRYPTION_KEY debe estar configurada")
	}

	return NewClientCredentialsEncryption(conf.CLIENT_CREDENTIALS_ENCRYPTION_KEY)
}
