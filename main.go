package main

import (
	"flag"
	"fmt"
	"terraform-workspace-cleanup/cmd"
)

func main()  {
	flag.Usage = func() {
		fmt.Println("Usage: tf-workspace-cleanup -register-workspace=<workspace> -aws-account-id=12345678 -aws-iam-role=sirius-ci")
		fmt.Println("Usage: tf-workspace-cleanup -protected-workspaces=true -aws-account-id=12345678 -aws-iam-role=sirius-ci")
		flag.PrintDefaults()
	}
	var workspaceName string
	var protectedWorkspaces bool
	var awsAccountId string
	var awsIAMRoleName string

	flag.StringVar(&workspaceName, "register-workspace", "", "Register a workspace to be deleted at a later point")
	flag.StringVar(&awsAccountId, "aws-account-id", "", "Account ID for IAM Role")
	flag.StringVar(&awsIAMRoleName, "aws-iam-role", "", "AWS IAM Role Name ")
	flag.BoolVar(&protectedWorkspaces, "protected-workspaces", false, "get list of protected workspaces for deletion")
	flag.Parse()

	if awsAccountId == "" {
		fmt.Println("Error: You have not provided an AWS Account ID")
		flag.Usage()
	}

	if awsIAMRoleName == "" {
		fmt.Println("Error: You have not provided an AWS IAM Role Name")
		flag.Usage()
	}

	if protectedWorkspaces {
		cmd.RetrieveProtectedWorkspaces(&awsAccountId, &awsIAMRoleName)
	}

	if workspaceName != "" {
		cmd.RegisterWorkspace(&workspaceName, &awsAccountId, &awsIAMRoleName)
	} else if protectedWorkspaces != true {
		fmt.Println("Error: Workspace not passed")
		flag.Usage()
	}
}
