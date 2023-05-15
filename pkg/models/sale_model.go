package models

import "gorm.io/gorm"

type SaleModel struct {
	gorm.Model
	StoreID    uint
	Store      StoreModel `gorm:"foreignKey:StoreID;"`
	EmployeeID uint
	Employee   Employee `gorm:"foreignKey:EmployeeID"`
	TotalPrice float64  `json:"TotalPrice"`
}

type SaleItem struct {
	gorm.Model
	SaleID    uint
	Sale      SaleModel `gorm:"foreignKey:SaleID"`
	ProductID uint      `json:"ProductID,omitempty"`
	Product   Product   `gorm:"foreignKey:ProductID"`
	Quantity  uint      `json:"Quantity"`
}

type SaleService struct {
}

func createSale(storeID uint, employeeID uint, totalPrice float64) (uint, error) {
	db.AutoMigrate(&SaleModel{})

	sale := SaleModel{StoreID: storeID, EmployeeID: employeeID, TotalPrice: totalPrice}
	if err := db.Create(&sale).Error; err != nil {
		return 0, err
	}
	return sale.ID, nil
}

func (ss *SaleService) PurchaseItems(items []SaleItem, employeeID uint, storeID uint) error {
	db.AutoMigrate(&SaleItem{})
	totalPrice := calculateTotalPrice(items)
	saleid, err := createSale(storeID, employeeID, totalPrice)
	if err != nil {
		return err
	}
	for index := range items {
		items[index].SaleID = saleid
	}

	if err = db.Create(&items).Error; err != nil {
		return err
	}
	return nil
}

func calculateTotalPrice(items []SaleItem) float64 {
	var total float64 = 0
	for _, item := range items {
		var product Product
		if err := db.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			return 0
		}
		total += product.Price * float64(item.Quantity)
	}
	return total
}
