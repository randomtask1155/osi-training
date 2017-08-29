package supportdemossl

import(
  "fmt"
  "net/http"
  "crypto/tls"
  "time"
  "crypto/x509"
  "osi-training/colors"
  "encoding/pem"
  "strings"
)

var (  
  rootCA = []byte(`-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`)

serverBadHostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDTjCCAjagAwIBAgIJALPcPP2kFjq3MA0GCSqGSIb3DQEBCwUAMFoxCzAJBgNV
BAYTAlVTMQswCQYDVQQIEwJJTDEQMA4GA1UEBxMHQ2hpY2FnbzEMMAoGA1UEChMD
VVBSMQwwCgYDVQQLEwNFbmcxEDAOBgNVBAMTB3Vwci5sbGMwHhcNMTcwODI5MDE0
MzA0WhcNNDUwMTE0MDE0MzA0WjBOMQswCQYDVQQGEwJVUzELMAkGA1UECAwCSUwx
EDAOBgNVBAcMB0NoaWNhZ28xDDAKBgNVBAoMA3VwcjESMBAGA1UEAwwJMTI3LjAu
MC4yMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2Ie8KLSowaEdpVtd
ENVKxHL05djwNe32y/+Z+lJwb5W0s6DK5earVUorgmqkhqEBmCRL7iRBG9tDKw0Z
tdpV2PHADdwYX00eMQCuR90rBTzteGhqUtaniXf4zAwQ4lnic3tVzeRRTbpjjMds
ml2nf0nQI++kuskb2uRNvcP/gi7gCatX1KJUjdeQVaeqG52c+ZDik+PDDRHmFCOm
NXXTwUclF8P1Zm+Rh+1oBFUzrcJYmSzXbnUFYl2T45gU+uOwFPPF5lrlA59LpH5V
eSrzOm7dYM5jD0i+iqQn9lXmiaUFCXxbhmFpo+d0WQFtzlklDBzoRbB11vAQdT/H
gAVwQwIDAQABoyMwITAfBgNVHREEGDAWhwR/AAACgg5lcm5lc3QudXByLmxsYzAN
BgkqhkiG9w0BAQsFAAOCAQEAIrVtPBs3ON0grQ9p0G8HN2Htxzp+XMb7PXqOPtTY
4jKkIoMEO+o9ejJwCoCC8txXTcwmugJz4Qip+tf7jtL2J43b0dCgF9Xl0ObX6t2r
tXoEfOUV1ELaZ8uzS/DMEvpWKTOdovmYRQTbIN3LQlota9nT7Z6/KqZq5QrvGuOI
mdd2lRE1CRXQmlCO4+P3Xbq7oySDjCA8vwRPp1aZ0nbaNFXT/INwJhuXpP+C3WE7
VL8fjOXoFAD1XtOi6gGDB+Ixe+dCc46foRoGwZ724UrcPdTLoHUlXCYQv4useuSO
0AFkX4dnmDLyWHaEY+1R+xM0Cd39URCpTRj08HnyjRKyYg==
-----END CERTIFICATE-----`)

serverBadHostKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA2Ie8KLSowaEdpVtdENVKxHL05djwNe32y/+Z+lJwb5W0s6DK
5earVUorgmqkhqEBmCRL7iRBG9tDKw0ZtdpV2PHADdwYX00eMQCuR90rBTzteGhq
UtaniXf4zAwQ4lnic3tVzeRRTbpjjMdsml2nf0nQI++kuskb2uRNvcP/gi7gCatX
1KJUjdeQVaeqG52c+ZDik+PDDRHmFCOmNXXTwUclF8P1Zm+Rh+1oBFUzrcJYmSzX
bnUFYl2T45gU+uOwFPPF5lrlA59LpH5VeSrzOm7dYM5jD0i+iqQn9lXmiaUFCXxb
hmFpo+d0WQFtzlklDBzoRbB11vAQdT/HgAVwQwIDAQABAoIBAQCTiLity5qIhCzc
9EmPJKVTATjYO15AgKl+CeRyaSVvAkQjeFWHHAp1jJnkvKDpkl6xuVl7I3yxbE5x
0PVJsUB5Fw9S4jpBBRyjKuGU2Z9sbD0po4t+cf+pbPM6pKYEdFYVdV3bccvr+CnI
TEE1VMbBtht5wNku48R0+sX1nMxFeEUvwdw19fatt+fgPKWmt1847fDjVJHTBjN8
V866hxzEjfO02EhMRw1PQNpNCIwxxnxZomBPF0umJ+E3XPXhm4AEv6h3vEkMgeDb
4q4VxZeDPtxSpa8riAFbTO2/pLmncFDICk2VbMnO0fQoe4N4c5/8qCUIbz9VMKx9
hXshqGrhAoGBAP31UhCnFuaW9HFC4VT8G2i4lJt1y4nAs2tlOkRAGYgsN5WMpZlM
cwTtyVJKfX9KSCvhDwa5bTi1L1NSHpL3aXEehMr3dBfOt1nKhNrmOHxtu6r4fcBV
MhESxg5QLMEEJloXFSLpqLzutg+ZW/Vtgts56H10WFIrXQjqhXy7rCppAoGBANpF
Ye9cfyVXShlreORntClcBuV0cp1aOEKlZiMrRjdEO7bsou8RlfFizi9C4O18ssry
znLH/1VahoJXWlA5qFpF5+GtXpH99LivA6klZFYd7CxysQ3EFIRVGel3RhczFtqY
58SXBUtgLKQbXUINh97L4uJ+8s4BdESnfARNYXfLAoGBALY3QiW+rdjPsR5fHWyr
40K1RbNxxpPVVycUn8T+tRMWnt9H5j3jM5fywYOw/PDLLJOHTQ/HFg4dOpKvFk7f
BKHsgt1axcqbQMVFYSPrEdZ9qazX1Oqedmj3rnHwptvrDVWQwFKnkrHVGX6BzBCB
/eUsXJs/UmvKv1s5YoQH6FUZAoGBALXHWsyxx+IdVETfUItIIq7fSY8G7/1ECd7X
SLy5aXzVgrXRVdCsYKfrlbOeixy94VlL5cuJ7If1IFikzz5JBoBH/9nfPQdw7MJy
XyYbUfi2at6JlhLU3hFnnnDlOKZRkhLzjPIYo/wmWIdDkuNTyXgvbQEcoLIYW/Bh
W9NpbmcJAoGASpJSX8qDkmxbCIOldtIomdlzJATDByVILfrTGVOIUSh+FVHKtDgs
yhWYP0vvpmxIKACgWhhSdJHfsrN8SQOkgX9LrI+wsTCm+RubbWa/6LmgOmwqUqjU
aNVUmZp/6Jb9Izqa/Jo2Bisk0cE4Q/+msZqaU/h9jKcFHUhXuyU5sFU=
-----END RSA PRIVATE KEY-----`)

serverLocalHostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDTjCCAjagAwIBAgIJALPcPP2kFjq2MA0GCSqGSIb3DQEBCwUAMFoxCzAJBgNV
BAYTAlVTMQswCQYDVQQIEwJJTDEQMA4GA1UEBxMHQ2hpY2FnbzEMMAoGA1UEChMD
VVBSMQwwCgYDVQQLEwNFbmcxEDAOBgNVBAMTB3Vwci5sbGMwHhcNMTcwODI5MDEz
NTMzWhcNNDUwMTE0MDEzNTMzWjBOMQswCQYDVQQGEwJVUzELMAkGA1UECAwCSUwx
EDAOBgNVBAcMB0NoaWNhZ28xDDAKBgNVBAoMA3VwcjESMBAGA1UEAwwJMTI3LjAu
MC4xMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqb2NniDEYPRy2t95
YUgO1GEyKsILhxcfPCcFYnXRiEV3clAIqVcXf3JiQA2Jm6vGPXijznbrTKUUtpaI
/mfwFepUfGUwlYrVyjvE5T6L0c+nWfrRMI4u8XS2yUXtQkpaW35SfsgsrPkw7csx
qxGaOCWRNL9tnulJFJAaDC/6gfQ9x42fZICD64dF/YdrR5MbrhSnMbYKKhiF1HcE
pE8oDIm1mI5CiMOtANDyPv12SPSBy6AHcPcP9hpX3rCx5REl/0AVBtrvc0CJzDSP
yUeR4BJWS/ETMBVWlvf8k2Ri3ZCwAG3Iv2nzYQWZKrlwbilVg3q+xmimDcHZVCHb
28Y0qQIDAQABoyMwITAfBgNVHREEGDAWhwR/AAABgg5lcm5lc3QudXByLmxsYzAN
BgkqhkiG9w0BAQsFAAOCAQEAHgtMyrlgE972wYS34r2BF1ODA/OUFZ3RdiJZ5MtF
XAx4JNEKS9v0DBnSQSNFMevGdPnWwv2LAU33dxGbm2kT2bEkmquwJc8zyu9MgUxv
Rw6q3i3rt9Wt2oCLSdZDN1iXgjHYCaHBWho6UNI+pwELnFdICLZLBOkHeB/vAI0K
sL/gbRZLWbgn+jMR7JIaDgR0lSaTIac2PUvMnGeAB4WW1ptOz7HEC7Q5A8zrPOx7
0h8etkwqIIA0oajUXok9A7RqNNOf2s3TF6993xShOdO+/hO/nJ0P1NAgZG4eYkUz
76lMEh9XI8IDhV5Mx4gsyuAWB51Xsc3J/PORpB4335Qi1w==
-----END CERTIFICATE-----`)

serverLocalHostKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAqb2NniDEYPRy2t95YUgO1GEyKsILhxcfPCcFYnXRiEV3clAI
qVcXf3JiQA2Jm6vGPXijznbrTKUUtpaI/mfwFepUfGUwlYrVyjvE5T6L0c+nWfrR
MI4u8XS2yUXtQkpaW35SfsgsrPkw7csxqxGaOCWRNL9tnulJFJAaDC/6gfQ9x42f
ZICD64dF/YdrR5MbrhSnMbYKKhiF1HcEpE8oDIm1mI5CiMOtANDyPv12SPSBy6AH
cPcP9hpX3rCx5REl/0AVBtrvc0CJzDSPyUeR4BJWS/ETMBVWlvf8k2Ri3ZCwAG3I
v2nzYQWZKrlwbilVg3q+xmimDcHZVCHb28Y0qQIDAQABAoIBAHWalqdLrqQ6WeWZ
1kB3q4asfRSw80m61HJZ2to4JV/UOYkjOI0TgX1U5AvbFU+dOTHYBy6CqE6nRe2n
6OzTWL3maHjzpzDFS5DdznLC3h8VT68BH7UTniS/J/HWGCfI2kfwAvpCeOmzkZoA
Ib6T6JUkOjIgu7PRkEfF+7Yb++XTEar9W7l8MkWgB7nEClR/r5yYNS8RzqLSM9dx
ay3XaXWeS/zW4IgKy9BpQt6kOQB9ID1tSCtx3bFVcgUFi++YwNe4yb2SptEdCRFf
0uML3Q6XsLTWReSP8T+0emKk8JRNPmNv4qwu7l3JJbTn0fNTUTiv3xzVsy5SPxUa
SiXHuIECgYEA0HuBUWandBPFd3q+6yJ2YnRRq13RsuqUrPVW5Ht8PGw/P7EiWBzU
uo6yTGlUsxehtPZVppdJ+yo2L+UkQZGlgUSNx6LFCVXbjQ963cHl3V+gruJ4IMcS
HXXZ8ofU+KexKfBLPuDRDLYiQR1USzBROvOvWa7NvnI09EIWg5JystECgYEA0G2J
epaXGwHAhmAIU6AGHt6kmjHWFwGg9SAFmIbwi74D5hvQtq1tT1/DJhYGzfrGq/YL
NTywS0erpUB7PlEgqxG5o+cHdUG4ZK5JjmGSkAFLPFKjmtIdfb6bn3XD3OxiV3B9
nM6MZKGLDTy0kq8BqzxaZUkZ67qSxGh+ZWID6lkCgYBFQ1LyRUWKxerLyAbXvYpR
KU3cvf2mEFM2pweoPvJGbLwSn/nGBkYSeMf5pODT4x0BLvnDr+2POTXpcZo7AnVW
3fywf34wnMqlMahjNkD07AlJMMoyMZDuIrI25jO2LJgqU7/b8vrg2z3EHkdb94B1
MnJmPDH0fKLlJ3OtYKEkkQKBgEmxsR5LCHpgEDZy1f7bYc6gYgqy/EN+K+7/t3rK
m1qNgMtnolA02aVq8pEQ0K8bsAs1H5lfL+YuHR58whaykJ5r0fuFwDlRV2Uhypgx
H6UTEArwHTCsggjn0BZ9iRcf7VWFTKSY00Lxazzu7dm/TxPAbyXIxwV2Hlabq7Ul
BsVhAoGBALUhvg1Z41hxieK7Y0EN8Xzu6NtNRwR6N8q4a+l+C7gh+E/SQtj1onlp
0zNuLyDvzjw9lzHKvA7J7gCJ2yyHBMmetD0pKC2ssp1d6ixmcUTS9nlKJl+Ku+GB
8omqiy9IPQGNYZERg+9Mw8fC5LinH6FU6Ni6I1cFtKuC53MMLYO9
-----END RSA PRIVATE KEY-----`)

)

