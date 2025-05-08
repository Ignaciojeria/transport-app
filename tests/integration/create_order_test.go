package integration

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

//go:embed TEST00001-CreateOrder.json
var TEST00001CreateOrder []byte

var _ = Describe("TEST00001 - Create Account, Organization and Order", func() {
	It("creates all resources with dynamic organizationKey", func() {
		type EmbeddedRequest struct {
			Headers map[string]string `json:"headers"`
			Body    any               `json:"body"`
		}

		type OrganizationResponse struct {
			OrganizationKey string `json:"organizationKey"`
			Message         string `json:"message"`
		}

		email := "ignaciovl.j@gmail.com"

		// ---------- Crear cuenta ----------
		By("Creating account")
		createAccountRequest := EmbeddedRequest{
			Headers: map[string]string{},
			Body: map[string]any{
				"email": email,
			},
		}

		accountBody, err := json.Marshal(createAccountRequest.Body)
		Expect(err).To(BeNil())

		accountReq, err := http.NewRequest("POST", "http://localhost:8080/register", bytes.NewReader(accountBody))
		Expect(err).To(BeNil())
		accountReq.Header.Set("Content-Type", "application/json")

		accountResp, err := http.DefaultClient.Do(accountReq)
		Expect(err).To(BeNil())
		Expect(accountResp.StatusCode).To(Equal(http.StatusOK))

		// ---------- Crear organización ----------
		By("Creating organization")
		createOrganizationRequest := EmbeddedRequest{
			Headers: map[string]string{},
			Body: map[string]any{
				"name":    "MY EINAR ORGANIZATION",
				"email":   email,
				"country": "CL",
			},
		}

		orgBody, err := json.Marshal(createOrganizationRequest.Body)
		Expect(err).To(BeNil())

		orgReq, err := http.NewRequest("POST", "http://localhost:8080/organizations", bytes.NewReader(orgBody))
		Expect(err).To(BeNil())
		orgReq.Header.Set("Content-Type", "application/json")

		orgResp, err := http.DefaultClient.Do(orgReq)
		Expect(err).To(BeNil())
		Expect(orgResp.StatusCode).To(Equal(http.StatusOK))

		defer orgResp.Body.Close()
		respBody, err := io.ReadAll(orgResp.Body)
		Expect(err).To(BeNil())

		var orgRespData OrganizationResponse
		err = json.Unmarshal(respBody, &orgRespData)
		Expect(err).To(BeNil())
		Expect(orgRespData.OrganizationKey).ToNot(BeEmpty())

		// ---------- Crear orden ----------
		By("Creating order with hardcoded headers")
		orderReq, err := http.NewRequest("POST", "http://localhost:8080/orders", bytes.NewReader(TEST00001CreateOrder))
		Expect(err).To(BeNil())

		// Headers seteados manualmente
		orderReq.Header.Set("organization", orgRespData.OrganizationKey)
		orderReq.Header.Set("commerce", "TEST_COMMERCE")
		orderReq.Header.Set("consumer", "TEST_CONSUMER")
		orderReq.Header.Set("Content-Type", "application/json")

		orderResp, err := http.DefaultClient.Do(orderReq)
		Expect(err).To(BeNil())

		// Si necesitas esperar procesamiento asíncrono, hazlo después
		time.Sleep(1 * time.Second)

		Expect(orderResp.StatusCode).To(Equal(http.StatusOK))

	})
})
