package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
	"strconv"
	"time"
)

func main()  {
	flag.Usage = func() {
		fmt.Println("Usage: tf-workspace-cleanup -put <workspace>")
		fmt.Println("Usage: tf-workspace-cleanup -expired-workspaces true")
		flag.PrintDefaults()
	}
	var workspaceName string
	var expiredWorkspaces bool

	flag.StringVar(&workspaceName, "put", "", "workspace to register for deletion at later time")
	flag.BoolVar(&expiredWorkspaces, "expired-workspaces", false, "get list of expired workspaces for deletion")
	flag.Parse()

	if workspaceName == "" {
		fmt.Println("Error: Workspace not passed")
		flag.Usage()
	} else {
		PutWorkspace(&workspaceName)
	}

	if expiredWorkspaces {
		GetExpiredWorkspaces()
	}
}

func PutWorkspace(w *string) {

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalln(err)
	}
	RoleArn := ""
	// TODO - Allow AWS Account ID and Role name to be passed in via cli flag
	if len(os.Getenv("CI")) > 0 {
		RoleArn = "arn:aws:iam::288342028542:role/sirius-ci"
	} else {
		RoleArn = "arn:aws:iam::288342028542:role/operator"
	}

	creds := stscreds.NewCredentials(sess, RoleArn)
	awsConfig := aws.Config{Credentials: creds, Region: aws.String("eu-west-1")}

	svc := dynamodb.New(sess, &awsConfig)

	type Workspace struct {
		WorkspaceName string
		ExpiresTTL    int64
	}

	item := Workspace{
		WorkspaceName: *w,
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

func GetExpiredWorkspaces() {
	fmt.Println("Not yet implemented")
}