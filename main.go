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
		fmt.Println("Usage: tf-workspace-cleanup -put=<workspace>")
		fmt.Println("Usage: tf-workspace-cleanup -expired-workspaces=true")
		fmt.Println("Usage: tf-workspace-cleanup -put=<workspace> -aws-account-id=12345678 -aws-iam-role=sirius-ci")
		flag.PrintDefaults()
	}
	var workspaceName string
	var expiredWorkspaces bool
	var awsAccountId string
	var awsIAMRoleName string

	flag.StringVar(&workspaceName, "put", "", "workspace to register for deletion at later time")
	flag.StringVar(&awsAccountId, "aws-account-id", "", "Account ID for IAM Role")
	flag.StringVar(&awsIAMRoleName, "aws-iam-role", "", "AWS IAM Role Name ")
	flag.BoolVar(&expiredWorkspaces, "expired-workspaces", false, "get list of expired workspaces for deletion")
	flag.Parse()

	if awsAccountId == "" {
		fmt.Println("Error: You have not provided an AWS Account ID")
		flag.Usage()
	}

	if awsIAMRoleName == "" {
		fmt.Println("Error: You have not provided an AWS IAM Role Name")
		flag.Usage()
	}

	if expiredWorkspaces {
		GetExpiredWorkspaces(&awsAccountId, &awsIAMRoleName)
	}

	if workspaceName != "" && expiredWorkspaces == true{
		PutWorkspace(&workspaceName, &awsAccountId, &awsIAMRoleName)
	} else {
		fmt.Println("Error: Workspace not passed")
		flag.Usage()
	}
}

func PutWorkspace(workspace *string, accountId *string, iamRoleName *string) {

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

func GetExpiredWorkspaces(accountId *string, iamRoleName *string) {
	type Item struct {
		WorkspaceName string
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalln(err)
	}

	RoleArn := "arn:aws:iam::" + *accountId + ":role/" + *iamRoleName

	creds := stscreds.NewCredentials(sess, RoleArn)
	awsConfig := aws.Config{Credentials: creds, Region: aws.String("eu-west-1")}

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
