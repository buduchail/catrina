package rest

import (
	"errors"
	"strconv"
	"net/http"
	"github.com/valyala/fasthttp"
	"github.com/buduchail/catrina"
)

type (
	FastAPI struct {
		root      *pathHandler
		prefix    string
		prefixLen int
	}
)

func NewFast(prefix string) (api FastAPI) {
	api = FastAPI{}
	api.prefix = normalizePrefix(prefix)
	api.prefixLen = len(api.prefix)
	api.root = NewPathHandler(api.prefix)
	return api
}

func (api FastAPI) getBody(ctx *fasthttp.RequestCtx) catrina.Payload {
	return ctx.Request.Body()
}

func (api FastAPI) getQueryParameters(ctx *fasthttp.RequestCtx) catrina.QueryParameters {
	params := catrina.QueryParameters{}
	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		params[string(key)] = []string{string(value)}
	})
	return params
}

func (api FastAPI) sendResponse(ctx *fasthttp.RequestCtx, code int, body catrina.Payload, err error) error {

	if code == http.StatusOK {
		_, err = ctx.Write(body)
	} else {
		if err == nil {
			err = getHttpError(code)
		}
		ctx.Error(err.Error(), code)
	}

	return err
}

func (api FastAPI) handleResource(method string, id string, parentIds []string, ctx *fasthttp.RequestCtx, handler catrina.ResourceHandler) (code int, body catrina.Payload, err error) {

	switch method {
	case "POST":
		if id != "" {
			return http.StatusBadRequest, catrina.EmptyBody, errors.New("POST requests must not provide an ID")
		}
		return handler.Post(parentIds, api.getBody(ctx))
	case "GET":
		if id != "" {
			return handler.Get(id, parentIds)
		} else {
			return handler.GetMany(parentIds, api.getQueryParameters(ctx))
		}
	case "PUT":
		if id == "" {
			return http.StatusBadRequest, catrina.EmptyBody, errors.New("PUT method must provide an ID")
		}
		return handler.Put(id, parentIds, api.getBody(ctx))
	case "DELETE":
		if id == "" {
			return http.StatusBadRequest, catrina.EmptyBody, errors.New("DELETE method must provide an ID")
		}
		return handler.Delete(id, parentIds)
	}

	return http.StatusMethodNotAllowed, catrina.EmptyBody, errors.New("Method not allowed")
}

func (api FastAPI) handle(ctx *fasthttp.RequestCtx) {

	path := string(ctx.Request.URI().Path())

	if len(path) > api.prefixLen+1 {

		handler, id, parentIds := api.root.findHandler(path[api.prefixLen:])
		if handler == nil {
			api.sendResponse(ctx, http.StatusNotFound, catrina.EmptyBody, nil)
			return
		}

		code, body, err := api.handleResource(string(ctx.Method()), id, parentIds, ctx, handler)
		api.sendResponse(ctx, code, body, err)

	} else {
		api.sendResponse(ctx, http.StatusNotFound, catrina.EmptyBody, nil)
	}
}

func (api FastAPI) AddResource(name string, handler catrina.ResourceHandler) {
	api.root.addHandler(name, handler)
}

func (api FastAPI) AddMiddleware(m catrina.Middleware) {
	// NOT IMPLEMENTED
}

func (api FastAPI) Run(port int) {

	fasthttp.ListenAndServe(":"+strconv.Itoa(port), api.handle)
}
