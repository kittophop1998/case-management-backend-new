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
		lib.HandleError(ctx, lib.InternalServer.WithDetails("Failed to create user: "+err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusCreated, uid)
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
	departmentIdStr := ctx.Query("departmentId")
	isActiveStr := ctx.Query("isActive")
	var isActive *bool = nil
	if isActiveStr != "" {
		val := isActiveStr == "true"
		isActive = &val
	}

	var roleID, sectionIdID, centerID, departmentID uuid.UUID
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
	if departmentIdStr != "" {
		if id, err := uuid.Parse(departmentIdStr); err == nil {
			departmentID = id
		}
	}

	filter := model.UserFilter{
		Keyword:      keyword,
		Sort:         sort,
		IsActive:     isActive,
		RoleID:       roleID,
		SectionID:    sectionIdID,
		CenterID:     centerID,
		DepartmentID: departmentID,
	}

	users, total, err := h.UseCase.GetAll(ctx, page, limit, filter)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails("Failed to fetch users: "+err.Error()))
		return
	}

	lib.HandlePaginatedResponse(ctx, page, limit, total, users)
}

func (h *UserHandler) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid user ID"))
		return
	}

	user, err := h.UseCase.GetById(ctx, userId)
	if err != nil {
		lib.HandleError(ctx, lib.NotFound.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, user)
}

func (h *UserHandler) UpdateUserByID(ctx *gin.Context) {
	var userUpdate model.CreateUpdateUserRequest
	if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	id := ctx.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid user ID"))
		return
	}

	err = h.UseCase.UpdateUserById(ctx, userId, userUpdate)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, gin.H{"message": "User updated successfully"})
}
