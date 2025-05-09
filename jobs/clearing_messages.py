import logging
import os
import pymysql
import boto3
import csv
import zipfile
from datetime import datetime, timedelta
from tqdm import tqdm

# Logging setup
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[logging.FileHandler("dump_to_s3.log"), logging.StreamHandler()]
)

# Database connection setup
DB_HOST = 'your_db_host'
DB_USER = 'your_db_user'
DB_PASSWORD = 'your_db_password'
DB_NAME = 'your_db_name'

# S3 setup
AWS_ACCESS_KEY = 'your_aws_access_key'
AWS_SECRET_KEY = 'your_aws_secret_key'
S3_BUCKET_NAME = 'your_s3_bucket_name'
S3_FOLDER = 'backup_folder/'  # Folder in the bucket

# Number of days to retain records in the database
RECORDS_OLDER_THAN_DAYS = 5

# Convert datetime to epoch (milliseconds)
def datetime_to_epoch_ms(dt):
    return int(dt.timestamp() * 1000)

# Connect to MySQL database
def get_db_connection():
    return pymysql.connect(
        host=DB_HOST,
        user=DB_USER,
        password=DB_PASSWORD,
        database=DB_NAME,
        cursorclass=pymysql.cursors.DictCursor
    )

# Fetch records older than 5 days based on epoch timestamp
def fetch_old_records():
    connection = get_db_connection()
    try:
        with connection.cursor() as cursor:
            five_days_ago = datetime.now() - timedelta(days=RECORDS_OLDER_THAN_DAYS)
            five_days_ago_epoch = datetime_to_epoch_ms(five_days_ago)
            sql = "SELECT * FROM messages WHERE ts < %s"
            cursor.execute(sql, (five_days_ago_epoch,))
            records = cursor.fetchall()
            return records
    finally:
        connection.close()

# Dump records to CSV
def dump_to_csv(records, csv_filename):
    if records:
        logging.info(f'Dumping {len(records)} records to {csv_filename}')
        with open(csv_filename, mode='w', newline='') as file:
            writer = csv.DictWriter(file, fieldnames=records[0].keys())
            writer.writeheader()
            for record in tqdm(records, desc="Writing records to CSV"):
                writer.writerow(record)

# Zip the CSV file
def zip_file(csv_filename, zip_filename):
    logging.info(f'Zipping {csv_filename} to {zip_filename}')
    with zipfile.ZipFile(zip_filename, 'w', zipfile.ZIP_DEFLATED) as zipf:
        zipf.write(csv_filename, os.path.basename(csv_filename))

# Upload the file to S3
def upload_to_s3(zip_filename):
    s3_client = boto3.client(
        's3',
        aws_access_key_id=AWS_ACCESS_KEY,
        aws_secret_access_key=AWS_SECRET_KEY
    )
    s3_key = S3_FOLDER + os.path.basename(zip_filename)
    logging.info(f'Uploading {zip_filename} to S3 bucket {S3_BUCKET_NAME}/{s3_key}')
    with tqdm(total=os.path.getsize(zip_filename), unit='B', unit_scale=True, desc="Uploading to S3") as pbar:
        s3_client.upload_file(zip_filename, S3_BUCKET_NAME, s3_key, Callback=lambda bytes_transferred: pbar.update(bytes_transferred))

# Delete records older than 5 days based on epoch timestamp
def delete_old_records():
    connection = get_db_connection()
    try:
        with connection.cursor() as cursor:
            five_days_ago = datetime.now() - timedelta(days=RECORDS_OLDER_THAN_DAYS)
            five_days_ago_epoch = datetime_to_epoch_ms(five_days_ago)
            sql = "DELETE FROM messages WHERE ts < %s"
            result = cursor.execute(sql, (five_days_ago_epoch,))
            connection.commit()
            logging.info(f'Deleted {result} records older than {RECORDS_OLDER_THAN_DAYS} days')
    finally:
        connection.close()

# Main job function
def dump_and_upload_job():
    logging.info("Job started")

    records = fetch_old_records()
    if not records:
        logging.info("No records found older than 5 days")
        return

    # Filenames for CSV and Zip
    csv_filename = f"messages_{datetime.now().strftime('%Y%m%d')}.csv"
    zip_filename = csv_filename.replace(".csv", ".zip")

    dump_to_csv(records, csv_filename)
    zip_file(csv_filename, zip_filename)
    upload_to_s3(zip_filename)

    delete_old_records()

    # Clean up local files
    os.remove(csv_filename)
    os.remove(zip_filename)

    logging.info("Job completed successfully")

if __name__ == "__main__":
    dump_and_upload_job()
