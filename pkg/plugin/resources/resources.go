package resources

import (
	"encoding/json"
	"net/http"

	"github.com/grafana/netlify-datasource/pkg/plugin/client"
)

type ResourceHandler struct {
	client client.Client
	Router http.Handler
}

func getRoutes(h *ResourceHandler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/sites", h.HandleGetSites)

	return router
}

func NewResourcesHandler(client client.Client) ResourceHandler {
	r := ResourceHandler{
		client: client,
	}

	r.Router = getRoutes(&r)

	return r
}

func (h *ResourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

func (h *ResourceHandler) HandleGetSites(w http.ResponseWriter, r *http.Request) {
	sites, err := h.client.GetSites()
	if err != nil {
		// handle error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	res := make([]string, len(sites))

	for i, site := range sites {
		res[i] = site.ID
	}

	response, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
