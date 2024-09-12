package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Intiqo/jobadder-go-sdk/api"
)

const (
	authenticationUrl = "https://id.jobadder.com/connect/token"
)

type jobAdderTokenResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Api          string `json:"api"`
}

type JobAdderClientParams struct {
	APIBaseURL        string
	ClientID          string
	ClientSecret      string
	AuthorizationCode string
	RedirectUri       string
	AccessToken       string
	RefreshToken      string
	TokenExpiryTime   time.Time
}

type JobAdderClient struct {
	Api    *api.ClientWithResponses
	Params *JobAdderClientParams
}

// NewJobAdderClient initializes the client and checks token expiry for every request
func NewJobAdderClient(params *JobAdderClientParams) (*JobAdderClient, error) {
	jc := &JobAdderClient{Params: params}

	// Authenticate initially
	err := jc.Authenticate()
	if err != nil {
		return nil, err
	}

	// Create the API client
	apiClient, err := api.NewClientWithResponses(jc.Params.APIBaseURL,
		api.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			// Set the Authorization header with the valid token
			req.Header.Set("Authorization", "Bearer "+jc.Params.AccessToken)
			return nil
		}),
	)
	if err != nil {
		return nil, err
	}
	jc.Api = apiClient

	return jc, nil
}

// authenticate authenticates the client and refreshes the token if necessary
func (c *JobAdderClient) Authenticate() (err error) {
	var payload *strings.Reader
	if c.Params.AccessToken == "" && c.Params.RefreshToken == "" {
		// Get a new access token using the authorization code
		payload = strings.NewReader("client_id=" + c.Params.ClientID + "&client_secret=" + c.Params.ClientSecret + "&grant_type=authorization_code&code=" + c.Params.AuthorizationCode + "&redirect_uri=" + c.Params.RedirectUri)
	} else {
		// Refresh the access token using the refresh token
		payload = strings.NewReader("client_id=" + c.Params.ClientID + "&client_secret=" + c.Params.ClientSecret + "&grant_type=refresh_token&refresh_token=" + c.Params.RefreshToken)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, authenticationUrl, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	var result jobAdderTokenResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	// Update the tokens and expiry time in the params
	c.Params.AccessToken = result.AccessToken
	c.Params.RefreshToken = result.RefreshToken
	c.Params.APIBaseURL = result.Api
	c.Params.TokenExpiryTime = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)

	return nil
}

func (c *JobAdderClient) RefreshToken() (err error) {
	return c.Authenticate()
}
