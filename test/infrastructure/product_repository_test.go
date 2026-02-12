package infrastructure

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nurullahgd/product-app/common/postgresql"
	"github.com/nurullahgd/product-app/domain"
	"github.com/nurullahgd/product-app/persistence"
	"github.com/stretchr/testify/assert"
)

var ProductRepository persistence.IProductRepository
var dbPool *pgxpool.Pool
var ctx context.Context

func TestMain(m *testing.M) {
	ctx = context.Background()

	dbPool = postgresql.GetConnectionPool(ctx, postgresql.Config{
		Host:                  "localhost",
		Port:                  "6432",
		DbName:                "productapp",
		UserName:              "postgres",
		Password:              "postgres",
		MaxConnections:        "10",
		MaxConnectionIdleTime: "30s",
	})
	ProductRepository = persistence.NewProductRepository(dbPool)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setup(ctx context.Context, dbpool *pgxpool.Pool) {
	TestDataInitialize(ctx, dbpool)

}

func clear(ctx context.Context, dbpool *pgxpool.Pool) {
	TruncateTestData(ctx, dbpool)
}
func TestGetAllProducts(t *testing.T) {
	setup(ctx, dbPool)
	fmt.Println("TestAllProduct")

	expectedProduct := []domain.Product{
		{
			Id:       1,
			Name:     "AirFryer",
			Price:    3000.0,
			Discount: 22.0,
			Store:    "ABC TECH",
		},
		{
			Id:       2,
			Name:     "Ütü",
			Price:    1500.0,
			Discount: 10.0,
			Store:    "ABC TECH",
		},
		{
			Id:       3,
			Name:     "Çamaşır Makinesi",
			Price:    10000.0,
			Discount: 15.0,
			Store:    "ABC TECH",
		},
		{
			Id:       4,
			Name:     "Lambader",
			Price:    2000.0,
			Discount: 0.0,
			Store:    "Dekorasyon Sarayı",
		},
	}
	t.Run("GetAllProduct", func(t *testing.T) {
		actualProducts := ProductRepository.GetAllProducts()
		assert.Equal(t, 4, len(actualProducts))
		assert.Equal(t, expectedProduct, actualProducts)
	})
	clear(ctx, dbPool)
}

func TestGetAllProductsByStore(t *testing.T) {
	setup(ctx, dbPool)
	fmt.Println("TestAllProduct")

	expectedProduct := []domain.Product{
		{
			Id:       1,
			Name:     "AirFryer",
			Price:    3000.0,
			Discount: 22.0,
			Store:    "ABC TECH",
		},
		{
			Id:       2,
			Name:     "Ütü",
			Price:    1500.0,
			Discount: 10.0,
			Store:    "ABC TECH",
		},
		{
			Id:       3,
			Name:     "Çamaşır Makinesi",
			Price:    10000.0,
			Discount: 15.0,
			Store:    "ABC TECH",
		},
	}
	t.Run("GetAllProductsByStore", func(t *testing.T) {
		actualProducts := ProductRepository.GetAllProductsByStore("ABC TECH")
		assert.Equal(t, 3, len(actualProducts))
		assert.Equal(t, expectedProduct, actualProducts)
	})
	clear(ctx, dbPool)
}
func TestAddProduct(t *testing.T) {

	fmt.Println("TestAllProduct")

	expectedProduct := []domain.Product{
		{
			Id:       1,
			Name:     "AirFryer",
			Price:    3000.0,
			Discount: 22.0,
			Store:    "ABC TECH",
		},
	}
	newProduct := domain.Product{
		Name:     "AirFryer",
		Price:    3000.0,
		Discount: 22.0,
		Store:    "ABC TECH",
	}
	t.Run("GetAllProductsByStore", func(t *testing.T) {
		ProductRepository.AddProduct(newProduct)
		actualProducts := ProductRepository.GetAllProducts()
		assert.Equal(t, 1, len(actualProducts))
		assert.Equal(t, expectedProduct, actualProducts)
	})
	clear(ctx, dbPool)
}

func TestGetProductById(t *testing.T) {
	setup(ctx, dbPool)
	t.Run("GetProductById", func(t *testing.T) {
		product, _ := ProductRepository.GetById(1)
		_, err := ProductRepository.GetById(5)
		assert.Equal(t, "Product not found with id 5", err.Error())
		fmt.Println("pr", product)
		assert.Equal(t, domain.Product{
			Id:       1,
			Name:     "AirFryer",
			Price:    3000.0,
			Discount: 22.0,
			Store:    "ABC TECH"}, product)
	})

	clear(ctx, dbPool)
}
