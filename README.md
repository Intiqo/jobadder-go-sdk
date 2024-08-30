# JobAdder Golang SDK

This repository contains a Golang SDK for interacting with the JobAdder API. The SDK is auto-generated using `oapi-codegen` based on the OpenAPI specification provided by JobAdder.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [Initializing the Client](#initializing-the-client)
  - [Examples](#examples)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install the SDK, you can simply run:

```sh
go get github.com/Intiqo/jobadder-go-sdk
```

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
    apiClient, err := client.NewClientWithResponses("https://api.jobadder.com/v2", client.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
        req.Header.Set("Authorization", "Bearer YOUR_ACCESS_TOKEN")
        return nil
    }))

    if err != nil {
        log.Fatalf("Failed to create JobAdder API client: %v", err)
    }

    // Now you can start using the client
}
```

### Examples

List Jobs

```go
jobs, err := apiClient.GetJobsWithResponse(context.Background())
if err != nil {
    log.Fatalf("Failed to retrieve jobs: %v", err)
}

for _, job := range jobs.JSON200.Items {
    log.Printf("Job Title: %s", job.Title)
}
```

Create a New Job

```go
newJob := client.JobAdderCreateJobRequest{
    Title:       "Software Engineer",
    Description: "Join our growing engineering team!",
    Location:    "Remote",
}

jobResponse, err := apiClient.CreateJobWithResponse(context.Background(), newJob)
if err != nil {
    log.Fatalf("Failed to create job: %v", err)
}

log.Printf("Created Job ID: %d", jobResponse.JSON201.JobId)
```

## Configuration

The SDK is generated using oapi-codegen with the following configuration:

 • Package Name: api
 • Output Directory: ./api
 • Models and Methods: Separated based on OpenAPI specification.

If you need to regenerate the SDK, make sure you have the OpenAPI specification file (e.g., jobadder.json) in the root of the repository. You can then run the following command:

Generating the models (types):

```sh
oapi-codegen --package models --generate types -o ./models/models.gen.go openapi.json
```

Generating the APIs:

```sh
oapi-codegen --package api --generate client -o ./api/api.gen.go openapi.json
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

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
