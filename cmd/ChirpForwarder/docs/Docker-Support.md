ChirpForwarder is available as a docker image, so you can `docker pull 0xdeeb/chirpforwarder` it.  
# Providing Parameters 
To run ChirpForwarder in docker, all parameters must be provided using environment variables, which can be through the following.
## The `.canary/lastcheck` file
To keep track of last time the tool checked for incidents between runs, it automatically creates a directory `.canary` and maintains that through a file `lastcheck` created in it.  
This directory resides in `/.canary` in the docker image, and a local directory needs to be mapped to the docker container to maintain `lastcheck` between docker runs "e.g. `docker run -v "/home/user/.canary:/.canary" ...`".  
To override this behavior, delete the file or provide the `CANARY_SINCE` environment variable in the following format `yyyy-MM-dd HH:mm:ss`.  
If the file does not exist, and no `CANARY_SINCE` has been provided, the tool will default to getting incidents from the last seven days.
## `docker run` Example
Choose file output module, set output filename, then mount volume to get the file out of the container.
```
$ sudo docker run \
-e CANARY_APIKEY=4c...10 \
-e CANARY_DOMAIN=f...a \
-e CANARY_OUTPUT=file \
-e CANARY_FILENAME=/chirps/chirps.json \
-v "chirps:/chirps"  \
-v "/home/user/.canary:/.canary"  \
0xdeeb/chirpforwarder:latest
```
## `docker-compose` Example
A sample docker-compose file; just pass the parameters as environment variables and mount volumes as needed.
```yaml
version: '3'
services:
  chirpforwarder:
    image: "0xdeeb/chirpforwarder"
    environment: 
        CANARY_APIKEY: 4c8...b10
        CANARY_DOMAIN: f...a
        CANARY_INTERVAL: 60
        CANARY_OUTPUT: file
        CANARY_FILENAME: /chirps/chirps.json
    volumes:
        - "./chirps:/chirps"
        - "/home/user/.canary:/.canary"
```
## `docker-compose` without hard-coding sensitive info
`docker-compose` can set values from environment variables as well; this example shows how to set parameters without hard-coding sensitive info into the file by passing arguments as environment variables.  
One important note, if you're using sudo, make sure you pass the environment variables using `-E` as well to copy them over to the `sudo` command, e.g.:  
`$ CANARY_DOMAIN=f...a CANARY_APIKEY=4c8...b10 sudo -E docker-compose up`  
```yaml
version: '3'
services:
  chirpforwarder:
    image: "0xdeeb/chirpforwarder"
    environment: 
        CANARY_APIKEY: ${CANARY_APIKEY}
        CANARY_DOMAIN: ${CANARY_DOMAIN}
        CANARY_INTERVAL: 60
        CANARY_OUTPUT: file
        CANARY_FILENAME: /chirps/chirps.json
    volumes:
        - "./chirps:/chirps"
```