package mapper

import (
	"context"
	"transport-app/app/adapter/out/vroom/model"
	"transport-app/app/domain/optimization"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MapOptimizationResponse", func() {
	var ctx context.Context
	var originalFleet optimization.FleetOptimization

	BeforeEach(func() {
		ctx = context.Background()
		originalFleet = optimization.FleetOptimization{
			Vehicles: []optimization.Vehicle{
				{
					Plate: "ABC123",
					StartLocation: optimization.Location{
						Latitude:  -33.4567,
						Longitude: -70.6789,
						NodeInfo: optimization.NodeInfo{
							ReferenceID: "start-node-1",
						},
					},
					EndLocation: optimization.Location{
						Latitude:  -33.4568,
						Longitude: -70.6790,
						NodeInfo: optimization.NodeInfo{
							ReferenceID: "end-node-1",
						},
					},
				},
			},
			Visits: []optimization.Visit{
				{
					Pickup: optimization.VisitLocation{
						Coordinates: optimization.Coordinates{
							Latitude:  -33.4567,
							Longitude: -70.6789,
						},
						NodeInfo: optimization.NodeInfo{
							ReferenceID: "pickup-node-1",
						},
					},
					Delivery: optimization.VisitLocation{
						Coordinates: optimization.Coordinates{
							Latitude:  -33.4568,
							Longitude: -70.6790,
						},
						NodeInfo: optimization.NodeInfo{
							ReferenceID: "delivery-node-1",
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
		}
	})

	Describe("Mapeo de respuesta exitosa", func() {
		Context("Cuando VROOM devuelve rutas optimizadas", func() {
			It("debe mapear correctamente la respuesta al modelo OptimizedFleet", func() {
				// Arrange
				vroomResponse := model.VroomOptimizationResponse{
					Code: 0,
					Routes: []model.Route{
						{
							Vehicle:  1,
							Cost:     1500,
							Duration: 3600,
							Steps: []model.Step{
								{
									Type:     "start",
									Location: [2]float64{-70.6789, -33.4567},
								},
								{
									Type:     "pickup",
									Shipment: 1,
									Location: [2]float64{-70.6789, -33.4567},
								},
								{
									Type:     "delivery",
									Shipment: 1,
									Location: [2]float64{-70.6790, -33.4568},
								},
								{
									Type:     "end",
									Location: [2]float64{-70.6790, -33.4568},
								},
							},
						},
					},
					Unassigned: []model.UnassignedJob{},
				}

				// Act
				result, err := MapOptimizationResponse(ctx, vroomResponse, originalFleet)

				// Assert
				Expect(err).To(BeNil())
				Expect(result.Routes).To(HaveLen(1))
				Expect(result.Routes[0].VehiclePlate).To(Equal("ABC123"))
				Expect(result.Routes[0].Cost).To(Equal(int64(1500)))
				Expect(result.Routes[0].Duration).To(Equal(int64(3600)))
				Expect(result.Routes[0].Steps).To(HaveLen(4))
				Expect(result.Unassigned.Orders).To(HaveLen(0))
			})
		})

		Context("Cuando VROOM devuelve trabajos no asignados", func() {
			It("debe mapear correctamente los trabajos no asignados", func() {
				// Arrange
				vroomResponse := model.VroomOptimizationResponse{
					Code:   0,
					Routes: []model.Route{},
					Unassigned: []model.UnassignedJob{
						{
							ID:       1,
							Location: [2]float64{-70.6790, -33.4568},
							Reason:   "No vehicle available",
						},
					},
				}

				// Act
				result, err := MapOptimizationResponse(ctx, vroomResponse, originalFleet)

				// Assert
				Expect(err).To(BeNil())
				Expect(result.Routes).To(HaveLen(0))
				Expect(result.Unassigned.Orders).To(HaveLen(1))
				Expect(result.Unassigned.Orders[0].ReferenceID).To(Equal("ORD-001"))
			})
		})
	})

	Describe("Mapeo de tipos de pasos", func() {
		It("debe mapear correctamente los tipos de pasos de VROOM", func() {
			// Arrange
			vroomResponse := model.VroomOptimizationResponse{
				Code: 0,
				Routes: []model.Route{
					{
						Vehicle: 1,
						Steps: []model.Step{
							{Type: "start"},
							{Type: "job", Job: 1},
							{Type: "pickup", Shipment: 1},
							{Type: "delivery", Shipment: 1},
							{Type: "end"},
						},
					},
				},
			}

			// Act
			result, err := MapOptimizationResponse(ctx, vroomResponse, originalFleet)

			// Assert
			Expect(err).To(BeNil())
			Expect(result.Routes[0].Steps).To(HaveLen(5))
			Expect(result.Routes[0].Steps[0].Type).To(Equal("start"))
			Expect(result.Routes[0].Steps[1].Type).To(Equal("delivery"))
			Expect(result.Routes[0].Steps[2].Type).To(Equal("pickup"))
			Expect(result.Routes[0].Steps[3].Type).To(Equal("delivery"))
			Expect(result.Routes[0].Steps[4].Type).To(Equal("end"))
		})
	})

	Describe("Manejo de errores", func() {
		Context("Cuando la respuesta de VROOM tiene un código de error", func() {
			It("debe manejar correctamente el error", func() {
				// Arrange
				vroomResponse := model.VroomOptimizationResponse{
					Code:  1,
					Error: "Invalid input",
				}

				// Act
				result, err := MapOptimizationResponse(ctx, vroomResponse, originalFleet)

				// Assert
				Expect(err).To(BeNil()) // El mapper no valida códigos de error, solo mapea
				Expect(result.Routes).To(HaveLen(0))
			})
		})
	})
})
