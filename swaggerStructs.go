package main

type SwaggerTemplate struct {
  Swagger string `json:"swagger"`
  Info InfoStruct `json:"info"`
  Host string `json:"host"`
  BasePath string `json:"basePath"`
  Tags []TagStruct `json:"tags"`
  Schemes []string `json:"schemes"`
  //Paths map[string]map[string]VerbStruct `json:"paths"`
  Paths PathsStruct `json:"paths"`
  //SecurityDefinitions string `json:"securityDefinitions"`
  Definitions DefinitionStruct `json:"definitions"`
}

type InfoStruct struct {
  Description string `json:"description"`
  Version string `json:"version"`
  Title string `json:"title"`
  Contact struct {
    Email string `json:"email"`
  } `json:"contact"`
  License struct {
    Name string `json:"name"`
    URL string `json:"url"`
  } `json:"license"`
} 

type Definition struct {
  Type string `json:"type"`
  Properties string `json:"properties"`
  Xml struct {
    Name string `json:"name"`
  } `json:"xml"`
}

type DefinitionStruct struct {
  Definitions []Definition
}

type TagStruct struct {
  Name string `json: "name"`
  Description string `json: "description"`
}

type PathsStruct struct {
  Paths []PathStruct 
}

type PathStruct struct {
  Verbs []VerbStruct
}

type VerbStruct struct {
  Tags []string `json:"tags"`
  Summary string `json:"summary"`
  Description string `json:"description"`
  OperationID string `json:"operationId"`
  Consumes []string `json:"consumes"`
  Produces []string `json:"produces"`
 // Port string `json:"port"`
  Connection ConnectionStruct `json:"connection"`
  //Parameters []VerbParameters `json:"parameters"`
  Parameters string `json:"parameters"`
}

type ConnectionStruct struct {
  Port int `json:"port"`
  Type string `json:"type"`
  RemoteType string `json:"remoteType"`
}

type VerbParameters struct {
  In string `json:"in"`
    Name string `json:"name"`
    Description string `json:"description"`
    Required bool `json:"required"`
    Schema struct {
      Ref string `json:"$ref"`
    } `json:"schema"`
}
