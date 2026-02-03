package infrastructer

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nurullahgd/product-app/common/postgresql"
	"github.com/nurullahgd/product-app/persistence"
)

var ProductRepository persistence.IProductRepository
var dbPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbPool = postgresql.GetConnectionPool(ctx, postgresql.Config{
		Host:                  "localhost",
		Port:                  "6432",
		DbName:                "productapp",
		UserName:              "postgres",
		Password:              "postgres",
		MaxConnections:        "10",
		MaxConnectionIdleTime: "30s",
	})
	fmt.Println("Before All Tests")
	ProductRepository = persistence.NewProductRepository(dbPool)
	exitCode := m.Run()
	fmt.Println("After all tests")
	os.Exit(exitCode)
}

func TestGetAllProducts(t *testing.T) {
	fmt.Println("testo")
	fmt.Println(ProductRepository)
}
