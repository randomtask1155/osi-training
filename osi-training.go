package main

import(
  "fmt"
  "flag"
  "osi-training/supportdemoip"
  "osi-training/supportdemotcp"
  "osi-training/supportdemossl"
  "osi-training/supportdemohttp"
  //"osi-training/options"
)

var (
  // TestType default test type is IP
  TestType = flag.String("t", "ip", "[ ip | tcp | ssl | http ]")
  
  // Tests is a map of test names to their function
  Tests = map[string]func(){
      "ip": supportdemoip.Run, 
      "tcp": supportdemotcp.Run,
      "ssl": supportdemossl.Run,
      "http": supportdemohttp.Run}
)

func runTCP(){
  fmt.Printf("running tcp tests")
}

func main(){
  flag.Parse()
  Tests[*TestType]()
}
