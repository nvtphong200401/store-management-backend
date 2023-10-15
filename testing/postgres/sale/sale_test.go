package sale_test

import (
	"encoding/json"
	"testing"

	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
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
	saleRepo := respository.NewSaleRepository(&s.TxStore)
	products := []models.Product{{
		ID:          "123",
		ProductName: "Phongne",
		PriceIn:     20000,
		PriceOut:    30000,
		Stock:       20,
		StoreID:     1,
	}}

	saleItems, err := parseModel(products)
	require.NoError(s.T(), err)
	statusCode, _ := saleRepo.BuyItems(saleItems, 1, 1)
	require.Equal(s.T(), 200, statusCode)

	saleItems[0].Stock = 8

	statusCode, _ = saleRepo.SellItems(saleItems, 1, 1)
	require.Equal(s.T(), 200, statusCode)

	code, header := productRepo.GetProducts(1, 1, 10)
	require.Equal(s.T(), 200, code)

	productAfter := header["data"].([]models.Product)
	require.NotEmpty(s.T(), productAfter)

	require.Equal(s.T(), uint(12), productAfter[0].Stock)
}

func (s *saleSuiteTest) TestInvalidStock() {

	saleRepo := respository.NewSaleRepository(&s.TxStore)
	products := []models.Product{{
		ID:          "123",
		ProductName: "Phongne",
		PriceIn:     20000,
		PriceOut:    30000,
		Stock:       20,
		StoreID:     1,
	}}
	saleItems, err := parseModel(products)
	require.NoError(s.T(), err)

	statusCode, _ := saleRepo.BuyItems(saleItems, 1, 1)
	require.Equal(s.T(), 200, statusCode)

	saleItems[0].Stock = 30

	statusCode, _ = saleRepo.SellItems(saleItems, 1, 1)
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

	saleItems, err := parseModel(products)
	require.NoError(s.T(), err)
	saleItems[0].Stock = 8

	saleRepo := respository.NewSaleRepository(&s.TxStore)
	statusCode, _ := saleRepo.BuyItems(saleItems, 1, 1)
	require.Equal(s.T(), 200, statusCode)

	productRepo := respository.NewProductRepository(&s.TxStore)
	code, header := productRepo.SearchProduct("20042001", 1, 1, 10)
	require.Equal(s.T(), 200, code)
	productAfter := header["data"].([]models.Product)
	require.NotEmpty(s.T(), productAfter)
}

func parseModel(products []models.Product) ([]models.SaleItem, error) {
	bytes, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}
	var saleItems []models.SaleItem
	err = json.Unmarshal(bytes, &saleItems)
	if err != nil {
		return nil, err
	}

	for i := range saleItems {
		saleItems[i].Product = products[i]
	}
	return saleItems, nil
}
