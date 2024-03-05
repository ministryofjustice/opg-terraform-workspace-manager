# terraform-workspace-cleanup

Forked from [this code base](https://github.com/TomTucka/terraform-workspace-manager).

## Use Case

Sometimes when building products, we want environments for developers to be built on-demand or on the opening of pull requests. However,
we need an automated way to destroy these environments, so they don't sit around costing us $$$$

So what is the purpose of this? Well, when you create a new environment using a terraform workspace using this you can add the workspace you have
created in your automated pipeline into a DynamoDB table with a TTL value. You can then set up an automated job to also use this tool to
pull protected workspaces from DynamoDB

## Creating the DynamoDB table

We have provided a Terraform Module that will create a DynamoDB table in the configuration required by our tool.
To use our module you can use the following snippet:

```
module "workspace-cleanup" {
  source  = "git@github.com:ministryofjustice/opg-terraform-workspace-manager/terraform/workspace_cleanup"
  enabled = true
}
```

## Using this tool

```
Usage: terraform-workspace-manager -register-workspace=<workspace> -time-to-protect=2 -aws-account-id=12345678 -aws-iam-role=operator
Usage: terraform-workspace-manager -protected-workspaces=true -aws-account-id=12345678 -aws-iam-role=operator
  -aws-account-id string
    	Account ID for IAM Role
  -aws-iam-role string
    	AWS IAM Role Name
  -protected-workspaces
    	get list of protected workspaces for deletion
  -register-workspace string
    	Register a workspace to be deleted at a later point
  -time-to-protect
        Time in hours to protect workspace for
```
