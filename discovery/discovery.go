package discovery

import (
	"context"

	billing "cloud.google.com/go/billing/apiv1"
	logging "cloud.google.com/go/logging/apiv2"
	monitoring "cloud.google.com/go/monitoring/apiv3"
	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
)

// GcpServices holds the initialized Google Cloud API clients
type GcpServices struct {
	Billing              *billing.CloudBillingClient
	CloudResourceManager *resourcemanager.ProjectsClient
	Monitoring           *monitoring.MetricClient
	Logging              *logging.ConfigClient
}

// InitGcpServices initializes and returns a GcpServices instance
func InitGcpServices(ctx context.Context) (*GcpServices, error) {
	// Initialize Logging client
	loggingClient, err := logging.NewConfigClient(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize Logging client
	billingClient, err := billing.NewCloudBillingClient(ctx)
	if err != nil {
		return nil, err
	}

	// Return the GcpServices instance
	return &GcpServices{
		Logging: loggingClient,
		Billing: billingClient,
	}, nil
}
