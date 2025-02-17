# go_s3_buckets


## Clean, Build, Run, and other Stuff

  ```
	$: go build -o bin/s3_list_objects cmd/s3_list_objects/main.go
	$:
	$:
	$: bin/s3_list_objects --bucket <bucket_name> --region us-east-1 --profile default
	S3 Objects:
  ```

## Utilities

  1. bin/s3_list_objects
    - requires s3 bucket name
    - optionally accepts s3 bucket prefix
    - optional filter on filetype (e.g. ".txt", ".pdf")
    - optional region, defaults to us-east-1
    - optional profile name, specify an aws profile name or default

## file hierarchy and more stuff

  ```
	$: tree
	.
	├── bin
	│   ├── s3_list_objects
	│   └── s3_upload_tracker
	├── cmd
	│   ├── s3_list_objects
	│   │   └── main.go
	│   └── s3_upload_tracker
	│       └── main.go
	├── db
	│   └── database.go
	├── go.mod
	├── go_s3_buckets_list_objects
	├── go.sum
	├── internal
	│   └── s3utils
	│       └── s3.go
	├── policy.json
	├── python_auth_hdrs.py
	├── python_credentials.py
	├── README.md
	├── sqlite_extraction_tracker.go
	└── uploads.db
	
	7 directories, 15 files
  ```

## Database Information

I'm using sqlite to test.  Move to postres later.

Just a simple, single table (for now) to track meta data of local files pushed to s3 bucket. 
Information such as file name, bucket name, date_time stamp, and hash of file (for dedup) is maintained here.

### uploads table

  ```
    $: sqlite3 uploads.db .tables
    uploads

    $: sqlite3 uploads.db .schema
    CREATE TABLE uploads (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    file_name TEXT,
		    s3_prefix TEXT,
		    s3_key TEXT UNIQUE,
		    file_hash TEXT UNIQUE,
		    timestamp TEXT,
		    allowed_users TEXT
	    );
    CREATE TABLE sqlite_sequence(name,seq);
  ```

## Ancillary Stuff

  - python_auth_hdrs.py
    simple, debug, used to view aws auth header when trying to compare go sdk with python boto3

  - python_credentials.py
    debug and validation, used to view user credentials, from an aws perspective.  What are my credentials when running



