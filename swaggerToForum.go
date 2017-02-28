package main

import (
  "fmt"
  "github.com/golang/glog"
  "io/ioutil"
  "strings"
  "os"
  "encoding/json"
)

var defaults map[string]interface{}
var config Config

func SwaggerToForum(swagger SwaggerTemplate, conf Config) {
  config = conf
  LoadDefaults()

  fmt.Println(defaults)

  host := swagger.Host
  basePath := swagger.BasePath
  projectName := strings.ToLower(strings.Replace(swagger.Info.Title, " ", "_", -1))

  for i := range swagger.Definitions.Definitions {
    var defMap map[string]json.RawMessage
    json.Unmarshal([]byte(swagger.Definitions.Definitions[i].Properties), &defMap)
    var defName string

    for key, _ := range defMap {
      defName = key
    } 
    
    CreateDefinitionDocument(projectName, swagger.Definitions.Definitions[i], defName)
  }

  for i := range swagger.Paths.Paths {
    for j := range swagger.Paths.Paths[i].Verbs {
      glog.V(1).Info("Verb in swaggerToForum is: " + toJson(swagger.Paths.Paths[i].Verbs))
      //need to check all of these to see if they exist before creating them, an then update them instead
      listenerName := CreateListener(projectName, swagger.Paths.Paths[i].Verbs[j])
      remoteName := CreateRemote(projectName, swagger.Paths.Paths[i].Verbs[j], host)
      contentPolicyName := CreateContentPolicy(projectName, swagger.Paths.Paths[i].Verbs[j].Tags[0], swagger.Tags, listenerName)
      CreateVirtualDirectory(projectName, swagger.Paths.Paths[i].Verbs[j], basePath, listenerName, remoteName, contentPolicyName)
      //need to check if the vd is there before deleting it.
      callRest("policies/jsonPolicies/" + contentPolicyName + "/virtualDirectories/New Virtual Directory", "DELETE", "", config)
    }
  }

}

func CreateParametersDocument(projectName string, parameters string, opId string) {
  escapedParameters := strings.Replace(toJson(parameters), "\"", "\\\"", -1)

  jsonString := "{\"name\":\"" + projectName + "_document_policy_parameters_" + opId + "\",\"document\":[\""  + escapedParameters + "\"]}"

  glog.V(1).Info("jsonString of parameter doc before post is: " + jsonString)

  callRest("policies/documents", "POST", jsonString, config)
}

func CreateDefinitionDocument(projectName string, definition Definition, opId string) {
  escapedDefinition := strings.Replace(definition.Properties, "\"", "\\\"", -1)

  jsonString := "{\"name\":\"" + projectName + "_document_policy_definition_" + opId + "\",\"document\":[\""  + escapedDefinition + "\"]}"

  callRest("policies/documents", "POST", jsonString, config)
}

func CreateVirtualDirectory(projectName string, verb VerbStruct, basePath string, listenerName string, remoteName string, contentPolicyName string) {
  var vd virtualDirectory
  vd = defaults["virtual_directory"].(virtualDirectory)

  opId := strings.Replace(verb.OperationID, " ", "_", -1)

  vd.Name = projectName + "_virtual_directory_" + opId 
  vd.VirtualPath, _ = getPathAndVerbFromJson(verb.Tags[0], opId)
  if config.ContextAllLowerCase {
    vd.VirtualPath = strings.ToLower(vd.VirtualPath)
  }
  vd.RemotePath = vd.VirtualPath
  vd.ListenerPolicy = listenerName
  vd.RemotePolicy =  remoteName
  vd.Description = verb.Description
  callRest("policies/jsonPolicies/" + contentPolicyName + "/virtualDirectories", "POST", toJson(vd), config)
  
  glog.V(2).Info("verb.Parameters is: " + verb.Parameters)
  CreateParametersDocument(projectName, verb.Parameters, opId)

  fmt.Println(toJson(vd))
}

func CreateContentPolicy(projectName string, tag string, tags []TagStruct, listenerName string) string {
  var content contentPolicy
  content = defaults["json_content_policy"].(contentPolicy)

  content.Name = projectName + "_json_content_policy_" + tag
  content.ListenerPolicy = listenerName
  for _, thisTag := range tags {
    if thisTag.Name == tag {
      content.Description = thisTag.Description
    }
  }

  callRest("policies/jsonPolicies", "POST", toJson(content), config)

  fmt.Println(toJson(content))

  return content.Name
}

