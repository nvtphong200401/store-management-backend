package sale_test

import (
	"encoding/json"
	"fmt"
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

func (s *saleSuiteTest) TestSaleProducts() {

	productRepo := respository.NewProductRepository(&s.TxStore)
	products := []models.Product{{
		ID:          "123",
		ProductName: "Phongne",
		PriceIn:     20000,
		PriceOut:    30000,
		Stock:       20,
		StoreID:     1,
	}}

	err := productRepo.AddProduct(products)
	require.NoError(s.T(), err)

	saleRepo := respository.NewSaleRepository(&s.TxStore)
	bytes, err := json.Marshal(products)
	require.NoError(s.T(), err)
	var saleItems []models.SaleItem
	err = json.Unmarshal(bytes, &saleItems)
	require.NoError(s.T(), err)
	saleItems[0].Stock = 8
	statusCode, _ := saleRepo.SellItems(saleItems, 1, 1)
	require.Equal(s.T(), 200, statusCode)
	code, header := productRepo.GetProducts(1, 1, 10)
	require.Equal(s.T(), 200, code)
	productAfter := header["data"].([]models.Product)
	fmt.Println(productAfter)
	require.NotEmpty(s.T(), productAfter)
	require.Equal(s.T(), 12, productAfter[0].Stock)
}

func (s *saleSuiteTest) TestInvalidStock() {

	productRepo := respository.NewProductRepository(&s.TxStore)
	products := []models.Product{{
		ID:          "123",
		ProductName: "Phongne",
		PriceIn:     20000,
		PriceOut:    30000,
		Stock:       20,
		StoreID:     1,
	}}

	err := productRepo.AddProduct(products)
	require.NoError(s.T(), err)

	saleRepo := respository.NewSaleRepository(&s.TxStore)
	bytes, err := json.Marshal(products)
	require.NoError(s.T(), err)
	var saleItems []models.SaleItem
	err = json.Unmarshal(bytes, &saleItems)
	require.NoError(s.T(), err)
	saleItems[0].Stock = 30
	statusCode, _ := saleRepo.SellItems(saleItems, 1, 1)
	require.NotEqual(s.T(), 200, statusCode)
}

func (s *saleSuiteTest) TestBuyProduct() {
	products := []models.Product{{
		ID:          "123",
		ProductName: "Phongne",
		PriceIn:     20000,
		PriceOut:    30000,
		Stock:       20,
		StoreID:     1,
	}}

	bytes, err := json.Marshal(products)
	require.NoError(s.T(), err)
	var saleItems []models.SaleItem
	err = json.Unmarshal(bytes, &saleItems)
	require.NoError(s.T(), err)

}
