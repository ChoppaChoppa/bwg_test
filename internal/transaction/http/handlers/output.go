package handlers

import (
	"bwg_test/internal/transaction/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) Output(c echo.Context) error {
	var transaction models.Transaction
	if err := c.Bind(&transaction); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Error:     true,
			ErrorText: err.Error(),
			Code:      http.StatusBadRequest,
		})
	}

	if err := h.service.Output(c.Request().Context(), &transaction); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Error:     true,
			ErrorText: err.Error(),
			Code:      http.StatusBadRequest,
		})
	}

	return c.JSON(http.StatusBadRequest, models.Response{
		Data: "транзакция принята в обработку",
		Code: http.StatusBadRequest,
	})
}
