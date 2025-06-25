package mapper

import (
	"context"
	"time"
	"transport-app/app/domain"
	"transport-app/app/domain/optimization"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MapOptimizedFleetToPlan", func() {
	var ctx context.Context
	var originalPlan domain.Plan
	var optimizedFleet optimization.OptimizedFleet

	BeforeEach(func() {
		ctx = context.Background()
		originalPlan = domain.Plan{
			Headers: domain.Headers{
				Commerce: "test-commerce",
				Consumer: "test-consumer",
				Channel:  "test-channel",
			},
			ReferenceID: "PLAN-001",
		}

		optimizedFleet = optimization.OptimizedFleet{
			PlannedDate: time.Now(),
			Routes: []optimization.OptimizedRoute{
				{
					VehiclePlate: "ABC123",
					Cost:         1500,
					Duration:     3600,
					Steps: []optimization.OptimizedStep{
						{
							Type: "start",
							Location: optimization.Location{
								Latitude:  -33.4567,
								Longitude: -70.6789,
								NodeInfo: optimization.NodeInfo{
									ReferenceID: "start-node",
								},
							},
						},
						{
							Type: "delivery",
							Location: optimization.Location{
								Latitude:  -33.4568,
								Longitude: -70.6790,
								NodeInfo: optimization.NodeInfo{
									ReferenceID: "delivery-node",
								},
							},
							Orders: []optimization.Order{
								{
									ReferenceID: "ORD-001",
									DeliveryUnits: []optimization.DeliveryUnit{
										{
											Lpn:       "LPN-001",
											Weight:    500000,
											Insurance: 25000,
											Items: []optimization.Item{
												{Sku: "SKU-001"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Unassigned: optimization.OptimizedUnassigned{
				Orders: []optimization.Order{
					{
						ReferenceID: "ORD-002",
						DeliveryUnits: []optimization.DeliveryUnit{
							{
								Lpn:       "LPN-002",
								Weight:    300000,
								Insurance: 15000,
								Items: []optimization.Item{
									{Sku: "SKU-002"},
								},
							},
						},
					},
				},
				Vehicles: []optimization.Vehicle{
					{Plate: "XYZ789"},
				},
				Origins: []optimization.NodeInfo{
					{ReferenceID: "origin-node"},
				},
			},
		}
	})

	Describe("Conversión exitosa", func() {
		It("debe convertir correctamente OptimizedFleet a domain.Plan", func() {
			// Act
			result := MapOptimizedFleetToPlan(ctx, optimizedFleet)

			// Assert
			Expect(result.Headers).To(Equal(originalPlan.Headers))
			Expect(result.ReferenceID).To(Equal(originalPlan.ReferenceID))
			Expect(result.PlannedDate).To(Equal(optimizedFleet.PlannedDate))
			Expect(result.Routes).To(HaveLen(1))
			Expect(result.UnassignedOrders).To(HaveLen(1))
			Expect(result.UnassignedVehicles).To(HaveLen(1))
			Expect(result.UnassignedOrigins).To(HaveLen(1))
		})

		It("debe mapear correctamente las rutas optimizadas", func() {
			// Act
			result := MapOptimizedFleetToPlan(ctx, optimizedFleet)

			// Assert
			Expect(result.Routes[0].ReferenceID).To(Equal("ROUTE-ABC123"))
			Expect(result.Routes[0].Vehicle.Plate).To(Equal("ABC123"))
			Expect(result.Routes[0].Orders).To(HaveLen(1))
			Expect(result.Routes[0].Orders[0].ReferenceID).To(Equal(domain.ReferenceID("ORD-001")))
		})

		It("debe mapear correctamente los elementos no asignados", func() {
			// Act
			result := MapOptimizedFleetToPlan(ctx, optimizedFleet)

			// Assert
			Expect(result.UnassignedOrders[0].ReferenceID).To(Equal(domain.ReferenceID("ORD-002")))
			Expect(result.UnassignedVehicles[0].Plate).To(Equal("XYZ789"))
			Expect(result.UnassignedOrigins[0].ReferenceID).To(Equal(domain.ReferenceID("origin-node")))
		})
	})

	Describe("Casos edge", func() {
		It("debe manejar flotas optimizadas vacías", func() {
			// Arrange
			emptyFleet := optimization.OptimizedFleet{
				PlannedDate: time.Now(),
				Routes:      []optimization.OptimizedRoute{},
				Unassigned:  optimization.OptimizedUnassigned{},
			}

			// Act
			result := MapOptimizedFleetToPlan(ctx, emptyFleet)

			// Assert
			Expect(result.Routes).To(HaveLen(0))
			Expect(result.UnassignedOrders).To(HaveLen(0))
			Expect(result.UnassignedVehicles).To(HaveLen(0))
			Expect(result.UnassignedOrigins).To(HaveLen(0))
		})

		It("debe preservar los headers y referenceID del plan original", func() {
			// Act
			result := MapOptimizedFleetToPlan(ctx, optimizedFleet)

			// Assert
			Expect(result.Headers.Commerce).To(Equal("test-commerce"))
			Expect(result.Headers.Consumer).To(Equal("test-consumer"))
			Expect(result.Headers.Channel).To(Equal("test-channel"))
			Expect(result.ReferenceID).To(Equal("PLAN-001"))
		})
	})
})
