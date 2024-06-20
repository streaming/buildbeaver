# Go API client for client

This is the BuildBeaver Dynamic Build API.

## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 0.3.00
- Package version: 1.0.0
- Build package: org.openapitools.codegen.languages.GoClientCodegen

## Installation

Install the following dependencies:

```shell
go get github.com/stretchr/testify/assert
go get golang.org/x/net/context
```

Put the package under your project folder and add the following in import:

```golang
import client "github.com/GIT_USER_ID/GIT_REPO_ID/client"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```golang
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `sw.ContextServerIndex` of type `int`.

```golang
ctx := context.WithValue(context.Background(), client.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `sw.ContextServerVariables` of type `map[string]string`.

```golang
ctx := context.WithValue(context.Background(), client.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), client.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), client.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to *http://localhost:3003/api/v1/dynamic*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*BuildApi* | [**GetArtifact**](docs/BuildApi.md#getartifact) | **Get** /artifacts/{artifactId} | Reads information about an artifact.
*BuildApi* | [**GetArtifactData**](docs/BuildApi.md#getartifactdata) | **Get** /artifacts/{artifactId}/data | Reads the data for an artifact.
*BuildApi* | [**GetBuild**](docs/BuildApi.md#getbuild) | **Get** /builds/{buildId} | Reads the current build graph for a build.
*BuildApi* | [**GetJob**](docs/BuildApi.md#getjob) | **Get** /jobs/{jobId} | Reads information about a job.
*BuildApi* | [**GetJobGraph**](docs/BuildApi.md#getjobgraph) | **Get** /jobs/{jobId}/graph | Reads information about a job&#39;s graph.
*BuildApi* | [**GetLogData**](docs/BuildApi.md#getlogdata) | **Get** /logs/{logDescriptorId}/data | Reads part of a log.
*BuildApi* | [**GetLogDescriptor**](docs/BuildApi.md#getlogdescriptor) | **Get** /logs/{logDescriptorId} | Fetches a Log Descriptor containing information about part of a log.
*BuildApi* | [**ListArtifacts**](docs/BuildApi.md#listartifacts) | **Get** /builds/{buildId}/artifacts | Reads information about all or some artifacts from a build.
*BuildApi* | [**Ping**](docs/BuildApi.md#ping) | **Get** /ping | Checks for connectivity to the Dynamic API.
*EventsApi* | [**GetEvents**](docs/EventsApi.md#getevents) | **Get** /builds/{buildId}/events | Reads events relating to a build.
*JobsApi* | [**CreateJobs**](docs/JobsApi.md#createjobs) | **Post** /builds/{buildId}/jobs | Creates and add a set of jobs to a build.


## Documentation For Models

 - [Artifact](docs/Artifact.md)
 - [ArtifactDefinition](docs/ArtifactDefinition.md)
 - [ArtifactDependency](docs/ArtifactDependency.md)
 - [ArtifactsPaginatedResponse](docs/ArtifactsPaginatedResponse.md)
 - [Build](docs/Build.md)
 - [BuildDefinition](docs/BuildDefinition.md)
 - [BuildGraph](docs/BuildGraph.md)
 - [BuildOptions](docs/BuildOptions.md)
 - [Commit](docs/Commit.md)
 - [DockerAWSAuth](docs/DockerAWSAuth.md)
 - [DockerAWSAuthDefinition](docs/DockerAWSAuthDefinition.md)
 - [DockerBasicAuth](docs/DockerBasicAuth.md)
 - [DockerBasicAuthDefinition](docs/DockerBasicAuthDefinition.md)
 - [DockerConfig](docs/DockerConfig.md)
 - [DockerConfigDefinition](docs/DockerConfigDefinition.md)
 - [EnvVar](docs/EnvVar.md)
 - [Event](docs/Event.md)
 - [ExternalResourceID](docs/ExternalResourceID.md)
 - [Job](docs/Job.md)
 - [JobDefinition](docs/JobDefinition.md)
 - [JobDependency](docs/JobDependency.md)
 - [JobGraph](docs/JobGraph.md)
 - [LogDescriptor](docs/LogDescriptor.md)
 - [NodeFQN](docs/NodeFQN.md)
 - [Repo](docs/Repo.md)
 - [RunnerApiEndpoints](docs/RunnerApiEndpoints.md)
 - [SecretString](docs/SecretString.md)
 - [SecretStringDefinition](docs/SecretStringDefinition.md)
 - [Service](docs/Service.md)
 - [ServiceDefinition](docs/ServiceDefinition.md)
 - [Step](docs/Step.md)
 - [StepDefinition](docs/StepDefinition.md)
 - [StepDependency](docs/StepDependency.md)
 - [WorkflowTimings](docs/WorkflowTimings.md)


## Documentation For Authorization



### jwt_build_token

- **Type**: API key
- **API key parameter name**: Authorization
- **Location**: HTTP header

Note, each API key must be added to a map of `map[string]APIKey` where the key is: Authorization and passed in as the auth context for each request.


## Documentation for Utility Methods

Due to the fact that model structure members are all pointers, this package contains
a number of utility functions to easily obtain pointers to values of basic types.
Each of these functions takes a value of the given basic type and returns a pointer to it:

* `PtrBool`
* `PtrInt`
* `PtrInt32`
* `PtrInt64`
* `PtrFloat`
* `PtrFloat32`
* `PtrFloat64`
* `PtrString`
* `PtrTime`

## Author


