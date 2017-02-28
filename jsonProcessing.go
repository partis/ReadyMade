package main

import (
  "encoding/json"
  "io/ioutil"
  "bytes"
  "strings"
  "github.com/golang/glog"
)

func swaggerToJson(template *SwaggerTemplate) string {
  bytes, err := json.Marshal(template)
  if err != nil {
    glog.Fatal(err.Error())
  }

  return string(bytes)
}

func toJson(template interface{}) string {
  bytes, err := json.Marshal(template)
  if err != nil {
    glog.Fatal(err.Error())
  }

  return string(bytes)
}

func isJson(s string) (map[string]interface{}, bool) {
  var js map[string]interface{}

  err := json.Unmarshal([]byte(s), &js)

  return js, err == nil
}

func (ver *VerbStruct) MarshalJSON() ([]byte, error) {
  var jsonBytes []byte
  comma := []byte(",")

  jsonBytes = append(jsonBytes, []byte("{")...)
  
  //Tags
  tags, err := json.Marshal(ver.Tags)
  if err != nil {
    glog.Error(err)
  }
  jsonBytes = append(jsonBytes, []byte("\"tags\":")...)
  jsonBytes = append(jsonBytes, tags...)
  jsonBytes = append(jsonBytes, comma...)
  //jsonBytes = append(jsonBytes, []byte("\"tags\":[")...)
  //tagString := ""
  //for i := range ver.Tags {
  //  tagString = "\"" + ver.Tags[i] + "\","
  //}
  //tagString = strings.TrimSuffix(tagString, ",")
  //jsonBytes = append(jsonBytes, []byte(tagString)...)
  //jsonBytes = append(jsonBytes, []byte("],")...)

  //Summary
  summary, err := json.Marshal(ver.Summary)
  if err != nil {
    glog.Error(err)
  }
  jsonBytes = append(jsonBytes, []byte("\"summary\":")...)
  jsonBytes = append(jsonBytes, summary...)
  jsonBytes = append(jsonBytes, comma...)

  //Description
  description, err := json.Marshal(ver.Description)
  if err != nil {
    glog.Error(err)
  }
  jsonBytes = append(jsonBytes, []byte("\"description\":")...)
  jsonBytes = append(jsonBytes, description...)
  jsonBytes = append(jsonBytes, comma...)

  //OperationID
  operationId, err := json.Marshal(ver.OperationID)
  if err != nil {
    glog.Error(err)
  }
  jsonBytes = append(jsonBytes, []byte("\"operationId\":")...)
  jsonBytes = append(jsonBytes, operationId...)
  jsonBytes = append(jsonBytes, comma...)

  //Consumes
  consumes, err := json.Marshal(ver.Consumes)
  if err != nil {
    glog.Error(err)
  }
  jsonBytes = append(jsonBytes, []byte("\"consumes\":")...)
  jsonBytes = append(jsonBytes, consumes...)
  jsonBytes = append(jsonBytes, comma...)

  //Produces
  produces, err := json.Marshal(ver.Produces)
  if err != nil {
    glog.Error(err)
  }
  jsonBytes = append(jsonBytes, []byte("\"produces\":")...)
  jsonBytes = append(jsonBytes, produces...)
  jsonBytes = append(jsonBytes, comma...)

  //Connection
  connection, err := json.Marshal(ver.Connection)
  if err != nil {
    glog.Error(err)
  }
  jsonBytes = append(jsonBytes, []byte("\"connection\":")...)
  jsonBytes = append(jsonBytes, connection...)
  jsonBytes = append(jsonBytes, comma...)

  //Parameters
  jsonBytes = append(jsonBytes, []byte("\"parameters\":[")...)
  parameters, err := json.Marshal(ver.Parameters)
  if err != nil {
    glog.Error(err)
  }
  parameters = bytes.Replace(parameters, []byte("\\\""), []byte("\""), -1)
  parameters = []byte(removeQuotes(string(parameters)))
  jsonBytes = append(jsonBytes, parameters...)

  jsonBytes = append(jsonBytes, []byte("]")...)
   

  jsonBytes = append(jsonBytes, []byte("}")...)
  glog.V(2).Info("json bytes for verb struct: %s", string(jsonBytes))
  return jsonBytes, nil
}

