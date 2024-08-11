package roles

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	dynamodb "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/jsii-runtime-go"
)

func CreatePBRHandlerRole(stack awscdk.Stack, log_group awslogs.LogGroup, user_table dynamodb.Table, discount_table dynamodb.Table) awsiam.Role {
	role := awsiam.NewRole(stack, jsii.String("PBR-Role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})

	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("logs:CreateLogGroup"),
			jsii.String("logs:PutLogEvents"),
			jsii.String("logs:DescribeLogStreams"),
			jsii.String("logs:CreateLogStream"),
			jsii.String("dynamodb:BatchGet*"),
			jsii.String("dynamodb:DescribeStream"),
			jsii.String("dynamodb:DescribeTable"),
			jsii.String("dynamodb:Get*"),
			jsii.String("dynamodb:Query"),
			jsii.String("dynamodb:Scan"),
			jsii.String("dynamodb:BatchWrite*"),
			jsii.String("dynamodb:CreateTable"),
			jsii.String("dynamodb:Delete*"),
			jsii.String("dynamodb:Update*"),
			jsii.String("dynamodb:PutItem"),
		},
		Resources: &[]*string{
			jsii.String(*user_table.TableArn()),
			jsii.String(*discount_table.TableArn()),
			jsii.String(*log_group.LogGroupArn()),
			jsii.String("*"),
		},
	}))
	return role
}
