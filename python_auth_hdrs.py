import boto3
import logging

logging.basicConfig(level=logging.DEBUG)
s3 = boto3.client('s3', region_name='us-east-1')
s3.list_objects_v2(Bucket="<bucket_name>")