func (pss *PathsStruct) MarshalJSON() ([]byte, error) {
  var jsonBytes []byte
  comma := []byte(",")
  currentPath := ""

  jsonBytes = append(jsonBytes, []byte("{")...)

  for i := range pss.Paths {
    tag, err := json.Marshal(string(pss.Paths[i].Verbs[0].Tags[0]))
    operationId, err := json.Marshal(string(pss.Paths[i].Verbs[0].OperationID))
    if err != nil {
      glog.V(2).Infof("Unable to marshal tag or operationid: %s %s", string(pss.Paths[i].Verbs[0].Tags[0]), string(pss.Paths[i].Verbs[0].OperationID))
      return nil, err
    }

    glog.V(2).Info(string(jsonBytes))

    if bytes.HasPrefix(jsonBytes, []byte("{\"Paths\":[")) {
      jsonBytes = bytes.TrimPrefix(jsonBytes, []byte("{\"Paths\":["))
    }

    path, verb := getPathAndVerbFromJson(string(tag), string(operationId))

    if path == currentPath {
      jsonBytes = bytes.TrimSuffix(jsonBytes, comma)
      jsonBytes = bytes.TrimSuffix(jsonBytes, []byte("}"))
      jsonBytes = append(jsonBytes, comma...)
    } else {
      jsonBytes = append(jsonBytes, []byte("\"" + path + "\": {")...)
      currentPath = path
    }

    for j := range pss.Paths[i].Verbs {
      jsonBytes = append(jsonBytes, []byte("\"" + verb + "\":")...)
      bytes, err := json.Marshal(pss.Paths[i].Verbs[j])
      if err != nil {
        glog.V(2).Infof("Unable to marshal verb: %s", pss.Paths[i].Verbs[j])
        return nil, err
      }
      jsonBytes = append(jsonBytes, bytes...)
      jsonBytes = append(jsonBytes, comma...)
      glog.V(2).Infof("json bytes after path %s: %s", path, string(jsonBytes))
    }
    jsonBytes = bytes.TrimSuffix(jsonBytes, comma)
    jsonBytes = append(jsonBytes, []byte("}")...)
    jsonBytes = append(jsonBytes, comma...)
  }

  jsonBytes = bytes.TrimSuffix(jsonBytes, comma)
  jsonBytes = append(jsonBytes, []byte("}")...)
  glog.V(2).Info("json bytes after all paths: %s", string(jsonBytes))
  return jsonBytes, nil
}

func (def *DefinitionStruct) MarshalJSON() ([]byte, error) {
  var jsonBytes []byte = []byte("{")

  for i := range def.Definitions {

    properties := def.Definitions[i].Properties
    //strip the escape characters from properties
    //properties = strings.Replace(properties, "\\", "", -1)


    xml, err := json.Marshal(def.Definitions[i].Xml)
    if err != nil {
      return nil, err
    }


    properties = strings.Replace(properties, "\"properties\":{", "\"type\":\"" + string(def.Definitions[i].Type) + "\",\"xml\":" + string(xml) + "," +  "\"properties\":{", 1)
    properties = strings.TrimPrefix(properties, "{")
    properties = strings.TrimSuffix(properties, "}")
    glog.Info("Properties are: " + properties)
    jsonBytes = append(jsonBytes, []byte(properties)...)
    jsonBytes = append(jsonBytes, []byte(",")...)
  }

  jsonBytes = bytes.TrimSuffix(jsonBytes, []byte(","))
  jsonBytes = append(jsonBytes, []byte("}")...)

  glog.V(2).Info("json from definiaton structs is " +  string(jsonBytes))

  if string(jsonBytes) == "" {
    jsonBytes = []byte("{}")
  }
  return jsonBytes, nil
}

/**func (swag *SwaggerTemplate) UnmarshalJSON(b []byte) (err error) {
  var entries map[string]json.RawMessage
  json.Unmarshal(b, &entries)
  
  glog.V(1)Info(entries["swagger"])
  
  swag.Swagger = removeQuotes(string(entries["swagger"]))
  var info InfoStruct
  json.Unmarshal(entries["info"], &info)
  swag.Info = info

  swag.Host = removeQuotes(string(entries["host"]))
  swag.BasePath = removeQuotes(string(entries["basePath"]))

  var tags []TagStruct
  json.Unmarshal(entries["tags"], &tags)
  swag.Tags = tags

  var schemes []string
  json.Unmarshal(entries["schemes"], &schemes)
  swag.Schemes = schemes

  glog.V(1).Info("Paths in swagger unmarsahl is: " + string(entries["paths"]))

  paths := &PathsStruct{}
  json.Unmarshal(entries["paths"], &paths)

  glog.V(1).Info("Paths struct is: " + paths.Paths[0].Verbs[0].Summary)
  swag.Paths = *paths

  var definitions DefinitionStruct
  json.Unmarshal(entries["definitions"], &definitions)
  swag.Definitions = definitions

  return
}**/

