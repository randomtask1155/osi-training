package supportdemotcp 

import(
  "net"
  "fmt"
  "time"
  "osi-training/colors"
  "osi-training/supportdemoip"
)

var (
  
  // ServerPort port used by tcpserver
  ServerPort = "8081"
  
  // BadServerPort is where the tcp server will not be listening on
  BadServerPort = "8082"
)




func startTCPServer(laddr *net.TCPAddr){
  
  ln, err := net.ListenTCP("tcp", laddr)
  if err != nil {
    panic(fmt.Sprintf("Could not start server on port %s: %s", ServerPort, err))
  }
  for {
    conn, err := ln.AcceptTCP()
    if err != nil {
        fmt.Printf("TCP SERVER:CONN ERROR: %s\n", err)
    }
    conn.Close()
  }
}

// Run performs the tcp tests
func Run() {
  supportdemoip.FindBadIP()
  badip := supportdemoip.IPAddress
  localAddress := fmt.Sprintf("127.0.0.1:%s", ServerPort)
  localBadAddress := fmt.Sprintf("127.0.0.1:%s", BadServerPort)
  
  tcpAddr, err := net.ResolveTCPAddr("tcp", localAddress)
  if err != nil {
    panic(fmt.Sprintf("Could not resolve tcp address: %s", err))
  }
  
  tcpBadAddr, err := net.ResolveTCPAddr("tcp", localBadAddress)
  if err != nil {
    panic(fmt.Sprintf("Could not resolve tcp address: %s", err))
  }
  
  tcpTimeoutAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", badip, BadServerPort))
  if err != nil {
    panic(fmt.Sprintf("Could not resolve tcp address: %s", err))
  }
  
  go startTCPServer(tcpAddr)
  time.Sleep(1*time.Second) // give time for the server to start
  
  // connection refused test 
  fmt.Printf("Performing TCP CONNECTION REFUSED test to %s: ", tcpBadAddr)
  _, err = net.DialTCP("tcp", nil,  tcpBadAddr)
  if err != nil {
    fmt.Printf("%s\n", colors.ColorizeError(err))
  } else {
    fmt.Printf("\nOops that was not supposed to work! Try changing the bad port\n")
  }
  fmt.Printf(colors.ColorizeText(fmt.Sprintf(`
Connection refused usually happens when the remote server is not listening on 
the given server port.  For example in this case we are connecting to %s
and host 127.0.0.1 is not not listening on port %s

The TCP communication will look like this

IP 127.0.0.1.63329 > %[1]s: Flags [S.]
IP %[1]s < 127.0.0.1.63329 : Flags [R.]

NOTE: "S" stands for TCP SYN Flag and "R" Stands for TCP RESET Flag

`, tcpBadAddr, BadServerPort)))
  
  // Borken Pipe Error Test
  fmt.Printf("Performing BROKEN PIPE Error test to %s: ", tcpAddr)
  conn, err := net.DialTCP("tcp", nil,  tcpAddr)
  if err != nil {
    fmt.Printf("%s\n", colors.ColorizeError(err))
  }
  time.Sleep(1* time.Second) // give time for session to be created and closed
  conn.Write([]byte("hello")) // go is buggy here as it should have errored
  time.Sleep(1* time.Second) // give time for session to be created and closed
  _, err = conn.Write([]byte("hello"))
  if err != nil {
    fmt.Printf("%s\n", colors.ColorizeError(err))
  }
  fmt.Printf(colors.ColorizeText(fmt.Sprintf(`
Broken Pipe error occurred becaused the client attempted to write data to the 
TCP socket that was already closed

The TCP communication looks like this 

TCP handshake completes meaning Session is Esablished
IP 127.0.0.1.63330 > %[1]s: Flags [S.]
IP %[1]s < 127.0.0.1.63330 : Flags [SA.]
IP 127.0.0.1.63330 > %[1]s: Flags [A.]

Server closes TCP session and client acknowleded
IP %[1]s < 127.0.0.1.63330 : Flags [F.]
IP 127.0.0.1.63330 > %[1]s: Flags [FA.]
IP %[1]s < 127.0.0.1.63330 : Flags [A.]

But application tries to send data on the closed TCP socket.  Server responsds 
with TCP RESET
IP 127.0.0.1.63330 > %[1]s: Flags [P.]
IP %[1]s < 127.0.0.1.63330 : Flags [R.]

NOTE: 
"A" stand fro TCP ACK Flag which means receiver has acknowleded receipt of 
packet or packetsand 
"P" is for TCP PUSH flag which means source is sending data
"F" is for TCP FIN flag which means sender would like to close the TCP session

`, tcpAddr )))
  
  fmt.Printf("Performing TCP RESET test to %s: ", tcpAddr)
  b := make([]byte,1024)
  _, err = conn.Read(b)
  if err != nil {
    fmt.Printf("%s\n", colors.ColorizeError(err))
  }
  fmt.Printf(colors.ColorizeText(fmt.Sprintf(`
TCP RESET error is observed when an application attempts to read from a 
TCP socket that recieved a TCP RESET flag. 
For example the application was previously reading and writing from a socket and
all of sudden closed aburptly

The way we reproduce this error is via the following TCP communication 

TCP handshake completes meaning Session is Esablished
IP 127.0.0.1.63330 > %[1]s: Flags [S.]
IP %[1]s < 127.0.0.1.63330 : Flags [SA.]
IP 127.0.0.1.63330 > %[1]s: Flags [A.]

Server closes TCP session and client acknowleded
IP %[1]s < 127.0.0.1.63330 : Flags [F.]
IP 127.0.0.1.63330 > %[1]s: Flags [FA.]
IP %[1]s < 127.0.0.1.63330 : Flags [A.]

But application tries to send data on the closed TCP socket.  Server responsds 
with TCP RESET. 
IP 127.0.0.1.63330 > %[1]s: Flags [P.]
IP %[1]s < 127.0.0.1.63330 : Flags [R.]

For purpose of this example we ignore this write failure and proceed to read 
from the socket anway and go returns the connection reset error


`, tcpAddr )))
  // dial timeout test
  fmt.Printf("Performing TCP TIMEOUT test to %s: ", tcpTimeoutAddr)
  _, err = net.DialTCP("tcp", nil,  tcpTimeoutAddr)
  if err != nil {
    fmt.Printf("%s\n", colors.ColorizeError(err))
  } else {
    fmt.Printf("\nOpps that was suppose to fail! try a different bad ip range")
  }
  fmt.Printf(colors.ColorizeText(fmt.Sprintf(`
In this example we reproduce a socket timeout error by attempting to connect
to a non-routable network and port.  This essentially simulates packet loss 

TCP communication will look like this.  The client will continue to send SYN 
packets until the OS kernel times out the TCP session. 
IP 127.0.0.1.63331 > %[1]s: Flags [S.]
IP 127.0.0.1.63331 > %[1]s: Flags [S.]
IP 127.0.0.1.63331 > %[1]s: Flags [S.]
IP 127.0.0.1.63331 > %[1]s: Flags [S.]
IP 127.0.0.1.63331 > %[1]s: Flags [S.]
IP 127.0.0.1.63331 > %[1]s: Flags [S.]

`, tcpTimeoutAddr )))
}
