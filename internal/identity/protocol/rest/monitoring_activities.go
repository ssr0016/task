package rest

import (
	"task/internal/api/errors"
	"task/internal/api/response"
	"task/internal/identity/monitoringactivities"
	"task/internal/identity/monitoringactivities/logsmonitoring"

	"github.com/gofiber/fiber/v2"
)

type monitorinActivitiesHandler struct {
	s monitoringactivities.Service
	m logsmonitoring.Service
}

func NewMonitoringActivitiesHandler(s monitoringactivities.Service, m logsmonitoring.Service) *monitorinActivitiesHandler {
	return &monitorinActivitiesHandler{
		s: s,
		m: m,
	}
}

func (h *monitorinActivitiesHandler) GetMonitoringActivities(ctx *fiber.Ctx) error {
	var query monitoringactivities.SearchLogActivityQuery

	if err := ctx.QueryParser(&query); err != nil {
		return errors.ErrorBadRequest(err)
	}

	result, err := h.s.SearchLogActivities(ctx.Context(), &query)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"result": result,
	})
}

func (h *monitorinActivitiesHandler) MonitoringLogs(ctx *fiber.Ctx) error {
	var query logsmonitoring.MonitoringLogsQuery

	if err := ctx.QueryParser(&query); err != nil {
		return errors.ErrorBadRequest(err)
	}

	result, err := h.m.MonotoringLogs(ctx.Context(), &query)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"result": result,
	})
}
