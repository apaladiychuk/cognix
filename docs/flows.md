

# Orchestrator
TBD: 
 - Define genral behaviour 
 - How to have multiple orchestrator isntances 

The orchestrator is resposible of monitoring the connector table in the relational database and decide when it's time for a new analysis of the single connector
A single connector row represents a knowledge source the user decieded he wants to be analyzed by CogniX and stored in the vetcor database.

The orchestrator is resposible deciede if to scan again the source defined in the connecotr, and eventually to send analysys requests to both Connector or chaunker dependig on the spefic use case:
A new scan (connector) or analysis (chunker) shall be started only if there are no nother process running


The user shall decide it to disable the the connector, if so no futher scans will be started for that connector. 
Disabled is set in statuses
All date time are defined as UTC

## Status
The status of the connector for the moment is determined by the filed (verify status field name in the connector table)
In the future the status will be retrieved from NATS

Statuses:
- Active (connector just created)
- Pending scan (orchestrator sent a message to Connector or Chunker to scan the source)
- Working by Connector. The first action the Connector will do once it reveives the NATS message is to set this status. Once it sends the message to Chunker, it will be set back to Pending scan (after sending the message to NATS)
- Working by Chunker. The first action the Chunker will do once it reveives the NATS message is to set this status. Once finished it will set the appriate status (scan completed succesully, with error, unable to process)
- Scan completed successfully, set by chunker
- Scan completed with errors, set by chunker 
- Unable to process, set by chunker. This status indicates there is a problem with the source. The oarchestrator shall not request any more analysis for this resource. Only one (important only one) emmail shall be sent to the user reporting the issue

IMPORTANT:
TBD what to do when the row is in the pending or working status? 
We need to find our if there's a way for NATS to notify that a new item has been added to the dead letter queue. 
If so the orchestrator, or any other service, shall subscribe to this message (like we do for chunker, only one subscriber shall receive the message) and set the status to unable to process.
Again, only one (important only one) emmail shall be sent to the user reporting the issue

## Rules to scan again any given connector line
- no active scan 
- no pending requested scan
- the time of the last activity shall be greater than the refresh_freq (in seconds) 
- refresh_freq is coming form config (env) and different by filetype

## URL 
The user can decide to connect a URL as knowledge source
It shall provide:
- URL (mandatory)
- Sitemap URL (optiana)
- Scan all the links found in the page (bit, optiona)
- If the system shall look for a sitemap even if not provided

The Orchestrator, for this particular file type, it will forward the request to Chunker.

For this particular file type no file is stored in MinIO



If the user will delete the souce, the API receiving the request will be responsable to physically (soft delete) the row from the relational database and all the entities related to that documentid inside the vecotr database
Soft delete means that the 

## File 
The user can decide to upload a file as as knowledge source
It shall provide:
- the file to be analyzed (mandatory)

The user will be able to upload a file to our internal MinIO
The file will be scanned only Once.
The Orchestrator will send a request to Chunker as soon as the file is uploaded correctly
The URL property of the proto message will contain all the information for Chunker to open the file from MinIO
The uploaded file needs to be analyzed by chunker only onvce

## OneDrive - Google Drive and other cloud drives
The user can decide to coonect a cloud drive as knowledge source
It shall provide:
- the path to be analyzed (mandatory)
- if we need to scan the path only or all subfolders
- All the needed credentials to access the provided path

## MS Teams - Slack
The user can decide to upload a file as as knowledge source
It shall provide:
- the file to be analyzed (mandatory)



# Connector
- The connector will not perform any action to the vector database
- Max allowed single file size 1Gb
- All files, from any type of connector, are stored inside the MinIO instance

Connector table definition shall have a field status, with the following statuses
- 

## OneDrive - Google Drive and other cloud drives




# Chat