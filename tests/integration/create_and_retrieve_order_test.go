package integration

import (
	"bytes"
	_ "embed"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

//go:embed TEST00001-CreateAndRetrievePackageWithLpns.json
var rawTest00001 []byte
var _ = Describe("create order using post and retrieve order data using graphql", func() {
	It("TEST00001-CreateAndRetrievePackageWithLpns.json", func() {
		payload, err := loadRequest(rawTest00001)
		Expect(err).To(BeNil())

		// Simulación: puedes reemplazar esto con una llamada real a POST
		req, err := http.NewRequest("POST", "http://localhost:8080/orders", bytes.NewReader(payload.Body))
		Expect(err).To(BeNil())

		for k, v := range payload.Headers {
			req.Header.Set(k, v)
		}

		resp, err := http.DefaultClient.Do(req)
		Expect(err).To(BeNil())
		Expect(resp.StatusCode).To(Equal(http.StatusCreated))

		// Luego haces la query GraphQL aquí...
	})
})
