package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	writerInstanceID = os.Getenv("WRITER_INSTANCE_IDENTIFIER")
	clusterID        = os.Getenv("CLUSTER_IDENTIFIER")
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, snsEvent events.SNSEvent) {

	if !strings.Contains(snsEvent.Records[0].SNS.Message, FailoverEndedCode) {
		fmt.Println("skip")
		return
	}

	isWriter, err := IsWriterInstance(clusterID, writerInstanceID)
	if err != nil {
		panic(err)
	}

	if !isWriter {
		FailoverDBCluster(clusterID)
	}
}

func main() {
	lambda.Start(Handler)
}
