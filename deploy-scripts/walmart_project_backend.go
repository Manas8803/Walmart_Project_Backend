package main

import (
	"fmt"
	"os"

	"github.com/Manas8803/Walmart_Project_Backend/roles"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	dynamodb "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type WalmartProjectBackendStackProps struct {
	awscdk.StackProps
}

const stack_name = "WalmartProjectStack"

func NewWalmartProjectBackendStack(scope constructs.Construct, id string, props *WalmartProjectBackendStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	//^ Discount-TABLE
	discount_table := dynamodb.NewTable(stack, jsii.String(fmt.Sprintf("%s-Discount-Table", stack_name)), &dynamodb.TableProps{
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("beacon_id"),
			Type: dynamodb.AttributeType_STRING,
		},
		TableName: jsii.String(fmt.Sprintf("%s-Discount_Table", stack_name)),
	})

	//^ User-TABLE
	user_table := dynamodb.NewTable(stack, jsii.String(fmt.Sprintf("%s-User-Table", stack_name)), &dynamodb.TableProps{
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("user_id"),
			Type: dynamodb.AttributeType_STRING,
		},
		TableName: jsii.String(fmt.Sprintf("%s-User_Table", stack_name)),
	})

	//^ Log group of Proximity_Beacon_Service handler
	logGroup_pbr := awslogs.NewLogGroup(stack, jsii.String("PBR_Service-LogGroup"), &awslogs.LogGroupProps{
		LogGroupName:  jsii.String(fmt.Sprintf("/aws/lambda/%s-PBR_Service", stack_name)),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	//^ Proximity_Beacon_Service handler
	pbr_handler := awslambda.NewFunction(stack, jsii.String("PBR_Service-Lambda"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("../pbr-service"), nil),
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Timeout: awscdk.Duration_Seconds(jsii.Number(10)),
		Role:    roles.CreatePBRHandlerRole(stack, logGroup_pbr, user_table, discount_table),
		Environment: &map[string]*string{
			"REGION":               jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
			"DISCOUNT_TABLE_ARN":   jsii.String(*discount_table.TableArn()),
			"USER_TABLE_ARN":       jsii.String(*user_table.TableArn()),
			"NOTIFY_WEBSOCKET_URL": jsii.String(os.Getenv("NOTIFY_WEBSOCKET_URL")),
		},
		FunctionName: jsii.String(fmt.Sprintf("%s-PBR_Service-Lambda", stack_name)),
		LogGroup:     logGroup_pbr,
	})

	awsapigateway.NewLambdaRestApi(stack, jsii.String(fmt.Sprintf("%s-PBR_Service-Gateway", stack_name)), &awsapigateway.LambdaRestApiProps{
		Handler: pbr_handler,
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins: awsapigateway.Cors_ALL_ORIGINS(),
			AllowMethods: awsapigateway.Cors_ALL_METHODS(),
			AllowHeaders: awsapigateway.Cors_DEFAULT_HEADERS(),
		},
	})

	//~ WEBSOCKET API :
	//^ Connect Route
	// logGroup_connect := awslogs.NewLogGroup(stack, jsii.String("Connect-LogGroup"), &awslogs.LogGroupProps{
	// 	LogGroupName:  jsii.String(fmt.Sprintf("/aws/lambda/%s-Connect", stack_name)),
	// 	RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	// })
	awslambda.NewFunction(stack, jsii.String("Connect-Lambda"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("../websocket/connect"), nil),
		Runtime: awslambda.Runtime_NODEJS_16_X(),
		Handler: jsii.String("index.handler"),
		Timeout: awscdk.Duration_Seconds(jsii.Number(10)),
		Environment: &map[string]*string{
			"REGION":             jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
			"DISCOUNT_TABLE_ARN": jsii.String(*discount_table.TableArn()),
			"USER_TABLE_ARN":     jsii.String(*user_table.TableArn()),
		},
		FunctionName: jsii.String(fmt.Sprintf("%s-Connect-Lambda", stack_name)),
		Role:         roles.CreateWebSocketLambdaRole(stack, "Connect", stack_name),
		// LogGroup:     logGroup_connect,
	})

	//^ Disconnect Route
	// logGroup_disconnect := awslogs.NewLogGroup(stack, jsii.String("Disconnect-LogGroup"), &awslogs.LogGroupProps{
	// 	LogGroupName:  jsii.String(fmt.Sprintf("/aws/lambda/%s-Connect", stack_name)),
	// 	RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	// })
	awslambda.NewFunction(stack, jsii.String("Disconnect-Lambda"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("../websocket/disconnect"), nil),
		Runtime: awslambda.Runtime_NODEJS_16_X(),
		Handler: jsii.String("index.handler"),
		Timeout: awscdk.Duration_Seconds(jsii.Number(10)),
		Environment: &map[string]*string{
			"REGION":             jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
			"DISCOUNT_TABLE_ARN": jsii.String(*discount_table.TableArn()),
			"USER_TABLE_ARN":     jsii.String(*user_table.TableArn()),
		},
		FunctionName: jsii.String(fmt.Sprintf("%s-Disconnect-Lambda", stack_name)),
		Role:         roles.CreateWebSocketLambdaRole(stack, "Disconnect", stack_name),
		// LogGroup:     logGroup_disconnect,
	})

	//^ Notify Route
	logGroup_notify := awslogs.NewLogGroup(stack, jsii.String("Notify-LogGroup"), &awslogs.LogGroupProps{
		LogGroupName:  jsii.String(fmt.Sprintf("/aws/lambda/%s-Notify", stack_name)),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
	awslambda.NewFunction(stack, jsii.String("Notify-Lambda"), &awslambda.FunctionProps{
		Code:    awslambda.Code_FromAsset(jsii.String("../websocket/notify"), nil),
		Runtime: awslambda.Runtime_NODEJS_16_X(),
		Handler: jsii.String("index.handler"),
		Timeout: awscdk.Duration_Seconds(jsii.Number(10)),
		Environment: &map[string]*string{
			"REGION":             jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
			"DISCOUNT_TABLE_ARN": jsii.String(*discount_table.TableArn()),
			"USER_TABLE_ARN":     jsii.String(*user_table.TableArn()),
		},
		FunctionName: jsii.String(fmt.Sprintf("%s-Notify-Lambda", stack_name)),
		Role:         roles.CreateWebSocketLambdaRole(stack, "Notify", stack_name),
		LogGroup:     logGroup_notify,
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewWalmartProjectBackendStack(app, stack_name, &WalmartProjectBackendStackProps{
		awscdk.StackProps{
			Env:       env(),
			StackName: jsii.Sprintf(stack_name),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {

	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