func (pss *PathsStruct) UnmarshalJSON(b []byte) (err error) {
  s := strings.Trim(string(b), "\"")
  if s == "null" {
    glog.Warning("bytes are null, something isn't right")
  }

  glog.V(2).Info(s)
  var paths map[string]json.RawMessage
  json.Unmarshal(b, &paths)

  for i := range paths {
    glog.V(2).Info(i)
    //glog.V(2).Info(paths[i])
    var path PathStruct
    err := json.Unmarshal(paths[i], &path)
    if err != nil {
      glog.V(2).Info("Error unmarshalling " + i + ":", err.Error)
    }
    glog.V(2).Info("Path " + i + " is: " + toJson(path.Verbs))

    pss.Paths = append(pss.Paths, path)
    glog.V(2).Info(pss.Paths)
    glog.V(2).Info("Paths JSON is : " + toJson(pss))
  }

  return
}

func (ps *PathStruct) UnmarshalJSON(b []byte) (err error) {
  s := strings.Trim(string(b), "\"")
  if s == "null" {
    glog.Warning("bytes are null, something isn't right")
  }

  glog.V(2).Info(s)
  var path map[string]json.RawMessage
  json.Unmarshal(b, &path)

  for i := range path {
    glog.V(2).Info(i)
    //glog.V(2).Info(path[i])
    var verb VerbStruct
    json.Unmarshal(path[i], &verb)

    glog.V(2).Info("Verb " + i + " is: " + toJson(verb))

    ps.Verbs = append(ps.Verbs, verb)
    glog.V(2).Info("Path.verb is: " + toJson(ps.Verbs))
    glog.V(2).Info("Path is: " + toJson(ps))
  }

  return
}

func (defs *DefinitionStruct) UnmarshalJSON(b []byte) (err error) {
   s := strings.Trim(string(b), "\"")
  if s == "null" {
    glog.Warning("bytes are null, something isn't right")
  }

  glog.V(2).Info(s)
  var definitions map[string]json.RawMessage
  json.Unmarshal(b, &definitions)
 
  for i := range definitions {
    glog.V(2).Info(i)
    var definition Definition
    json.Unmarshal(definitions[i], &definition)

    glog.V(2).Info("Definition " + i + " is: " + toJson(definition))

    definition.Properties = "{\"" + i + "\":" + definition.Properties + "}"

    defs.Definitions = append(defs.Definitions, definition)
    glog.V(2).Info("Defs.Definition is: " + toJson(defs.Definitions))
    glog.Flush()
    glog.V(2).Info("Definitions is: " + toJson(defs))
  }

  return
}

func (ver *VerbStruct) UnmarshalJSON(b []byte) (err error) {
  s := strings.Trim(string(b), "\"")
  if s == "null" {
    glog.Warning("bytes are null, something isn't right")
  }

  glog.V(2).Info(s)
  var verb map[string]json.RawMessage
  json.Unmarshal(b, &verb)

  var tags []string
  json.Unmarshal(verb["tags"], &tags)  
  //ver.Tags = removeQuotes(string(verb["tags"]))
  ver.Tags = tags
  ver.Summary = removeQuotes(string(verb["summary"]))
  ver.Description = removeQuotes(string(verb["description"]))
  ver.OperationID = removeQuotes(string(verb["operationId"]))

  var consumes []string
  json.Unmarshal(verb["consumes"], &consumes)
  //ver.Consumes = removeQuotes(string(verb["consumes"]))
  ver.Consumes = consumes

  var produces []string
  json.Unmarshal(verb["produces"], &produces)
  //ver.Produces = removeQuotes(string(verb["produces"]))
  ver.Produces = produces
  
  var conn ConnectionStruct
  json.Unmarshal(verb["connection"], &conn)

  ver.Connection = conn

  glog.V(1).Info("Parameters from bytes in unmarshal: " + string(verb["parameters"]))

  params := removeWhiteSpace(string(verb["parameters"]))

  params = removeQuotes(params) 

  params = strings.TrimPrefix(params, "[")

  ver.Parameters = strings.TrimSuffix(params, "]")

  glog.V(1).Info("Parameters from verbStruct in unmarshal: " + ver.Parameters)

  return
}

