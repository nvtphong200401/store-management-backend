package sale_test

import (
	"testing"

	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
	"github.com/nvtphong200401/store-management/pkg/handlers/usecases"
	postgres_test "github.com/nvtphong200401/store-management/testing/postgres"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type saleSuiteTest struct {
	postgres_test.PostgresSuite
}

func TestSaleSuite(t *testing.T) {

	saleSuite := &saleSuiteTest{
		postgres_test.PostgresSuite{},
	}
	suite.Run(t, saleSuite)
}

func (s *saleSuiteTest) TearDownTest() {

	err := s.TearDownSuite()
	require.NoError(s.T(), err)

}

func (s *saleSuiteTest) TestBuyAndSellProducts() {

	productRepo := respository.NewProductRepository(&s.TxStore)
	productUseCase := usecases.NewProductUseCases(productRepo)
	saleRepo := respository.NewSaleRepository(&s.TxStore)
	saleUseCase := usecases.NewSaleUseCases(saleRepo)
	products := []models.Product{{
		ID:          "123",
		ProductName: "Phongne",
		PriceIn:     20000,
		PriceOut:    30000,
		Stock:       20,
		StoreID:     1,
	}}

	statusCode, _ := saleUseCase.BuyItems(products, 1, 1)
	require.Equal(s.T(), 200, statusCode)

	products[0].Stock = 8

	statusCode, _ = saleUseCase.SellItems(products, 1, 1)
	require.Equal(s.T(), 200, statusCode)

	code, header := productUseCase.GetProducts(1, 1, 10)
	require.Equal(s.T(), 200, code)

	productAfter := header["data"].([]models.Product)
	require.NotEmpty(s.T(), productAfter)

	require.Equal(s.T(), uint(12), productAfter[0].Stock)
}

func (s *saleSuiteTest) TestInvalidStock() {

	saleRepo := respository.NewSaleRepository(&s.TxStore)
	saleUseCase := usecases.NewSaleUseCases(saleRepo)
	products := []models.Product{{
		ID:          "123",
		ProductName: "Phongne",
		PriceIn:     20000,
		PriceOut:    30000,
		Stock:       20,
		StoreID:     1,
	}}

	statusCode, _ := saleUseCase.BuyItems(products, 1, 1)
	require.Equal(s.T(), 200, statusCode)

	products[0].Stock = 30

	statusCode, _ = saleUseCase.SellItems(products, 1, 1)
	require.NotEqual(s.T(), 200, statusCode)
}

func (s *saleSuiteTest) TestBuyProductNotExist() {
	products := []models.Product{{
		ID:          "20042001",
		ProductName: "Phongne",
		PriceIn:     20000,
		PriceOut:    30000,
		Stock:       20,
		StoreID:     1,
	}}

	saleRepo := respository.NewSaleRepository(&s.TxStore)
	saleUseCase := usecases.NewSaleUseCases(saleRepo)
	statusCode, _ := saleUseCase.BuyItems(products, 1, 1)
	require.Equal(s.T(), 200, statusCode)

	productRepo := respository.NewProductRepository(&s.TxStore)
	productUseCase := usecases.NewProductUseCases(productRepo)
	code, header := productUseCase.SearchProducts("20042001", 1, 1, 10)
	require.Equal(s.T(), 200, code)
	productAfter := header["data"].([]models.Product)
	require.NotEmpty(s.T(), productAfter)
}
