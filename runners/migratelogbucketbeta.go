package runners

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	loggingpb "cloud.google.com/go/logging/apiv2/loggingpb"
	"github.com/adrlyx/goclouder/discovery"
	"github.com/adrlyx/goclouder/helpers"
	"google.golang.org/grpc/status"
)

var (
	filepath = "files/input/change_log_bucket_input"
)

func processPartProjects(part map[int]string, results chan<- map[int]string, id int, billingAccount string, newBucketName string, newBucketLocation string) {
	// Initialize context
	ctx := context.Background()

	// Initialize GCP services
	gcpServices, err := discovery.InitGcpServices(ctx)
	if err != nil {
		log.Fatalf("ERROR > Failed to initialize GcpServices: %v", err)
	}

	processedProjectsPart := make(map[int]string)

	// Ensure that the function returns the results regardless of the outcome
	defer func() {
		defer gcpServices.Logging.Close()
		defer gcpServices.Billing.Close()
		results <- processedProjectsPart
	}()

	for k, v := range part {
		testBilling := helpers.CheckBillingAccount(ctx, v, gcpServices)
		if !testBilling {
			fmt.Printf("INFO > %s > No billing account set\n", v)
			// Loop until the billing account is enabled
			for {
				if helpers.CheckBillingAccount(ctx, v, gcpServices) {
					fmt.Printf("INFO > %s > Billing account enabled\n", v)
					break
				}
				err := helpers.EnableBillingAccount(ctx, v, billingAccount, gcpServices)
				if err != nil {
					fmt.Printf("ERROR > Failed to enable billing account for %s: %v\n", v, err)
				}
				fmt.Printf("INFO > %s > Billing account is not yet enabled. Checking again in 2 seconds...\n", v)
				time.Sleep(2 * time.Second) // Wait for 10 seconds before checking again
			}
		} else {
			fmt.Printf("INFO > %s > Billing account check passed\n", v)
		}
		processedProjectsPart[k], err = migrateBucket(v, ctx, gcpServices, billingAccount, newBucketName, newBucketLocation)
		if err != nil {
			fmt.Printf("ERROR > %s > Quota exceeded\n", v)
			break
		}
		fmt.Printf("INFO > %s > Processed project\n", v)
	}
	results <- processedProjectsPart
}

func migrateBucket(projectID string, ctx context.Context, gcpServices *discovery.GcpServices, billingAccount string, newBucketName string, newBucketLocation string) (str string, e error) {
	// Create log bucket
	parent := "projects/" + projectID + "/locations/" + newBucketLocation
	req := &loggingpb.CreateBucketRequest{
		Parent:   parent,
		BucketId: newBucketName,
	}
	resp, err := gcpServices.Logging.CreateBucket(ctx, req)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Message() == ("Valid linked billing account is required\n") {
				fmt.Printf("ERROR > %s >  Billing account is required\n", projectID)
				return "Billing account is required", nil
			} else {
				fmt.Printf("ERROR > %v\n", s.Message())
			}
		}
	} else {
		fmt.Printf("INFO > %v > Bucket created\n", resp.Name)
	}

	// Update sink destination
	updateSink, err := helpers.UpdateLogSink(ctx, projectID, gcpServices, newBucketLocation, newBucketName)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == 8 { // 5 = NotFound found here "google.golang.org/grpc/codes"
				return projectID, err
			} else {
				fmt.Printf("ERROR > %v\n", s.Message())
			}
		}
	}
	fmt.Printf("INFO > %s > Sink updated for project %s\n", projectID, updateSink.Destination)

	return projectID, nil
}

func MigrateLogBucket(newBucketName string, newBucketLocation string, billingAccount string) {
	fmt.Println("LOG > Running change log bucket function with 4 channels")

	// Read project IDs from file
	projectIds, err := helpers.ReadFile(filepath)
	if err != nil {
		fmt.Println("ERROR > :", err)
		return
	}

	// Extract keys and sort them
	keys := make([]int, 0, len(projectIds))
	for k := range projectIds {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Determine size of each part
	n := len(projectIds)
	partSize := (n + 3) / 4

	// Create channels
	results := make(chan map[int]string, 4)

	// Start 4 goroutines
	for i := 0; i < 4; i++ {
		start := i * partSize
		end := start + partSize
		if end > n {
			end = n
		}

		part := make(map[int]string)
		for _, k := range keys[start:end] {
			part[k] = projectIds[k]
		}

		go processPartProjects(part, results, i, billingAccount, newBucketName, newBucketLocation)
	}

	// Collect results
	finalResults := make(map[int]string)
	for i := 0; i < 4; i++ {
		partResult := <-results
		for k, v := range partResult {
			finalResults[k] = v
		}
	}

	fmt.Println("Processed Projects:")
	for _, v := range finalResults {
		fmt.Println(v)
	}
}
