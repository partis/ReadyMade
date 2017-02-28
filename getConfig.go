package main

import (
  "encoding/json"
  "io/ioutil"
  "github.com/golang/glog"
  "fmt"
)

type Config struct {
  SwaggerFile string `json:"swaggerFile"`
  ForumURL string `json:"forumURL"`
  ForumUsername string `json:"forumUsername"`
  ForumPass string `json:"forumPass"`
  ContextAllLowerCase bool `json:"contextAllLowerCase"`
}

func ReadConfig(configfile string) Config {
  raw, err := ioutil.ReadFile(configfile)
  if err != nil {
    glog.Fatal(err)
  }

  var config Config
  json.Unmarshal(raw, &config)

  glog.Info(config)
  return config
}

func DisplayConfig() {
  fmt.Println(toJson(config))
}
