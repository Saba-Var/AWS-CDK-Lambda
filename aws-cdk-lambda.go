package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AwsCdkLambdaStackProps struct {
	awscdk.StackProps
}

func NewAwsCdkLambdaStack(scope constructs.Construct, id string, props *AwsCdkLambdaStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	table := awsdynamodb.NewTable(stack, jsii.String("userTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("username"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName:     jsii.String("userTable"),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	lambdaFunc := awslambda.NewFunction(stack, jsii.String("lambdaFunc"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_NODEJS_20_X(),
		Handler: jsii.String("main"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"), nil),
	})

	table.GrantReadWriteData(lambdaFunc)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewAwsCdkLambdaStack(app, "AwsCdkLambdaStack", &AwsCdkLambdaStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
