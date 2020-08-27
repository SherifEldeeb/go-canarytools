The `elastic` output module saves alerts to elasticsearch.   
Each event/alert will be indexed as a new document.  
It supports on-prem Elasticsearch clusters and the managed elastic cloud service.  
No auth, basic auth & SSL/TLS are all supported.    

|Flag|Type|Environment Variable|Description
|---|---|---|---|
|`-eshost`|string| CANARY_ESHOST|elasticsearch host:port|
|`-esindex`|string| CANARY_ESINDEX|elasticsearch index; defaults to `canarychirps` if not specified|
|`-esuser`|string| CANARY_ESUSER|elasticsearch user 'basic auth'|
|`-espass`|string| CANARY_ESPASS|elasticsearch password 'basic auth'|
|`-escloudapikey`|string| CANARY_ESCLOUDAPIKEY|elasticsearch Base64-encoded token for authorization; if set, overrides username and password|
|`-escloudid`|string|CANARY_ESCLOUDID|endpoint for the Elastic Cloud Service 'https://elastic.co/cloud'|

# SSL/TLS Support
To send alerts encrypted using SSL/TLS, you have to specify the parameters using the [SSL/TLS Settings](SSL-TLS-Settings)
