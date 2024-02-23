package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/netlify-datasource/pkg/plugin/models"
)

type Client struct {
	models.Settings
	client *http.Client
}

func NewClient(settings models.Settings) Client {
	c := Client{}

	c.client = &http.Client{}
	c.BaseUrl = settings.BaseUrl
	c.AccessToken = settings.AccessToken
	c.AccountId = settings.AccountId
	c.SiteId = settings.SiteId

	return c
}

type ErrorResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (c Client) doGet(url string, response any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+c.AccessToken)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			err = fmt.Errorf("error reading error body code: %d response: %s", res.StatusCode, err.Error())
			return err
		}

		err = fmt.Errorf("error: code: %d, response: %s", res.StatusCode, string(b))
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		backend.Logger.Info("Unmarshal error", "err", string(err.Error()))

		return err
	}

	return nil
}

type Doer[T any] func(s string) (T, error)

func httpGetter[T any](ctx context.Context, doer Doer[T], variable string, ch chan T, error_ch chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	backend.Logger.Info("httpGetter", "doing request")
	res, err := doer(variable)
	if err != nil {
		backend.Logger.Info("httpGetter", "found error", "err", err.Error())
		error_ch <- err
		return
	}

	backend.Logger.Info("httpGetter", "got result", 0, "res", res)
	ch <- res
}

func DoGets[T any](ctx context.Context, doer Doer[T], variables []string) ([]T, []error) {
	backend.Logger.Info("DoGets", "variables", variables, "len", len(variables))

	var wg sync.WaitGroup
	ch := make(chan T, len(variables))
	error_ch := make(chan error, len(variables))

	for _, variable := range variables {
		wg.Add(1)
		go httpGetter[T](ctx, doer, variable, ch, error_ch, &wg)
	}

	backend.Logger.Info("DoGets", "waiting for results")
	go func() {
		backend.Logger.Info("DoGets", "waiting on waitgroup")
		wg.Wait()
		backend.Logger.Info("DoGets", "done waiting closing channel")
		close(ch)
		close(error_ch)
	}()

	backend.Logger.Info("DoGets", "receiving results")

	results := make([]T, 0, len(variables))
	for r := range ch {
		backend.Logger.Info("DoGets", "looping over results", "r", r, "len", len(results))
		results = append(results, r)
	}

	errors := make([]error, 0, len(variables))
	for err := range error_ch {
		errors = append(errors, err)
	}

	backend.Logger.Info("DoGets", "returning results", "results", results, "errors", errors)
	return results, errors
}

func (c Client) buildUrl(pattern string, siteId string) string {
	sid := ""
	if siteId != "" {
		sid = siteId
	} else {
		sid = c.SiteId
	}

	pattern = strings.Replace(pattern, "{site_id}", sid, -1)

	return pattern
}