func (def *Definition) UnmarshalJSON(b []byte) (err error) {
  s := strings.Trim(string(b), "\"")
  if s == "null" {
    glog.Warning("bytes are null, something isn't right")
  }

  glog.V(2).Info(s)
  var definition map[string]json.RawMessage
  json.Unmarshal(b, &definition)

  for i := range definition {
    glog.V(2).Info("Definition unmarshal: " + i)
    switch i {
    case "type":
      def.Type = removeQuotes(string(definition[i]))
    case "properties":
      glog.V(2).Info(removeQuotes(string(definition[i])))
      def.Properties = "{\"properties\":" + removeQuotes(string(definition[i])) + "}" 
      glog.V(2).Info(def.Properties)
    case "xml":
      glog.V(2).Info("xml is: " +removeQuotes(string(definition[i])))
      //create mapping here, if more fields are needed in xml then consider using a new struct
      var xmlTemp map[string]json.RawMessage
      json.Unmarshal(definition[i], &xmlTemp)

      for j := range xmlTemp {
        glog.V(2).Info("xmlTemp unmarshal: " + j)
        switch j {
        case "name":
          def.Xml.Name = removeQuotes(string(xmlTemp[j])) 
        }
      }
    }
  }

  return
}


func getPathFromJson(pathJson []byte) (paths []string) {
  pathString := string(pathJson)

  glog.V(1).Info(pathString) 

  return paths
}

func getPathAndVerbFromJson(tag string, operationID string) (path string, verb string) {
  glog.V(1).Info(operationID)
  context, verb := getContextAndVerb(removeQuotes(string(operationID)), removeQuotes(string(tag)))
  path = "/" + removeQuotes(string(tag)) + context
  return path, verb
}

func getContextAndVerb(operationID string, tag string) (context string, verb string) {

  switch true {
    case strings.HasPrefix(operationID, "add"):
      context = strings.TrimPrefix(operationID, "add")
      verb = "post"
    case strings.HasPrefix(operationID, "post"):
      context = strings.TrimPrefix(operationID, "post")
      verb = "post"
    case strings.HasPrefix(context, "upload"):
      context = strings.TrimPrefix(operationID, "upload")
      verb = "post"
    case strings.HasPrefix(operationID, "update"):
      context = strings.TrimPrefix(operationID, "update")
      verb = "put"
    case strings.HasPrefix(operationID, "get"):
      context = strings.TrimPrefix(operationID, "get")
      verb = "get"
    case strings.HasPrefix(operationID, "find"):
      context = operationID
      verb = "get"
    default:
      context = operationID
      verb = "get"
    }
    glog.V(1).Info("Verb is : " + verb)
    glog.V(1).Info("Tag is : " + tag)

    context = strings.Replace(context, strings.Title(tag), "", 1)
    context = strings.Replace(context, " ", "_", -1)
    glog.V(1).Info("Context is : " + context)
    if context != "" {
      context = "/" + context
    }

  return context, verb
}

func getSwaggerTemplate() SwaggerTemplate {
  raw, err := ioutil.ReadFile("./Swagger_UI_Poc.json")
  if err != nil {
    glog.Fatal(err.Error())
  }

  var c SwaggerTemplate
  json.Unmarshal(raw, &c)
  return c
}

func removeQuotes(quotedString string) (unquotedString string) {
  unquotedString = strings.TrimPrefix(quotedString, "\"")
  return strings.TrimSuffix(unquotedString, "\"")
}

func removeWhiteSpace(stringWithSpaces string) (stringWithoutSpaces string) {
  stringWithoutSpaces = strings.Replace(stringWithSpaces, "\n", "", -1)
  stringWithoutSpaces = strings.Replace(stringWithoutSpaces, "\r", "", -1)
  stringWithoutSpaces = strings.Replace(stringWithoutSpaces, "\t", "", -1)
  return stringWithoutSpaces
}
