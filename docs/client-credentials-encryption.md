# 🔐 Encriptación de Client Credentials

## Descripción

El sistema implementa encriptación AES-256-GCM para proteger los `ClientSecret` de las credenciales de cliente antes de almacenarlos en la base de datos.

## 🛡️ Características de Seguridad

- **Algoritmo**: AES-256-GCM
- **Tamaño de clave**: 32 bytes (256 bits)
- **Codificación**: Base64
- **Nonce aleatorio**: Generado automáticamente para cada encriptación
- **Autenticación**: GCM proporciona autenticación además de encriptación

## 📋 Configuración

### 1. Generar Clave de Encriptación

Ejecuta el script para generar una clave segura:

```bash
go run scripts/generate_encryption_key.go
```

### 2. Configurar Variable de Entorno

Agrega la clave generada como variable de entorno:

```bash
export CLIENT_CREDENTIALS_ENCRYPTION_KEY="tu-clave-generada-en-base64"
```

### 3. Para Producción

En entornos de producción, usa un gestor de secretos:

#### Google Cloud Secret Manager
```bash
# Crear el secreto
gcloud secrets create client-credentials-encryption-key --data-file=encryption-key.txt

# Configurar en Cloud Run
gcloud run services update transport-app \
  --set-secrets=CLIENT_CREDENTIALS_ENCRYPTION_KEY=client-credentials-encryption-key:latest
```

#### AWS Secrets Manager
```bash
# Crear el secreto
aws secretsmanager create-secret \
  --name client-credentials-encryption-key \
  --secret-string "tu-clave-generada-en-base64"
```

## 🔧 Uso en el Código

### Servicio de Encriptación

```go
import "transport-app/app/shared/infrastructure/encryption"

// El servicio se inyecta automáticamente en los repositorios
encryptionService := encryption.NewClientCredentialsEncryptionFromConfig(conf)

// Encriptar
encrypted, err := encryptionService.Encrypt("mi-secreto-cliente")

// Desencriptar
decrypted, err := encryptionService.Decrypt(encrypted)
```

### Repositorios

Los repositorios manejan la encriptación automáticamente:

```go
// Al guardar - el ClientSecret se encripta automáticamente
credentials := domain.ClientCredentials{
    ClientID: "mi-client-id",
    ClientSecret: "mi-secreto-sin-encriptar", // Se encripta automáticamente
}
savedCredentials, err := upsertClientCredentials(ctx, credentials)

// Al leer - el ClientSecret se desencripta automáticamente
foundCredentials, err := findClientCredentialsByID(ctx, id)
// foundCredentials.ClientSecret ya está desencriptado
```

## 🧪 Testing

Ejecuta los tests de encriptación:

```bash
go test ./app/shared/infrastructure/encryption/...
```

## ⚠️ Consideraciones Importantes

### 1. Rotación de Claves
- **NO** cambies la clave de encriptación sin migrar los datos existentes
- Implementa un proceso de migración si necesitas cambiar la clave

### 2. Backup de Claves
- Mantén un backup seguro de la clave de encriptación
- Sin la clave, los datos encriptados no se pueden recuperar

### 3. Logs y Debugging
- Los secretos encriptados se muestran como `[ENCRYPTED]` en los logs
- Nunca logees secretos desencriptados

### 4. Migración de Datos Existentes

Si tienes datos sin encriptar, crea un script de migración:

```go
// Ejemplo de migración
func migrateExistingCredentials(ctx context.Context) error {
    // 1. Leer todas las credenciales sin encriptar
    // 2. Encriptar cada ClientSecret
    // 3. Actualizar en la base de datos
    return nil
}
```

## 🔍 Monitoreo

### Métricas Recomendadas
- Número de encriptaciones/desencriptaciones por minuto
- Errores de encriptación/desencriptación
- Tiempo de respuesta del servicio de encriptación

### Alertas
- Errores de desencriptación (podrían indicar clave incorrecta)
- Fallos en el servicio de encriptación

## 📚 Referencias

- [AES-GCM](https://en.wikipedia.org/wiki/Galois/Counter_Mode)
- [Go crypto/aes](https://golang.org/pkg/crypto/aes/)
- [Go crypto/cipher](https://golang.org/pkg/crypto/cipher/) 