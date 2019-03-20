package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

const (
	FailoverEndedCode = "RDS-EVENT-0071"
)

// IsWriterInstance 特定のDBインスタンスがクラスター内でWriterの役割を持っているか？
func IsWriterInstance(clusterID string, instanceID string) (bool, error) {
	rdsClient := rds.New(session.New())
	input := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(clusterID),
	}

	clustersOutput, err := rdsClient.DescribeDBClusters(input)
	if err != nil {
		return false, err
	}

	cluster := clustersOutput.DBClusters[0]
	for _, instance := range cluster.DBClusterMembers {
		if aws.StringValue(instance.DBInstanceIdentifier) == instanceID {
			return aws.BoolValue(instance.IsClusterWriter), nil
		}
	}

	return false, nil
}

// FailoverDBCluster 特定のDBクラスターのフェイルオーバーを開始する
func FailoverDBCluster(clusterID string) error {
	rdsClient := rds.New(session.New())
	failInput := &rds.FailoverDBClusterInput{
		DBClusterIdentifier: aws.String(clusterID),
	}
	_, err := rdsClient.FailoverDBCluster(failInput)
	return err
}
