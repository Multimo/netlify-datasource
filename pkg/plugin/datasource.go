package plugin

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"

	"github.com/grafana/netlify-datasource/pkg/plugin/client"
	"github.com/grafana/netlify-datasource/pkg/plugin/models"
	"github.com/grafana/netlify-datasource/pkg/plugin/query"
	"github.com/grafana/netlify-datasource/pkg/plugin/resources"
)

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct {
	client       client.Client
	queryHandler query.QueryHandler
	backend.CallResourceHandler
}

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
//
// type assertion to validate that the datasource implements the interface
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
	// _ backend.CallResourceHandler   = (*resources.ResourcesHandler.Router)(nil)
)

// NewDatasource creates a new datasource instance.
func NewDatasource(ctx context.Context, config backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	settings, err := models.LoadSettings(ctx, config)
	if err != nil {
		return nil, err
	}

	client := client.NewClient(settings)
	query := query.NewQueryHandler(client)

	resourcesHandler := resources.NewResourcesHandler(client)

	ds := Datasource{
		client:              client,
		queryHandler:        query,
		CallResourceHandler: httpadapter.New(resourcesHandler.Router),
	}

	return &ds, nil
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return d.queryHandler.HandleQueries(ctx, req)
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Data source is working"

	// TODO:
	// validate settings on health check

	_, err := d.client.GetSites()
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
