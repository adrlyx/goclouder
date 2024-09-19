// rename to service (this "service" concept I talk about is a
// reference to "a service" in DDD, domain driven design)
package discovery

// GcpService
// The purpose of this struct + methods would be to enable functionalty
// where you interract with GCP.
//
// helpers.CheckBillingAccount should be part of this "class"
//
// InitGcpServices should be renamed to NewGcpService and "initialize"
// the GcpService struct.
//

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
