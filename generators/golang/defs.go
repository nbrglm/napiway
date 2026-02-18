package golang

// SDK FILES DEFINITIONS

type GoSdkClientAndTypesFileTemplateData struct {
	// Package name for the SDK client
	PackageName string

	// Name of the SDK client struct
	ClientName string

	// Version of the SDK client
	//
	// Matches the spec.version from the Spec File.
	ClientVersion string

	Endpoints []GoSdkEndpointDef
}

type GoSdkEndpointDef struct {
	// Name of the Endpoint
	Name string

	Request GoRequestStructDef

	Responses map[int]GoResponseStructDef
}

// SERVER FILE DEFINITIONS

type GoServerFileTemplateData struct {
	PackageName string

	Request GoRequestStructDef

	// Response definitions
	//
	// These structs do NOT need a New() function to be created, but
	// require a WriteResponse function to be created.
	//
	// Only populated if ContentType is application/json.
	// Responses, and therefor WriteResponse functions will exist even if no headers, or body is defined.
	Responses []GoResponseStructDef
}

type GoRequestStructDef struct {
	// Name of the request
	Name        string
	Description *string
	ContentType string
	Method      string
	Path        string

	// Maximum allowed request body size in bytes, if any
	//
	// Nil, a default limit will be applied (256 KB).
	MaxBodyBytes *int64

	HeaderParams []GoParamDef
	QueryParams  []GoParamDef
	PathParams   []GoParamDef

	// Supporting structs that need a New() function to be created
	//
	// Example: RequestBody, if any, will be here, and any other structs
	// that are defined in the spec as a result of being the children of RequestBody.
	SupportingStructs []GoStructDef

	// Name of the request body struct, if any
	//
	// The definition is in SupportingStructs.
	//
	// Nil if there is no request body.
	RequestBodyName *string

	// Authentication methods which are ALL required for this request
	AuthAll []GoAuthMethod

	// Authentication methods of which ANY one is required for this request
	AuthAny []GoAuthMethod
}

type GoResponseStructDef struct {
	// HTTP status code of the response
	StatusCode int

	// Name of the response struct
	//
	// The definition is in SupportingStructs.
	Name string

	Description *string

	// content type of the response
	ContentType string

	Headers []GoParamDef

	// Other supporting structs needed for this response
	//
	// Example: If the response body has nested objects, those structs will be here.
	SupportingStructs []GoStructDef

	// Name of the response body struct, if any
	//
	// The definition is in SupportingStructs.
	//
	// Nil if there is no response body.
	ResponseBodyName *string
}

// Definition of a Go struct
type GoStructDef struct {
	Name   string
	Fields []GoFieldDef
}

// Definition of a field in a Go struct
type GoFieldDef struct {
	// e.g. "User", "Age", "IsActive" etc.
	Name string

	// Field description from the spec
	Description *string

	// Go type (e.g. "string", "float64", "boolean" etc.)
	//
	// If a type should be pointer, it should NOT have a `*` here, the template will add that based on Required.
	//
	// For arrays, this is the array type (e.g. `[]string`, `[]MyStruct` etc.)
	//
	// For objects, this is the struct name (e.g. `MyStruct`)
	GoType string

	// e.g. `json:"fieldName,omitempty"`
	//
	// fieldName will just be Name as is, pascal case.
	Tag string

	// Whether the type has a Validate() method defined on it or on the element type (for arrays)
	RecurseValidate bool

	// Whether the field is an array
	IsArray bool

	// Element type if the field is an array
	ElemGoType string

	// Whether the field is required, detects presence, not emptiness
	Required bool

	// Only considered for strings and arrays
	//
	// The template assumes this is only true for strings and arrays.
	//
	// If it is a []string, i.e an array of strings, this means the array elements must be non-empty strings.
	NonEmpty bool
}

// Definition of a parameter (header, query, path) in Go
type GoParamDef struct {
	// Parameter name
	Name string

	// Parameter name in transport (e.g. header name, query param name, path param name)
	//
	// This should be the exact name as it appears in the HTTP request.
	TransportName string

	// Parameter type
	//
	// This should not have a `*` even if the parameter is optional, the template will add that based on Required.
	GoType string

	// Whether the parameter is required, for strings means non-empty
	Required bool

	// Parameter description
	Description *string
}

type GoAuthMethodType string

const (
	GoAuthMethodTypeHeader GoAuthMethodType = "header"
)

// Definition of an authentication method in Go
type GoAuthMethod struct {
	// ID of the authentication method.
	//
	// This ID is referenced by endpoints to specify which authentication methods to use.
	//
	// Note: This must be unique across all authentication methods in the API spec.
	ID string

	// Name of the header to be used for authentication.
	//
	// This name is referenced by endpoints to specify which authentication methods to use.
	//
	// Note: This must be the name you want to use for the header in the actual HTTP request.
	Name string

	// Name of the transport (e.g. header name)
	//
	// This must be the name you want to use for the header in the actual HTTP request.
	TransportName string

	// Type of the auth method
	Type GoAuthMethodType

	// Description of the auth method
	//
	// Freeform text describing the auth method.
	Description *string

	// Optional format of the auth method
	//
	// Freeform text, e.g. "Bearer {token}", "jwt", etc.
	Format *string
}
