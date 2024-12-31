package tidbrepository

import (
	"context"
	"log"
	"testing"
	"time"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/tidb"

	"github.com/testcontainers/testcontainers-go"
	tidbmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSaveTransportOrder(t *testing.T) {
	ctx := context.Background()

	// Start MySQL container
	mysqlContainer, err := tidbmysql.Run(ctx,
		"mysql:8.0.36",
		tidbmysql.WithDatabase("foo"),
		tidbmysql.WithUsername("root"),
		tidbmysql.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("3306/tcp").WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("Failed to start MySQL container: %v", err)
	}
	defer func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate MySQL container: %v", err)
		}
	}()

	// Get connection string
	dsn, err := mysqlContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	// Initialize GORM with the MySQL connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize GORM: %v", err)
	}
	log.Printf("Connected to database: %s", dsn)

	// Set up TiDB connection
	tiDBConn := tidb.TIDBConnection{
		DB: db,
	}

	// Migrate tables
	err = table.MigrateTables(tiDBConn, configuration.TiDBConfiguration{
		TIDB_RUN_MIGRATIONS: "true",
	})
	if err != nil {
		t.Fatalf("Failed running migrations: %v", err)
	}

	// Run subtests
	t.Run("InsertTransportOrder", func(t *testing.T) {
		order := domain.TransportOrder{
			Tenant: domain.Tenant{
				Organization: "my-org-name",
				Commerce:     "UNIMARC",
				Consumer:     "CROSS-COMMERCE-API",
			},
			ReferenceID: "1234",
		}
		// Save transport order
		saveOrderFunc := NewSaveTransportOrder(tiDBConn)
		_, err := saveOrderFunc(ctx, order)
		if err != nil {
			t.Fatalf("Failed to save transport order: %v", err)
		}
	})

}
