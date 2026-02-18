package spec

import (
	"fmt"
	"slices"
)

type ServerLanguage string

type Config struct {
	GoServer *GoServerGeneration `yaml:"goServer,omitempty"`

	GoSDK *GoSDKGeneration `yaml:"goSdk,omitempty"`

	TsSDK *TsSDKGeneration `yaml:"tsSdk,omitempty"`

	// Spec defining the API to be generated
	Spec Specification `yaml:"spec"`
}

func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config is required")
	}

	if c.GoServer == nil && c.GoSDK == nil && c.TsSDK == nil {
		return fmt.Errorf("at least one of goServer, goSdk, or tsSdk generation must be specified")
	}

	if c.GoServer != nil {
		if err := c.GoServer.Validate(); err != nil {
			return fmt.Errorf("invalid goServer configuration: %w", err)
		}
	}

	if c.GoSDK != nil {
		if err := c.GoSDK.Validate(); err != nil {
			return fmt.Errorf("invalid goSdk configuration: %w", err)
		}
	}

	if c.TsSDK != nil {
		if err := c.TsSDK.Validate(); err != nil {
			return fmt.Errorf("invalid tsSdk configuration: %w", err)
		}
	}

	if err := c.Spec.Validate(); err != nil {
		return fmt.Errorf("invalid spec: %w", err)
	}

	return nil
}

// GoServerGeneration contains configuration for generating Go server code.
//
// It generates models and helper functions for the specified API endpoints.
//
// The generated code will be saved in the specified output directory.
//
// The generated files:
// - helpers.go -> Contains the Per-Endpoint Request DecodeAndValidate functions for the API endpoints.
// - models.go -> Contains the models for request and per-status response bodies.
type GoServerGeneration struct {
	// OutputDir is the directory where the generated server code will be saved
	OutputDir string `yaml:"outputDir"`

	// PackageName is the Go module name for the generated server code
	//
	// Usually, this is same as the last part of the output directory path.
	PackageName string `yaml:"packageName"`
}

func (g *GoServerGeneration) Validate() error {
	if g == nil {
		return fmt.Errorf("goServer configuration is required")
	}
	if g.OutputDir == "" {
		return fmt.Errorf("outputDir is required for goServer generation")
	}
	if g.PackageName == "" {
		return fmt.Errorf("moduleName is required for goServer generation")
	}
	return nil
}

// GoSDKGeneration contains configuration for generating Go SDK code.
//
// It generates client code for the specified API endpoints.
//
// The generated code will be saved in the specified output directory.
//
// The generated files:
// - client.go -> Contains the client struct and methods for calling the API endpoints.
// - models.go -> Contains the models for request and per-status response bodies.
// - go.mod -> The Go module file for the generated SDK.
// - .gitignore -> The gitignore file for the generated SDK.
// - README.md -> The readme file for the generated SDK.
type GoSDKGeneration struct {
	// OutputDir is the directory where the generated SDK will be saved
	OutputDir string `yaml:"outputDir"`

	// The module name for the generated SDK
	//
	// E.g. "github.com/username/project/sdk"
	ModuleName string `yaml:"moduleName"`
}

func (g *GoSDKGeneration) Validate() error {
	if g == nil {
		return fmt.Errorf("goSdk configuration is required")
	}

	if g.OutputDir == "" {
		return fmt.Errorf("outputDir is required for goSdk generation")
	}

	if g.ModuleName == "" {
		return fmt.Errorf("moduleName is required for goSdk generation")
	}
	return nil
}

// TsSDKGeneration contains configuration for generating TypeScript SDK code.
//
// It generates client code for the specified API endpoints.
//
// The generated code will be saved in the specified output directory.
//
// The generated files:
// - client.ts -> Contains the client class and methods for calling the API endpoints.
// - models.ts -> Contains the models for request and per-status response bodies.
// - package.json -> The npm package file for the generated SDK.
// - tsconfig.json -> The TypeScript configuration file for the generated SDK.
// - .gitignore -> The gitignore file for the generated SDK.
// - README.md -> The readme file for the generated SDK.
type TsSDKGeneration struct {
	// OutputDir is the directory where the generated SDK will be saved
	OutputDir string `yaml:"outputDir"`

	// The npm package name for the generated SDK (package.json name field)
	PackageName string `yaml:"packageName"`

	// A brief description of the SDK (package.json description field)
	Description *string `yaml:"description,omitempty"`

	// The author of the SDK (package.json author field)
	Author *string `yaml:"author,omitempty"`

	// The license for the SDK (package.json license field)
	License *string `yaml:"license,omitempty"`

	// The repository URL for the SDK (package.json repository.url field)
	Repository *string `yaml:"repository,omitempty"`

	// The website or documentation URL for the SDK (package.json homepage field)
	Website *string `yaml:"website,omitempty"`

	// Keywords for the SDK (package.json keywords field)
	Keywords []string `yaml:"keywords,omitempty"`
}

