package v1

//go:generate mockgen -source=router.go -destination=mock/mock_router.go -package=mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type LogsetService interface {
	Get(ctx context.Context, logsetID uuid.UUID) ([]map[string]any, error)
	Ingest(ctx context.Context, logs []map[string]any) (uuid.UUID, error)
}

type Router struct {
	logsetService LogsetService
}

func NewRouter(logsetService LogsetService) *Router {
	return &Router{logsetService: logsetService}
}

func (r *Router) Register(group *echo.Group) {
	group.POST("/logsets", r.createLogset)
	group.GET("/logsets/:uuid", r.getLogset)
}
