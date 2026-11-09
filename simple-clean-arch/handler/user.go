package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"simple-http/model"
	"simple-http/repository"
	"simple-http/usecase"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: userUsecase}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users := h.usecase.GetUsers()
	sendResponse(c, http.StatusOK, "00", "success", users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := getUserID(c)
	if err != nil {
		sendResponse(c, http.StatusBadRequest, "01", err.Error(), nil)
		return
	}

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		sendError(c, err)
		return
	}

	sendResponse(c, http.StatusOK, "00", "success", user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var request model.User
	if err := c.ShouldBindJSON(&request); err != nil {
		sendResponse(c, http.StatusBadRequest, "01", "invalid request body", nil)
		return
	}

	if err := validateRequest(request); err != nil {
		sendResponse(c, http.StatusBadRequest, "01", err.Error(), nil)
		return
	}

	user, err := h.usecase.CreateUser(request)
	if err != nil {
		sendError(c, err)
		return
	}

	sendResponse(c, http.StatusCreated, "00", "user created successfully", user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := getUserID(c)
	if err != nil {
		sendResponse(c, http.StatusBadRequest, "01", err.Error(), nil)
		return
	}

	var request model.UserUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		sendResponse(c, http.StatusBadRequest, "01", "invalid request body", nil)
		return
	}

	if err := validateRequest(request); err != nil {
		sendResponse(c, http.StatusBadRequest, "01", err.Error(), nil)
		return
	}

	user, err := h.usecase.UpdateUser(id, request)
	if err != nil {
		sendError(c, err)
		return
	}

	sendResponse(c, http.StatusOK, "00", "user updated successfully", user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := getUserID(c)
	if err != nil {
		sendResponse(c, http.StatusBadRequest, "01", err.Error(), nil)
		return
	}

	user, err := h.usecase.DeleteUser(id)
	if err != nil {
		sendError(c, err)
		return
	}

	sendResponse(c, http.StatusOK, "00", "user deleted successfully", user)
}

func getUserID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid user id")
	}
	return id, nil
}

func sendResponse(c *gin.Context, status int, code, message string, data any) {
	c.JSON(status, model.Response{
		RespCode:    code,
		RespMessage: message,
		RespData:    data,
	})
}

func sendError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		sendResponse(c, http.StatusNotFound, "02", err.Error(), nil)
	case errors.Is(err, repository.ErrDuplicateID):
		sendResponse(c, http.StatusConflict, "01", err.Error(), nil)
	default:
		sendResponse(c, http.StatusBadRequest, "01", err.Error(), nil)
	}
}
