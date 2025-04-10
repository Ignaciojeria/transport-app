package firebaseauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/firebaseadminsdk"

	firebase "firebase.google.com/go/v4"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type Login func(context.Context, domain.UserCredentials) (domain.ProviderToken, error)

const firebaseAuthURL = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key="

func init() {
	ioc.Registry(
		NewLogin,
		configuration.NewConf,
		firebaseadminsdk.NewFirebaseAdmin)
}

func NewLogin(
	conf configuration.Conf,
	app *firebase.App) Login {
	apiKey := conf.FIREBASE_API_KEY
	return func(ctx context.Context, uc domain.UserCredentials) (domain.ProviderToken, error) {
		// 🔹 Paso 1: Autenticar usuario con Firebase REST API
		payload := map[string]string{
			"email":             uc.Email,
			"password":          uc.Password,
			"returnSecureToken": "true",
		}

		payloadBytes, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", firebaseAuthURL+apiKey, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return domain.ProviderToken{}, err
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return domain.ProviderToken{}, err
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return domain.ProviderToken{}, err
		}

		// 🔹 Validar credenciales incorrectas
		if _, exists := result["error"]; exists {
			return domain.ProviderToken{}, errors.New("invalid credentials")
		}

		// Obtener UID del usuario autenticado
		uid := result["localId"].(string)
		refreshToken := result["refreshToken"].(string) // 🔥 Obtener Refresh Token

		// 🔹 Paso 2: Asignar Custom Claims en Firebase
		authClient, err := app.Auth(ctx)
		if err != nil {
			return domain.ProviderToken{}, err
		}

		customClaims := map[string]interface{}{} // TODO

		err = authClient.SetCustomUserClaims(ctx, uid, customClaims)
		if err != nil {
			return domain.ProviderToken{}, err
		}

		// 🔹 Paso 3: Forzar regeneración del ID Token con Custom Claims
		idToken, err := refreshIDToken(apiKey, refreshToken)
		if err != nil {
			return domain.ProviderToken{}, err
		}

		// 🔹 Retornar el ID Token con Custom Claims y Refresh Token
		return domain.ProviderToken{
			TokenType:    "idToken",
			TokenValue:   idToken,
			RefreshToken: refreshToken, // 🔥 Incluir el Refresh Token en la respuesta
			ExpiresIn:    3600,
			Provider:     "firebase",
		}, nil
	}
}

// 🔹 Función para regenerar el ID Token con los nuevos Custom Claims
func refreshIDToken(apiKey, refreshToken string) (string, error) {
	refreshPayload := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}

	refreshPayloadBytes, _ := json.Marshal(refreshPayload)
	req, err := http.NewRequest("POST", "https://securetoken.googleapis.com/v1/token?key="+apiKey, bytes.NewBuffer(refreshPayloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var refreshResult map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&refreshResult); err != nil {
		return "", err
	}

	// 🔹 Retornar el nuevo ID Token
	return refreshResult["id_token"].(string), nil
}
