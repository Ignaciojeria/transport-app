package tidbrepository

import (
	"context"
	"log"
	"testing"
	"time"
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
	defer mysqlContainer.Terminate(ctx)

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

	// Pass the GORM DB instance to your repository function
	NewSaveTransportOrder(tidb.TIDBConnection{
		DB: db,
	})
}
