package rest

import (
	"net/http"
	"github.com/buduchail/catrina"
)

type (
	// Base implementation of ResourceHandler interface, to be used by
	// RestAPI instances. This implementation can be used by concrete
	// resource handlers to provide default behaviour for HTTP verbs
	// that are not implemented by the handler. By embedding this type,
	// concrete handlers will only need to implement methods for verbs
	// they handle.
	ResourceHandler struct {
	}
)

func (s ResourceHandler) Options() (code int, body catrina.Payload, err error) {
	return http.StatusOK, catrina.EmptyBody, nil
}

func (s ResourceHandler) Post(parentIds []string, payload catrina.Payload) (code int, body catrina.Payload, err error) {
	return http.StatusMethodNotAllowed, catrina.EmptyBody, nil
}

func (s ResourceHandler) Get(id string, parentIds []string) (code int, body catrina.Payload, err error) {
	return http.StatusMethodNotAllowed, catrina.EmptyBody, nil
}

func (s ResourceHandler) GetMany(parentIds []string, params catrina.QueryParameters) (code int, body catrina.Payload, err error) {
	return http.StatusMethodNotAllowed, catrina.EmptyBody, nil
}

func (s ResourceHandler) Put(id string, parentIds []string, payload catrina.Payload) (code int, body catrina.Payload, err error) {
	return http.StatusMethodNotAllowed, catrina.EmptyBody, nil
}

func (s ResourceHandler) Delete(id string, parentIds []string) (code int, body catrina.Payload, err error) {
	return http.StatusMethodNotAllowed, catrina.EmptyBody, nil
}
