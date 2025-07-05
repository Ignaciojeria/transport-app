package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService maneja la firma y validación de tokens JWT
type JWTService struct {
	secretKey     []byte
	privateKey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
	issuer        string
	signingMethod jwt.SigningMethod
}

// Claims personalizadas para el JWT
type Claims struct {
	Sub     string            `json:"sub"`
	Scopes  []string          `json:"scopes"`
	Context map[string]string `json:"context"`
	Tenant  string            `json:"tenant"`
	jwt.RegisteredClaims
}

// JWK representa una clave JSON Web Key
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKS representa un JSON Web Key Set
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// NewJWTService crea una nueva instancia del servicio JWT con HMAC
func NewJWTService(secretKey string, issuer string) *JWTService {
	return &JWTService{
		secretKey:     []byte(secretKey),
		issuer:        issuer,
		signingMethod: jwt.SigningMethodHS256,
	}
}

// NewJWTServiceWithRSA crea una nueva instancia del servicio JWT con RSA
func NewJWTServiceWithRSA(privateKeyPEM, publicKeyPEM, issuer string) (*JWTService, error) {
	// Parsear clave privada
	privateKeyBlock, _ := pem.Decode([]byte(privateKeyPEM))
	if privateKeyBlock == nil {
		return nil, fmt.Errorf("error decodificando clave privada PEM")
	}

	// Intentar parsear como PKCS#1 primero, luego PKCS#8
	var privateKey *rsa.PrivateKey
	var err error

	// Intentar PKCS#1 (formato tradicional RSA)
	privateKey, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		// Si falla, intentar PKCS#8 (formato moderno)
		key, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("error parseando clave privada (PKCS#1 y PKCS#8): %w", err)
		}

		// Verificar que es una clave RSA
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("clave privada no es RSA")
		}
	}

	// Parsear clave pública
	publicKeyBlock, _ := pem.Decode([]byte(publicKeyPEM))
	if publicKeyBlock == nil {
		return nil, fmt.Errorf("error decodificando clave pública PEM")
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parseando clave pública: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("clave pública no es RSA")
	}

	return &JWTService{
		privateKey:    privateKey,
		publicKey:     rsaPublicKey,
		issuer:        issuer,
		signingMethod: jwt.SigningMethodRS256,
	}, nil
}

// GenerateToken genera un nuevo token JWT
func (j *JWTService) GenerateToken(sub string, scopes []string, context map[string]string, tenant string, audience string, expirationMinutes int) (string, int64, error) {
	// Calcular tiempo de expiración
	expirationTime := time.Now().Add(time.Duration(expirationMinutes) * time.Minute)
	expiresAt := expirationTime.Unix()

	// Crear claims
	claims := &Claims{
		Sub:     sub,
		Scopes:  scopes,
		Context: context,
		Tenant:  tenant,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
			Subject:   sub,
			Audience:  []string{audience},
		},
	}

	// Crear token
	token := jwt.NewWithClaims(j.signingMethod, claims)

	// Agregar kid al header si es RSA
	if j.signingMethod == jwt.SigningMethodRS256 && j.publicKey != nil {
		// Generar el mismo kid que se usa en GetJWKS
		kid := fmt.Sprintf("key-%x", j.publicKey.N.Bytes()[:8])
		token.Header["kid"] = kid
	}

	// Firmar token
	var tokenString string
	var err error

	if j.signingMethod == jwt.SigningMethodRS256 {
		tokenString, err = token.SignedString(j.privateKey)
	} else {
		tokenString, err = token.SignedString(j.secretKey)
	}

	if err != nil {
		return "", 0, fmt.Errorf("error firmando token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateToken valida y parsea un token JWT
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	// Parsear token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar método de firma
		if j.signingMethod == jwt.SigningMethodRS256 {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return j.publicKey, nil
		} else {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return j.secretKey, nil
		}
	})

	if err != nil {
		return nil, fmt.Errorf("error parseando token: %w", err)
	}

	// Verificar si el token es válido
	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	// Extraer claims
	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("error extrayendo claims del token")
}

// GetPublicKeyPEM retorna la clave pública en formato PEM para Zuplo
func (j *JWTService) GetPublicKeyPEM() (string, error) {
	if j.publicKey == nil {
		return "", fmt.Errorf("no hay clave pública configurada")
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(j.publicKey)
	if err != nil {
		return "", fmt.Errorf("error serializando clave pública: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(publicKeyPEM), nil
}

// GetJWKS retorna el JSON Web Key Set para Zuplo
func (j *JWTService) GetJWKS() (*JWKS, error) {
	if j.publicKey == nil {
		return nil, fmt.Errorf("no public key configured")
	}

	// Generar un Key ID único basado en el módulo de la clave
	kid := fmt.Sprintf("key-%x", j.publicKey.N.Bytes()[:8])

	// Convertir la clave pública RSA a formato JWK
	jwk := JWK{
		Kty: "RSA",
		Use: "sig",
		Kid: kid,
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(j.publicKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(j.publicKey.E)).Bytes()),
	}

	return &JWKS{
		Keys: []JWK{jwk},
	}, nil
}

// GenerateRSAKeyPair genera un par de claves RSA
func GenerateRSAKeyPair(bits int) (privateKeyPEM, publicKeyPEM string, err error) {
	// Generar clave privada
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", fmt.Errorf("error generando clave privada: %w", err)
	}

	// Serializar clave privada
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}))

	// Serializar clave pública
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("error serializando clave pública: %w", err)
	}

	publicKeyPEM = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}))

	return privateKeyPEM, publicKeyPEM, nil
}

// GenerateSecretKey genera una clave secreta aleatoria
func GenerateSecretKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("error generando clave secreta: %w", err)
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// RefreshToken genera un nuevo token basado en uno existente
func (j *JWTService) RefreshToken(tokenString string, expirationMinutes int) (string, int64, error) {
	// Validar token existente
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", 0, fmt.Errorf("error validando token para refresh: %w", err)
	}

	// Obtener el audience del token original
	var audience string
	if len(claims.Audience) > 0 {
		audience = claims.Audience[0]
	}

	// Generar nuevo token con los mismos claims pero nueva expiración
	return j.GenerateToken(claims.Sub, claims.Scopes, claims.Context, claims.Tenant, audience, expirationMinutes)
}
