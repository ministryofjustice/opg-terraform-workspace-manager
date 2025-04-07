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

func RegisterWorkspace(workspace *string, accountId *string, iamRoleName *string, timeToProtect *int64, assumeRole *bool) {

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

	type Workspace struct {
		WorkspaceName string
		ExpiresTTL    int64
	}


	item := Workspace{
		WorkspaceName: *workspace,
		ExpiresTTL:    time.Now().Add(time.Hour * time.Duration(*timeToProtect)).Unix(),

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
