package seed

import (
	"case-management/internal/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DispositionSubMap map[string]uuid.UUID

func SeedDispositionSub(db *gorm.DB, dispositionMain map[string]uuid.UUID) DispositionSubMap {
	dispositionSubMap := make(DispositionSubMap)

	dispositionSub := []model.DispositionSub{
		{NameEN: "Check available credit limit.", NameTH: "เช็ควงเงินใช้ได้ในบัตร", MainID: dispositionMain["Check Limit"]},
		{NameEN: "Check cash advance limit.", NameTH: "เช็ควงเงินในการกดเงินสด", MainID: dispositionMain["Check Limit"]},
		{NameEN: "Check purchase limit.", NameTH: "เช็ควงเงินในการรูดชื้อสินค้า", MainID: dispositionMain["Check Limit"]},
		{NameEN: "Check available limit after payment.", NameTH: "เช็ควงเงินเมื่อชำระเข้าแล้ววงเงินบัตรใช้ได้เท่าไหร่", MainID: dispositionMain["Check Limit"]},
		{NameEN: "What is the maximum credit limit?", NameTH: "วงเงินสูงสุดของบัตรวงเงินเท่าไหร่ / วงเงินอนุมัติเท่าไหร่", MainID: dispositionMain["Check Limit"]},
	}

	for _, dispoSub := range dispositionSub {
		if err := db.FirstOrCreate(&dispoSub, model.DispositionSub{NameEN: dispoSub.NameEN}).Error; err != nil {
			panic("Failed to seed disposition sub: " + err.Error())
		}
		dispositionSubMap[dispoSub.NameEN] = dispoSub.ID
	}

	return dispositionSubMap
}
