package http

import (
	"case-management/infrastructure/config"
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"case-management/utils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	UseCase usecase.DashboardUseCase
	Config  *config.Config
}

func (h *DashboardHandler) GetCustInfo(ctx *gin.Context) {
	id := ctx.Param("aeon_id")

	reqID := ctx.GetHeader("X-Request-ID")

	c := context.WithValue(ctx.Request.Context(), utils.CtxKeyApisKey, h.Config.Headers.ApiKey)
	c = context.WithValue(c, utils.CtxKeyApiLang, h.Config.Headers.ApiLanguage)
	c = context.WithValue(c, utils.CtxKeyDeviceOS, h.Config.Headers.ApiDeviceOS)
	c = context.WithValue(c, utils.CtxKeyChannel, h.Config.Headers.ApiChannel)
	c = context.WithValue(c, utils.CtxKeyRequestID, reqID)

	_, err := h.UseCase.CustInfo(c, id)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	lib.HandleError(ctx, lib.NotFound.WithDetails("Customer not found"))

	mock := &model.GetCustInfoResponse{
		NationalID:      "1234",
		CustomerNameEng: "Jane Doe",
		CustomerNameTH:  "เจน โด",
		MobileNO:        "0812345678",
		MailToAddress:   "123/45 หมู่บ้านสุขสันต์ ถ.สุขสวัสดิ์ แขวงบางปะกอก เขตราษฎร์บูรณะ กทม. 10140",
		MailTo:          "Jane Doe",
	}

	lib.HandleResponse(ctx, http.StatusOK, mock)
}

func (h *DashboardHandler) GetCustProfile(ctx *gin.Context) {
	// id := ctx.Param("id")
	// customer, err := h.UseCase.CustProfile(ctx, id)
	// if err != nil {
	// 	lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
	// 	return
	// }

	customerMock := &model.GetCustProfileResponse{
		LastCardApplyDate:             "25 Aug 2023",
		CustomerSentiment:             "",
		PhoneNoLastUpdateDate:         "01 Aug 2024",
		LastIncreaseCreditLimitUpdate: "",
		LastReduceCreditLimitUpdate:   "",
		LastIncomeUpdate:              "29 Aug 2023",
		SuggestedAction:               "Update salary slip",
		TypeOfJob:                     "PRIVATE COMPANY",
		MaritalStatus:                 "Single",
		Gender:                        "Female",
		LastEStatementSentDate:        "",
		EStatementSentStatus:          "",
		StatementChannel:              "Paper",
		ConsentForDisclose:            "Incomplete",
		BlockMedia:                    "No blocked",
		ConsentForCollectUse:          "Incomplete",
	}
	lib.HandleResponse(ctx, http.StatusOK, customerMock)
}

func (h *DashboardHandler) GetCustSegment(ctx *gin.Context) {
	// id := ctx.Param("id")
	// customer, err := h.UseCase.CustSegment(ctx, id)
	// if err != nil {
	// 	lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
	// 	return
	// }

	customerMock := &model.GetCustSegmentResponse{
		Sweetheart:      "",
		ComplaintLevel:  "Complaint Level: ",
		CustomerGroup:   "NORMAL",
		ComplaintGroup:  "",
		CustomerType:    "MD",
		MemberStatus:    "BAD CREDIT",
		CustomerSegment: "Existing Customer - Active",
		UpdateData:      "01 Jan 0001",
	}
	lib.HandleResponse(ctx, http.StatusOK, customerMock)
}

func (h *DashboardHandler) GetCustSuggestion(ctx *gin.Context) {
	// id := ctx.Param("id")
	// customer, err := h.UseCase.CustSuggestion(ctx, id)
	// if err != nil {
	// 	lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
	// 	return
	// }

	customerMock := model.GetCustSuggestionResponse{
		SuggestCards: []string{"No suggestion"},
		SuggestPromotions: []model.GetCustSuggestionPromotionResponse{
			{
				PromotionCode:            "P24099EEBE",
				PromotionName:            "BIC CAMERA Coupon with Aeon Credit Card",
				PromotionDetails:         "ซื้อสินค้าปลอดภาษี สูงสุด 10%  และ รับส่วนลด สูงสุด 7% เมื่อซื้อสินค้าที่ร้าน BicCamera ประเทศญี่ปุ่น, ร้าน Air BicCamera และ ร้าน KOJIMA ด้วยบัตรเครดิตอิออนทุกประเภท (ยกเว้นบัตรเครดิตเพื่อองค์กร) ซึ่ง BicCamera เป็นห้างสรรพสินค้าในประเทศญี่ปุ่น จำหน่ายสินค้าหลากหลายประเภท เช่น เครื่องใช้ไฟฟ้า ยา เครื่องสำอาง และของใช้ในชีวิตประจำวัน โปรดแสดงภาพบาร์โค้ดบนสื่อประชาสัมพันธ์นี้ ที่แคชเชียร์",
				Action:                   "Apply",
				PromotionResultTimestamp: "10 Feb 2025, 16.07",
				Period:                   "4 Sep 2024 - 31 Aug 2025",
				EligibleCard:             []string{"BIG C WORLD MASTERCARD"},
			},
		},
	}

	lib.HandleResponse(ctx, http.StatusOK, customerMock)
}
