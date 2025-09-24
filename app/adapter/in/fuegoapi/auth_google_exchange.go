package fuegoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

type GoogleExchangeRequest struct {
	Code        string `json:"code"`
	State       string `json:"state"`
	RedirectURI string `json:"redirect_uri"`
}

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type GoogleExchangeResponse struct {
	Token string         `json:"token"`
	User  GoogleUserInfo `json:"user"`
	Error string         `json:"error,omitempty"`
}

func init() {
	ioc.Registry(authGoogleExchange, httpserver.New, configuration.NewConf)
}

func authGoogleExchange(s httpserver.Server, conf configuration.Conf) {
	fuego.Post(s.Manager, "/auth/google/exchange",
		func(c fuego.ContextWithBody[GoogleExchangeRequest]) (GoogleExchangeResponse, error) {
			req, err := c.Body()
			if err != nil {
				return GoogleExchangeResponse{
					Error: "Error leyendo request body",
				}, nil
			}

			fmt.Printf("Google Exchange - Code: %s, State: %s\n", req.Code, req.State)

			// 1. Intercambiar código por token con Google
			tokenData, err := exchangeCodeForToken(req.Code, req.RedirectURI, conf)
			if err != nil {
				fmt.Printf("Error intercambiando código: %v\n", err)
				return GoogleExchangeResponse{
					Error: "Error intercambiando código por token",
				}, nil
			}

			// 2. Obtener información del usuario
			userInfo, err := getUserInfo(tokenData.AccessToken)
			if err != nil {
				fmt.Printf("Error obteniendo info de usuario: %v\n", err)
				return GoogleExchangeResponse{
					Error: "Error obteniendo información del usuario",
				}, nil
			}

			// 3. Validar dominio (opcional)
			// if !isValidDomain(userInfo.Email) {
			//     return &map[string]interface{}{
			//         "error": "Dominio no autorizado",
			//     }, nil
			// }

			// 4. Generar tu token personalizado (JWT)
			// TODO: Implementar generación de JWT real
			customToken := fmt.Sprintf("jwt_%s_%s", userInfo.ID, req.State)

			fmt.Printf("Usuario autenticado: %s (%s)\n", userInfo.Name, userInfo.Email)

			return GoogleExchangeResponse{
				Token: customToken,
				User:  *userInfo,
			}, nil
		}, option.Summary("Exchange Google OAuth code for tokens"))
}

func exchangeCodeForToken(code, redirectURI string, conf configuration.Conf) (*GoogleTokenResponse, error) {
	clientID := conf.GOOGLE_OAUTH_CLIENT_ID         // Desde variable de entorno
	clientSecret := conf.GOOGLE_OAUTH_CLIENT_SECRET // Desde variable de entorno

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)

	resp, err := http.Post("https://oauth2.googleapis.com/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Google token endpoint error: %s", string(body))
	}

	var tokenResp GoogleTokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func getUserInfo(accessToken string) (*GoogleUserInfo, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Google userinfo endpoint error: %s", string(body))
	}

	var userInfo GoogleUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// func isValidDomain(email string) bool {
//     allowedDomains := []string{"transportapp.com", "empresa.com"}
//     parts := strings.Split(email, "@")
//     if len(parts) != 2 {
//         return false
//     }
//     domain := parts[1]
//     for _, allowed := range allowedDomains {
//         if domain == allowed {
//             return true
//         }
//     }
//     return false
// }
