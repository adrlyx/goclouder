package helpers

// move to discovery/discovery.go (but renamed service/gcp.go?)

import (
	"context"
	"fmt"

	billingpb "cloud.google.com/go/billing/apiv1/billingpb"
	"github.com/adrlyx/goclouder/discovery"
)

// This is a client towards gcp
func CheckBillingAccount(ctx context.Context, projectID string, gcpService *discovery.GcpServices) bool {
	projectID = "projects/" + projectID
	req := &billingpb.GetProjectBillingInfoRequest{
		Name: projectID,
	}
	resp, err := gcpService.Billing.GetProjectBillingInfo(ctx, req)
	if err != nil {
		fmt.Printf("ERROR > Failed to get billing accounts : %v\n", err)
		return false
	}

	return resp.BillingEnabled
}
