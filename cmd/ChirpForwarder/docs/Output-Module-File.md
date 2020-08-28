The `file` output module saves alerts to a file.   
Events/Alerts will be separated by newline `\n`.  
It supports automatic file rotation, and retention.  

|Flag|Type|Environment Variable|Description
|---|---|---|---|
|`-filename`|string| CANARY_FILENAME|file name, defaults to `canaryChirps.json` if not specified|
|`-maxsize`|int| CANARY_MAXSIZE|file max size in megabytes, file will rotate wen reaches this size|
|`-maxbackups`|int| CANARY_MAXBACKUPS|max number of rotated files to keep|
|`-maxage`|int| CANARY_MAXAGE|files max age in days 'older than this will be deleted'|
|`-compress`|bool| CANARY_COMPRESS|if set to true, files will be compressed|
