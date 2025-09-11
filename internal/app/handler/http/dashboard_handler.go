package http

import (
	"case-management/infrastructure/config"
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"case-management/utils"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	UseCase usecase.DashboardUseCase
	Config  *config.Config
}

// func (h *DashboardHandler) GetCustInfo(ctx *gin.Context) {
// 	log.Println("[Handler] => Entering GetCustInfo")

// 	Id := ctx.Param("id")
// 	if Id == "" {
// 		lib.HandleError(ctx, fmt.Errorf("missing aeon_id parameter"))
// 		return
// 	}

// 	var req model.ConnectorCustomerInfoRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		log.Println("[Handler] => Failed to bind JSON:", err)
// 		lib.HandleError(ctx, fmt.Errorf("invalid request body: %w", err))
// 		return
// 	}

// 	cfg, err := config.Load("sit")
// 	if err != nil {
// 		lib.HandleError(ctx, fmt.Errorf("internal server error"))
// 		return
// 	}

// 	// เตรียม context.Context พร้อมค่าต่างๆ
// 	c := ctx.Request.Context()
// 	c = context.WithValue(c, utils.CtxKeyApisKey, cfg.Headers.ApiKey)
// 	c = context.WithValue(c, utils.CtxKeyApiLang, cfg.Headers.ApiLanguage)
// 	c = context.WithValue(c, utils.CtxKeyDeviceOS, cfg.Headers.ApiDeviceOS)
// 	c = context.WithValue(c, utils.CtxKeyChannel, cfg.Headers.ApiChannel)

// 	reqID := ctx.GetHeader("X-Request-ID")
// 	if reqID != "" {
// 		c = context.WithValue(c, utils.CtxKeyRequestID, reqID)
// 	}

// 	resp, err := h.UseCase.CustInfo(c, ctx, req)
// 	if err != nil {
// 		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
// 		return
// 	}

// 	lib.HandleResponse(ctx, http.StatusOK, resp)
// }

func (h *DashboardHandler) GetCustInfo(ctx *gin.Context) {
	log.Println("[Handler] => Entering GetCustInfo")

	id := ctx.Param("id")
	if id == "" {
		lib.HandleError(ctx, fmt.Errorf("missing customer id parameter"))
		return
	}

	cfg, err := config.Load("sit")
	if err != nil {
		lib.HandleError(ctx, fmt.Errorf("internal server error"))
		return
	}

	// เตรียม context.Context พร้อมค่าต่างๆ
	c := ctx.Request.Context()
	c = context.WithValue(c, utils.CtxKeyApisKey, cfg.Headers.ApiKey)
	c = context.WithValue(c, utils.CtxKeyApiLang, cfg.Headers.ApiLanguage)
	c = context.WithValue(c, utils.CtxKeyDeviceOS, cfg.Headers.ApiDeviceOS)
	c = context.WithValue(c, utils.CtxKeyChannel, cfg.Headers.ApiChannel)

	reqID := ctx.GetHeader("X-Request-ID")
	if reqID != "" {
		c = context.WithValue(c, utils.CtxKeyRequestID, reqID)
	}

	// Pass id ไปที่ UseCase แทน struct
	resp, err := h.UseCase.CustInfo(c, ctx, id)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, resp)
}

func (h *DashboardHandler) GetCustProfile(ctx *gin.Context) {
	log.Println("[Handler] => Entering GetCustProfile")

	aeonId := ctx.Query("aeon_id")
	if aeonId == "" {
		lib.HandleError(ctx, fmt.Errorf("missing aeon_id parameter"))
		return
	}

	cfg, err := config.Load("sit")
	if err != nil {
		lib.HandleError(ctx, fmt.Errorf("internal server error"))
		return
	}

	c := ctx.Request.Context()
	c = context.WithValue(c, utils.CtxKeyApisKey, cfg.Headers.ApiKey)
	c = context.WithValue(c, utils.CtxKeyApiLang, cfg.Headers.ApiLanguage)
	c = context.WithValue(c, utils.CtxKeyDeviceOS, cfg.Headers.ApiDeviceOS)
	c = context.WithValue(c, utils.CtxKeyChannel, cfg.Headers.ApiChannel)

	reqID := ctx.GetHeader("X-Request-ID")
	if reqID != "" {
		c = context.WithValue(c, utils.CtxKeyRequestID, reqID)
	}

	resp, err := h.UseCase.CustProfile(c, ctx, aeonId)
	if err != nil {
		log.Println("[Handler] Error from UseCase:", err)
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, resp)
}

func (h *DashboardHandler) GetCustSegment(ctx *gin.Context) {
	log.Println("[Handler] => Entering GetCustSegment")

	aeonId := ctx.Query("aeon_id")
	if aeonId == "" {
		lib.HandleError(ctx, fmt.Errorf("missing aeon_id parameter"))
		return
	}

	cfg, err := config.Load("sit")
	if err != nil {
		lib.HandleError(ctx, fmt.Errorf("internal server error"))
		return
	}

	c := ctx.Request.Context()
	c = context.WithValue(c, utils.CtxKeyApisKey, cfg.Headers.ApiKey)
	c = context.WithValue(c, utils.CtxKeyApiLang, cfg.Headers.ApiLanguage)
	c = context.WithValue(c, utils.CtxKeyDeviceOS, cfg.Headers.ApiDeviceOS)
	c = context.WithValue(c, utils.CtxKeyChannel, cfg.Headers.ApiChannel)

	reqID := ctx.GetHeader("X-Request-ID")
	if reqID != "" {
		c = context.WithValue(c, utils.CtxKeyRequestID, reqID)
	}

	resp, err := h.UseCase.CustSegment(c, ctx, aeonId)
	if err != nil {
		log.Println("[Handler] Error from UseCase:", err)
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, resp)
}

func (h *DashboardHandler) GetCustSuggestion(ctx *gin.Context) {
	log.Println("[Handler] => Entering GetCustSuggestion")

	aeonId := ctx.Query("aeon_id")
	if aeonId == "" {
		lib.HandleError(ctx, fmt.Errorf("missing aeon_id parameter"))
		return
	}

	cfg, err := config.Load("sit")
	if err != nil {
		lib.HandleError(ctx, fmt.Errorf("internal server error"))
		return
	}

	c := ctx.Request.Context()
	c = context.WithValue(c, utils.CtxKeyApisKey, cfg.Headers.ApiKey)
	c = context.WithValue(c, utils.CtxKeyApiLang, cfg.Headers.ApiLanguage)
	c = context.WithValue(c, utils.CtxKeyDeviceOS, cfg.Headers.ApiDeviceOS)
	c = context.WithValue(c, utils.CtxKeyChannel, cfg.Headers.ApiChannel)

	reqID := ctx.GetHeader("X-Request-ID")
	if reqID != "" {
		c = context.WithValue(c, utils.CtxKeyRequestID, reqID)
	}

	resp, err := h.UseCase.CustSuggestion(c, ctx, aeonId)
	if err != nil {
		log.Println("[Handler] Error from UseCase:", err)
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, resp)
}
