package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adrlyx/goclouder/runners"
)

var (
	changeBucketInfo = "Execute change log bucket function\n" +
		"Add GCP Project IDs in a list in files/input/change_log_bucket_input.\n" +
		"The program will create a bucket with the name you gave it in location of your choice.\n" +
		"It will then update the _Default sink in the project to this new bucket\n" +
		"This functionality is because we had to migrate all of our _Default log buckets to eu from global in GCP.\n" +
		"Required flags:\n" +
		"-billing-account: Name of billing account, e.g. billingAccounts/XXXXX-XXXXX-XXXXX\n" +
		"-new-bucket-name: Name of the new bucket\n" +
		"-new-bucket-location: Location of the new bucket\n\n"
	verifyInfo = "Execute verify function.\n" +
		"Add GCP Project IDs in a list in files/input/verify_input.\n" +
		"The program will go through all projects in the file and get all log sinks in this format:\n" +
		"<project-id> <sink-name>, <sink-destination> <sink-status>\n\n" +
		"The status will give true if the sink destination is pointing to a sink that has global as location\n" +
		"The data will be outputed to files/output/verifydata.\n\n"
	getProjectInfoChannelTest = "Execute get project id channel test function\n"
	newBucketName             string
	newBucketLoc              string
	billingAccount            string
)

func init() {
	flag.StringVar(&billingAccount, "billing-account", "", "Name of billing account, e.g. billingAccounts/XXXXX-XXXXX-XXXXX\n")
	flag.StringVar(&newBucketName, "new-bucket-name", "", "Location of the new bucket, e.g. eu (Used with -change-log-bucket)\n")
	flag.StringVar(&newBucketLoc, "new-bucket-location", "", "Location of the new bucket, e.g. my-bucket (Used with -change-log-bucket)\n")
	flag.Bool("change-log-bucket", false, changeBucketInfo)
	flag.Bool("verify-log-sinks", false, verifyInfo)
	flag.Bool("get-project-info-channel-test", false, getProjectInfoChannelTest)
	flag.Parse()
}

func main() {
	if flag.Lookup("change-log-bucket").Value.(flag.Getter).Get().(bool) {
		if newBucketName == "" || newBucketLoc == "" || billingAccount == "" {
			fmt.Println("Please provide both --new-bucket-name and --new-bucket-location and --billing-account")
			fmt.Println("Use -h for help")
			os.Exit(1)
		}
		runners.MigrateLogBucket(newBucketName, newBucketLoc, billingAccount)
	} else if flag.Lookup("verify-log-sinks").Value.(flag.Getter).Get().(bool) {
		runners.VerifyDataFunc()
	} else if flag.Lookup("get-project-info-channel-test").Value.(flag.Getter).Get().(bool) {
		runners.GetProjectInfoChannelTest()
	} else {
		fmt.Println("Please provide a flag")
		flag.Usage()
		os.Exit(1)
	}
}
