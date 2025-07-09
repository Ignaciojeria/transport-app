package usecase

import (
	"context"
	"testing"
	"transport-app/app/domain"

	"github.com/biter777/countries"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock del repositorio
type MockUpsertClientCredentials struct {
	mock.Mock
}

func (m *MockUpsertClientCredentials) UpsertClientCredentials(ctx context.Context, credentials domain.ClientCredentials) (domain.ClientCredentials, error) {
	args := m.Called(ctx, credentials)
	return args.Get(0).(domain.ClientCredentials), args.Error(1)
}

func TestCreateClientCredentials(t *testing.T) {
	// Arrange
	mockRepo := &MockUpsertClientCredentials{}
	createClientCredentials := NewCreateClientCredentials(mockRepo.UpsertClientCredentials)

	tenantID := uuid.New()
	tenantCountry := countries.CL
	scopes := []string{"orders:read", "orders:write"}

	// Configurar el mock para retornar las credenciales creadas
	mockRepo.On("UpsertClientCredentials", mock.Anything, mock.AnythingOfType("domain.ClientCredentials")).
		Return(func(ctx context.Context, credentials domain.ClientCredentials) (domain.ClientCredentials, error) {
			// Simular que las credenciales se guardaron correctamente
			return credentials, nil
		})

	// Act
	credentials, err := createClientCredentials(context.Background(), tenantID, tenantCountry, scopes)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, credentials.ID)
	assert.Equal(t, tenantID, credentials.TenantID)
	assert.Equal(t, tenantCountry, credentials.TenantCountry)
	assert.NotEmpty(t, credentials.ClientID)
	assert.NotEmpty(t, credentials.ClientSecret)
	assert.Equal(t, scopes, credentials.AllowedScopes)
	assert.Equal(t, "active", credentials.Status)

	// Verificar que el ClientID tiene el formato esperado (24 caracteres)
	assert.Len(t, credentials.ClientID, 24)

	// Verificar que el ClientSecret tiene el formato esperado (43 caracteres)
	assert.Len(t, credentials.ClientSecret, 43)

	mockRepo.AssertExpectations(t)
}

func TestGenerateClientID(t *testing.T) {
	// Act
	clientID1, err1 := generateClientID()
	clientID2, err2 := generateClientID()

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, clientID1)
	assert.NotEmpty(t, clientID2)
	assert.Len(t, clientID1, 24)
	assert.Len(t, clientID2, 24)

	// Verificar que son diferentes (aleatorios)
	assert.NotEqual(t, clientID1, clientID2)
}

func TestGenerateClientSecret(t *testing.T) {
	// Act
	secret1, err1 := generateClientSecret()
	secret2, err2 := generateClientSecret()

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, secret1)
	assert.NotEmpty(t, secret2)
	assert.Len(t, secret1, 43)
	assert.Len(t, secret2, 43)

	// Verificar que son diferentes (aleatorios)
	assert.NotEqual(t, secret1, secret2)
}

func TestCreateClientCredentialsWithError(t *testing.T) {
	// Arrange
	mockRepo := &MockUpsertClientCredentials{}
	createClientCredentials := NewCreateClientCredentials(mockRepo.UpsertClientCredentials)

	tenantID := uuid.New()
	tenantCountry := countries.CL
	scopes := []string{"orders:read"}

	// Configurar el mock para retornar un error
	mockRepo.On("UpsertClientCredentials", mock.Anything, mock.AnythingOfType("domain.ClientCredentials")).
		Return(domain.ClientCredentials{}, assert.AnError)

	// Act
	credentials, err := createClientCredentials(context.Background(), tenantID, tenantCountry, scopes)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, credentials.ID)
	assert.Contains(t, err.Error(), "error guardando client credentials")

	mockRepo.AssertExpectations(t)
}
