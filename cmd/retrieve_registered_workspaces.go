package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
)

type Item struct {
	WorkspaceName string
}

func RetrieveProtectedWorkspaces(accountId *string, iamRoleName *string, assumeRole *bool) {

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalln(err)
	}

	awsConfig := aws.Config{Region: aws.String("eu-west-1")}

	if *assumeRole {
		RoleArn := "arn:aws:iam::" + *accountId + ":role/" + *iamRoleName
		creds := stscreds.NewCredentials(sess, RoleArn)
		awsConfig.Credentials = creds
	}

	svc := dynamodb.New(sess, &awsConfig)

	params := &dynamodb.ScanInput{
		TableName: aws.String("WorkspaceCleanup"),
	}

	result, err := svc.Scan(params)
	if err != nil {
		exitWithError(fmt.Errorf("failed to make Query API call, %v", err))
	}

	items := []Item{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		exitWithError(fmt.Errorf("failed to unmarshal Query result items, %v", err))
	}

	for i, item := range items {
		fmt.Print(item.WorkspaceName, " ")
		i++
	}
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
