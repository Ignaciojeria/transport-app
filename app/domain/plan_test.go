package domain

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Plan", func() {
	var (
		ctx1, ctx2 context.Context
		now        = time.Now()
		tomorrow   = now.AddDate(0, 0, 1)
	)

	BeforeEach(func() {
		ctx1 = buildCtx("org1", "CL")
		ctx2 = buildCtx("org2", "AR")
	})

	Describe("DocID", func() {
		It("should generate different document IDs for different contexts", func() {
			plan1 := Plan{
				ReferenceID: "PLAN-001",
			}
			plan2 := Plan{
				ReferenceID: "PLAN-001",
			}

			Expect(plan1.DocID(ctx1)).ToNot(Equal(plan2.DocID(ctx2)))
		})

		It("should generate different document IDs for different reference IDs", func() {
			plan1 := Plan{
				ReferenceID: "PLAN-001",
			}
			plan2 := Plan{
				ReferenceID: "PLAN-002",
			}

			Expect(plan1.DocID(ctx1)).ToNot(Equal(plan2.DocID(ctx1)))
		})

		It("should generate the same document ID for same context and reference ID", func() {
			plan1 := Plan{
				ReferenceID: "PLAN-001",
			}
			plan2 := Plan{
				ReferenceID: "PLAN-001",
			}

			Expect(plan1.DocID(ctx1)).To(Equal(plan2.DocID(ctx1)))
		})
	})

	Describe("UpdateIfChanged", func() {
		It("should update planned date if provided and different", func() {
			original := Plan{
				ReferenceID: "PLAN-001",
				PlannedDate: now,
			}
			newPlan := Plan{
				PlannedDate: tomorrow,
			}

			updated, changed := original.UpdateIfChanged(newPlan)
			Expect(changed).To(BeTrue())
			Expect(updated.ReferenceID).To(Equal("PLAN-001")) // No debe cambiar
			Expect(updated.PlannedDate).To(Equal(tomorrow))
		})

		It("should not update planned date when empty value is provided", func() {
			original := Plan{
				ReferenceID: "PLAN-001",
				PlannedDate: now,
			}
			newPlan := Plan{
				PlannedDate: time.Time{}, // zero value
			}

			updated, changed := original.UpdateIfChanged(newPlan)
			Expect(changed).To(BeFalse())
			Expect(updated.ReferenceID).To(Equal("PLAN-001"))
			Expect(updated.PlannedDate).To(Equal(now))
		})

		It("should not mark as changed if same value is provided", func() {
			original := Plan{
				ReferenceID: "PLAN-001",
				PlannedDate: now,
			}
			newPlan := Plan{
				PlannedDate: now, // mismo valor
			}

			updated, changed := original.UpdateIfChanged(newPlan)
			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(original))
		})

		It("should ignore ReferenceID even if provided in newPlan", func() {
			original := Plan{
				ReferenceID: "PLAN-001",
				PlannedDate: now,
			}
			newPlan := Plan{
				ReferenceID: "PLAN-002", // Esto no debería afectar
				PlannedDate: tomorrow,
			}

			updated, changed := original.UpdateIfChanged(newPlan)
			Expect(changed).To(BeTrue())
			Expect(updated.ReferenceID).To(Equal("PLAN-001")) // Debe mantener el original
			Expect(updated.PlannedDate).To(Equal(tomorrow))
		})
	})

	Describe("AssignIndexesToAllOrders", func() {
		It("should assign indexes to unassigned orders without LPN", func() {
			// Crear órdenes sin asignar con unidades de entrega sin LPN
			unassignedOrder1 := Order{
				ReferenceID: "ORDER-001",
				DeliveryUnits: []DeliveryUnit{
					{
						Lpn: "", // Sin LPN
						Items: []Item{
							{Sku: "SKU-001"},
							{Sku: "SKU-002"},
						},
					},
				},
			}

			unassignedOrder2 := Order{
				ReferenceID: "ORDER-002",
				DeliveryUnits: []DeliveryUnit{
					{
						Lpn: "LPN-001", // Con LPN - no debería cambiar
						Items: []Item{
							{Sku: "SKU-003"},
						},
					},
					{
						Lpn: "", // Sin LPN
						Items: []Item{
							{Sku: "SKU-004"},
						},
					},
				},
			}

			plan := &Plan{
				ReferenceID:      "PLAN-001",
				UnassignedOrders: []Order{unassignedOrder1, unassignedOrder2},
			}

			// Ejecutar el método
			plan.AssignIndexesToAllOrders()

			// Verificar que las unidades sin LPN ahora tienen noLPNReference
			Expect(plan.UnassignedOrders[0].DeliveryUnits[0].noLPNReference).To(Equal("ORDER-001"))
			Expect(plan.UnassignedOrders[1].DeliveryUnits[0].noLPNReference).To(Equal("")) // Con LPN, no cambia
			Expect(plan.UnassignedOrders[1].DeliveryUnits[1].noLPNReference).To(Equal("ORDER-002"))
		})

		It("should assign indexes to orders in routes without LPN", func() {
			// Crear rutas con órdenes
			routeOrder1 := Order{
				ReferenceID: "ROUTE-ORDER-001",
				DeliveryUnits: []DeliveryUnit{
					{
						Lpn: "", // Sin LPN
						Items: []Item{
							{Sku: "SKU-005"},
						},
					},
				},
			}

			routeOrder2 := Order{
				ReferenceID: "ROUTE-ORDER-002",
				DeliveryUnits: []DeliveryUnit{
					{
						Lpn: "LPN-002", // Con LPN
						Items: []Item{
							{Sku: "SKU-006"},
						},
					},
				},
			}

			route := Route{
				ReferenceID: "ROUTE-001",
				Orders:      []Order{routeOrder1, routeOrder2},
			}

			plan := &Plan{
				ReferenceID: "PLAN-002",
				Routes:      []Route{route},
			}

			// Ejecutar el método
			plan.AssignIndexesToAllOrders()

			// Verificar que las unidades sin LPN en las rutas ahora tienen noLPNReference
			Expect(plan.Routes[0].Orders[0].DeliveryUnits[0].noLPNReference).To(Equal("ROUTE-ORDER-001"))
			Expect(plan.Routes[0].Orders[1].DeliveryUnits[0].noLPNReference).To(Equal("")) // Con LPN, no cambia
		})

		It("should handle plan with both unassigned orders and routes", func() {
			// Orden sin asignar
			unassignedOrder := Order{
				ReferenceID: "UNASSIGNED-001",
				DeliveryUnits: []DeliveryUnit{
					{
						Lpn: "", // Sin LPN
						Items: []Item{
							{Sku: "SKU-007"},
						},
					},
				},
			}

			// Orden en ruta
			routeOrder := Order{
				ReferenceID: "ROUTE-ORDER-003",
				DeliveryUnits: []DeliveryUnit{
					{
						Lpn: "", // Sin LPN
						Items: []Item{
							{Sku: "SKU-008"},
						},
					},
				},
			}

			route := Route{
				ReferenceID: "ROUTE-002",
				Orders:      []Order{routeOrder},
			}

			plan := &Plan{
				ReferenceID:      "PLAN-003",
				UnassignedOrders: []Order{unassignedOrder},
				Routes:           []Route{route},
			}

			// Ejecutar el método
			plan.AssignIndexesToAllOrders()

			// Verificar que ambas órdenes tienen sus índices asignados
			Expect(plan.UnassignedOrders[0].DeliveryUnits[0].noLPNReference).To(Equal("UNASSIGNED-001"))
			Expect(plan.Routes[0].Orders[0].DeliveryUnits[0].noLPNReference).To(Equal("ROUTE-ORDER-003"))
		})

		It("should handle empty plan", func() {
			plan := &Plan{
				ReferenceID:      "PLAN-004",
				UnassignedOrders: []Order{},
				Routes:           []Route{},
			}

			// No debería causar error
			Expect(func() {
				plan.AssignIndexesToAllOrders()
			}).ToNot(Panic())
		})

		It("should handle orders with empty delivery units", func() {
			order := Order{
				ReferenceID:   "EMPTY-ORDER-001",
				DeliveryUnits: []DeliveryUnit{},
			}

			plan := &Plan{
				ReferenceID:      "PLAN-005",
				UnassignedOrders: []Order{order},
			}

			// No debería causar error
			Expect(func() {
				plan.AssignIndexesToAllOrders()
			}).ToNot(Panic())
		})
	})

	Describe("AssignSequenceNumbersToAllOrders", func() {
		It("should assign sequence numbers to all orders in the plan", func() {
			// Crear órdenes de prueba
			order1 := Order{ReferenceID: "ORD-001"}
			order2 := Order{ReferenceID: "ORD-002"}
			order3 := Order{ReferenceID: "ORD-003"}
			order4 := Order{ReferenceID: "ORD-004"}

			// Crear rutas con órdenes
			route1 := Route{
				ReferenceID: "ROUTE-001",
				Orders:      []Order{order1, order2},
			}
			route2 := Route{
				ReferenceID: "ROUTE-002",
				Orders:      []Order{order3},
			}

			// Crear plan con órdenes no asignadas y rutas
			plan := Plan{
				ReferenceID:      "PLAN-001",
				UnassignedOrders: []Order{order4},
				Routes:           []Route{route1, route2},
			}

			// Asignar números de secuencia
			plan.AssignSequenceNumbersToAllOrders()

			// Verificar que las órdenes no asignadas tienen secuencia
			Expect(plan.UnassignedOrders[0].GetSequenceNumber()).To(Equal(1))

			// Verificar que las órdenes en las rutas tienen secuencia
			Expect(plan.Routes[0].Orders[0].GetSequenceNumber()).To(Equal(2))
			Expect(plan.Routes[0].Orders[1].GetSequenceNumber()).To(Equal(3))
			Expect(plan.Routes[1].Orders[0].GetSequenceNumber()).To(Equal(4))
		})

		It("should handle empty plan", func() {
			plan := Plan{
				ReferenceID: "PLAN-EMPTY",
			}

			// No debería causar errores
			Expect(func() {
				plan.AssignSequenceNumbersToAllOrders()
			}).NotTo(Panic())
		})

		It("should return correct total order count", func() {
			plan := Plan{
				ReferenceID:      "PLAN-COUNT",
				UnassignedOrders: []Order{{ReferenceID: "ORD-001"}, {ReferenceID: "ORD-002"}},
				Routes: []Route{
					{Orders: []Order{{ReferenceID: "ORD-003"}}},
					{Orders: []Order{{ReferenceID: "ORD-004"}, {ReferenceID: "ORD-005"}}},
				},
			}

			Expect(plan.GetTotalOrderCount()).To(Equal(5))
		})
	})

	Describe("AssignSequenceNumbersToRouteOrders", func() {
		It("should assign sequence numbers to orders in a specific route", func() {
			order1 := Order{ReferenceID: "ORD-001"}
			order2 := Order{ReferenceID: "ORD-002"}
			order3 := Order{ReferenceID: "ORD-003"}

			route := Route{
				ReferenceID: "ROUTE-001",
				Orders:      []Order{order1, order2, order3},
			}

			plan := Plan{
				ReferenceID: "PLAN-ROUTE",
				Routes:      []Route{route},
			}

			// Asignar secuencia solo a la primera ruta
			plan.AssignSequenceNumbersToRouteOrders(0)

			// Verificar que las órdenes tienen secuencia
			Expect(plan.Routes[0].Orders[0].GetSequenceNumber()).To(Equal(1))
			Expect(plan.Routes[0].Orders[1].GetSequenceNumber()).To(Equal(2))
			Expect(plan.Routes[0].Orders[2].GetSequenceNumber()).To(Equal(3))
		})

		It("should handle invalid route index", func() {
			plan := Plan{
				ReferenceID: "PLAN-INVALID",
				Routes:      []Route{},
			}

			// No debería causar errores
			Expect(func() {
				plan.AssignSequenceNumbersToRouteOrders(-1)
				plan.AssignSequenceNumbersToRouteOrders(10)
			}).NotTo(Panic())
		})
	})
})
