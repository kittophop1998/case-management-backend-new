package http

import (
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	UseCase usecase.UserUseCase
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	user := &model.CreateUpdateUserRequest{}
	if err := ctx.ShouldBindJSON(user); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	uid, err := h.UseCase.Create(ctx, user)
	if err != nil {
		lib.HandleError(ctx, lib.NewAppError(http.StatusInternalServerError, lib.MessageError{
			Th: "ไม่สามารถสร้างผู้ใช้ได้",
			En: "Failed to create user",
		}, err.Error()))
		return
	}

	lib.NewResponse(ctx, http.StatusCreated, gin.H{"id": uid})
}

func (h *UserHandler) GetAllUsers(ctx *gin.Context) {
	limit, err := getLimit(ctx)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	page, err := getPage(ctx)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	fmt.Printf("Fetching users with pagination: page=%d, limit=%d\n", page, limit)

	sort := ctx.DefaultQuery("sort", "is_active desc")
	keyword := ctx.Query("keyword")
	roleIdStr := ctx.Query("roleId")
	sectionIdStr := ctx.Query("sectionId")
	centerIdStr := ctx.Query("centerId")
	isActiveStr := ctx.Query("isActive")
	var isActive *bool = nil
	if isActiveStr != "" {
		val := isActiveStr == "true"
		isActive = &val
	}

	var roleID, sectionIdID, centerID uuid.UUID
	if roleIdStr != "" {
		if id, err := uuid.Parse(roleIdStr); err == nil {
			roleID = id
		}
	}
	if sectionIdStr != "" {
		if id, err := uuid.Parse(sectionIdStr); err == nil {
			sectionIdID = id
		}
	}
	if centerIdStr != "" {
		if id, err := uuid.Parse(centerIdStr); err == nil {
			centerID = id
		}
	}

	filter := model.UserFilter{
		Keyword:   keyword,
		Sort:      sort,
		IsActive:  isActive,
		RoleID:    roleID,
		SectionID: sectionIdID,
		CenterID:  centerID,
	}

	users, total, err := h.UseCase.GetAll(ctx, page, limit, filter)
	if err != nil {
		lib.HandleError(ctx, lib.NewAppError(http.StatusInternalServerError, lib.MessageError{
			Th: "ไม่สามารถดึงข้อมูลผู้ใช้ได้",
			En: "Failed to retrieve users",
		}, err.Error()))
		return
	}

	lib.HandlePaginatedResponse(ctx, users, page, limit, total)
}

func (h *UserHandler) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.UseCase.GetById(ctx, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	lib.NewResponse(ctx, http.StatusOK, user)
}

func (h *UserHandler) UpdateUserByID(ctx *gin.Context) {
	var input model.CreateUpdateUserRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	id := ctx.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.UseCase.UpdateUserById(ctx, userId, input)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Can't update user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}
