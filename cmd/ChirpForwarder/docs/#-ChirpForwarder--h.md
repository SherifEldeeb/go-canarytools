```~$ ChirpForwarder -h```
```
+++++++++++++++++++++++++++.ooo.`  -/-`          -/-`
++++++++++++++++++++++++++.ooooo` -+++:.        -+++:.
+++++++++++++++++++++++++++.ooo.` -+++++/:.     -+++++/:`
++++++++++++++/:---:/+++++/+++++` `-/++++++/-`  `-/++++++/-`
++++++++++++/-        .---:+++++`    `-/++++++:.   `-/++++++:.
+++++++++++:`          .:/++++++`       .:++++++/:`   .:++++++/-`
++++++++++/           -+++++++++`        `-/++++++/-`  `-/++++++/-`
+++++++++/`          `++++++++++`            `-/++++++:.   `-/++++++:.
+++++++/.            .++++++++++`            `-/+++++/:.   `-/+++++/:.
+++++/.              `++++++++++`        `-/++++++/-`  `-/++++++/-`
+++/-                 /+++++++++`       .:++++++/-`   .:++++++/-`
++/`                  /+++++++++`    `-/++++++:.   `-/++++++:.
:.                   :++++++++++` `-/++++++/-`  `-/++++++/-`
                   `:+++++++++++` -+++++/:`     -+++++/:`
                  -/++++++++++++` -+++:.        -+++:.
               `-/++++++++++++++`  -/-`          -/-`

INFO[0000] starting canary ChirpForwarder
Usage of ChirpForwarder:
  -apikey string
        API Key
  -compress
        [OUT|FILE] file compress log files?
  -domain string
        canarytools domain
  -escloudapikey string
        [OUT|ELASTIC] elasticsearch Base64-encoded token for authorization; if set, overrides username and password
  -escloudid string
        [OUT|ELASTIC] endpoint for the Elastic Cloud Service 'https://elastic.co/cloud'
  -eshost string
        [OUT|ELASTIC] elasticsearch host
  -esindex string
        [OUT|ELASTIC] elasticsearch index (default "canarychirps")
  -espass string
        [OUT|ELASTIC] elasticsearch password 'basic auth'
  -esuser string
        [OUT|ELASTIC] elasticsearch user 'basic auth'
  -feeder string
        input module (default "consoleapi")
  -filename string
        [OUT|FILE] file name
  -filter string
        filter to apply to incident ('none', or 'dropevents')
  -host string
        [OUT|TCP] host
  -insecure
        [SSL/TLS CLIENT] ignore cert errors
  -interval int
        alert fetch interval 'in seconds'
  -kafkabrokers string
        [OUT|KAFKA] kafka brokers "broker:port"
                        for multiple brokers, separate using semicolon "broker1:9092;broker2:9092"
  -kafkatopic string
        [OUT|KAFKA] elasticsearch user 'basic auth'
  -loglevel string
        set loglevel, can be one of ('info', 'warning' or 'debug')
  -maxage int
        [OUT|FILE] file max age in days 'older than this will be deleted'
  -maxbackups int
        [OUT|FILE] file max number of files to keep
  -maxsize int
        [OUT|FILE] file max size in megabytes
  -output string
        output module ('tcp', 'file', 'elastic' or 'kafka')
  -port int
        [OUT|TCP] TCP/UDP port
  -since string
        get events newer than this time.
                        format has to be like this: 'yyyy-MM-dd HH:mm:ss'
                        if nothing provided, it will check value from '.canary.lastcheck' file,
                        if .canary.lastcheck file does not exist, it will default to events from last 7 days
  -ssl
        [SSL/TLS CLIENT] are we using SSL/TLS? setting this to true enables encrypted clinet configs
  -sslclientca string
        [SSL/TLS CLIENT] path to client rusted CA certificate file
  -sslclientcert string
        [SSL/TLS CLIENT] path to client SSL/TLS cert  file
  -sslclientkey string
        [SSL/TLS CLIENT] path to client SSL/TLS Key  file
  -then string
        what to do after getting an incident? can be one of ('nothing', or 'ack')
  -tokenfile string
        the token file 'canarytools.config' which contains api token and the domain
  -which string
        which incidents to fetch? can be one of ('all', or 'unacknowledged')```