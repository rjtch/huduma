package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/huduma/config"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/sirupsen/logrus"
)

//API is the api's definition
type API struct {
	Name    string `json:"name"`
	Active  bool   `json:"active"`
	Log     *logrus.Entry
	Config  *config.Config
	Handler http.Handler
	Port    int `json:"port"`
	//	db     *gorm.DB
	Version string `json:"version"`
}

//NewAPI is the instance of our API
func NewAPI(version string, conf *config.Config) *API {

	api := &API{
		Active:  true,
		Version: version,
		Port:    conf.Port,
		Config:  conf,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Get("/", api.Info)
	return api
}

func sendJSON(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Set("Content-Type", "appliction/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.Encode(obj)
}

//Serve is used to listen api's port
func (api *API) Serve() error {
	p := fmt.Sprintf(":%d", api.Port)
	api.Log.Infof("API started on: %s", p)
	return http.ListenAndServe(p, api.Handler)
}
