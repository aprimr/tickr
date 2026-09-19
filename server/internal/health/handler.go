package health

import (
	"context"
	"net/http"
	"time"

	"github.com/aprimr/tickr/internal/pkg/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler interface {
	HandleHealth(w http.ResponseWriter, r *http.Request)
}

type healthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) HealthHandler {
	return &healthHandler{
		db: db,
	}
}

// HandleHealth performs the database ping and returns system health status
func (h *healthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		res := map[string]string{"health": "error", "database": "unreachable"}
		response.Error(w, http.StatusServiceUnavailable, "database ping failed", res)
		return
	}

	res := map[string]string{"health": "ok", "database": "connected"}
	response.JSON(w, http.StatusOK, "database ping successfull", res)
}
