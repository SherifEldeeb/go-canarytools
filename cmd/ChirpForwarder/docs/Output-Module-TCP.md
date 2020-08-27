The `tcp` output module sends alerts to TCP listeners (like logstash tcp input, or splunk TCP receiver).  
Events/Alerts will be separated by newline `\n`.  

|Flag|Type|Environment Variable|Description
|---|---|---|---|
|`-host`|string| CANARY_HOST|tcp listener host (IP or Hostname)|
|`-port`|int| CANARY_PORT|tcp listener port number|


# SSL/TLS Support
To send alerts over tcp encrypted using SSL/TLS, you have to specify the parameters using the [SSL/TLS Settings](SSL-TLS-Settings)