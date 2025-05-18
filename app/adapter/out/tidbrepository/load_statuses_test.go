package tidbrepository

/*
import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoadStatuses", func() {
	var (
		loadStatuses LoadStatuses
		statuses     Statuses
	)

	BeforeEach(func() {
		loadStatuses = NewLoadStatuses(connection)
		// Primero llamamos a loadStatuses() para asegurarnos que los registros se guarden
		statuses = loadStatuses()
	})

	It("should load all statuses and verify they can be found by document ID", func() {
		ctx := context.Background()

		// Test each status
		testCases := []struct {
			name     string
			status   domain.Status
			expected string
		}{
			{"Available", statuses.Available(), domain.StatusAvailable},
			{"Scanned", statuses.Scanned(), domain.StatusScanned},
			{"Picked", statuses.Picked(), domain.StatusPicked},
			{"Planned", statuses.Planned(), domain.StatusPlanned},
			{"InTransit", statuses.InTransit(), domain.StatusInTransit},
			{"Cancelled", statuses.Cancelled(), domain.StatusCancelled},
			{"Finished", statuses.Finished(), domain.StatusFinished},
		}

		for _, tc := range testCases {
			By("Testing " + tc.name + " status")

			// Verify the status exists in the map
			Expect(tc.status.Status).To(Equal(tc.expected))

			// Verify the status exists in the database
			var dbStatus table.Status
			err := connection.DB.WithContext(ctx).
				Table("statuses").
				Where("document_id = ?", tc.status.DocID().String()).
				First(&dbStatus).Error
			Expect(err).ToNot(HaveOccurred(), "Failed to find status in database: %v", err)
			Expect(dbStatus.Status).To(Equal(tc.expected))
			Expect(dbStatus.ID).To(Equal(tc.status.ID))
		}
	})

	It("should fail if database has no statuses table", func() {
		loadStatuses = NewLoadStatuses(noTablesContainerConnection)
		Expect(func() {
			loadStatuses()
		}).To(Panic())
	})
})*/
