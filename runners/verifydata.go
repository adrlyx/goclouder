package runners

import (
	"context"
	"fmt"
	"log"

	"github.com/adrlyx/goclouder/discovery"
	"github.com/adrlyx/goclouder/helpers"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/status"

	loggingpb "cloud.google.com/go/logging/apiv2/loggingpb"
)

type VerifyData struct {
	projectID string
	sinkName  string
	sinkDest  string
	status    bool
}

func VerifyDataFunc() {
	ctx := context.Background()

	gcpServices, err := discovery.InitGcpServices(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize GcpServices: %v", err)
	}
	filepath := "files/input/verify_input"
	projectIds, err := helpers.ReadFile(filepath)
	if err != nil {
		fmt.Println("ERROR > :", err)
		return
	}

	var verifyData []VerifyData

	for _, projectID := range projectIds {
		parent := "projects/" + projectID
		req := &loggingpb.ListSinksRequest{
			Parent: parent,
		}
		it := gcpServices.Logging.ListSinks(ctx, req)
		defer gcpServices.Logging.Close()

		for {
			sink, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				if s, ok := status.FromError(err); ok {
					if s.Code() == 5 { // 5 = NotFound
						fmt.Printf("WARNING > Project does not exist %s\n", projectID)
						break
					} else {
						fmt.Printf("ERROR > %v\n", s.Message())
						return
					}
				}
			}
			var strToCheck string = "global"
			status := helpers.DoesStringContain(strToCheck, sink.Destination)

			verifyData = append(verifyData, VerifyData{
				projectID: projectID,
				sinkName:  sink.Name,
				sinkDest:  sink.Destination,
				status:    status,
			})
		}
		fmt.Printf("INFO > Done for %s\n", projectID)
	}
	var dataToWrite []string
	for _, sink := range verifyData {
		dataToWrite = append(dataToWrite, sink.projectID+" "+sink.sinkName+" "+sink.sinkDest+" "+fmt.Sprintf("%t", sink.status))
		fmt.Printf("%s %s %s %t\n", sink.projectID, sink.sinkName, sink.sinkDest, sink.status)
	}

	outputFile := "files/output/verifydata"
	test := helpers.WriteMapToFile(outputFile, dataToWrite)
	if test != nil {
		fmt.Println("ERROR > :", test)
		return
	}
	fmt.Println("LOG > Data written to file successfully")
}
