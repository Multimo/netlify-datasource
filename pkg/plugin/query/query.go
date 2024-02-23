package query

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/framestruct"

	"github.com/grafana/netlify-datasource/pkg/plugin/client"
)

type QueryHandler struct {
	client client.Client
}

func NewQueryHandler(client client.Client) QueryHandler {
	return QueryHandler{
		client: client,
	}
}

func (q QueryHandler) HandleQueries(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, query := range req.Queries {
		res := q.Query(ctx, req.PluginContext, query)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[query.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	Entity         string `json:"entity"` // builds, deployments
	SiteId         string `json:"siteId"` // uuid
	ParsingOptions struct {
		SelectedFields []string `json:"selectedFields"`
	} `json:"parsingOptions"`
}

func parseSiteIdsAsVariables(siteIds string) ([]string, error) {
	siteIdSlice := make([]string, 0)
	// expecting string as "siteId" or "{siteId, siteId, siteId}"
	if strings.HasPrefix(siteIds, "{") && strings.HasSuffix(siteIds, "}") {
		siteIds = strings.TrimPrefix(siteIds, "{")
		siteIds = strings.TrimSuffix(siteIds, "}")
		siteIdSlice = strings.Split(siteIds, ",")
		return siteIdSlice, nil
	}

	siteIdSlice = append(siteIdSlice, siteIds)

	return siteIdSlice, nil
}

func (q QueryHandler) Query(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	// Unmarshal the JSON into our queryModel.
	var qm queryModel
	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal failed on query: %v", err.Error()))
	}

	backend.Logger.Info("queryParams", "SelectedFields", qm.ParsingOptions.SelectedFields)

	selectedFields := qm.ParsingOptions.SelectedFields

	sitesIds, err := parseSiteIdsAsVariables(qm.SiteId)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed on parsing siteIds: %v", err.Error()))
	}

	backend.Logger.Info("query", "entity", qm.Entity, "siteId", qm.SiteId, "sitesIds", sitesIds)

	switch qm.Entity {
	case "builds":
		return q.HandleBuildsQuery(ctx, sitesIds, selectedFields)
	case "deployments":
		return q.HandleDeploymentsQuery(ctx, sitesIds, selectedFields)
	case "forms":
		return q.HandleFormsQuery(ctx, sitesIds)
	case "form-submissions":
		return q.HandleFormSubmissionsQuery(ctx, sitesIds)
	case "builds-account":
		return q.HandleBuildAccountDetails(ctx)
	case "sites":
		return q.HandleSitesQuery(ctx)
	case "accounts":
		return q.HandleAccounts(ctx)
	case "":
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("missing query param entity: %v", err.Error()))
	default:
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("Unidentified query param entity: %v", err.Error()))
	}
}

func contains(slice []string, item string) bool {
	for _, value := range slice {
		if value == item {
			return true
		}
	}
	return false
}

func (q QueryHandler) HandleBuildsQuery(ctx context.Context, siteIds []string, selectedFields []string) backend.DataResponse {
	var response backend.DataResponse

	res, errors := client.DoGets[client.BuildsResponse](ctx, q.client.GetBuilds, siteIds)
	if len(errors) > 0 {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to get builds: %v", errors[0].Error()))
	}

	builds := make(client.BuildsResponse, 0, len(res))
	for _, r := range res {
		builds = append(builds, r...)
	}

	// ID Sha CreatedAt State
	// if len(selectedFields) > 0 {
	// 	for _, build := range builds {
	// 		buildValue := reflect.ValueOf(build).Elem()

	// 		// for _, v := range selectedFields {
	// 		// 	fieldValue := buildValue.FieldByName(v)
	// 		// 	if !fieldValue.IsValid() {
	// 		// 		continue
	// 		// 	}

	// 		// 	fieldValue.Set(reflect.Zero(fieldValue.Type()))
	// 		// }

	// 		for i := 0; i < buildValue.NumField(); i++ {
	// 			fieldName := buildValue.Type().Field(i)

	// 			// Check if the field should be kept
	// 			if !contains(selectedFields, fieldName.Name) {
	// 				// Reset the field's value to its zero value
	// 				fieldValue := buildValue.Field(i)
	// 				fieldValue.Set(reflect.Zero(fieldValue.Type()))
	// 			}
	// 		}
	// 	}
	// }

	backend.Logger.Info("HandleBuildsQuery", "len", len(res), "builds", builds)

	dataFrames, err := framestruct.ToDataFrame("builds", builds)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed Builds to frame conversion: %v", err.Error()))
	}

	// dataFrames.

	response.Frames = append(response.Frames, dataFrames)

	return response
}

