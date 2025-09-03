package seed

import (
	"case-management/internal/domain/model"

	"gorm.io/gorm"
)

func SeedProducts(db *gorm.DB) {
	products := []model.Products{
		{Name: "Product 1"},
		{Name: "Product 2"},
	}

	for _, product := range products {
		if err := db.FirstOrCreate(&product, model.Products{Name: product.Name}).Error; err != nil {
			panic("failed to seed product: " + err.Error())
		}
	}
}
