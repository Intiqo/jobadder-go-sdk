# JobAdder Golang SDK

This repository contains a Golang SDK for interacting with the JobAdder API. The SDK is auto-generated using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen/) based on the [OpenAPI specification provided by JobAdder](https://api.jobadder.com/v2/docs).

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [Initializing the Client](#initializing-the-client)
  - [API Examples](#api-examples)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install the SDK, execute the below command:

```sh
go get github.com/Intiqo/jobadder-go-sdk
```

To install a specific version of the SDK, execute the below command:

```sh
go get github.com/Intiqo/jobadder-go-sdk@v0.0.1
```

Replace the version as required

## Usage

### Initializing the client

```go
package main

import (
    "context"
    "log"
    "github.com/Intiqo/jobadder-go-sdk/client"
)

func main() {
 // Get a new client
 params := client.JobAdderClientParams{
  APIBaseURL:        "https://api.jobadder.com/v2",
  ClientID:          "<your_client_id>",
  ClientSecret:      "<your_client_secret>",
  AuthorizationCode: "<your_authorization_code>",
  AccessToken:       "",
  TokenExpiryTime:   time.Time{},
  RefreshToken:      "",
  RedirectUri:       "https://<my_app.com>/jobadder",
 }
 cl, err := client.NewJobAdderClient(&params)
 if err != nil {
  panic(err)
 }

// Now you can start using the client
}
```

#### Obtain Access Token from Authorization Code

When you are connecting with JobAdder the very first time, [you need to obtain an authorization code](https://api.jobadder.com/v2/docs#section/Getting-Started/Authentication) for your app, and pass that authorization code along with the Client ID & Client Secret to the SDK. At this point, pass empty string for access & refresh token parameters so that the SDK can use the authorization code as the grant type to get an access token.

The SDK will then automatically update the params passed with the new Access Token, Refresh Token, API Url & Token Expiry time. Store the obtained access & refresh tokens securely to reuse it later if you restart your server or other such instances.

- If the Client ID or Client Secret is invalid, JobAdder will respond with `{"error":"invalid_client"}`.
- If the Authorization Code is invalid, JobAdder will respond with `{"error":"invalid_grant"}`.
- In both of these cases, the error is transparently passed on by the SDK from JobAdder.

#### Obtain Access Token from Refresh Token

When you initially obtain the access token from authorization code, it is important to store the refresh token somewhere safe to be reused later. If you don't store it, you'll need to ask your user to pass through the authentication process every `x` minutes (currently 60 minutes) as the access token gets expired.

If you've stored the refresh token, and want to get a new access token, use the function `RefreshToken()` with the parameters you had securely stored. If the refresh token is valid, the SDK will get a new access token and update the API client for further usage.

```go
package main

import (
    "context"
    "log"
    "github.com/Intiqo/jobadder-go-sdk/client"
)

func main() {
 // Get a new client
 params := client.JobAdderClientParams{
  APIBaseURL:        "https://api.jobadder.com/v2",
  ClientID:          "<your_client_id>",
  ClientSecret:      "<your_client_secret>",
  AuthorizationCode: "<your_authorization_code>",
  AccessToken:       "<your_access_token>",
  TokenExpiryTime:   time.Time{},
  RefreshToken:      "<your_refresh_token>",
  RedirectUri:       "https://<my_app.com>/jobadder",
 }
 cl, err := client.NewJobAdderClient(&params)
 if err != nil {
  panic(err)
 }

// Get a new access token from the refresh token
 err = cl.RefreshToken()
 if err != nil {
  return
 }
}
```

- If the refresh token is invalid, JobAdder will respond with `{"error":"invalid_grant"}`.
- So, you'll need to initiate the entire authorization process again in this case.

All of the above steps are standard [OAuth 2.0 protocol](https://oauth.net/2/) and is well explained in the [JobAdder Authentication Documentation](https://api.jobadder.com/v2/docs#section/Getting-Started/Authentication).

### Automatic Token Refresh

It is best advised to have an automatic token refresh process where you get a new access token 'x' minutes before the access token expiry time. This can be easily achieved using something like go-cron or writing a native version of cron jobs.

If you are using go-cron, you can refer to the below sample code:

```go
package main

import (
    "context"
    "log"
    "github.com/Intiqo/jobadder-go-sdk/client"
    "github.com/go-co-op/gocron"
)

func main() {
 // Get a new client
 params := client.JobAdderClientParams{
  APIBaseURL:        "https://api.jobadder.com/v2",
  ClientID:          "<your_client_id>",
  ClientSecret:      "<your_client_secret>",
  AuthorizationCode: "<your_authorization_code>",
  AccessToken:       "<your_access_token>",
  TokenExpiryTime:   time.Time{},
  RefreshToken:      "<your_refresh_token>",
  RedirectUri:       "https://<my_app.com>/jobadder",
 }
 cl, err := client.NewJobAdderClient(&params)
 if err != nil {
  panic(err)
 }

 // Change timezone as required
 sch := gocron.NewScheduler(time.UTC)

 // Schedule a job to refresh token
 c.scheduleTokenRefresh(cl.Params.TokenExpiryTime, sch)
}

func scheduleTokenRefresh(expiryTime time.Time, sch *gocron.Scheduler) {
 // Calculate the time 5 minutes before expiry
 refreshTime := expiryTime.Add(-5 * time.Minute)
 delay := refreshTime.Sub(time.Now())

 if delay > 0 {
  c.TokenJob, _ = sch.After(delay).Do(func() {
   err := c.RefreshToken()
   if err != nil {
    // Handle token refresh error (log or retry)
    return
   }
   // Reschedule the job after token is refreshed
   c.scheduleTokenRefresh(cl.Params.TokenExpiryTime, sch)
  })
 }
}
```

### API Examples

List Companies

```go
compParams := &api.FindCompaniesParams{}
companies, err := apiClient.FindCompaniesWithResponse(context.Background(), compParams)
if err != nil {
log.Fatalf("Failed to get company: %v", err)
}

for _, comp := range *companies.JSON200.Items {
log.Printf("Company: %s", *comp.Name)
}
```

## Configuration

The SDK is generated using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen/) with the following configuration:

- Package Name: api
- Output Directory: ./api

If you need to regenerate the SDK, make sure you have the OpenAPI specification file (e.g., openapi.json) in the root of the repository. You can then run the following command:

```sh
oapi-codegen --package=api --generate types,client openapi.json > ./api/api.gen.go
```

## Contributing

We welcome contributions to the JobAdder SDK! If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.

Steps to Contribute

 1. Fork the repository.
 2. Create a new branch (git checkout -b feature/your-feature).
 3. Make your changes.
 4. Commit your changes (git commit -m 'Add some feature').
 5. Push to the branch (git push origin feature/your-feature).
 6. Open a Pull Request.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
