package utils

import (
	"encoding/json"

	"github.com/nvtphong200401/store-management/pkg/handlers/models"
)

func ProductsToSaleItems(products []models.Product) ([]models.SaleItem, error) {
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
