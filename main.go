package main

import (
	"flag"
	"fmt"
	"terraform-workspace-cleanup/cmd"
)

func main()  {
	flag.Usage = func() {
		fmt.Println("Usage: tf-workspace-cleanup -register-workspace=<workspace> -time-to-protect=2 -aws-account-id=12345678 -aws-iam-role=sirius-ci")
		fmt.Println("Usage: tf-workspace-cleanup -protected-workspaces=true -aws-account-id=12345678 -aws-iam-role=sirius-ci")
		flag.PrintDefaults()
	}
	var workspaceName string
	var protectedWorkspaces bool
	var awsAccountId string
	var awsIAMRoleName string
	var timeToProtect int64
	var assumeRole bool

	flag.StringVar(&workspaceName, "register-workspace", "", "Register a workspace to be deleted at a later point")
	flag.StringVar(&awsAccountId, "aws-account-id", "", "Account ID for IAM Role")
	flag.StringVar(&awsIAMRoleName, "aws-iam-role", "", "AWS IAM Role Name")
	flag.Int64Var(&timeToProtect, "time-to-protect", 1 , "Time in hours to protect workspace for")
	flag.BoolVar(&protectedWorkspaces, "protected-workspaces", false, "get list of protected workspaces for deletion")
	flag.BoolVar(&assumeRole, "assume-role", true, "whether to assume the passed role rather than use calling creds")
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
		cmd.RetrieveProtectedWorkspaces(&awsAccountId, &awsIAMRoleName, &assumeRole)
	}

	if workspaceName != "" {
		cmd.RegisterWorkspace(&workspaceName, &awsAccountId, &awsIAMRoleName, &timeToProtect, &assumeRole)
	} else if protectedWorkspaces != true {
		fmt.Println("Error: Workspace not passed")
		flag.Usage()
	}
}
