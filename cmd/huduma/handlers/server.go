package handlers

import (
	"net/http"

	"github.com/huduma/internal/middleware"
	"github.com/huduma/internal/mongo"

	"github.com/huduma/internal/router"
)

//API is the huduma server a set of routes
func API(db *mongo.BooksDB) http.Handler {
	huduma := router.NewHuduma(middleware.ErrorHandler, middleware.Logger)

	b := Books{
		db: db,
	}

	huduma.HandleWrapper("GET", "/v1/books", b.Read)
	huduma.HandleWrapper("POST", "/v1/books", b.Create)
	huduma.HandleWrapper("GET", "/v1/books", b.retreive)
	huduma.HandleWrapper("PUT", "/v1/books", b.Update)
	huduma.HandleWrapper("DELETE", "/v1/books", b.Delete)

	h := HealthCkeck{
		db: db,
	}
	huduma.HandleWrapper("GET", "/v1/healthcheck", h.RunCheck)

	return huduma
}
