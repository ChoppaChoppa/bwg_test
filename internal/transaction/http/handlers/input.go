package handlers

import (
	"bwg_test/internal/transaction/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) Input(c echo.Context) error {
	var transaction models.Transaction
	if err := c.Bind(&transaction); err != nil {
		return err
	}

	err := h.service.Input(c.Request().Context(), &transaction)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
