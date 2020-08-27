The `kafka` output module publishes alerts to kafka.   
Each event/alert will be published as a new event. 
SSL/TLS auth is supported.    

|Flag|Type|Environment Variable|Description
|---|---|---|---|
|`-kafkabrokers`|string| CANARY_KAFKABROKERS|kafka brokers "broker:port"<br>- for multiple brokers, separate using semicolon "broker1:9092;broker2:9092"|
|`-kafkatopic`|string| CANARY_KAFKATOPIC|kafka topic; defaults to `canarychirps` if not set|


# SSL/TLS Support
To send alerts encrypted using SSL/TLS, you have to specify the parameters using the [SSL/TLS Settings](SSL-TLS-Settings.md)