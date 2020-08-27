# First, get `API Key` and `Domain`
To get started, you need two things: an **API token**, and your **Canary Console Domain** `https://**your_domain**.canary.tools`.  
To get the API token, in your canary console, go to `Settings`, then  turn `ON` the `API` switch; copy your API key  
You might also want to `Download Token File` which has those two pieces of info in it.  
![Get canary API](https://github.com/SherifEldeeb/ChirpForwarder/raw/master/assets/images/01-GetAPI.png)
***
# Examples
## Output to file, everything else on default settings.
This will get **unacknowledged** incidents from the **last 7 days**, write them to a file called 'canaryChirps.json'
```
~$ ChirpForwarder \
> -apikey 13c...d1f \
> -domain f...a \
> -output file
```
## Get all incidents, from certain time/date, output to file.
This will get **all** incidents, since **"2020-05-25 00:00:00"**, write them to a file called 'canaryChirps.json'
```
~$ ChirpForwarder \
> -apikey 13c...d1f \
> -domain f...a \
> -output file \
> -since "2020-05-25 00:00:00" \
> -which all
```

## Elasticsearch output, with SSL & Basic Auth
```
ChirpForwarder \
-apikey 13c...0d1f \
-domain f...a \
-output elastic \
-esuser elastic \
-espass elastic \
-eshost "https://192.168.0.165:9200" \
-ssl true \
-sslclientca /tmp/ca/ca.crt \
-sslclientkey /tmp/chirp/chirp.key \
-sslclientcert /tmp/chirp/chirp.crt
```

