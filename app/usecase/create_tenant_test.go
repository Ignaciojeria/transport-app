package usecase

import (
	"context"
	"fmt"
	"strings"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	"github.com/biter777/countries"
	goqu "github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tenant", func() {
	Describe("Default inner join readiness", func() {
		It("should allow inner join across all tenant-related tables", func() {
			// Create first tenant
			tenantID1 := uuid.New()
			ctx1 := buildCtx(tenantID1.String(), "CL")

			err := testCreateTenant(ctx1, tenantID1)
			Expect(err).ToNot(HaveOccurred())

			// Create second tenant with its own context
			tenantID2 := uuid.New()
			ctx2 := buildCtx(tenantID2.String(), "CL")

			err = testCreateTenant(ctx2, tenantID2)
			Expect(err).ToNot(HaveOccurred())

			// 1. Obtener las tablas que tienen columna tenant_id
			var tableNames []string
			err = connection.Raw(`
				SELECT table_name 
				FROM information_schema.columns 
				WHERE column_name = 'tenant_id' AND table_schema = 'public'
			`).Scan(&tableNames).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(tableNames).To(ContainElement("order_headers")) // asegurarse de que la base existe

			// 2. Preparar JOIN dinámico usando goqu
			baseTable := "order_headers"
			baseAlias := "a" + uuid.NewString()[0:7]
			ds := goqu.From(goqu.T(baseTable).As(baseAlias))

			for _, table := range tableNames {
				if table == baseTable {
					continue
				}

				if avoidJoin(table) {
					continue
				}

				alias := "a" + uuid.NewString()[0:7]
				ds = ds.Join(goqu.T(table).As(alias),
					goqu.On(goqu.Ex{
						fmt.Sprintf("%s.tenant_id", alias): goqu.I(fmt.Sprintf("%s.tenant_id", baseAlias)),
					}))
			}

			// 3. WHERE por tenant + SELECT COUNT(*)
			ds = ds.Where(goqu.Ex{
				fmt.Sprintf("%s.tenant_id", baseAlias): tenantID2.String(),
			}).Select(goqu.COUNT("*"))

			sql, args, err := ds.ToSQL()
			Expect(err).ToNot(HaveOccurred())

			// Eliminar todas las comillas dobles de la query
			sql = strings.ReplaceAll(sql, `"`, "")

			// Print query in a format easy to copy-paste into TiDB
			fmt.Printf("\n=== COPY-PASTE THIS QUERY INTO TiDB ===\n")
			fmt.Printf("%s\n", sql)
			fmt.Printf("=== WITH THESE ARGS ===\n")
			fmt.Printf("%v\n", args)
			fmt.Printf("=====================================\n\n")

			var count int64
			err = connection.Raw(sql, args...).Scan(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
		})
	})
})

func avoidJoin(table string) bool {
	// Whitelist temporal — eliminar esta función cuando todas las entidades estén habilitadas
	switch table {
	case
		"order_headers",
		"contacts",
		"address_infos",
		"node_infos",
		"delivery_units",
		"order_types",
		"order_references",
		"order_delivery_units",
		"vehicle_categories",
		"vehicles",
		"carriers",
		"drivers",
		"vehicle_headers",
		"plans",
		"plan_headers",
		"states",
		"provinces",
		"districts",
		"node_info_headers",
		"node_types",
		//"non_delivery_reasons",
		"routes":
		return false // se permiten
	default:
		return true // evitar join
	}
}

func testCreateTenant(ctx context.Context, tenantID uuid.UUID) error {
	email := "ignaciovl.j@gmail.com"
	err := NewRegister(nil, tidbrepository.NewUpsertAccount(connection))(ctx, domain.UserCredentials{
		Email: email,
	})
	if err != nil {
		return err
	}
	return NewCreateTenant(
		tidbrepository.NewUpsertOrderHeaders(connection),
		tidbrepository.NewUpsertContact(connection),
		tidbrepository.NewUpsertAddressInfo(connection),
		tidbrepository.NewUpsertNodeInfo(connection),
		tidbrepository.NewUpsertDeliveryUnits(connection),
		tidbrepository.NewUpsertOrderType(connection),
		tidbrepository.NewUpsertOrder(connection),
		tidbrepository.NewUpsertOrderReferences(connection),
		tidbrepository.NewUpsertOrderDeliveryUnits(connection),
		tidbrepository.NewSaveTenant(connection),
		tidbrepository.NewUpsertVehicleCategory(connection),
		tidbrepository.NewUpsertCarrier(connection),
		tidbrepository.NewUpsertDriver(connection),
		tidbrepository.NewUpsertVehicle(connection),
		tidbrepository.NewUpsertVehicleHeaders(connection),
		tidbrepository.NewUpsertPlan(connection),
		tidbrepository.NewUpsertPlanHeaders(connection),
		tidbrepository.NewUpsertRoute(connection),
		tidbrepository.NewUpsertState(connection),
		tidbrepository.NewUpsertProvince(connection),
		tidbrepository.NewUpsertDistrict(connection),
		tidbrepository.NewUpsertNodeInfoHeaders(connection),
		tidbrepository.NewUpsertNodeType(connection),
	)(ctx, domain.Tenant{
		ID: tenantID,
		Operator: domain.Operator{
			Contact: domain.Contact{
				PrimaryEmail: email,
			},
		},
		Country: countries.CL,
	})
}
