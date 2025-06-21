package mapper

import (
	"context"
	"transport-app/app/domain/optimization"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MapOptimizationRequest", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("Mapeo de vehículos", func() {
		Context("Cuando se proporciona un vehículo con todos los campos", func() {
			It("debe mapear correctamente todos los campos del vehículo", func() {
				// Arrange
				fleetOptimization := optimization.FleetOptimization{
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
							Skills: []string{"XL", "refrigerated"},
							TimeWindow: optimization.TimeWindow{
								Start: "08:00",
								End:   "18:00",
							},
							Capacity: optimization.Capacity{
								Weight:                1000000, // 1kg en gramos
								Volume:                5000,    // 5m³
								Insurance:             500000,  // 500k CLP
								DeliveryUnitsQuantity: 50,
							},
						},
					},
					Visits: []optimization.Visit{},
				}

				// Act
				result, err := MapOptimizationRequest(ctx, fleetOptimization)

				// Assert
				Expect(err).To(BeNil())
				Expect(result.Vehicles).To(HaveLen(1))

				vehicle := result.Vehicles[0]
				Expect(vehicle.ID).To(Equal(1))
				Expect(vehicle.Start).To(Equal(&[2]float64{-70.6789, -33.4567}))
				Expect(vehicle.End).To(Equal(&[2]float64{-70.6790, -33.4568}))
				Expect(vehicle.Capacity).To(Equal([]int64{1000000, 50, 500000}))
				Expect(vehicle.Skills).To(Equal([]int64{1, 2}))
				Expect(vehicle.TimeWindow).To(Equal([]int{28800, 64800})) // 08:00 y 18:00 en segundos
			})
		})

		Context("Cuando se proporciona un vehículo con campos opcionales vacíos", func() {
			It("debe omitir los campos opcionales vacíos", func() {
				// Arrange
				fleetOptimization := optimization.FleetOptimization{
					Vehicles: []optimization.Vehicle{
						{
							Plate: "XYZ789",
							StartLocation: optimization.Location{
								Latitude:  0,
								Longitude: 0,
							},
							EndLocation: optimization.Location{
								Latitude:  0,
								Longitude: 0,
							},
							Skills: []string{},
							TimeWindow: optimization.TimeWindow{
								Start: "",
								End:   "",
							},
							Capacity: optimization.Capacity{
								Weight:                0,
								Volume:                0,
								Insurance:             0,
								DeliveryUnitsQuantity: 0,
							},
						},
					},
					Visits: []optimization.Visit{},
				}

				// Act
				result, err := MapOptimizationRequest(ctx, fleetOptimization)

				// Assert
				Expect(err).To(BeNil())
				Expect(result.Vehicles).To(HaveLen(1))

				vehicle := result.Vehicles[0]
				Expect(vehicle.ID).To(Equal(1))
				Expect(vehicle.Start).To(BeNil())
				Expect(vehicle.End).To(BeNil())
				Expect(vehicle.Capacity).To(BeNil())
				Expect(vehicle.Skills).To(BeNil())
				Expect(vehicle.TimeWindow).To(BeNil())
			})
		})
	})

	Describe("Mapeo de jobs (solo delivery)", func() {
		Context("Cuando se proporciona una visita con solo delivery válido", func() {
			It("debe crear un job correctamente", func() {
				// Arrange
				fleetOptimization := optimization.FleetOptimization{
					Vehicles: []optimization.Vehicle{},
					Visits: []optimization.Visit{
						{
							Pickup: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  0,
									Longitude: 0,
								},
							},
							Delivery: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  -33.4567,
									Longitude: -70.6789,
								},
								ServiceTime: 300, // 5 minutos
								Contact: optimization.Contact{
									Email:      "cliente@test.com",
									Phone:      "+56912345678",
									NationalID: "12345678-9",
									FullName:   "Juan Pérez",
								},
								Skills: []string{"fragile"},
								TimeWindow: optimization.TimeWindow{
									Start: "09:00",
									End:   "17:00",
								},
							},
							Orders: []optimization.Order{
								{
									ReferenceID: "ORD-001",
									DeliveryUnits: []optimization.DeliveryUnit{
										{
											Lpn:       "LPN-001",
											Weight:    500000, // 500g
											Volume:    1000,   // 1m³
											Insurance: 25000,  // 25k CLP
											Items: []optimization.Item{
												{Sku: "SKU-001"},
												{Sku: "SKU-002"},
											},
										},
									},
								},
							},
						},
					},
				}

				// Act
				result, err := MapOptimizationRequest(ctx, fleetOptimization)

				// Assert
				Expect(err).To(BeNil())
				Expect(result.Jobs).To(HaveLen(1))
				Expect(result.Shipments).To(HaveLen(0))

				job := result.Jobs[0]
				Expect(job.ID).To(Equal(1)) // ID generado por el registry
				Expect(job.Location).To(Equal([2]float64{-70.6789, -33.4567}))
				Expect(job.Service).To(Equal(int64(300)))
				Expect(job.Amount).To(Equal([]int64{500000, 1, 25000}))
				Expect(job.Skills).To(Equal([]int64{1}))
				Expect(job.TimeWindows).To(Equal([][]int{{32400, 61200}})) // 09:00 y 17:00 en segundos

				// Verificar CustomUserData
				Expect(job.CustomUserData).To(HaveKey("orders"))
				Expect(job.CustomUserData).To(HaveKey("delivery_contact"))
				deliveryContact := job.CustomUserData["delivery_contact"].(optimization.Contact)
				Expect(deliveryContact.Email).To(Equal("cliente@test.com"))
				Expect(deliveryContact.FullName).To(Equal("Juan Pérez"))
			})
		})
	})

	Describe("Mapeo de jobs (solo pickup)", func() {
		Context("Cuando se proporciona una visita con solo pickup válido", func() {
			It("debe crear un job para pickup correctamente", func() {
				// Arrange
				fleetOptimization := optimization.FleetOptimization{
					Vehicles: []optimization.Vehicle{},
					Visits: []optimization.Visit{
						{
							Pickup: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  -33.4566,
									Longitude: -70.6788,
								},
								ServiceTime: 180, // 3 minutos
								Contact: optimization.Contact{
									Email:      "pickup@test.com",
									Phone:      "+56987654321",
									NationalID: "98765432-1",
									FullName:   "María González",
								},
								Skills: []string{"heavy"},
								TimeWindow: optimization.TimeWindow{
									Start: "08:00",
									End:   "12:00",
								},
							},
							Delivery: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  0,
									Longitude: 0,
								},
							},
							Orders: []optimization.Order{
								{
									ReferenceID: "ORD-003",
									DeliveryUnits: []optimization.DeliveryUnit{
										{
											Lpn:       "LPN-003",
											Weight:    300000, // 300g
											Volume:    800,    // 0.8m³
											Insurance: 15000,  // 15k CLP
											Items: []optimization.Item{
												{Sku: "SKU-004"},
											},
										},
									},
								},
							},
						},
					},
				}

				// Act
				result, err := MapOptimizationRequest(ctx, fleetOptimization)

				// Assert
				Expect(err).To(BeNil())
				Expect(result.Jobs).To(HaveLen(1))
				Expect(result.Shipments).To(HaveLen(0))

				job := result.Jobs[0]
				Expect(job.ID).To(Equal(1)) // ID generado por el registry
				Expect(job.Location).To(Equal([2]float64{-70.6788, -33.4566}))
				Expect(job.Service).To(Equal(int64(180)))
				Expect(job.Amount).To(Equal([]int64{300000, 1, 15000}))
				Expect(job.Skills).To(Equal([]int64{1}))
				Expect(job.TimeWindows).To(Equal([][]int{{28800, 43200}})) // 08:00 y 12:00 en segundos

				// Verificar CustomUserData
				Expect(job.CustomUserData).To(HaveKey("orders"))
				Expect(job.CustomUserData).To(HaveKey("pickup_contact"))
				pickupContact := job.CustomUserData["pickup_contact"].(optimization.Contact)
				Expect(pickupContact.Email).To(Equal("pickup@test.com"))
				Expect(pickupContact.FullName).To(Equal("María González"))
			})
		})
	})

	Describe("Mapeo de shipments (pickup + delivery)", func() {
		Context("Cuando se proporciona una visita con pickup y delivery válidos", func() {
			It("debe crear un shipment correctamente", func() {
				// Arrange
				fleetOptimization := optimization.FleetOptimization{
					Vehicles: []optimization.Vehicle{},
					Visits: []optimization.Visit{
						{
							Pickup: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  -33.4566,
									Longitude: -70.6788,
								},
								ServiceTime: 180, // 3 minutos
								Contact: optimization.Contact{
									Email:      "pickup@test.com",
									Phone:      "+56987654321",
									NationalID: "98765432-1",
									FullName:   "María González",
								},
								Skills: []string{"heavy"},
								TimeWindow: optimization.TimeWindow{
									Start: "08:00",
									End:   "12:00",
								},
							},
							Delivery: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  -33.4567,
									Longitude: -70.6789,
								},
								ServiceTime: 300, // 5 minutos
								Contact: optimization.Contact{
									Email:      "delivery@test.com",
									Phone:      "+56912345678",
									NationalID: "12345678-9",
									FullName:   "Juan Pérez",
								},
								Skills: []string{"fragile"},
								TimeWindow: optimization.TimeWindow{
									Start: "14:00",
									End:   "18:00",
								},
							},
							Orders: []optimization.Order{
								{
									ReferenceID: "ORD-002",
									DeliveryUnits: []optimization.DeliveryUnit{
										{
											Lpn:       "LPN-002",
											Weight:    750000, // 750g
											Volume:    1500,   // 1.5m³
											Insurance: 50000,  // 50k CLP
											Items: []optimization.Item{
												{Sku: "SKU-003"},
											},
										},
									},
								},
							},
						},
					},
				}

				// Act
				result, err := MapOptimizationRequest(ctx, fleetOptimization)

				// Assert
				Expect(err).To(BeNil())
				Expect(result.Jobs).To(HaveLen(0))
				Expect(result.Shipments).To(HaveLen(1))

				shipment := result.Shipments[0]
				Expect(shipment.ID).To(Equal(1))

				// Verificar pickup
				Expect(shipment.Pickup.ID).To(Equal(1)) // ID generado por el registry
				Expect(shipment.Pickup.Location).To(Equal(&[2]float64{-70.6788, -33.4566}))
				Expect(shipment.Pickup.TimeWindows).To(Equal([][]int{{28800, 43200}})) // 08:00 y 12:00 en segundos

				// Verificar delivery
				Expect(shipment.Delivery.ID).To(Equal(2)) // ID generado por el registry
				Expect(shipment.Delivery.Location).To(Equal(&[2]float64{-70.6789, -33.4567}))
				Expect(shipment.Delivery.TimeWindows).To(Equal([][]int{{50400, 64800}})) // 14:00 y 18:00 en segundos

				// Verificar otros campos
				Expect(shipment.Service).To(Equal(int64(480))) // 180 + 300
				Expect(shipment.Amount).To(Equal([]int64{750000, 1, 50000}))
				Expect(shipment.Skills).To(Equal([]int64{1, 2})) // heavy, fragile

				// Verificar CustomUserData
				Expect(shipment.CustomUserData).To(HaveKey("orders"))
				Expect(shipment.CustomUserData).To(HaveKey("pickup_contact"))
				Expect(shipment.CustomUserData).To(HaveKey("delivery_contact"))

				pickupContact := shipment.CustomUserData["pickup_contact"].(optimization.Contact)
				Expect(pickupContact.Email).To(Equal("pickup@test.com"))
				Expect(pickupContact.FullName).To(Equal("María González"))

				deliveryContact := shipment.CustomUserData["delivery_contact"].(optimization.Contact)
				Expect(deliveryContact.Email).To(Equal("delivery@test.com"))
				Expect(deliveryContact.FullName).To(Equal("Juan Pérez"))
			})
		})
	})

	Describe("Mapeo de skills", func() {
		Context("Cuando se usan skills duplicados", func() {
			It("debe asignar el mismo ID a skills idénticos", func() {
				// Arrange
				fleetOptimization := optimization.FleetOptimization{
					Vehicles: []optimization.Vehicle{
						{
							Plate:  "VEH-001",
							Skills: []string{"XL", "refrigerated"},
						},
						{
							Plate:  "VEH-002",
							Skills: []string{"XL", "heavy"}, // XL se repite
						},
					},
					Visits: []optimization.Visit{
						{
							Delivery: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  -33.4567,
									Longitude: -70.6789,
								},
								Skills: []string{"refrigerated", "fragile"}, // refrigerated se repite
							},
						},
					},
				}

				// Act
				result, err := MapOptimizationRequest(ctx, fleetOptimization)

				// Assert
				Expect(err).To(BeNil())
				Expect(result.Vehicles).To(HaveLen(2))
				Expect(result.Jobs).To(HaveLen(1))

				// Verificar que XL tiene el mismo ID en ambos vehículos
				Expect(result.Vehicles[0].Skills).To(ContainElement(int64(1))) // XL
				Expect(result.Vehicles[1].Skills).To(ContainElement(int64(1))) // XL

				// Verificar que refrigerated tiene el mismo ID en vehículo y job
				Expect(result.Vehicles[0].Skills).To(ContainElement(int64(2))) // refrigerated
				Expect(result.Jobs[0].Skills).To(ContainElement(int64(2)))     // refrigerated
			})
		})
	})

	Describe("Mapeo de contactos", func() {
		Context("Cuando se usan contactos duplicados en la misma ubicación", func() {
			It("debe asignar el mismo ID a contactos idénticos en la misma ubicación", func() {
				// Arrange
				fleetOptimization := optimization.FleetOptimization{
					Vehicles: []optimization.Vehicle{},
					Visits: []optimization.Visit{
						{
							Delivery: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  -33.4567,
									Longitude: -70.6789,
								},
								Contact: optimization.Contact{
									Email:      "cliente@test.com",
									Phone:      "+56912345678",
									NationalID: "12345678-9",
									FullName:   "Juan Pérez",
								},
							},
						},
						{
							Delivery: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  -33.4567, // Misma ubicación
									Longitude: -70.6789,
								},
								Contact: optimization.Contact{
									Email:      "cliente@test.com", // Mismo contacto
									Phone:      "+56912345678",
									NationalID: "12345678-9",
									FullName:   "Juan Pérez",
								},
							},
						},
					},
				}

				// Act
				result, err := MapOptimizationRequest(ctx, fleetOptimization)

				// Assert
				Expect(err).To(BeNil())
				Expect(result.Jobs).To(HaveLen(2))

				// Verificar que ambos jobs tienen el mismo ID (misma ubicación + mismo contacto)
				Expect(result.Jobs[0].ID).To(Equal(result.Jobs[1].ID))
			})
		})
	})

	Describe("Opciones de VROOM", func() {
		It("debe configurar las opciones correctamente", func() {
			// Arrange
			fleetOptimization := optimization.FleetOptimization{
				Vehicles: []optimization.Vehicle{},
				Visits:   []optimization.Visit{},
			}

			// Act
			result, err := MapOptimizationRequest(ctx, fleetOptimization)

			// Assert
			Expect(err).To(BeNil())
			Expect(result.Options).ToNot(BeNil())
			Expect(result.Options.G).To(BeTrue())
			Expect(result.Options.Steps).To(BeTrue())
			Expect(result.Options.Overview).To(BeTrue())
			Expect(result.Options.MinimizeVehicles).To(BeTrue())
		})
	})

	Describe("Casos edge", func() {
		Context("Cuando no hay vehículos ni visitas", func() {
			It("debe retornar un request vacío pero válido", func() {
				// Arrange
				fleetOptimization := optimization.FleetOptimization{
					Vehicles: []optimization.Vehicle{},
					Visits:   []optimization.Visit{},
				}

				// Act
				result, err := MapOptimizationRequest(ctx, fleetOptimization)

				// Assert
				Expect(err).To(BeNil())
				Expect(result.Vehicles).To(HaveLen(0))
				Expect(result.Jobs).To(HaveLen(0))
				Expect(result.Shipments).To(HaveLen(0))
				Expect(result.Options).ToNot(BeNil())
			})
		})

		Context("Cuando una visita no tiene coordenadas válidas", func() {
			It("debe omitir la visita", func() {
				// Arrange
				fleetOptimization := optimization.FleetOptimization{
					Vehicles: []optimization.Vehicle{},
					Visits: []optimization.Visit{
						{
							Pickup: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  0,
									Longitude: 0,
								},
							},
							Delivery: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  0,
									Longitude: 0,
								},
							},
						},
					},
				}

				// Act
				result, err := MapOptimizationRequest(ctx, fleetOptimization)

				// Assert
				Expect(err).To(BeNil())
				Expect(result.Jobs).To(HaveLen(0))
				Expect(result.Shipments).To(HaveLen(0))
			})
		})

		Context("Cuando hay múltiples tipos de visitas en el mismo request", func() {
			It("debe mapear correctamente cada tipo de visita", func() {
				// Arrange
				fleetOptimization := optimization.FleetOptimization{
					Vehicles: []optimization.Vehicle{},
					Visits: []optimization.Visit{
						{
							// Visita 1: Solo pickup válido -> Job
							Pickup: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  -33.4566,
									Longitude: -70.6788,
								},
								ServiceTime: 180,
								Contact: optimization.Contact{
									Email:    "pickup@test.com",
									FullName: "María González",
								},
								Skills: []string{"heavy"},
							},
							Delivery: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  0,
									Longitude: 0,
								},
							},
							Orders: []optimization.Order{
								{
									ReferenceID: "ORD-001",
									DeliveryUnits: []optimization.DeliveryUnit{
										{
											Lpn:    "LPN-001",
											Weight: 300000,
										},
									},
								},
							},
						},
						{
							// Visita 2: Solo delivery válido -> Job
							Pickup: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  0,
									Longitude: 0,
								},
							},
							Delivery: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  -33.4567,
									Longitude: -70.6789,
								},
								ServiceTime: 300,
								Contact: optimization.Contact{
									Email:    "delivery@test.com",
									FullName: "Juan Pérez",
								},
								Skills: []string{"fragile"},
							},
							Orders: []optimization.Order{
								{
									ReferenceID: "ORD-002",
									DeliveryUnits: []optimization.DeliveryUnit{
										{
											Lpn:    "LPN-002",
											Weight: 500000,
										},
									},
								},
							},
						},
						{
							// Visita 3: Pickup y delivery válidos -> Shipment
							Pickup: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  -33.4568,
									Longitude: -70.6790,
								},
								ServiceTime: 120,
								Contact: optimization.Contact{
									Email:    "pickup2@test.com",
									FullName: "Ana López",
								},
								Skills: []string{"refrigerated"},
							},
							Delivery: optimization.VisitLocation{
								Coordinates: optimization.Coordinates{
									Latitude:  -33.4569,
									Longitude: -70.6791,
								},
								ServiceTime: 240,
								Contact: optimization.Contact{
									Email:    "delivery2@test.com",
									FullName: "Carlos Ruiz",
								},
								Skills: []string{"fragile"},
							},
							Orders: []optimization.Order{
								{
									ReferenceID: "ORD-003",
									DeliveryUnits: []optimization.DeliveryUnit{
										{
											Lpn:    "LPN-003",
											Weight: 750000,
										},
									},
								},
							},
						},
					},
				}

				// Act
				result, err := MapOptimizationRequest(ctx, fleetOptimization)

				// Assert
				Expect(err).To(BeNil())
				Expect(result.Jobs).To(HaveLen(2))      // 2 jobs (pickup + delivery)
				Expect(result.Shipments).To(HaveLen(1)) // 1 shipment

				// Verificar que los jobs tienen IDs diferentes
				Expect(result.Jobs[0].ID).ToNot(Equal(result.Jobs[1].ID))

				// Verificar que el shipment tiene pickup y delivery con IDs diferentes
				shipment := result.Shipments[0]
				Expect(shipment.Pickup.ID).ToNot(Equal(shipment.Delivery.ID))

				// Verificar que todos los skills se mapean correctamente
				// heavy, fragile, refrigerated, fragile (duplicado)
				allSkills := []int64{}
				for _, job := range result.Jobs {
					allSkills = append(allSkills, job.Skills...)
				}
				allSkills = append(allSkills, shipment.Skills...)
				Expect(allSkills).To(ContainElements(int64(1), int64(2), int64(3))) // 3 skills únicos
			})
		})
	})
})
