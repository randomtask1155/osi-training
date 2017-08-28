package supportdemossl

import(
  "fmt"
  "net/http"
  "crypto/tls"
  "time"
  //"osi-training/colors"
)

var (
  // HTTPServerPort listeng port for http server
  HTTPServerPort = "8081"
  
  rootCA = `-----BEGIN CERTIFICATE-----
MIID9zCCAt+gAwIBAgIJAM2Us+7ONCpKMA0GCSqGSIb3DQEBCwUAMFoxCzAJBgNV
BAYTAlVTMQswCQYDVQQIEwJJTDEQMA4GA1UEBxMHQ2hpY2FnbzEMMAoGA1UEChMD
VVBSMQwwCgYDVQQLEwNFbmcxEDAOBgNVBAMTB3Vwci5sbGMwHhcNMTcwNTE2MDAz
NDM3WhcNMjAwMzA1MDAzNDM3WjBaMQswCQYDVQQGEwJVUzELMAkGA1UECBMCSUwx
EDAOBgNVBAcTB0NoaWNhZ28xDDAKBgNVBAoTA1VQUjEMMAoGA1UECxMDRW5nMRAw
DgYDVQQDEwd1cHIubGxjMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
xAUoqGbGBjp9FcOLSkCMc5CBLJ1kqGZ6y7YgwrSMwviBGnrSu0M7lPYf1THBUUjr
LZZTRLVscnfEx1K06prurjYrovijZNj0jmMZFwGzU1ulAerMwnHQ42PjYBKHz7Zn
xnH/aftBGyjX2iu2jI/xpATIDbNk5WSE+TkZJoY7E7our6VVsD3EP+SYNZ7gGNw9
4XqNXlLjcsNDJLopHAEeItyQBHe71ZLotTBSDQvq6NxA0AOypB9WgZrBv3lSW7wS
Ddmzehzag0ZdEyzl9wXGtiQjCwVg3p+GDqo6s+JTXK8OWKLABsUowT7ubzG+fzhk
sVJYaBdE903ZiK5cnjoCtwIDAQABo4G/MIG8MB0GA1UdDgQWBBSur1N8K1Acm7kO
2uw9/zLlCfoB7zCBjAYDVR0jBIGEMIGBgBSur1N8K1Acm7kO2uw9/zLlCfoB76Fe
pFwwWjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAklMMRAwDgYDVQQHEwdDaGljYWdv
MQwwCgYDVQQKEwNVUFIxDDAKBgNVBAsTA0VuZzEQMA4GA1UEAxMHdXByLmxsY4IJ
AM2Us+7ONCpKMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAFoQkQlL
1xCKBMZO38AYkUW6wPdybgS6rnPXSeYUxbIuYZ3UFhSYG7GgVolumU7FJltye4ZY
OSevRCPCdUI3PgCZzInnohXsOhhg8fniJJEFj6s3zyEuhGmut79LgjPURrN87f/o
nTUIbg6a7J3pKxGv47W73GdetsYhcjLnzEkuyK7FLPggQPS/vtRrMzqvEic6KGOJ
1d2AX42dGaJO2iSY7P2Z4sQXUl4v5655qNgLeDe7i+BupSHX1HDIeendjRZc0DfJ
0ptGQnuWMpu00Gk2JWVT25ehfjZ2dAJFt+6qi82VrzXZvLmkUm10OqdS4PIx52M5
CqfuiOLUvw4prUk=
-----END CERTIFICATE-----`

  serverCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDNzCCAh8CCQC6+1HnfWbvijANBgkqhkiG9w0BAQsFADBaMQswCQYDVQQGEwJV
UzELMAkGA1UECBMCSUwxEDAOBgNVBAcTB0NoaWNhZ28xDDAKBgNVBAoTA1VQUjEM
MAoGA1UECxMDRW5nMRAwDgYDVQQDEwd1cHIubGxjMB4XDTE3MDUxNjAwNDAxNVoX
DTE4MDkyODAwNDAxNVowYTELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAklMMRAwDgYD
VQQHEwdDaGljYWdvMQwwCgYDVQQKEwNVUFIxDDAKBgNVBAsTA0VuZzEXMBUGA1UE
AxMOZXJuZXN0LnVwci5sbGMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQC0V2b+V3aH2MANQdECedtFQoy05JRwUrMcA0FUD21e2XOAfCmwEHWiyRidGRd1
5T2jj7mXNKibmXB7WD5G+2NaniF2vFykex8xOoQBBu0WgXBld026vnAWhKi9H8qB
H2XPJmO37Nk6zPVJZmxccjWwH/ORGEIEu+l9va9dVE7BXhwMo9FIm0xA48sOC15m
OA84Wj3/g/FxHP+KoBlbMo5NxerR0u0P/t73v9CGQxiLObzc2pJamy2X4v4VyMLe
FBIDOV8+n8GgzKIJjbZifbzFzBlVesdH9gVswm0rJt3wr0CD+DDfvIhLiIrw+EkG
lkY9vPpyyJlh6Zwkqb/foUQDAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAIKtW4Mz
6fqf3jDoc8KSELUzBjx+5kQYfIb7T0nWlu8dandtnaV0+yWzggQjSJvl9Cu81SaJ
8GZZ+7aLhi+ooYfuzBj3B9lgkyvPO5MSqZihml/L7hdR8EVw/ChUBlr2ECMGJ45S
e9eyCNtJvN2jnYGHy07kFc5k5leCWF2/bVD4HyX6dxHmjvf1YNpaoT2x3bLybwJt
v9Ns2Ywl3Mexqt5+JgU0Wn9fSvFmtXPYUWmt+N30MQBb+dNQGK1m+PvceZpQUMJn
QEIIjPsPreZ4+nshI7/vvx4P4INFsJRqFU2lQ0V7oZ3pqiOo2ZhqXptIMD7W4Pbl
qk5nHsrG3ZQXcPw=
-----END CERTIFICATE-----`)

  serverKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAtFdm/ld2h9jADUHRAnnbRUKMtOSUcFKzHANBVA9tXtlzgHwp
sBB1oskYnRkXdeU9o4+5lzSom5lwe1g+RvtjWp4hdrxcpHsfMTqEAQbtFoFwZXdN
ur5wFoSovR/KgR9lzyZjt+zZOsz1SWZsXHI1sB/zkRhCBLvpfb2vXVROwV4cDKPR
SJtMQOPLDgteZjgPOFo9/4PxcRz/iqAZWzKOTcXq0dLtD/7e97/QhkMYizm83NqS
Wpstl+L+FcjC3hQSAzlfPp/BoMyiCY22Yn28xcwZVXrHR/YFbMJtKybd8K9Ag/gw
37yIS4iK8PhJBpZGPbz6csiZYemcJKm/36FEAwIDAQABAoIBAFnM8/JOpy06F1gC
oAs1lqREYUBqWigcZifazxsGm5WRflxKd4++gnVO4bzIk9AXGLxrgyTyCcuBemtR
I6HA6FZlS5COdytIS939n2HLix+b0NYVtFT7I2FzKXsTa4kkk1O1cA5UGE+ZY8Lr
B6Z5EJblMincBHPLBpegbsRwiM7sMYFjNFnbKxGCmNDmsffFC2Jk3iyJ7EBfSDdL
h50rS6Qlsdh1/YfOHBoM15JdDtTucugYFToB0UAATjcSD22/9f4K9F5ezwHhVo7l
qFnWzH/zII+xHRiRJL1l3sy2tmFxK5f5HYWJsBSzk8XqxW6ZZUkiZ0l81wSWdfxX
R/G+MZECgYEA431U+mxFk55S48lN5YRG0FgJuJL1kHV3mk10J2jUEzqq01GcuM/e
Vs8lPaMe8JSd7LNjl13ZlELyGhlBajRJZyT9qPkSmBg9moQ3VRQrlzIugX8M7MkD
6igJHoiELFmWSHnb8F8lNsm21F21AruDB7+aMjVZGegG7TW38NDmPu0CgYEAyvFk
GaQf1mBYP03oC7a5+rYTzLhhkHzJheIxNaPWcwvQybp7iNkZZN9X/fDNMomzPQJU
B8Hx8W2HhT6X9Q3krp4mu3xcL4Fu4PyGuWuSeDynTPpsLyw7CQCEryaBuOJsuaaC
kJ8o5q0z0FM7JxerxmdbjtWOfupxj1xJvotlQK8CgYEApLQER38GbNdPGh6QKGE5
x/RFrX2xDxMNMglr1HCgQv/R2EeZWXEef+lgBB1K9FIVeN90do37Ts2dbWnlo5gR
oPz2LlWOsdGzIGEjkpSU7tXlN9qdp/3tuKWd3J5oW5fNJ9IafBDW57Dpjx39ROov
9vcxE7LuPMRx52JiN5d3bA0CgYAv6fgZOa+unIaZQ5qCkXytXamlDu8x/tRRgMrf
gDQUa3i6+AVMlP0y3KxYry6zPOGNiOwv/LFTr+lsIxAbclFIjNxWLZFSQfvcsKJa
SrSFIMTbHtDF/mpdwLqS48OC1CqZVl/DJ2Cvvbra95uiqisLJ8HtIiyHeHChSUe7
7gtUzwKBgQDFI5FVx8LcBeieILa19DGTOZVVeLS3b5f9V6EMRa5Qu2NXS2DMro4w
zy4zr1EP2V2posZnrMuuWT+Gr+B0vtDIKQkBW/dUwUTfbLqkJ+HV+yswIBUkecm8
4KhyOAtNrhbCZSfERncgHiGPbSmogauD/VUu3mXEvVHrJTht7e41pg==
-----END RSA PRIVATE KEY-----`)
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("<html>hello</html>"))
}

func runHTTPServer(port string){
  http.HandleFunc("/", handleRoot) 

  
  certs := make([]tls.Certificate,1)
  cert, err := tls.X509KeyPair(serverCert, serverKey)
  if err != nil {
    panic(fmt.Sprintf("Could not load certs: %s", err))
  }
  certs[0] = cert
  
  cfg := &tls.Config{
        //MinVersion:               tls.VersionTLS12,
        //CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
            tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
        },
        Certificates:  certs,      
  }
  
  
  srv := &http.Server{
            Addr:         ":"+port,
            TLSConfig:    cfg,
            //TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),            
  }
  
  err = srv.ListenAndServeTLS("", "")
  if err != nil {
     panic(fmt.Sprintf("Could not listen on port %s: %s", HTTPServerPort, err))
  }
}

// Run performs the ssl tests
func Run() {
  go runHTTPServer("8081")
  time.Sleep(1 * time.Second)
  // host verification failure 
  client := &http.Client{}
  _, err := client.Get("https://127.0.0.1:8081")
  if err != nil {
    fmt.Println(err)
  }
  // ca verification failure  
  // cipher verification failure 
}