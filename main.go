package main

import (
	"flag"
	"fmt"
	"terraform-workspace-cleanup/cmd"
)

func main()  {
	flag.Usage = func() {
		fmt.Println("Usage: tf-workspace-cleanup -register-workspace=<workspace> -aws-account-id=12345678 -aws-iam-role=sirius-ci")
		fmt.Println("Usage: tf-workspace-cleanup -expired-workspaces=true -aws-account-id=12345678 -aws-iam-role=sirius-ci")
		flag.PrintDefaults()
	}
	var workspaceName string
	var expiredWorkspaces bool
	var awsAccountId string
	var awsIAMRoleName string

	flag.StringVar(&workspaceName, "register-workspace", "", "Register a workspace to be deleted at a later point")
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
		cmd.RetrieveExpiredWorkspaces(&awsAccountId, &awsIAMRoleName)
	}

	if workspaceName != "" {
		cmd.RegisterWorkspace(&workspaceName, &awsAccountId, &awsIAMRoleName)
	} else if expiredWorkspaces != true {
		fmt.Println("Error: Workspace not passed")
		flag.Usage()
	}
}
