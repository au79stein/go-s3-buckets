# go_s3_buckets


## Clean, Build, Run, and other Stuff

  ```
	$: go build -o bin/s3_list_objects cmd/s3_list_objects/main.go
	$:
	$:
	$: bin/s3_list_objects --bucket <bucket_name> --region us-east-1 --profile default
	S3 Objects:
  ```

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