func CreateListener(projectName string, verb VerbStruct) string {
  if verb.Connection.Type == "" {
    verb.Connection.Type = "http"
  }

  tag := strings.ToLower(strings.Replace(verb.Tags[0], " ", "_", -1))
  switch verb.Connection.Type {
  case "http":
    var listener httpListener
    listener = defaults["http_listener"].(httpListener)
    
    if verb.Connection.Port != 0 {
      	listener.Port = verb.Connection.Port
    }
    listener.Name = projectName + "_http_listener_policy_" + tag
  
    callRest("policies/httpListenerPolicies", "POST", toJson(listener), config)
    
    fmt.Println(toJson(listener))
    return listener.Name
  case "amqp10":
    var listener amqp10Listener
    listener = defaults["amqp10_listener"].(amqp10Listener)

    if verb.Connection.Port != 0 {
        listener.Port = verb.Connection.Port
    }
    listener.Name = projectName + "_amqp10_listener_policy_" + tag

    callRest("policies/amqp10ListenerPolicies", "POST", toJson(listener), config)

    fmt.Println(toJson(listener))
    return listener.Name
  }

  return ""
}

func CreateRemote(projectName string, verb VerbStruct, host string) string {
  if verb.Connection.RemoteType == "" {
    verb.Connection.RemoteType = "http"
  }

  tag := strings.ToLower(strings.Replace(verb.Tags[0], " ", "_", -1))

  switch verb.Connection.RemoteType {
  case "http":
    var remote httpRemote
    remote = defaults["http_remote"].(httpRemote)
    
    if verb.Connection.Port != 0 {
        remote.RemotePort = verb.Connection.Port
    }
    remote.Name = projectName + "_http_remote_policy_" + tag
    remote.RemoteServer = host
  
    callRest("policies/httpRemotePolicies", "POST", toJson(remote), config)
    
    fmt.Println(toJson(remote))
    return remote.Name
  case "amqp10":
    var remote amqp10Remote
    remote = defaults["amqp10_remote"].(amqp10Remote)

    if verb.Connection.Port != 0 {
        remote.RemotePort = verb.Connection.Port
    }
    remote.Name = projectName + "_amqp10_remote_policy_" + tag
    remote.RemoteServer = host

    callRest("policies/amqp10RemotePolicies", "POST", toJson(remote), config)

    fmt.Println(toJson(remote))
    return remote.Name
  }

  return ""
}

func LoadDefaults() {
  files, err := ioutil.ReadDir("./defaults")
  if err != nil {
    glog.Error(err)
  }

  defaults = make(map[string]interface{})

  for _, file := range files {
    fileName := file.Name()
    fmt.Println(fileName)
    if strings.HasSuffix(fileName, ".json") {
      switch true {
      case strings.Contains(fileName, "listener"):
        loadListenerDefaults(file)
      case strings.Contains(fileName, "remote"):
        loadRemoteDefaults(file)
      case strings.Contains(fileName, "policy"):
        loadContentPolicyDefaults(file)
      case strings.Contains(fileName, "virtual_directory"):
        loadVirtualDirectoryDefaults(file)
      default:
        glog.Warning(fileName + " does not match a known policy type")
      }
    } else {
      glog.Warning("Defaults files should be json")
    }
  }
}

func loadListenerDefaults(file os.FileInfo) {
  raw, err := ioutil.ReadFile("./defaults/" + file.Name())
  if err != nil {
    glog.Fatal(err)
  }

  fileName := file.Name()
  
  switch true {
  case strings.Contains(fileName, "http"):
    var def httpListener
    json.Unmarshal(raw, &def)
    defaults["http_listener"] = def
  case strings.Contains(fileName, "amqp10"):
    var def amqp10Listener
    json.Unmarshal(raw, &def)
    defaults["amqp10_listener"] = def
  default:
    glog.Warning("Unknown listener type")
  }
}

func loadRemoteDefaults(file os.FileInfo) {
  raw, err := ioutil.ReadFile("./defaults/" + file.Name())
  if err != nil {
    glog.Fatal(err)
  }

  fileName := file.Name()

  switch true {
  case strings.Contains(fileName, "http"):
    var def httpRemote
    json.Unmarshal(raw, &def)
    defaults["http_remote"] = def
  case strings.Contains(fileName, "amqp10"):
    var def amqp10Remote
    json.Unmarshal(raw, &def)
    defaults["amqp10_remote"] = def
  default:
    glog.Warning("Unknown remote type")
  }
}

func loadContentPolicyDefaults(file os.FileInfo) {
  raw, err := ioutil.ReadFile("./defaults/" + file.Name())
  if err != nil {
    glog.Fatal(err)
  }

  fileName := file.Name()

  var def contentPolicy
  err = json.Unmarshal(raw, &def)
  if err != nil {
    glog.Warning("Unable to parse default file: " + fileName)
    glog.Error(err)
  }

  switch true {
  case strings.HasPrefix(fileName, "json"):
    defaults["json_content_policy"] = def
  case strings.HasPrefix(fileName, "xml"):
    defaults["xml_content_policy"] = def
  case strings.HasPrefix(fileName, "html"):
    defaults["html_content_policy"] = def
  }
}

func loadVirtualDirectoryDefaults(file os.FileInfo) {
  raw, err := ioutil.ReadFile("./defaults/" + file.Name())
  if err != nil {
    glog.Fatal(err)
  }

  var def virtualDirectory
  json.Unmarshal(raw, &def)
  defaults["virtual_directory"] = def

}
