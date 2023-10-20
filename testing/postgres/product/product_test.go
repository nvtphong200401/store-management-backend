package postgres_test

import (
	"net/http"
	"testing"

	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
	"github.com/nvtphong200401/store-management/pkg/handlers/usecases"
	postgres_test "github.com/nvtphong200401/store-management/testing/postgres"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type productSuiteTest struct {
	postgres_test.PostgresSuite
}

func TestProductSuite(t *testing.T) {

	productSuite := &productSuiteTest{
		postgres_test.PostgresSuite{},
	}
	suite.Run(t, productSuite)
}

func (s *productSuiteTest) TearDownTest() {

	err := s.TearDownSuite()
	require.NoError(s.T(), err)

}

func (s *productSuiteTest) TestReAdd() {
	repo := respository.NewProductRepository(&s.TxStore)
	products := []models.Product{{
		ID:          "123",
		ProductName: "Phongne",
		PriceIn:     20000,
		PriceOut:    30000,
		Stock:       20,
		StoreID:     1,
	}}
	err := repo.AddProduct(products)
	require.NoError(s.T(), err)
	err = repo.DeleteProduct(products)
	require.NoError(s.T(), err)
	err = repo.AddProduct(products)
	require.NoError(s.T(), err)
}

func (s *productSuiteTest) TestAddDuplicate() {
	repo := respository.NewProductRepository(&s.TxStore)
	productUseCase := usecases.NewProductUseCases(repo)
	products := []models.Product{{
		ID:          "123",
		ProductName: "Phongne",
		PriceIn:     20000,
		PriceOut:    30000,
		Stock:       20,
		StoreID:     1,
	}}
	err := repo.AddProduct(products)
	require.NoError(s.T(), err)
	err = repo.AddProduct(products)
	require.NoError(s.T(), err)
	status, header := productUseCase.GetProducts(1, 1, 10)
	require.Equal(s.T(), status, http.StatusOK)
	totalItems := header["totalItems"]
	require.Equal(s.T(), 1, totalItems)
}

func (s *productSuiteTest) TestSearch() {
	repo := respository.NewProductRepository(&s.TxStore)
	productUseCase := usecases.NewProductUseCases(repo)
	products := []models.Product{{
		ID:          "123",
		ProductName: "Đây là rap việt",
		PriceIn:     20000,
		PriceOut:    30000,
		Stock:       20,
		StoreID:     1,
	}}
	err := repo.AddProduct(products)
	require.NoError(s.T(), err)
	status, header := productUseCase.SearchProducts("day", 1, 1, 10)
	require.Equal(s.T(), status, http.StatusOK)
	data := header["data"].([]models.Product)
	require.Equal(s.T(), 1, len(data))
	require.Equal(s.T(), "Đây là rap việt", data[0].ProductName)
}