type DeploysResponse []struct {
	ID           string    `json:"id"`
	Build_id     string    `json:"build_id"`
	State        string    `json:"state"` // ready, error, retrying
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	PublishedAt  time.Time `json:"published_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	DeployTime   int64     `json:"deploy_time"`
	ManualDeploy bool      `json:"manual_deploy"`
	ErrorMessage string    `json:"error_message"`
	Branch       string    `json:"branch"`
	Context      string    `json:"context"`
}

func (c Client) GetDeployments(siteId string) (DeploysResponse, error) {
	deploys := DeploysResponse{}
	url := c.buildUrl("https://api.netlify.com/api/v1/sites/{site_id}/deploys", siteId)

	err := c.doGet(url, &deploys)
	if err != nil {
		return deploys, err
	}

	return deploys, nil
}

type BuildsResponse []struct {
	ID        string    `json:"id"`
	DeployID  string    `json:"deploy_id"`
	Sha       string    `json:"sha"`
	Done      bool      `json:"done"`
	Error     string    `json:"error"`
	CreatedAt time.Time `json:"created_at"`
}

type MapResponse = []map[string]any

// state
// "new" "pending_review" "accepted" "rejected" "enqueued" "building" "uploading" "uploaded" "preparing" "prepared" "processing" "processed" "ready" "error" "retrying"
func (c Client) GetBuilds(siteId string) (BuildsResponse, error) {
	backend.Logger.Info("GetBuilds", "siteId", siteId)

	builds := BuildsResponse{}
	url := c.buildUrl("https://api.netlify.com/api/v1/sites/{site_id}/builds", siteId)

	err := c.doGet(url, &builds)
	if err != nil {
		return builds, err
	}

	return builds, nil
}

type SitesResponse []struct {
	ID                        string    `json:"id"`
	State                     string    `json:"state"`
	Plan                      string    `json:"plan"`
	Name                      string    `json:"name"`
	CustomDomain              string    `json:"custom_domain"`
	DomainAliases             []string  `json:"domain_aliases"`
	BranchDeployCustomDomain  string    `json:"branch_deploy_custom_domain"`
	DeployPreviewCustomDomain string    `json:"deploy_preview_custom_domain"`
	Password                  string    `json:"password"`
	NotificationEmail         string    `json:"notification_email"`
	URL                       string    `json:"url"`
	SslURL                    string    `json:"ssl_url"`
	AdminURL                  string    `json:"admin_url"`
	ScreenshotURL             string    `json:"screenshot_url"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
	UserID                    string    `json:"user_id"`
	SessionID                 string    `json:"session_id"`
	Ssl                       bool      `json:"ssl"`
	ForceSsl                  bool      `json:"force_ssl"`
	ManagedDNS                bool      `json:"managed_dns"`
	DeployURL                 string    `json:"deploy_url"`
	PublishedDeploy           struct {
		ID            string `json:"id"`
		SiteID        string `json:"site_id"`
		UserID        string `json:"user_id"`
		BuildID       string `json:"build_id"`
		State         string `json:"state"`
		Name          string `json:"name"`
		URL           string `json:"url"`
		SslURL        string `json:"ssl_url"`
		AdminURL      string `json:"admin_url"`
		DeployURL     string `json:"deploy_url"`
		DeploySslURL  string `json:"deploy_ssl_url"`
		ScreenshotURL string `json:"screenshot_url"`
		ReviewID      int64  `json:"review_id"`
		Draft         bool   `json:"draft"`
		// Required          []string `json:"required"`
		// RequiredFunctions []string `json:"required_functions"`
		ErrorMessage string    `json:"error_message"`
		Branch       string    `json:"branch"`
		CommitRef    string    `json:"commit_ref"`
		CommitURL    string    `json:"commit_url"`
		Skipped      bool      `json:"skipped"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		PublishedAt  time.Time `json:"published_at"`
		Title        string    `json:"title"`
		Context      string    `json:"context"`
		Locked       bool      `json:"locked"`
		ReviewURL    string    `json:"review_url"`
		Framework    string    `json:"framework"`
		// FunctionSchedules []struct {
		// 	Name string `json:"name"`
		// 	Cron string `json:"cron"`
		// } `json:"function_schedules"`
	} `json:"published_deploy"`
	AccountName  string `json:"account_name"`
	AccountSlug  string `json:"account_slug"`
	GitProvider  string `json:"git_provider"`
	DeployHook   string `json:"deploy_hook"`
	Capabilities struct {
		Property1 struct {
		} `json:"property1"`
		Property2 struct {
		} `json:"property2"`
	} `json:"capabilities"`
	ProcessingSettings struct {
		HTML struct {
			PrettyUrls bool `json:"pretty_urls"`
		} `json:"html"`
	} `json:"processing_settings"`
	BuildSettings struct {
		ID           int64  `json:"id"`
		Provider     string `json:"provider"`
		DeployKeyID  string `json:"deploy_key_id"`
		RepoPath     string `json:"repo_path"`
		RepoBranch   string `json:"repo_branch"`
		Dir          string `json:"dir"`
		FunctionsDir string `json:"functions_dir"`
		Cmd          string `json:"cmd"`
		// AllowedBranches []string `json:"allowed_branches"`
		PublicRepo  bool   `json:"public_repo"`
		PrivateLogs bool   `json:"private_logs"`
		RepoURL     string `json:"repo_url"`
		Env         struct {
			Property1 string `json:"property1"`
			Property2 string `json:"property2"`
		} `json:"env"`
		InstallationID int64 `json:"installation_id"`
		StopBuilds     bool  `json:"stop_builds"`
	} `json:"build_settings"`
	IDDomain         string `json:"id_domain"`
	DefaultHooksData struct {
		AccessToken string `json:"access_token"`
	} `json:"default_hooks_data"`
	BuildImage      string `json:"build_image"`
	Prerender       string `json:"prerender"`
	FunctionsRegion string `json:"functions_region"`
}

func (c Client) GetSites() (SitesResponse, error) {
	sites := SitesResponse{}

	err := c.doGet("https://api.netlify.com/api/v1/sites", &sites)
	if err != nil {
		return sites, err
	}

	return sites, nil
}

type FormsResponse []struct {
	ID              string    `json:"id"`
	SiteId          string    `json:"site_id"`
	Name            string    `json:"name"`
	Paths           []string  `json:"paths"`
	SubmissionCount int64     `json:"submission_count"`
	CreatedAt       time.Time `json:"created_at"`
}

func (c Client) GetForms(siteId string) (FormsResponse, error) {
	forms := FormsResponse{}
	url := c.buildUrl("https://api.netlify.com/api/v1/sites/{site_id}/forms", siteId)

	err := c.doGet(url, &forms)
	if err != nil {
		return forms, err
	}

	return forms, nil
}

type FormSubmissionsResponse []struct {
	Id        string            `json:"id"`
	Number    int64             `json:"number"`
	Email     string            `json:"email"`
	Name      string            `json:"name"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Company   string            `json:"company"`
	Summary   string            `json:"summary"`
	Body      string            `json:"body"`
	Data      map[string]string `json:"data"`
	CreatedAt time.Time         `json:"created_at"`
	SiteUrl   string            `json:"site_url"`
}

func (c Client) GetFormSubmittions(siteId string) (FormSubmissionsResponse, error) {
	submissions := FormSubmissionsResponse{}
	url := c.buildUrl("https://api.netlify.com/api/v1/sites/{site_id}/submissions", siteId)

	err := c.doGet(url, &submissions)
	if err != nil {
		return submissions, err
	}

	return submissions, nil
}

type BuildAccountResponse struct {
	Active             int64 `json:"active"`
	PendingConcurrency int64 `json:"pending_concurrency"`
	Enqueued           int64 `json:"enqueued"`
	// BuildCount         int64 `json:"build_count"`
	Minutes struct {
		Current int64 `json:"current"`
		// CurrentAverageSec        int64    `json:"current_average_sec"`
		Previous                 int64     `json:"previous"`
		PeriodStartDate          time.Time `json:"period_start_date"`
		PeriodEndDate            time.Time `json:"period_end_date"`
		LastUpdatedAt            time.Time `json:"last_updated_at"`
		IncludedMinutes          int64     `json:"included_minutes"`
		IncludedMinutesWithPacks int64     `json:"included_minutes_with_packs"`
	} `json:"minutes"`
}

func (c Client) GetBuildAccountDetails() (BuildAccountResponse, error) {
	accountDetails := BuildAccountResponse{}

	err := c.doGet("https://api.netlify.com/api/v1/"+c.AccountId+"/builds/status", &accountDetails)
	if err != nil {
		return accountDetails, err
	}

	return accountDetails, nil
}

type AccountResponse []struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
	Capabilities struct {
		Sites struct {
			Included int64 `json:"included"`
			Used     int64 `json:"used"`
		} `json:"sites"`
		Collaborators struct {
			Included int64 `json:"included"`
			Used     int64 `json:"used"`
		} `json:"collaborators"`
	} `json:"capabilities"`
	BillingName     string `json:"billing_name"`
	BillingEmail    string `json:"billing_email"`
	BillingDetails  string `json:"billing_details"`
	BillingPeriod   string `json:"billing_period"`
	PaymentMethodID string `json:"payment_method_id"`
	TypeName        string `json:"type_name"`
	TypeID          string `json:"type_id"`
	// OwnerIds        []string  `json:"owner_ids"`
	// RolesAllowed    []string  `json:"roles_allowed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c Client) GetAccounts() (AccountResponse, error) {
	accountDetails := AccountResponse{}

	err := c.doGet("https://api.netlify.com/api/v1/accounts", &accountDetails)
	if err != nil {
		return accountDetails, err
	}

	return accountDetails, nil
}
