package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	// Generar una clave aleatoria de 32 bytes para AES-256
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal("Error generando clave aleatoria:", err)
	}

	// Codificar en base64
	encodedKey := base64.StdEncoding.EncodeToString(key)

	fmt.Println("🔐 Clave de encriptación generada para CLIENT_CREDENTIALS_ENCRYPTION_KEY:")
	fmt.Println("")
	fmt.Println(encodedKey)
	fmt.Println("")
	fmt.Println("📝 Instrucciones:")
	fmt.Println("1. Copia esta clave y configúrala como variable de entorno CLIENT_CREDENTIALS_ENCRYPTION_KEY")
	fmt.Println("2. Asegúrate de que esta clave esté disponible en todos los entornos (desarrollo, staging, producción)")
	fmt.Println("3. Mantén esta clave segura y no la compartas en el código fuente")
	fmt.Println("4. Para producción, usa un gestor de secretos como Google Secret Manager o AWS Secrets Manager")
}
