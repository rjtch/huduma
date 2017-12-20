package oauth2

import (
	"encoding/json"
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/notifier"
	"github.com/hellofresh/janus/pkg/opentracing"
	"github.com/hellofresh/janus/pkg/render"
	"github.com/hellofresh/janus/pkg/router"
)

// Controller is the api rest controller
type Controller struct {
	repo     Repository
	notifier notifier.Notifier
}

// NewController creates a new instance of Controller
func NewController(repo Repository, notifier notifier.Notifier) *Controller {
	return &Controller{repo, notifier}
}

// Get is the find all handler
func (c *Controller) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := opentracing.FromContext(r.Context(), "datastore.FindAll")
		data, err := c.repo.FindAll()
		span.Finish()

		if err != nil {
			errors.Handler(w, err)
			return
		}

		render.JSON(w, http.StatusOK, data)
	}
}

// GetBy is the find by handler
func (c *Controller) GetBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := router.URLParam(r, "name")
		span := opentracing.FromContext(r.Context(), "datastore.FindByName")
		data, err := c.repo.FindByName(name)
		span.Finish()

		if err != nil {
			errors.Handler(w, err)
			return
		}

		render.JSON(w, http.StatusOK, data)
	}
}

// PutBy is the update handler
func (c *Controller) PutBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		name := router.URLParam(r, "name")

		span := opentracing.FromContext(r.Context(), "datastore.FindByName")
		oauth, err := c.repo.FindByName(name)
		span.Finish()

		if oauth.Name == "" {
			errors.Handler(w, ErrOauthServerNotFound)
			return
		}

		if err != nil {
			errors.Handler(w, err)
			return
		}

		err = json.NewDecoder(r.Body).Decode(oauth)
		if nil != err {
			errors.Handler(w, err)
			return
		}

		span = opentracing.FromContext(r.Context(), "datastore.Add")
		err = c.repo.Add(oauth)
		c.dispatch(notifier.NoticeOAuthServerUpdated)
		span.Finish()

		if nil != err {
			errors.Handler(w, errors.New(http.StatusBadRequest, err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Post is the create handler
func (c *Controller) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var oauth OAuth

		err := json.NewDecoder(r.Body).Decode(&oauth)
		if nil != err {
			errors.Handler(w, err)
			return
		}

		span := opentracing.FromContext(r.Context(), "datastore.Add")
		err = c.repo.Add(&oauth)
		c.dispatch(notifier.NoticeOAuthServerAdded)
		span.Finish()

		if nil != err {
			errors.Handler(w, errors.New(http.StatusBadRequest, err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// DeleteBy is the delete handler
func (c *Controller) DeleteBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := router.URLParam(r, "name")

		span := opentracing.FromContext(r.Context(), "datastore.Remove")
		err := c.repo.Remove(name)
		span.Finish()

		if err != nil {
			errors.Handler(w, err)
			return
		}

		c.dispatch(notifier.NoticeOAuthServerRemoved)
		w.WriteHeader(http.StatusNoContent)
	}
}

func (c *Controller) dispatch(cmd notifier.NotificationCommand) {
	if c.notifier != nil {
		c.notifier.Notify(notifier.Notification{Command: cmd})
	}
}
