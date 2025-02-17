import boto3
session = boto3.Session(profile_name="<profile_name>")
credentials = session.get_credentials().get_frozen_credentials()
print(f"AWS_ACCESS_KEY_ID={credentials.access_key}")
print(f"AWS_SECRET_ACCESS_KEY={credentials.secret_key}")
print(f"AWS_SESSION_TOKEN={credentials.token}")

