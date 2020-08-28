Those are used to set general options, like loglevel, what type of output module, which filter to use ... etc.  
|Flag|Type|Environment Variable|Description
|---|---|---|---|
|`-feeder` |string| CANARY_FEEDER|input/feeder module: where to read alerts from; currently we only support "consoleapi"| 
|`-outout` |string| CANARY_OUTPUT|output module: how to send/forward alerts; valid values are ('tcp', 'file', 'elastic' and 'kafka')|
|`-loglevel`|string| CANARY_LOGLEVEL |set loglevel; valid values are ('info', 'warning', 'debug' and 'trace')|
|`-then` |string|CANARY_THEN|what to do after reading an incident? valid values are 'nothing' which will do as advertised, 'ack' which will acknowledge the incidents, or 'delete' which will delete the incidents once they have been successfully forwarded to destination|
|`-since`|string|CANARY_SINCE|get events newer than this time.<br>- format has to be like this: 'yyyy-MM-dd HH:mm:ss'<br>- if nothing provided, it will check value from '.canary.lastcheck' file,<br>- if .canary.lastcheck file does not exist, it will default to events from last 7 day|
|`-which`|string|CANARY_WHICH |which incidents to fetch? valid values are ('all', and 'unacknowledged')|
|`-filter`|string|CANARY_FILTER|filter to apply to incident; valid values are ('none', and 'dropevents')|
|`-flock`|string|CANARY_FLOCK|flock name to process incidents for 'if left empty, all incidents will be processed'|