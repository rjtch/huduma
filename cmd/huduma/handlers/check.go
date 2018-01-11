package handlers

import (
	"context"
	"net/http"

	"github.com/huduma/internal/mongo"
	"github.com/huduma/internal/router"
	"github.com/pkg/errors"
)

//HealthCkeck is the book API method handler set
type HealthCkeck struct {
	db *mongo.BooksDB
}

//RunCheck uses to check if all config are loaded and if the system is ready to accept connection
func (h *HealthCkeck) RunCheck(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	dbcon, err := h.db.Copy()
	if err != nil {
		return errors.Wrap(router.ErrDBNotConfigured, "")
	}
	defer dbcon.Close()

	if err := checker(); err != nil {
		return err
	}

	router.Response(ctx, w, []byte("Status ok"), http.StatusOK)
	return nil
}

func checker() error {
	return nil
}
