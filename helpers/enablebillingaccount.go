// same as checkBillingAccount. It should be part of the service
package helpers

import (
	"context"
	"fmt"

	billingpb "cloud.google.com/go/billing/apiv1/billingpb"
	"github.com/adrlyx/goclouder/discovery"
)

func EnableBillingAccount(ctx context.Context, projectID string, billingAccount string, gcpService *discovery.GcpServices) error {
	projectID = "projects/" + projectID
	req := &billingpb.UpdateProjectBillingInfoRequest{
		Name: projectID,
		ProjectBillingInfo: &billingpb.ProjectBillingInfo{
			BillingAccountName: billingAccount,
		},
	}
	resp, err := gcpService.Billing.UpdateProjectBillingInfo(ctx, req)
	if err != nil {
		fmt.Printf("ERROR > Failed to set billing account : %v\n", err)
		return err
	}
	if !resp.BillingEnabled {
		return fmt.Errorf("billing not enabled")
	}
	fmt.Printf("INFO > Billing account enabled for project : %s %t\n", projectID, resp.BillingEnabled)
	return nil
}
