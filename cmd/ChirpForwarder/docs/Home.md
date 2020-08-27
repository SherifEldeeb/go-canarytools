# ChirpForwarder
![ChirpForwarder Logo](https://github.com/SherifEldeeb/ChirpForwarder/raw/master/assets/images/ChirpForwarder.png)
## What?
ChirpForwarder is a tool that pulls alerts and incidents from a thinkst canary console "https://canary.tools/", then forwards them to local destinations (e.g. SIEM/SOAR solutions) for further processing/indexing/archival/etc.  
It's a single binary that has no dependencies, runs on windows, linux, macOS, and it's also available as a docker image.

## Why?
Even though there are many ways to receive alerts from the canary console, canary owners who would like to integrate alerts with local/on-prem SIEM solutions had to either do some custom development (to integrate with the console API), and/or expose a service listener to the internet (syslog or webhook).  
This tool aims to lower the barrier for canary owners who have a need to integrate canary alerts into their SIEM/SOAR solutions, or are looking for an easy way to archive alerts & incidents in a standard format, for regulatory & compliance purposes, without the need for custom development, or the need to open a listener to the internet.
## Notable Features
- Fetches all incidents, or only unacknowledged ones, with time filtering (e.g. only get unacknowledged incidents since X)
- Automatically fetches incidents on configurable intervals, with "remember last time checked" on restart.
- Can mark fetched incidents as "acknolwedged" once successfully forwarded to destination.
- Supports the following outputs:
  - **Flat file**: stores incidents as JSON-lines in flat files, with built-in automatic log rotation and retention on configurable settings.
  - **TCP**: forward events to a TCP listener (very popular with SIEM solutions); incidents are JSON encoded and `\n` separated.
  - **elasticsearch**: each incident is indexed as a document; supports HTTP basic auth.
  - **Kafka**: publishes events to a kafka cluster/topic.
  - **TLS Support**: ChirpForwarder supports TLS with custom CAs, certificate based client auth, and it can ignore certifcate errors (e.g. self-signed certs); TLS is supported for TCP, Kafka and elasticsearch.
- Parameters can be provided using commnad-line flags, or environment variables 
  - Environment variables are always ("CANARY_" + "UPPERCASEFLAG"), for example:
    - "`-output kafka`" is same as "`CANARY_OUTPUT=kafka`"
    - "`-esuser elastic`" is same as "`CANARY_ESUSER=elastic`"
  - API key, and the canary domain can also be "additonally" specified using the token file downloaded from the console, either through `-tokenfile`, or simply placing it in the user's home directory, and ChirpForwarder will look for that file there on its own.
- docker support with a sample "docker-compose" file; set environment variables, `docker-compose up` and off you go.