func (t *TsSDKGeneration) Validate() error {
	if t == nil {
		return fmt.Errorf("tsSdk configuration is required")
	}
	if t.OutputDir == "" {
		return fmt.Errorf("outputDir is required for tsSdk generation")
	}
	if t.PackageName == "" {
		return fmt.Errorf("packageName is required for tsSdk generation")
	}
	return nil
}

type Specification struct {
	ApiName     string `yaml:"apiName"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`

	// The website or documentation URL for the API
	Website *string `yaml:"website,omitempty"`

	// Contact email for the API
	Contact *string `yaml:"contact,omitempty"`

	// License information for the API, e.g., "MIT", "Apache-2.0", etc.
	//
	// This is an optional and freeform field.
	License *string `yaml:"license,omitempty"`

	// All possible authentication methods for the API.
	Auth []AuthMethod `yaml:"auth,omitempty"`

	// Map of endpoint names to endpoint definitions.
	//
	// The key is a unique name for the endpoint, e.g., "GetUser", "CreatePost", etc.
	// This will be used in code generation and documentation.
	Endpoints map[string]*Endpoint `yaml:"endpoints"`
}

func (s *Specification) Validate() error {
	if s.ApiName == "" {
		return fmt.Errorf("apiName is required")
	}
	if s.Version == "" {
		return fmt.Errorf("version is required")
	}
	for name, endpoint := range s.Endpoints {
		if err := endpoint.Validate(s.Auth); err != nil {
			return fmt.Errorf("endpoint %s: %w", name, err)
		}
	}
	for i, am := range s.Auth {
		if err := am.Validate(); err != nil {
			return fmt.Errorf("auth method %d: %w", i, err)
		}
	}
	return nil
}

type EndpointAuthMode string

const (
	EndpointAuthModeRequireAll EndpointAuthMode = "all"
	EndpointAuthModeRequireAny EndpointAuthMode = "any"
)

// EndpointAuthentication defines the authentication requirements for an endpoint.
//
// It specifies which authentication methods must be satisfied to access the endpoint.
// If both `All` and `Any` are specified, both conditions must be satisfied.
// If neither is specified, the endpoint does not require authentication.
//
// Example:
//
//	auth:
//	  all:
//	    - "apiKey"
//	  any:
//	    - "refreshToken"
//	    - "sessionToken"
//
// In the above example, the endpoint requires the "apiKey" authentication method to be satisfied,
// and at least one of "refreshToken" or "sessionToken" to be satisfied.
type EndpointAuthentication struct {
	// List of IDs of authentication methods that must be satisfied for this endpoint.
	All []string `yaml:"all"`

	// List of IDs of authentication methods where at least one must be satisfied for this endpoint.
	Any []string `yaml:"any"`
}

func (ea *EndpointAuthentication) Validate(globalAuth []AuthMethod) error {
	if ea == nil {
		return nil
	}
	authIDs := make([]string, 0)
	for _, am := range globalAuth {
		authIDs = append(authIDs, am.ID)
	}
	for _, id := range ea.All {
		if !slices.Contains(authIDs, id) {
			return fmt.Errorf("all: unknown auth method: %s", id)
		}
	}
	for _, id := range ea.Any {
		if !slices.Contains(authIDs, id) {
			return fmt.Errorf("any: unknown auth method: %s", id)
		}
	}
	return nil
}

type AuthMethodType string

const (
	AuthMethodHeader AuthMethodType = "header"
)

type AuthMethod struct {
	// This id is referenced by endpoints to specify which authentication methods to use.
	ID string `yaml:"id"`

	// Name of the header to be used for authentication.
	Name string `yaml:"name"`

	// Name of the transport (e.g. header name)
	//
	// This must be the name you want to use for the header in the actual HTTP request.
	TransportName string `yaml:"transportName"`

	Type        AuthMethodType `yaml:"type"`
	Description *string        `yaml:"description,omitempty"`

	// Optional format string for the header value.
	//
	// Freeform text, e.g., "Bearer {token}", "Token {token}", "jwt" etc.
	//
	// NOT ENFORCED, just for documentation purposes.
	Format *string `yaml:"format,omitempty"`
}

