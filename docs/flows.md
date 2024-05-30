

# Orchestrator
TBD: 
 - Define general behavior 
 - How to have multiple orchestrator instances 

[test](https://www.pp.ch)

The orchestrator is responsible for monitoring the connector table in the relational database and deciding when it's time for a new analysis of the single connector
A single connector row represents a knowledge source the user decided to have analyzed by CogniX and stored in the Vetcor database.

The orchestrator is responsible for deciding if to scan again the source defined in the connector, and eventually sending analysis requests to both the Connector or Chunker, depending on the specific use case:
A new scan (connector) or analysis (chunker) shall be started only if there are no other processes running


The user shall decide to disable the connector; if so, no further scans will be started for that connector. 
Disabled is set in status
All date times are defined as UTC

## Status
The status of the connector for the moment is determined by the field (verify status field name in the connector table)
In the future, the status will be retrieved from NATS

Statuses:
- Active (connector just created)
- Pending scan (orchestrator sent a message to Connector or Chunker to scan the source)
- Working. Set by Connector and Orchestrator. This means that eighter the Connector or the Chunker is working on it
- Scan completed successfully, set by chunker
- Scan completed with errors, set by chunker
- Disabled  
- Unable to process, set by chunker. This status indicates there is a problem with the source. The orchestrator shall not request any more analysis for this resource. Only one (important only one) email shall be sent to the user reporting the issue

IMPORTANT:
TBD what to do when the row is in the pending or working status? 
We need to find our if there's a way for NATS to notify that a new item has been added to the dead letter queue. 
If so the orchestrator, or any other service, shall subscribe to this message (like we do for chunker, only one subscriber shall receive the message) and set the status to unable to process.
Again, only one (important, only one) email shall be sent to the user reporting the issue

## Rules to scan again any given connector line
- no active scan 
- no pending requested scan
- the time of the last activity shall be greater than the refresh_freq (in seconds) 
- refresh_freq is coming from config (env) and is different by filetype [see paragraph below]()

## URL 
The user can decide to connect a URL as a knowledge source
It shall provide:
- URL (mandatory)
- Sitemap URL (optional)
- Scan all the links found on the page (bit, optional)
- If the system shall look for a sitemap even if not provided

The Orchestrator, for this particular file type, will forward the request to Chunker.

For this particular file type, no file is stored in MinIO

If the user deletes the source, the API receiving the request will be responsible for soft deleting the row from the relational database and physically deleting all the entities related to that documentid inside the vector database (it is very important to hard delete from the vector database to reduce storage utilization) 
Soft delete means that in the connector table, we have a flag, deleted, and a deleted date


## File 
The user can decide to upload a file as a knowledge source
It shall provide:
- the file to be analyzed (mandatory)

The user will be able to upload a file to our internal MinIO
The file will be scanned only Once. Once the status of this row is scan completed (with or without errors) the Orchestrator will not issue any other analysis request. 

The Orchestrator will send a request to Chunker as soon as the file is uploaded correctly
The URL property of the proto message will contain all the information for Chunker to open the file from MinIO
The uploaded file needs to be analyzed by the chunker only once


## OneDrive - Google Drive and other cloud drives
The user can decide to connect a cloud drive as a knowledge source
It shall provide:
- the path to be analyzed (mandatory)
- if we need to scan the path only or all subfolders
- All the needed credentials to access the provided path

The Orchestrator, for this particular file type, will forward the request to Connector.

## MS Teams - Slack
The user can decide to upload a file as a knowledge source
It shall provide:
- the file to be analyzed (mandatory)

TBD:
The whole flow
How to deal with the table documents? 

## Refresh frequency
URL: one week
OneDrive - Google Drive and other cloud drives: one week
MS Teams - Slack: daily, only new messages

# Connector
- The connector will not perform any action to the vector database
- Max allowed single file size 1Gb
- All files, from any type of connector, are stored inside the MinIO instance

## URL
Connector shall never receive a message for a file, it shall be sent directly to Chuker

## File
Connector shall never receive a message for a file, it shall be sent directly to Chuker

## OneDrive - Google Drive and other cloud drives
When the Connector receives a new message from the orchestrator it will
- create a guid “chunking_session” that will be sent to each Chunking message sent by this operation. This way the Chunker will be able to understand when to set this process as completed
- Set the connector status from Pending scan to  Working by Connector. 
- It will scan the drive (given the rules from the orchestrator, all sub-folder or not) and get a list of path/file
- for each path/file item will check in the documents table if the item shall be sent to Chunker, depending on the hash comparison between database and file actually scanned.
- if the item needs to be scanned (because is new or updated) 
  - Update chunking_session with the new chunking_session, set the status to not done
  - send a message to chunked
If the document does not need to be scanned (because hash comparison identical) 
  - Update chunking_session with the new chunking_session, set the status to not done
  - send a message to chunked

- delete (physically) all the documents in the database that are not anymore present in the original source  

# Chunker
## URL
Connector shall never receive a message for a file, it shall be sent directly to Chuker

## File
Connector shall never receive a message for a file, it shall be sent directly to Chuker

## OneDrive - Google Drive and other cloud drives
guid “chunking_session”

Document Table
CREATE TABLE "document" (
  "id" varchar PRIMARY KEY NOT NULL
  “parent_id”, integer // Alow nulls used for URLs
  "connector_id" integer NOT NULL,
  "link" varchar,
  "last_update" timestamp,
  "signature" text
  “chunking_session” guid // allow nulls
  “chunking_status” chunking status done, not done (talk with gian about type, this is a lookup) 
);




