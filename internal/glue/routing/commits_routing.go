package routing

import (
	"net/http"

	h "go-github-tracker/internal/handlers"
	"go-github-tracker/platforms/routers"
)

func CommitsRouting(handler *h.RepositoriesHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodGet,
			Path:        "/commits/{repositoryName}",
			Handle:      handler.GetAllRepositories,
			MiddleWares: []http.HandlerFunc{},
		},
	}
}
