package main

type policyList struct {
  Policy []struct {
    Name string `json:"name"`
    URL string `json:"url"`
  } `json:"policy"`
}

type contentPolicy struct {
  RemotePath string `json:"remotePath"`
  ListenerPolicy string `json:"listenerPolicy"`
  Name string `json:"name"`
  VirtualPath string `json:"virtualPath"`
  RemotePolicy string `json:"remotePolicy"`
  Description string `json:"description"`
  RequestProcessType string `json:"requestProcessType"`
  ResponseProcessType string `json:"responseProcessType"`
  RequestProcess string `json:"requestProcess"`
  ResponseProcess string `json:"responseProcess"`
}

type virtualDirectory struct {
  Description string `json:"description"`
  AclPolicy string `json:"aclPolicy"`
  Name string `json:"name"`
  Enabled bool `json:"enabled"`
  ErrorTemplate string `json:"errorTemplate"`
  ListenerPolicy string `json:"listenerPolicy"`
  VirtualPath string `json:"virtualPath"`
  RemotePath string `json:"remotePath"`
  RemotePolicy string `json:"remotePolicy"`
  RequestFilterPolicy string `json:"requestFilterPolicy"`
  RequestProcess string `json:"requestProcess"`
  RequestProcessType string `json:"requestProcessType"`
  ResponseProcess string `json:"responseProcess"`
  ResponseProcessType string `json:"responseProcessType"`
  UseRemotePolicy bool `json:"useRemotePolicy"`
  VirtualHost string `json:"virtualHost"`
}

type httpListener struct {
  Name string `json:"name"`
  Enabled bool `json:"enabled"`
  ErrorTemplate string `json:"errorTemplate"`
  Interface string `json:"interface"`
  Port int `json:"port"`
  ReadTimeoutMillis int `json:"readTimeoutMillis"`
  UseDeviceIP bool `json:"useDeviceIp"`
  IPAclPolicy string `json:"ipAclPolicy"`
  ListenerSSLEnabled bool `json:"listenerSSLEnabled"`
  ListenerSSLPolicy string `json:"listenerSSLPolicy"`
  PasswordAuthenticationRealm string `json:"passwordAuthenticationRealm"`
  RequirePasswordAuthentication bool `json:"requirePasswordAuthentication"`
  UseBasicAuthentication bool `json:"useBasicAuthentication"`
  UseChunking bool `json:"useChunking"`
  UseCookieAuthentication bool `json:"useCookieAuthentication"`
  UseDigestAuthentication bool `json:"useDigestAuthentication"`
  UseFormPostAuthentication bool `json:"useFormPostAuthentication"`
  UseKerberosAuthentication bool `json:"useKerberosAuthentication"`
}

type amqp10Listener struct {
  AclPolicy string `json:"aclPolicy"`
  IpAclPolicy string `json:"ipAclPolicy"`
  ReadTimeoutMillis int `json:"readTimeoutMillis"`
  Ip string `json:"ip"`
  Port int `json:"port"`
  UseDeviceIp bool `json:"useDeviceIp"`
  Name string `json:"name"`
  SaslMechanism string `json:"saslMechanism"`
  UseSsl bool `json:"useSsl"`
  Description string `json:"description"`
  Enabled bool `json:"enabled"`
  Interface string `json:"interface"`
  ErrorTemplate string `json:"errorTemplate"`
  SslPolicy string `json:"sslPolicy"`
}

type httpRemote struct {
  UseBasicAuth bool `json:"useBasicAuth"`
  ProxyPolicy string `json:"proxyPolicy"`
  SSLInitiationPolicy string `json:"SSLInitiationPolicy"`
  TcpConnectionTimeout int `json:"tcpConnectionTimeout"`
  HttpAuthenticationUserPolicy string `json:"httpAuthenticationUserPolicy"`
  UseChunking bool `json:"useChunking"`
  TcpReadTimeout int `json:"tcpReadTimeout"`
  Name string `json:"name"`
  EnableSSL bool `json:"enableSSL"`
  RemotePort int `json:"remotePort"`
  Enabled bool `json:"enabled"`
  RemoteServer string `json:"remoteServer"`
  ProcessResponse bool `json:"processResponse"`
}

type amqp10Remote struct {
  UseSsl bool `json:"useSsl"`
  IdleTimeoutMillis int `json:"idleTimeoutMillis"`
  UserPolicy string `json:"userPolicy"`
  SslPolicy string `json:"sslPolicy"`
  TransferTimeoutMillis int `json:"transferTimeoutMillis"`
  Name string `json:"name"`
  SaslMechanism string `json:"saslMechanism"`
  Description string `json:"description"`
  RemotePort int `json:"remotePort"`
  Enabled bool `json:"enabled"`
  RemoteServer string `json:"remoteServer"`
  ProcessResponse bool `json:"processResponse"`
}

type rabbitMqRemote struct {
  Name string `json:"name"`
  Description string `json:"description"`
  Enabled bool `json:"enabled"`
}
