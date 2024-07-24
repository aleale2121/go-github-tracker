package routing

import (
	"net/http"

	h "commits-manager-service/internal/handlers"
	"commits-manager-service/platforms/routers"
)

func RepositoriesRouting(handler *h.RepositoriesHandler) []routers.Route {
	return []routers.Route{
		{
			Method:      http.MethodGet,
			Path:        "/repositories",
			Handle:      handler.GetAllRepositories,
			MiddleWares: []http.HandlerFunc{},
		},
	}
}
