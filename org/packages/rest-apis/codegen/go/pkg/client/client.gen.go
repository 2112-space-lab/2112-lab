// Package client provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

// HealthInfo defines model for HealthInfo.
type HealthInfo struct {
	ServiceName    *string    `json:"serviceName,omitempty"`
	ServiceStatus  *string    `json:"serviceStatus,omitempty"`
	UptimeSinceUtc *time.Time `json:"uptimeSinceUtc,omitempty"`
}

// VersionInfo defines model for VersionInfo.
type VersionInfo struct {
	BuildNumber string `json:"buildNumber"`
	Version     string `json:"version"`
}

// N400BadRequest defines model for N400BadRequest.
type N400BadRequest struct {
	Message *string `json:"message,omitempty"`
}

// N500InternalError defines model for N500InternalError.
type N500InternalError struct {
	Message *string `json:"message,omitempty"`
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetVersion request
	GetVersion(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetHealth request
	GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetVersion(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetVersionRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetHealthRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetVersionRequest generates requests for GetVersion
func NewGetVersionRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/app/api/v1/version")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetHealthRequest generates requests for GetHealth
func NewGetHealthRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/app/health")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetVersionWithResponse request
	GetVersionWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetVersionResponse, error)

	// GetHealthWithResponse request
	GetHealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error)
}

type GetVersionResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *VersionInfo
	JSON400      *N400BadRequest
	JSON500      *N500InternalError
}

// Status returns HTTPResponse.Status
func (r GetVersionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetVersionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetHealthResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *HealthInfo
	JSON400      *N400BadRequest
	JSON500      *N500InternalError
}

// Status returns HTTPResponse.Status
func (r GetHealthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetHealthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetVersionWithResponse request returning *GetVersionResponse
func (c *ClientWithResponses) GetVersionWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetVersionResponse, error) {
	rsp, err := c.GetVersion(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetVersionResponse(rsp)
}

// GetHealthWithResponse request returning *GetHealthResponse
func (c *ClientWithResponses) GetHealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error) {
	rsp, err := c.GetHealth(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetHealthResponse(rsp)
}

// ParseGetVersionResponse parses an HTTP response from a GetVersionWithResponse call
func ParseGetVersionResponse(rsp *http.Response) (*GetVersionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetVersionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest VersionInfo
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest N400BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest N500InternalError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGetHealthResponse parses an HTTP response from a GetHealthWithResponse call
func ParseGetHealthResponse(rsp *http.Response) (*GetHealthResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetHealthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest HealthInfo
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest N400BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest N500InternalError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RUTU/kOBD9K1btHkMny8IlN5DQbl9aiJ7hMuqDiSuJUWKbcqUlhPLfR5WPpqGDmDmM",
	"RnNLyuWqV+/51QsUvg3eoeMI+QsQxuBdxOFnc5Fl19rc4VOHkSVSeMfohk8dQmMLzda79DF6J7FY1Nhq",
	"+QrkAxLbsVCLMeoK5ZOfA0IOkcm6Cvo+mSP+4RELhl5CBmNBNkhtyOFaGzVj6BPYXGbZ2jGS080Nkaff",
	"BmxGobZIeyQ1opG8sd/Q4n/UDddrV/rT9hFpbwvc6HYJQjKfb1lzFxczusC2xa11BX7lQlJKT61myMFo",
	"xjM5heTz4RK4R4rWu2WgD51tzKZrH5AWYezHy8s8Ej51ltBA/u2QmLwpuVsk205Q3pL+pUZ1dbtWBkvr",
	"rASVLxXXqASyrjR7UoTaIKmJPyHAciP1r0JQd+Ph9nB4QA/Z6mKVyUA+oNPBQg7/rrLVOSQQNNcDF6kO",
	"IdXBpvt/0qO5K+RTsP8hD9CmPCUjiTxHqG9fUZ8AEwmG5LUZi90f6Htj1PMs+ykT/E1YQg5/pa/mT6cn",
	"mx6/g4U3L/TTaEYltxtkNCp2RYExll3TPAt9FyOgpT4H4Om79dIncPlD1078Pziua1tNzxPnx3yL+rqK",
	"8vpGL8JOLservicey1mPgM/nGtFmx+LE+U4NfKM/ROvlj1Zno/FCcafMhSfy9JrfkTVcMDhqTIIGOGsihZg4x",
	"T1Pn+WzYDmig3x0avC9040zw1nFUpSdVIbN11SDwsUsJGy0csp91n+FrZ9TrOnPDDp+H6Hf99wAAAP//",
	"dm/Ux1oHAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
