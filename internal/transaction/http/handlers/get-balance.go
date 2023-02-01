package handlers

import (
	"bwg_test/internal/transaction/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *Handler) GetBalance(c echo.Context) error {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Error:     true,
			ErrorText: err.Error(),
			Data:      id,
			Code:      http.StatusBadRequest,
		})
	}

	transactions, err := h.service.GetBalance(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Error:     true,
			ErrorText: err.Error(),
			Data:      id,
			Code:      http.StatusBadRequest,
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Data: transactions,
		Code: http.StatusOK,
	})
}
