package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"github.com/nurullahgd/product-app/domain"
	constants "github.com/nurullahgd/product-app/persistence/common"
)

type IProductRepository interface {
	GetAllProducts() []domain.Product
	GetAllProductsByStore(storeName string) []domain.Product
	AddProduct(product domain.Product) error
	GetById(productId int64) (domain.Product, error)
}

func NewProductRepository(dbPool *pgxpool.Pool) IProductRepository {
	return &ProductRepository{
		dbPool: dbPool,
	}
}

type ProductRepository struct {
	dbPool *pgxpool.Pool
}

func (productRepository *ProductRepository) GetAllProducts() []domain.Product {
	ctx := context.Background()
	productRows, err := productRepository.dbPool.Query(ctx, "Select * from products")
	if err != nil {
		log.Errorf("Error While Getting All Products %v", err)
		return []domain.Product{}
	}

	return ExtractProductsFromRows(productRows)
}

func (productRepository *ProductRepository) GetAllProductsByStore(storeName string) []domain.Product {
	ctx := context.Background()

	getProductsByStoreNameSql := `Select * from products where store = $1`

	productRows, err := productRepository.dbPool.Query(ctx, getProductsByStoreNameSql, storeName)

	if err != nil {
		log.Errorf("Error While Getting All Products %v", err)
		return []domain.Product{}
	}

	return ExtractProductsFromRows(productRows)
}

func (productRepository *ProductRepository) AddProduct(product domain.Product) error {
	ctx := context.Background()
	insertSql := `INSERT INTO products (name,price,discount,store) VALUES ($1,$2,$3,$4)`
	addNewProduct, err := productRepository.dbPool.Exec(ctx, insertSql, product.Name, product.Price, product.Discount, product.Store)
	if err != nil {
		log.Error("Failed to add new product", err)
		return err
	}
	log.Info(fmt.Printf("Product added with %v", addNewProduct))
	return nil

}
func ExtractProductsFromRows(productRows pgx.Rows) []domain.Product {
	var products = []domain.Product{}

	var id int64
	var name string
	var price float32
	var discount float32
	var store string

	for productRows.Next() {
		productRows.Scan(&id, &name, &price, &discount, &store)
		products = append(products, domain.Product{
			Id:       id,
			Name:     name,
			Price:    price,
			Discount: discount,
			Store:    store,
		})
	}

	return products
}

func (ProductRepository *ProductRepository) GetById(productId int64) (domain.Product, error) {
	ctx := context.Background()

	getByIdSql := `Select * from products where id=$1`

	queryRow := ProductRepository.dbPool.QueryRow(ctx, getByIdSql, productId)

	var id int64
	var name string
	var price float32
	var discount float32
	var store string

	err := queryRow.Scan(&id, &name, &price, &discount, &store)
	if err != nil && err.Error() == constants.NOT_FOUND {
		return domain.Product{}, errors.New(fmt.Sprintf("Product not found with id %d", productId))
	}
	if err != nil {
		return domain.Product{}, errors.New(fmt.Sprintf("Something Wrong, Check To ID: %d", productId))
	}
	return domain.Product{
		Id:       id,
		Name:     name,
		Price:    price,
		Discount: discount,
		Store:    store,
	}, nil
}
