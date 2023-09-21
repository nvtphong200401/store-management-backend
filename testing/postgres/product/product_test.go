package postgres_test

import (
	"net/http"
	"testing"

	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
	postgres_test "github.com/nvtphong200401/store-management/testing/postgres"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type productSuiteTest struct {
	postgres_test.PostgresSuite
}

func TestAuthSuite(t *testing.T) {

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
	status, header := repo.GetProducts(1, 1, 10)
	require.Equal(s.T(), status, http.StatusOK)
	totalItems := header["totalItems"]
	require.Equal(s.T(), int64(1), totalItems)
}
