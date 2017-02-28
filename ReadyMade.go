package main

import (
  "fmt"
  "github.com/golang/glog"
  "os"
  "flag"
)

func usage() {
  //print usage to iniate logging
  fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string]\n", )
  flag.PrintDefaults()
  os.Exit(2)
}

func init() {
  //set the usage to the above func
  flag.Usage = usage
  //parse the flags from the command line to configre logging
  flag.Parse()
}

func main() {
  glog.Info("Reading in swagger file")

  config := ReadConfig("ReadyMade.cfg")
  fmt.Println(toJson(config))

  var fileName string

  if config.SwaggerFile != "" {
    fileName = config.SwaggerFile
  } else {
    fileName = "ReadyMade.json"
  }
  swagger := ReadSwagger(fileName)
  glog.V(2).Info(swaggerToJson(&swagger))

  SwaggerToForum(swagger, config)

  glog.Flush()
}
