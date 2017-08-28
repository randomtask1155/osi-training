package supportdemoip


import (
  "net"
  "fmt"
  "osi-training/colors"
)

var (
  
  // HostName is a Dummy hostname 
  HostName = "ivegotalovelybunchofcoconuts.diddily.dee"
  
  // IPAddress is Dummy ipaddress
  IPAddress = "192.123.123.1"
)

func FindBadIP(){
  ip := net.ParseIP(IPAddress)
  ip = ip.To4()
  for {
      if ip[3] > 255 {
        fmt.Printf("ERROR: Could not find a bad IP to test with\nTry a different range\n")
      }
      _, err := net.LookupAddr(IPAddress)
      if err == nil {
        ip[3]++
        IPAddress = fmt.Sprintf("%s", ip)
        continue
      }
      break
  }
}

// Run performs the ip tests
func Run() {
  FindBadIP()
  
  fmt.Printf("Perfroming a forward DNS Lookup of hostname %s: ", HostName)
  _, err := net.LookupHost(HostName)
  if err != nil {
    fmt.Printf("%s\n", colors.ColorizeError(err))
  }else {
    fmt.Printf("\nOpps hostname lookup was sucessful!\nTry a different hostname\n")
  }
  
  fmt.Printf("Perfroming a reverse DNS Lookup of ip address %s: ", IPAddress)
  _, err =  net.LookupAddr(IPAddress)
  if err != nil {
    fmt.Printf("%s\n", colors.ColorizeError(err))
  } else {
    fmt.Printf("\nOpps this was supposed to failed but it passed!!\nTry a different ip address?\n")
    return
  }
  
  fmt.Printf(colors.ColorizeText(fmt.Sprintf(`
Other Types of IP related errors include the following and they generally mean
there is a problem in the hosts routing table or a layer 3 issue where client
network can not route to server network
  * no route to host
  * destination host unreachable
    `)))

}