func (q QueryHandler) HandleDeploymentsQuery(ctx context.Context, siteIds []string, selectedFields []string) backend.DataResponse {
	var response backend.DataResponse

	res, errors := client.DoGets[client.DeploysResponse](ctx, q.client.GetDeployments, siteIds)
	if len(errors) > 0 {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to get deployments: %v", errors[0].Error()))
	}

	deployments := make(client.DeploysResponse, 0, len(res))
	for _, r := range res {
		deployments = append(deployments, r...)
	}

	dataFrames, err := framestruct.ToDataFrame("deployments", deployments)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed deployments to frame conversion: %v", err.Error()))
	}

	response.Frames = append(response.Frames, dataFrames)

	return response
}

func (q QueryHandler) HandleSitesQuery(ctx context.Context) backend.DataResponse {
	var response backend.DataResponse

	res, err := q.client.GetSites()
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to get deploys: %v", err.Error()))
	}

	dataFrames, err := framestruct.ToDataFrame("sites", res)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed Sites to frame conversion: %v", err.Error()))
	}

	// add the frames to the response.
	response.Frames = append(response.Frames, dataFrames)

	return response
}

func (q QueryHandler) HandleFormsQuery(ctx context.Context, siteIds []string) backend.DataResponse {
	var response = backend.DataResponse{}

	res, errors := client.DoGets[client.FormsResponse](ctx, q.client.GetForms, siteIds)
	if len(errors) > 0 {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to get forms: %v", errors[0].Error()))
	}

	forms := make(client.FormsResponse, 0, len(res))
	for _, r := range res {
		forms = append(forms, r...)
	}

	dataFrames, err := framestruct.ToDataFrame("forms", forms)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed forms to frame conversion: %v", err.Error()))
	}

	// add the frames to the response.
	response.Frames = append(response.Frames, dataFrames)

	return response
}

func (q QueryHandler) HandleFormSubmissionsQuery(ctx context.Context, siteIds []string) backend.DataResponse {
	var response = backend.DataResponse{}

	res, errors := client.DoGets[client.FormSubmissionsResponse](ctx, q.client.GetFormSubmittions, siteIds)
	if len(errors) > 0 {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to get forms submissions: %v", errors[0].Error()))
	}

	form_submissions := make(client.FormSubmissionsResponse, 0, len(res))
	for _, r := range res {
		form_submissions = append(form_submissions, r...)
	}

	// create data frame response.
	// For an overview on data frames and how grafana handles them:
	// https://grafana.com/developers/plugin-tools/introduction/data-frames
	dataFrames, err := framestruct.ToDataFrame("form_submissions", form_submissions)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed forms submissions to frame conversion: %v", err.Error()))
	}

	response.Frames = append(response.Frames, dataFrames)

	return response
}

func (q QueryHandler) HandleBuildAccountDetails(ctx context.Context) backend.DataResponse {
	var response backend.DataResponse

	res, err := q.client.GetBuildAccountDetails()
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to get build account details: %v", err.Error()))
	}

	dataFrames, err := framestruct.ToDataFrame("build_account_details", res)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed Build Account to frame conversion: %v", err.Error()))
	}

	response.Frames = append(response.Frames, dataFrames)

	return response
}

func (q QueryHandler) HandleAccounts(ctx context.Context) backend.DataResponse {
	var response backend.DataResponse

	res, err := q.client.GetAccounts()
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to get accounts: %v", err.Error()))
	}

	dataFrames, err := framestruct.ToDataFrame("accounts", res)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed Build Account to frame conversion: %v", err.Error()))
	}

	response.Frames = append(response.Frames, dataFrames)

	return response
}
