package routing

import (
	"net/http"

	h "go-github-tracker/internal/handlers"
	"go-github-tracker/platforms/routers"
)

func CommitsRouting(handler *h.CommitsHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodGet,
			Path:        "/commits/{repositoryName}",
			Handle:      handler.GetAllcommits,
			MiddleWares: []http.HandlerFunc{},
		},
	}
}