func (am *AuthMethod) Validate() error {
	if am.Name == "" {
		return fmt.Errorf("name is required")
	}
	switch am.Type {
	case AuthMethodHeader:
		// valid
	default:
		return fmt.Errorf("invalid type: %s", am.Type)
	}
	return nil
}

type EndpointMethod string

const (
	EndpointMethodGet    EndpointMethod = "GET"
	EndpointMethodPost   EndpointMethod = "POST"
	EndpointMethodPut    EndpointMethod = "PUT"
	EndpointMethodDelete EndpointMethod = "DELETE"
)

type Endpoint struct {
	// HTTP method for the endpoint, e.g., "GET", "POST", "PUT", "DELETE", etc.
	Method EndpointMethod `yaml:"method"`

	// Path of the endpoint, e.g., "/users", "/posts/{id}", etc.
	//
	// This is a required field.
	// It should be a valid URL path.
	Path string `yaml:"path"`

	// Content type for the request, if empty, defaults to application/json
	ContentType *string `yaml:"contentType,omitempty"`

	// Description of the endpoint
	Description *string `yaml:"description,omitempty"`

	// Path parameters for the request
	PathParams []Param `yaml:"pathParams,omitempty"`

	// Headers for the request
	//
	// No need to specify authentication headers here,
	// those will be handled via the `auth` field.
	Headers []Param `yaml:"headers,omitempty"`

	// Query parameters for the request
	QueryParams []Param `yaml:"queryParams,omitempty"`

	// Authentication requirements for this endpoint.
	//
	// If omitted, the endpoint does not require authentication.
	Auth *EndpointAuthentication `yaml:"auth,omitempty"`

	// RequestBody body for the endpoint
	RequestBody *HTTPBody `yaml:"requestBody,omitempty"`

	// Maximum allowed request body size in bytes, if any
	//
	// Nil, a default limit will be applied (256 KB).
	MaxBodyBytes *int64

	// Map of HTTP status codes to responses
	Responses map[int]*Response `yaml:"responses,omitempty"`
}

func (e *Endpoint) Validate(authMethods []AuthMethod) error {
	if e == nil {
		return fmt.Errorf("endpoint is nil")
	}
	switch e.Method {
	case EndpointMethodGet, EndpointMethodPost, EndpointMethodPut, EndpointMethodDelete:
		// valid
	default:
		return fmt.Errorf("invalid method: %s", e.Method)
	}
	if e.Path == "" {
		return fmt.Errorf("path is required")
	}
	if e.ContentType == nil || *e.ContentType == "" {
		defaultContentType := "application/json"
		e.ContentType = &defaultContentType
	}
	if e.RequestBody != nil {
		if err := e.RequestBody.Validate(); err != nil {
			return fmt.Errorf("request: %w", err)
		}
	}
	for code, resp := range e.Responses {
		if code < 100 || code > 599 {
			return fmt.Errorf("response %d: invalid HTTP status code", code)
		}
		if err := resp.Validate(); err != nil {
			return fmt.Errorf("response %d: %w", code, err)
		}
	}
	for i, pp := range e.PathParams {
		if err := pp.Validate(); err != nil {
			return fmt.Errorf("pathParam %d: %w", i, err)
		}
	}
	for i, header := range e.Headers {
		if err := header.Validate(); err != nil {
			return fmt.Errorf("header %d: %w", i, err)
		}
	}
	for i, qp := range e.QueryParams {
		if err := qp.Validate(); err != nil {
			return fmt.Errorf("queryParam %d: %w", i, err)
		}
	}
	if e.Auth != nil {
		// Note: global auth methods are not passed here for validation.
		// This should be handled in Specification.Validate.
		if err := e.Auth.Validate(authMethods); err != nil {
			return fmt.Errorf("auth: %w", err)
		}
	}
	return nil
}

type Response struct {
	// Description of the response
	Description *string `yaml:"description,omitempty"`

	// Response headers
	Headers []Param `yaml:"headers,omitempty"`

	ContentType *string `yaml:"contentType,omitempty"`

	// Response body
	Body *HTTPBody `yaml:"body,omitempty"`
}

