package typescript

type TsSdkPackageJsonTemplateData struct {
	PackageName string
	Version     string
	Description *string
	Website     *string
	Repository  *string
	License     *string
	Author      *string
	Keywords    []string
}

type TsSdkModelsFileData struct {
	ClientName string
	Types      []TypeData

	Requests []RequestData
}

type TsSdkApiFileData struct {
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

	RawBody bool
	Method  string

	Path        string
	ContentType string

	HeaderParams []ParamData
	QueryParams  []ParamData
	PathParams   []ParamData

	RequestBodyName *string

	AuthAll []AuthMethodData
	AuthAny []AuthMethodData

	Responses []ResponseData
}

type ResponseData struct {
	StatusCode  int
	Name        string
	Description *string

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
}

const (
	TypeStrString         = "string"
	TypeStrInteger        = "integer"
	TypeStrDouble         = "double"
	TypeStrBoolean        = "boolean"
	TypeStrFreeFormObject = "Record<string, any>"
)

type TypeFieldData struct {
	Name        string
	Description *string
	Type        string
	IsArray     bool
	Required    bool
	NonEmpty    bool
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
