There are few categories of options you can fiddle with "most are set to sane defaults":
* General Settings
* SSL/TLS Settings
* Input Modules
* Filter Modules
* Mapper Modules
* Output Modules 
# How to provide settings/parameters?
Parameters can be provided using **commnad-line flags**, or **environment variables**; if both have been provided, command line flags overrides environment variables.
- Environment variables are always ("CANARY_" + "UPPERCASEFLAG"), for example:
  - "`-output kafka`" is same as "`CANARY_OUTPUT=kafka`"
  - "`-esuser elastic`" is same as "`CANARY_ESUSER=elastic`"
- API key, and the canary domain can also be "additonally" specified using the token file downloaded from the console, either through `-tokenfile`, or simply placing it in the user's home directory, and ChirpForwarder will look for that file there on its own.
