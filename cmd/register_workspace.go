package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	_ "github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
	"strconv"
	"time"
)

type Session struct {
	AwsSession *session.Session
}

func RegisterWorkspace(workspace *string, accountId *string, iamRoleName *string) {

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalln(err)
	}

	RoleArn := "arn:aws:iam::" + *accountId + ":role/" + *iamRoleName

	creds := stscreds.NewCredentials(sess, RoleArn)
	awsConfig := aws.Config{Credentials: creds, Region: aws.String("eu-west-1")}

	svc := dynamodb.New(sess, &awsConfig)

	type Workspace struct {
		WorkspaceName string
		ExpiresTTL    int64
	}

	item := Workspace{
		WorkspaceName: *workspace,
		ExpiresTTL:    time.Now().AddDate(0, 0, 1).Unix(),
	}

	WorkspaceToPut, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling Workspace item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		Item:      WorkspaceToPut,
		TableName: aws.String("WorkspaceCleanup"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Successfully added '" + item.WorkspaceName + "' with TTL " + strconv.FormatInt(item.ExpiresTTL, 10) + " for workspace cleanup")
}
