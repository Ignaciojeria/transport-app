package fuegoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"transport-app/app/adapter/out/natspublisher"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/shared/infrastructure/jwt"
	"transport-app/app/shared/sharedcontext"

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
	AccessToken  string         `json:"access_token"`
	TokenType    string         `json:"token_type"`
	ExpiresIn    int            `json:"expires_in"`
	RefreshToken string         `json:"refresh_token"`
	User         GoogleUserInfo `json:"user"`
	Error        string         `json:"error,omitempty"`
}

type AuthenticationSubmittedEvent struct {
	UserID       string         `json:"user_id"`
	Email        string         `json:"email"`
	Name         string         `json:"name"`
	Provider     string         `json:"provider"`
	UserInfo     GoogleUserInfo `json:"user_info"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	Timestamp    time.Time      `json:"timestamp"`
}

func init() {
	ioc.Registry(authGoogleExchange, httpserver.New, configuration.NewConf, jwt.NewJWTServiceFromConfig, natspublisher.NewApplicationEvents)
}

func authGoogleExchange(s httpserver.Server, conf configuration.Conf, jwtService *jwt.JWTService, publish natspublisher.ApplicationEvents) {
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

			// 4. Generar access token JWT sin tenant (usuario autenticado pero sin tenant asignado)
			accessToken, err := jwtService.GenerateToken(
				fmt.Sprintf("google_%s", userInfo.ID), // sub = google_userID
				[]string{"profile", "select_tenant"},  // scopes para seleccionar tenant
				map[string]string{
					"provider": "google",
					"email":    userInfo.Email,
					"name":     userInfo.Name,
				}, // context con información del usuario
				"",              // sin tenant específico
				"zuplo-gateway", // audiencia
				60,              // 60 minutos de expiración
			)
			if err != nil {
				fmt.Printf("Error generando access token: %v\n", err)
				return GoogleExchangeResponse{
					Error: "Error generando token de acceso",
				}, nil
			}

			// 5. Generar refresh token para mantener la sesión
			refreshToken, err := jwtService.GenerateToken(
				fmt.Sprintf("google_%s", userInfo.ID), // sub = google_userID
				[]string{"refresh"},                   // scope para refresh
				map[string]string{
					"provider": "google",
					"email":    userInfo.Email,
					"name":     userInfo.Name,
				}, // context con información del usuario
				"",              // sin tenant específico
				"zuplo-gateway", // audiencia
				10080,           // 7 días de expiración
			)
			if err != nil {
				fmt.Printf("Error generando refresh token: %v\n", err)
				return GoogleExchangeResponse{
					Error: "Error generando refresh token",
				}, nil
			}

			// 6. Preparar evento de autenticación
			authEvent := AuthenticationSubmittedEvent{
				UserID:       fmt.Sprintf("google_%s", userInfo.ID),
				Email:        userInfo.Email,
				Name:         userInfo.Name,
				Provider:     "google",
				UserInfo:     *userInfo,
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				Timestamp:    time.Now(),
			}

			// 7. Publicar evento de autenticación
			eventPayload, _ := json.Marshal(authEvent)
			eventCtx := sharedcontext.AddEventContextToBaggage(c.Context(),
				sharedcontext.EventContext{
					EventType:  "authenticationSubmitted",
					EntityType: "authentication",
				},
			)

			if err := publish(eventCtx, domain.Outbox{
				Payload: eventPayload,
			}); err != nil {
				fmt.Printf("Error publicando evento de autenticación: %v\n", err)
				// No fallar el login por error en el evento, solo log
			}

			// 8. Extraer expiración para ExpiresIn
			claims, err := jwtService.ValidateToken(accessToken)
			if err != nil {
				fmt.Printf("Error validando token generado: %v\n", err)
				return GoogleExchangeResponse{
					Error: "Error validando token generado",
				}, nil
			}

			fmt.Printf("Usuario autenticado: %s (%s)\n", userInfo.Name, userInfo.Email)

			return GoogleExchangeResponse{
				AccessToken:  accessToken,
				TokenType:    "Bearer",
				ExpiresIn:    int(claims.ExpiresAt.Unix() - time.Now().Unix()),
				RefreshToken: refreshToken,
				User:         *userInfo,
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
