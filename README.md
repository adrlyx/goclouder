# GCP Log Bucket Management Tool

This tool is designed to help manage Google Cloud Platform (GCP) log buckets. It provides functionalities to migrate log buckets, verify log sinks, and test project information retrieval using channels.

## Features

1. **Change Log Bucket**:

   - Creates a new log bucket in a specified location.
   - Updates the `_Default` sink in the project to point to the new bucket.
   - Useful for migrating log buckets to a different location (e.g., from global to EU).

2. **Verify Log Sinks**:

   - Verifies the log sinks for a list of GCP projects.
   - Outputs the sink details and their status to a file.

3. **Get Project Info Channel Test**:
   - Tests the retrieval of project information using channels.

## Prerequisites

- Go 1.16 or later
- GCP credentials with appropriate permissions

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/adrlyx/goclouder.git
   cd goclouder
   ```

2. Build the program:
   ```sh
   go build -o gcp-log-bucket-tool ./main
   ```

## Usage

### Change Log Bucket

To change the log bucket for a list of GCP projects:

1. Create a file `files/input/change_log_bucket_input` and add the GCP Project IDs, one per line.

2. Run the command:

   ```sh
   ./gcp-log-bucket-tool -change-log-bucket -billing-account=<BILLING_ACCOUNT> -new-bucket-name=<NEW_BUCKET_NAME> -new-bucket-location=<NEW_BUCKET_LOCATION>
   ```

   Example:

   ```sh
   ./gcp-log-bucket-tool -change-log-bucket -billing-account=billingAccounts/12345-67890-ABCDE -new-bucket-name=my-new-bucket -new-bucket-location=eu
   ```

### Verify Log Sinks

To verify the log sinks for a list of GCP projects:

1. Create a file `files/input/verify_input` and add the GCP Project IDs, one per line.

2. Run the command:

   ```sh
   ./gcp-log-bucket-tool -verify-log-sinks
   ```

   The output will be saved to `files/output/verifydata`.

### Get Project Info Channel Test

To test the retrieval of project information using channels:

1. Run the command:
   ```sh
   ./gcp-log-bucket-tool -get-project-info-channel-test
   ```

## Flags

- `-billing-account`: Name of the billing account (required for `-change-log-bucket`).
- `-new-bucket-name`: Name of the new bucket (required for `-change-log-bucket`).
- `-new-bucket-location`: Location of the new bucket (required for `-change-log-bucket`).
- `-change-log-bucket`: Execute the change log bucket function.
- `-verify-log-sinks`: Execute the verify log sinks function.
- `-get-project-info-channel-test`: Execute the get project info channel test function.

## Help

For help, run:

```sh
./gcp-log-bucket-tool -h
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---