func createPemCert(cert, key []byte) *x509.Certificate {
  p, _ := pem.Decode(cert)
  if p == nil {
    panic(fmt.Sprintf("Could not parse pem cert"))
  }
  c, err := x509.ParseCertificate(p.Bytes)
  if err != nil {
    panic(fmt.Sprintf("Could not parse cert from decoded pem: %s", err))
  }
  return c
}

func createCert(cert, key []byte) tls.Certificate {
  c, err := tls.X509KeyPair(cert, key)
  if err != nil {
    panic(fmt.Sprintf("Could not load certs: %s", err))
  }
  return c
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("<html>hello</html>"))
}

func runHTTPServer(cert, key []byte, port string){ 

  certs := make([]tls.Certificate,1)
  certs[0] = createCert(cert, key)
  
  cfg := &tls.Config{
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
            tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
        },
        Certificates:  certs,     
  }
  
  srv := &http.Server{
            Addr:         ":"+port,
            TLSConfig:    cfg,
            //TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),            
  }
  
  err := srv.ListenAndServeTLS("", "")
  if err != nil {
     panic(fmt.Sprintf("Could not listen on port %s: %s", port, err))
  }
}

// Run performs the ssl tests
func Run() {
  http.HandleFunc("/", handleRoot)
  go runHTTPServer(serverBadHostCert, serverBadHostKey, "8081")
  time.Sleep(1 * time.Second)
  
  // host verification failure #################################################
  fmt.Printf("Performing Host verfication failure test sending GET request https://127.0.0.1:8081: ")
  client := &http.Client{}
  _, err := client.Get("https://127.0.0.1:8081")
  if err != nil {
    fmt.Printf("%s\n", colors.ColorizeError(err))
  }
  cert := createPemCert(serverBadHostCert, serverBadHostKey)
  fmt.Printf(colors.ColorizeText(fmt.Sprintf(`
Using cert:
Subjects Common Name = %[1]s
Issuer               = %[2]s

Alt Names 
%[3]s

In this case the Subject "%[1]s" does not match 127.0.0.1 used in the get request 

`, cert.Subject.CommonName, cert.Issuer.CommonName, strings.Join(cert.DNSNames, "\n") )))


  // ca verification failure ###################################################
  go runHTTPServer(serverLocalHostCert, serverLocalHostKey, "8082") 
  time.Sleep(1 * time.Second)
  fmt.Printf("Performing Certiicate Authority verfication failure test sending GET request https://127.0.0.1:8082: ")
  _, err = client.Get("https://127.0.0.1:8082")
  if err != nil {
    fmt.Printf("%s\n", colors.ColorizeError(err))
  }
  cert = createPemCert(serverLocalHostCert, serverLocalHostKey)
  fmt.Printf(colors.ColorizeText(fmt.Sprintf(`
Using cert:
Subjects Common Name = %[1]s
Issuer               = %[2]s

Alt Names 
%[3]s

In this case the Issue (Root Certificate Authority) "%[2]s" is not known to the client therefore can not be trusted 

`, cert.Subject.CommonName, cert.Issuer.CommonName, strings.Join(cert.DNSNames, "\n") )))


  // cipher verification failure ###############################################
  go runHTTPServer(serverLocalHostCert, serverLocalHostKey, "8083") 
  time.Sleep(1 * time.Second)
  
  // add the root certificate to the http client transport struct so we can verify the servers certificate
  roots := x509.NewCertPool()
  ok := roots.AppendCertsFromPEM(rootCA)
  if !ok {
    panic("Could not parse root cert")
  }
  
  // restrict the cipher to something that is not supported by the server and create the new http client
  cfg := &tls.Config{
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
        }, 
        RootCAs: roots,   
  }
  tr := &http.Transport{TLSClientConfig: cfg}
  client = &http.Client{Transport: tr}
  
  fmt.Printf("Performing Cipher verfication failure test sending GET request https://127.0.0.1:8083: ")
  _, err = client.Get("https://127.0.0.1:8083")
  if err != nil {
    fmt.Printf("%s\n", colors.ColorizeError(err))
  }
  cert = createPemCert(serverLocalHostCert, serverLocalHostKey)
  fmt.Printf(colors.ColorizeText(fmt.Sprintf(`
Using cert:
Subjects Common Name = %[1]s
Issuer               = %[2]s

Alt Names 
%[3]s

Supported Client Cipher 
TLS_RSA_WITH_AES_256_CBC_SHA

Support Server Ciphers 
TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
TLS_RSA_WITH_AES_256_GCM_SHA384,

`, cert.Subject.CommonName, cert.Issuer.CommonName, strings.Join(cert.DNSNames, "\n") )))
}