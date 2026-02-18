package typescript

// SDK FILES DEFINITIONS

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

type TsSdkApiAndModelsFilesTemplateData struct {
	ClientName    string
	ClientVersion string
	Endpoints     []TsSdkEndpointDef
}

type TsSdkEndpointDef struct {
	Name      string
	Request   TsRequestDef
	Responses map[int]TsResponseDef
}

type TsRequestDef struct {
	Name        string
	Description *string
	ContentType string
	Method      string
	Path        string

	MaxBodyBytes *int64

	HeaderParams []TsParamDef
	QueryParams  []TsParamDef
	PathParams   []TsParamDef

	SupportingTypes []TsTypeDef

	RequestBodyName *string

	AuthAll []TsAuthMethodDef
	AuthAny []TsAuthMethodDef
}

type TsResponseDef struct {
	StatusCode int

	Name string

	Description *string

	ContentType string

	Headers []TsParamDef

	SupportingTypes []TsTypeDef

	ResponseBodyName *string
}

type TsParamDef struct {
	Name          string
	TransportName string
	Type          string
	Required      bool
	Description   *string
}

type TsAuthMethodType string

const (
	TsAuthMethodTypeHeader TsAuthMethodType = "header"
)

type TsAuthMethodDef struct {
	ID            string
	Name          string
	TransportName string
	Type          TsAuthMethodType
	Description   *string
	Format        *string
}

type TsTypeDef struct {
	Name   string
	Fields []TsFieldDef
}

type TsFieldDef struct {
	Name            string
	Description     *string
	Type            string
	Required        bool
	RecurseValidate bool
	IsArray         bool
	ElemType        string
	NonEmpty        bool
}
