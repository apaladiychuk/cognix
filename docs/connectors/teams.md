#### Microsoft Teams

## Setup 
From CogniX UI navigate to connectors and create a new connector
Choose `Teams`
At step 2:
- Choose a name, it's just a description
- Fill the "Connector Specific Configration" with the json below filled with the corect data
- Refresh frequency in seconds is the delta of time that CogniX will use to start a new scan on your connected data source
- Connector credential, fill with a random number, it's not used

```json
{
  "channel": "",
  "topics": ["",""],
  "files": {
    "folder": "",
    "recursive": false,
  },
  "token": {
    "access_token": "",
    "expiry": "",
    "refresh_token": "",
    "token_type": ""
  }
}
```

- channel : name of channel for analyzing
- recursive :  false - scan only given folder , true - scan nested folders
- token : OAuth token for access to ```one drive```
- files : 
  - folder : optional, folder name for scanning
  - recursive :  false - scan only given folder , true - scan nested folders



