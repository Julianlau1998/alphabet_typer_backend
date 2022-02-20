package records

import (
	"alphabet_typer/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Delivery struct {
	recordService Service
}

func NewDelivery(recordService Service) Delivery {
	return Delivery{recordService: recordService}
}

func (d *Delivery) GetAll(c echo.Context) error {
	limit, _ := strconv.ParseInt(c.QueryParam("limit"), 0, 32)
	filter, _ := strconv.ParseInt(c.QueryParam("filter"), 0, 32)
	records, err := d.recordService.GetAll(limit, filter)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

func (d *Delivery) Post(c echo.Context) error {
	requestBody := new(models.Record)
	if err := c.Bind(requestBody); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	record, err := d.recordService.Post(requestBody)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, record.UUID)
}