func (r *Response) Validate() error {
	if r == nil {
		return fmt.Errorf("response is nil")
	}
	if r.ContentType == nil || *r.ContentType == "" {
		defaultContentType := "application/json"
		r.ContentType = &defaultContentType
	}
	if r.Body != nil {
		if err := r.Body.Validate(); err != nil {
			return fmt.Errorf("body: %w", err)
		}
	}
	for i, header := range r.Headers {
		if err := header.Validate(); err != nil {
			return fmt.Errorf("header %d: %w", i, err)
		}
	}
	return nil
}

type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeNumber  FieldType = "number"
	FieldTypeBoolean FieldType = "boolean"
	FieldTypeObject  FieldType = "object"
	FieldTypeArray   FieldType = "array"
)

type Field struct {
	// Type of the field
	Type FieldType `yaml:"type"`

	// Description of the field
	Description *string `yaml:"description,omitempty"`

	// For object type, the properties of the object
	Properties map[string]Field `yaml:"properties,omitempty"`

	// For array type, the items of the array
	Items *Field `yaml:"items,omitempty"`

	// Whether the field is required to be provided
	Required bool `yaml:"required,omitempty"`

	// For string and array types, whether the field must be non-empty
	//
	// If it is a []string, i.e an array of strings, this means the array elements must be non-empty strings.
	NonEmpty bool `yaml:"nonEmpty,omitempty"`
}

func (f *Field) Validate() error {
	if f == nil {
		return fmt.Errorf("field is nil")
	}
	if f.NonEmpty {
		switch f.Type {
		case FieldTypeString, FieldTypeArray:
			break
		default:
			return fmt.Errorf("nonEmpty is only applicable for string and array types, not for type %s", f.Type)
		}

		if !f.Required {
			return fmt.Errorf("nonEmpty requires the field to be required")
		}
	}

	switch f.Type {
	case FieldTypeString, FieldTypeNumber, FieldTypeBoolean:
		if len(f.Properties) != 0 || f.Items != nil {
			return fmt.Errorf("properties and items are not applicable for type %s", f.Type)
		}
	case FieldTypeObject:
		if len(f.Properties) == 0 {
			return fmt.Errorf("properties is required for object type")
		}
		for name, prop := range f.Properties {
			if err := prop.Validate(); err != nil {
				return fmt.Errorf("property %s: %w", name, err)
			}
		}
	case FieldTypeArray:
		if f.Items == nil {
			return fmt.Errorf("items is required for array type")
		}
		if err := f.Items.Validate(); err != nil {
			return fmt.Errorf("items: %w", err)
		}
	default:
		return fmt.Errorf("invalid field type: %s", f.Type)
	}

	return nil
}

type HTTPBody struct {
	Description *string `yaml:"description,omitempty"`

	// The properties of the body object
	Properties map[string]Field `yaml:"properties,omitempty"`
}

func (b *HTTPBody) Validate() error {
	if b == nil {
		return fmt.Errorf("http body is nil")
	}
	if len(b.Properties) == 0 {
		return fmt.Errorf("properties is required for http body")
	}
	for name, prop := range b.Properties {
		if err := prop.Validate(); err != nil {
			return fmt.Errorf("property %s: %w", name, err)
		}
	}
	return nil
}

type ParamType string

const (
	ParamTypeString  ParamType = "string"
	ParamTypeNumber  ParamType = "number"
	ParamTypeBoolean ParamType = "boolean"
)

type Param struct {
	// Name of the parameter
	Name string `yaml:"name"`

	// Name of the parameter in transport (e.g. header name, query param name, path param name)
	//
	// This should be the exact name as it appears in the HTTP request.
	TransportName string `yaml:"transportName"`

	// Type of the parameter
	Type ParamType `yaml:"type"`

	// Description of the parameter
	Description *string `yaml:"description,omitempty"`

	// Whether the parameter is required
	//
	// For string, required also means non-empty.
	Required bool `yaml:"required,omitempty"`
}

func (p *Param) Validate() error {
	if p == nil {
		return fmt.Errorf("param is nil")
	}
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.TransportName == "" {
		return fmt.Errorf("transportName is required")
	}
	// Required semantics:
	// - string: must be present and non-empty
	// - number/boolean: must be present
	switch p.Type {
	case ParamTypeString, ParamTypeNumber, ParamTypeBoolean:
		// valid
	default:
		return fmt.Errorf("invalid param type: %s", p.Type)
	}
	return nil
}
