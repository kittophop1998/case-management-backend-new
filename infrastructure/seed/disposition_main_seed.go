package seed

import (
	"case-management/internal/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DispositionMainMap map[string]uuid.UUID

func SeedDispositionMain(db *gorm.DB) DispositionMainMap {
	dispositionMainMap := make(DispositionMainMap)

	dispositionsMain := []model.DispositionMain{
		{NameEN: "Check Limit", NameTH: "ตรวจสอบขีดจำกัด", Description: "Check Limit"},
		{NameEN: "Payment", NameTH: "การชำระเงิน", Description: "Payment"},
		{NameEN: "Approval result", NameTH: "ผลการอนุมัติ", Description: "Approval result"},
		{NameEN: "Activate card", NameTH: "เปิดใช้งานบัตร", Description: "Activate card"},
		{NameEN: "Approval result", NameTH: "ผลการอนุมัติ", Description: "Approval result"},
		{NameEN: "Can not activate", NameTH: "ไม่สามารถเปิดใช้งาน", Description: "Can not activate"},
		{NameEN: "Aeon Happy Plan (Credit Card)", NameTH: "Aeon Happy Plan (บัตรเครดิต)", Description: "Aeon Happy Plan (Credit Card)"},
		{NameEN: "Aeon Happy Pay (Your Cash)", NameTH: "Aeon Happy Pay (เงินสดของคุณ)", Description: "Aeon Happy Pay (Your Cash)"},
		{NameEN: "Aeon Happy Point", NameTH: "Aeon Happy Point", Description: "Aeon Happy Point"},
		{NameEN: "Privileges Card", NameTH: "บัตรสิทธิพิเศษ", Description: "Privileges Card"},
		{NameEN: "Promotion", NameTH: "โปรโมชั่น", Description: "Promotion"},
		{NameEN: "Follow promotion", NameTH: "ติดตามโปรโมชั่น", Description: "Follow promotion"},
		{NameEN: "Re-Issue Card", NameTH: "ออกบัตรใหม่", Description: "Re-Issue Card"},
		{NameEN: "Re-Issue Pin", NameTH: "ออกรหัสใหม่", Description: "Re-Issue Pin"},
		{NameEN: "Annual Fee", NameTH: "ค่าธรรมเนียมรายปี", Description: "Annual Fee"},
		{NameEN: "Return Card", NameTH: "ส่งคืนบัตร", Description: "Return Card"},
		{NameEN: "Increase Limit", NameTH: "เพิ่มวงเงิน", Description: "Increase Limit"},
		{NameEN: "Decrease Credit Limit", NameTH: "ลดวงเงินบัตรเครดิต", Description: "Decrease Credit Limit"},
		{NameEN: "Request to adjust the income base", NameTH: "คำขอปรับฐานรายได้", Description: "Request to adjust the income base"},
		{NameEN: "Card sending Status", NameTH: "สถานะการจัดส่งบัตร", Description: "Card sending Status"},
		{NameEN: "Pin sending Status", NameTH: "สถานะการจัดส่งรหัส", Description: "Pin sending Status"},
		{NameEN: "Interest Rate", NameTH: "อัตราดอกเบี้ย", Description: "Interest Rate"},
		{NameEN: "ESM", NameTH: "ESM", Description: "ESM"},
		{NameEN: "Rabbit Function", NameTH: "ฟังก์ชัน Rabbit", Description: "Rabbit Function"},
		{NameEN: "SMS", NameTH: "ข้อความ SMS", Description: "SMS"},
		{NameEN: "Check card status", NameTH: "ตรวจสอบสถานะบัตร", Description: "Check card status"},
		{NameEN: "Cancel / Stop Insurance", NameTH: "ยกเลิก / หยุดการประกัน", Description: "Cancel / Stop Insurance"},
		{NameEN: "Detail Insurance", NameTH: "รายละเอียดประกัน", Description: "Detail Insurance"},
		{NameEN: "Transfer call to Authorize", NameTH: "โอนสายไปยัง Authorize", Description: "Transfer call to Authorize"},
		{NameEN: "Transfer call to CMS", NameTH: "โอนสายไปยัง CMS", Description: "Transfer call to CMS"},
		{NameEN: "Transfer call To Collection", NameTH: "โอนสายไปยัง Collection", Description: "Transfer call To Collection"},
		{NameEN: "Not Transfer call To Collection", NameTH: "ไม่โอนสายไปยัง Collection", Description: "Not Transfer call To Collection"},
		{NameEN: "Tranfer to other Dept.", NameTH: "โอนไปแผนกอื่น", Description: "Tranfer to other Dept."},
		{NameEN: "Telephone Connection", NameTH: "การเชื่อมต่อสายโทรศัพท์", Description: "Telephone Connection"},
		{NameEN: "NCB", NameTH: "ข้อมูลเครดิต (NCB)", Description: "NCB"},
		{NameEN: "Closing Account Letter", NameTH: "จดหมายปิดบัญชี", Description: "Closing Account Letter"},
		{NameEN: "Change Information", NameTH: "เปลี่ยนแปลงข้อมูล", Description: "Change Information"},
		{NameEN: "Cancel Card", NameTH: "ยกเลิกบัตร", Description: "Cancel Card"},
	}

	for _, disposition := range dispositionsMain {
		if err := db.FirstOrCreate(&disposition, model.DispositionMain{NameEN: disposition.NameEN}).Error; err != nil {
			panic("Failed to seed disposition main, " + err.Error())
		}
		dispositionMainMap[disposition.NameEN] = disposition.ID
	}
	return dispositionMainMap
}
