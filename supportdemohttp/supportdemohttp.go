package supportdemohttp

import(
  "fmt"
  "net/http"
  "time"
  "net"
  "osi-training/supportdemoip"
  "osi-training/colors"
)

var (
  //HTTPServerPort listening port for http server
  HTTPServerPort = "8081"
)

// handleServerTimeout sleeps for 2 minutes before responding
func handleServerTimeout(w http.ResponseWriter, r *http.Request) {
  time.Sleep(120 * time.Second)
  w.Write([]byte("<html>Hello</html>"))
}

func runHTTPServer(){
  http.HandleFunc("/sleep120", handleServerTimeout) 
  err := http.ListenAndServe(":"+HTTPServerPort, nil)
  if err != nil {
     panic(fmt.Sprintf("Could not listen on port %s: %s", HTTPServerPort, err))
  }
}

func myDialFunc(network, addr string) (net.Conn, error){
  return net.DialTimeout(network, addr, 1 * time.Second) 
}

// Run performs the http tests
func Run() {
  go runHTTPServer()
  supportdemoip.FindBadIP()
  badip := supportdemoip.IPAddress
  
  // dial timeout test 
  fmt.Printf("Performing dial timeout test to http://%s:%s: ", badip, HTTPServerPort)
  
  tr := &http.Transport{ Dial: myDialFunc} // this sets the request timeout for an http client
  client := &http.Client{Transport: tr}
  _, err := client.Get(fmt.Sprintf("http://%s:%s", badip, HTTPServerPort))
  if err !=nil {
    fmt.Printf("%s\n", colors.ColorizeError(err))
  } else {
    fmt.Printf("\nOops that was not supposed to work! Try changing the bad host range\n")
  }
  fmt.Printf(colors.ColorizeText(fmt.Sprintf(`
Dial timeout test is similar to the TCP dial timeout test. Basically we create 
and HTTP client that which will timeout if the TCP session is not established 
within 1 seconds as per "myDialFunc()"

`)))
  // request timeout test 
  fmt.Printf("Performing request timeout test to http://%s:%s: ", badip, HTTPServerPort)
  client = &http.Client{Timeout: 3 * time.Second}// timeout sets the request timeout
  _, err = client.Get(fmt.Sprintf("http://127.0.0.1:%s/sleep120", HTTPServerPort))
  if err != nil {
    fmt.Printf("%s\n", colors.ColorizeError(err))
  } else {
    fmt.Printf("\nOops that was not supposed to work!\n")
  }
  fmt.Printf(colors.ColorizeText(fmt.Sprintf(`
A request timeout can ocurr when the client sends an HTTP request and expects 
the server to respond within a certain amount of time.  In this case we 
configured the client to timeout after 3 seconds and configured the server to 
sleep for 120 seconds before responding.  That way a timeout error is always 
produced
 
`)))
  
}