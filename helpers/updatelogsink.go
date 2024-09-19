// move to service/gcp.go same as checkBillingAccount
package helpers

import (
	"context"
	"fmt"

	loggingpb "cloud.google.com/go/logging/apiv2/loggingpb"
	"github.com/adrlyx/goclouder/discovery"
)

func UpdateLogSink(ctx context.Context, project string, gcpService *discovery.GcpServices, newBucketLocation string, newBucketName string) (*loggingpb.LogSink, error) {
	dest := "logging.googleapis.com/projects/" + project + "/locations/" + newBucketLocation + "/buckets/" + newBucketName
	sinkName := "projects/" + project + "/sinks/_Default"
	req := &loggingpb.UpdateSinkRequest{
		SinkName: sinkName,
		Sink: &loggingpb.LogSink{
			Destination: dest,
		},
	}
	resp, err := gcpService.Logging.UpdateSink(ctx, req)
	if err != nil {
		fmt.Printf("ERROR > Failed to update sink : %v\n", err)
	}

	return resp, err
}
