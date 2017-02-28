package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "github.com/golang/glog"
)

func ReadSwagger(swaggerFile string) SwaggerTemplate {
  raw, err := ioutil.ReadFile(swaggerFile)
  if err != nil {
    log.Fatal(err)
  }

  _, isJson := isJson(string(raw))

  if isJson {
    var swagger SwaggerTemplate
    json.Unmarshal(raw, &swagger)

    glog.V(2).Info(swagger)
    return swagger
  } else {
    glog.Fatal("You swagger file is not valid JSON")
    return SwaggerTemplate{}
  }
}
