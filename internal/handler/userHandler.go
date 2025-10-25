package handler

import (
	"hackaton/internal/dto/request"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) Login(c echo.Context) error {
	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.service.User().Login(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateUser(c echo.Context) error {
	var req request.UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.service.User().CreateUser(c.Request().Context(), &req); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "user created"})
}

func (h *Handler) GetUsers(c echo.Context) error {
	users, err := h.service.User().GetUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUserByID(c echo.Context) error {
	id := c.Param("id")

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	resp, err := h.service.User().GetUserByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	id := c.Param("id")

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}
	var req request.UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.service.User().UpdateUser(c.Request().Context(), id, &req); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user updated"})
}

func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	if err := h.service.User().DeleteUser(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user deleted"})
}
