package golang

type GoTypesFileData struct {
	PackageName string

	Types []TypeData
}

type GoReqResFileData struct {
	RequestData
	PackageName string
}

type GoSdkClientFileData struct {
	PackageName   string
	ClientName    string
	ClientVersion string
	Endpoints     []EndpointData
}

type EndpointData struct {
	Name    string
	Request RequestData
}

type RequestData struct {
	Name        string
	Description *string

	// Indicates whether the request body should be ignored by the generated code
	// and only focus on other aspects of the request, such as headers, query parameters, path parameters, etc.
	RawBody bool

	Method string

	// URL path of the request, e.g. "/users/{userId}"
	Path string

	MaxBodyBytes *int64
	ContentType  string

	HeaderParams []ParamData
	QueryParams  []ParamData
	PathParams   []ParamData

	RequestBodyName *string

	AuthAll []AuthMethodData
	AuthAny []AuthMethodData

	// Responses
	Responses []ResponseData
}

type ResponseData struct {
	StatusCode  int
	Name        string
	Description *string

	// Indicates whether the response body should be ignored by the generated code
	//
	// This makes the generated code only set the status code and headers, without trying to include the response body.
	RawBody          bool
	ContentType      string
	Headers          []ParamData
	ResponseBodyName *string
}

type ParamData struct {
	Name          string
	TransportName string
	Type          string
	Required      bool
	Description   *string
}

type TypeData struct {
	Name        string
	Description *string

	Fields []TypeFieldData
	Enum   []string
}

// TypeStr is a string representation of a Go type, used for code generation purposes.
//
// It is not a full representation of all possible types, since the ones in the enum are
// primitive-ish types, while more complex types (e.g. structs) are represented by their name as a string, and the actual struct definition is in TypeData.
const (
	TypeStrString         = "string"
	TypeStrInteger        = "int64"
	TypeStrDouble         = "float64"
	TypeStrBoolean        = "bool"
	TypeStrFreeFormObject = "map[string]any"
)

type TypeFieldData struct {
	Name        string
	Description *string

	// Type is either a primitive type (TypeStr) or the name of a struct defined in Types.
	Type string

	// Whether the type is a pointer type. Cases:
	// True if the field is optional.
	// True if the field is a non-primitive type (e.g. struct), since we want to use pointer types for structs to allow for nil values and to avoid copying large structs.
	// False if IsArray is true.
	// False if the field is required and a primitive type, since we want to use value types for required primitive fields for better ergonomics.
	PtrType bool

	// Whether the type is a non-primitive type (e.g. struct).
	IsNonPrimitiveType bool

	// Tag is the struct field tag, used for JSON serialization, validation, etc.
	Tag string

	// Used to put the [] in the right place for array types. If true, the generated code will use []Type for this type.
	IsArray bool

	// Used to indicate that the field is an enum and should be parsed appropriately.
	IsEnum bool

	Required bool

	NonEmpty bool
}

type AuthMethodType string

const (
	AuthMethodTypeHeader AuthMethodType = "header"
)

type AuthMethodData struct {
	ID string

	Name string

	TransportName string

	Type AuthMethodType

	Description *string

	// Optional format of the auth method
	//
	// For example
	Format *string
}
