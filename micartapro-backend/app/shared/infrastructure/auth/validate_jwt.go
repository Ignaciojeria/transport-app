package auth

import (
	"errors"
	"fmt"
	"micartapro/app/shared/configuration"
	"time"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	keyfunc "github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

type SupabaseTokenValidator struct {
	jwks *keyfunc.JWKS
}

func init() {
	ioc.Registry(
		NewSupabaseTokenValidator,
		configuration.NewConf,
	)
}

func NewSupabaseTokenValidator(conf configuration.Conf) (SupabaseTokenValidator, error) {

	jwks, err := keyfunc.Get(conf.SUPABASE_JWKS_URL, keyfunc.Options{
		RefreshInterval: time.Hour,
		RefreshErrorHandler: func(err error) {
			fmt.Printf("Error refrescando JWKS: %v\n", err)
		},
	})
	if err != nil {
		return SupabaseTokenValidator{}, fmt.Errorf("error cargando JWKS: %w", err)
	}

	return SupabaseTokenValidator{jwks: jwks}, nil
}

// Validate valida el JWT y retorna las claims decodificadas
func (v SupabaseTokenValidator) ValidateJWT(tokenString string) (jwt.MapClaims, error) {

	if tokenString == "" {
		return nil, errors.New("token vacío")
	}

	// Parse completo con validación automática de firma usando JWKS
	token, err := jwt.Parse(tokenString, v.jwks.Keyfunc)
	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	// Verificar estructura + validez
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido o claims corruptas")
	}

	// Validación estándar del exp (jwt/v5 lo hace si se usa RegisteredClaims, pero aquí lo controlamos manual)
	if expRaw, ok := claims["exp"]; ok {
		if expFloat, ok := expRaw.(float64); ok {
			expTime := time.Unix(int64(expFloat), 0)
			if time.Now().After(expTime) {
				return nil, errors.New("token expirado")
			}
		}
	}

	return claims, nil
}
