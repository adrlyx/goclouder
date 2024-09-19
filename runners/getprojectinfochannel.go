package runners

import (
	"context"
	"fmt"
	"log"
	"sort"

	billingpb "cloud.google.com/go/billing/apiv1/billingpb"
	"github.com/adrlyx/goclouder/discovery"
	"github.com/adrlyx/goclouder/helpers"
)

func processPart(part map[int]string, results chan<- map[int]string, id int, CheckBillingAccount func(string) string) {
	processedPart := make(map[int]string)
	for k, v := range part {
		processedPart[k] = CheckBillingAccount(v)
		fmt.Println("Processed part: ", v)
	}
	results <- processedPart
}

func CheckBillingAccount(projectID string) string {
	ctx := context.Background()

	gcpService, err := discovery.InitGcpServices(ctx)
	if err != nil {
		log.Fatalf("ERROR > Failed to initialize GcpServices: %v", err)
	}

	projectID = "projects/" + projectID
	req := &billingpb.GetProjectBillingInfoRequest{
		Name: projectID,
	}
	resp, err := gcpService.Billing.GetProjectBillingInfo(ctx, req)
	defer gcpService.Billing.Close()
	if err != nil {
		fmt.Printf("ERROR > Failed to get billing accounts : %v\n", err)
	}
	if resp.BillingEnabled == true {
		return resp.BillingAccountName
	} else {
		return "No billing account enabled"
	}
}

func GetProjectInfoChannelTest() {
	fmt.Println("GetProjectInfoChannelTest")

	filepath := "files/input/channel_test_input"
	projectIds, err := helpers.ReadFile(filepath)
	if err != nil {
		fmt.Println("ERROR > :", err)
		return
	}

	for k, projectID := range projectIds {
		fmt.Println("Project ID: ", k, projectID)
	}

	// Extract keys and sort them
	keys := make([]int, 0, len(projectIds))
	for k := range projectIds {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// determine the size of each part
	n := len(projectIds)
	partSize := (n + 3) / 4 // +3 to handle cases where n is not a multiple of 4

	// Create channel
	results := make(chan map[int]string, 4)

	// start 4 goroutines
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

		go processPart(part, results, i, CheckBillingAccount)
	}

	// Collect results
	finalResult := make(map[int]string)
	for i := 0; i < 4; i++ {
		partResult := <-results
		for k, v := range partResult {
			finalResult[k] = v
		}
	}

	fmt.Println("Final result:")
	for _, k := range finalResult {
		fmt.Println(k)
	}
}
