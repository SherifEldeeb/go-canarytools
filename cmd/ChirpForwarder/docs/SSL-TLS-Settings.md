These settings apply to ALL output modules that supports SSL/TLS, like elasticsearch, kafka and TCP.
|Flag|Type|Environment Variable|Description
|---|---|---|---|
|`-ssl`|bool| CANARY_SSL|[SSL/TLS CLIENT] are we using SSL/TLS? setting this to true enables encrypted clinet configs|
|`-insecure`|bool| CANARY_INSECURE|[SSL/TLS CLIENT] ignore cert errors (e.g. expired or self signed certs)|
|`-sslclientca`|string| CANARY_SSLCLIENTCA|[SSL/TLS CLIENT] path to client rusted CA certificate file|
|`-sslclientkey`|string| CANARY_SSLCLIENTKEY|[SSL/TLS CLIENT] path to client SSL/TLS Key  file"|
|`-sslclientcert`|string| CANARY_SSLCLIENTCERT|[SSL/TLS CLIENT] path to client SSL/TLS cert  file|