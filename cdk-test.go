package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"

	"os"
)

type CdkTestStackProps struct {
	awscdk.StackProps
}

func NewCdkTestStack(scope constructs.Construct, id string, props *CdkTestStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here
	key := awsec2.NewCfnKeyPair(stack, jsii.String("Test"), &awsec2.CfnKeyPairProps{
		KeyName: jsii.String("Test"),
	})

	ami := awsec2.MachineImage_Lookup(&awsec2.LookupMachineImageProps{
		Name:   jsii.String("whslabs-*"),
		Owners: jsii.Strings("102933037533"),
	})

	vpc := awsec2.NewVpc(stack, jsii.String("Vpc"), &awsec2.VpcProps{
		Cidr: jsii.String("10.0.0.0/16"),
	})

	sg := awsec2.NewSecurityGroup(stack, jsii.String("Sg"), &awsec2.SecurityGroupProps{
		Vpc:              vpc,
		AllowAllOutbound: jsii.Bool(true),
	})

	sg.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(22)), jsii.String(""), jsii.Bool(false))

	instance := awsec2.NewInstance(stack, jsii.String("Instance"), &awsec2.InstanceProps{
		InstanceType: awsec2.NewInstanceType(jsii.String("t3.large")),
		MachineImage: ami,
		Vpc:          vpc,
		KeyName:      key.KeyName(),
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PUBLIC,
		},
		SecurityGroup: sg,
	})

	// example resource
	// queue := awssqs.NewQueue(stack, jsii.String("CdkTestQueue"), &awssqs.QueueProps{
	// 	VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	// })

	awscdk.NewCfnOutput(stack, jsii.String("Output1"), &awscdk.CfnOutputProps{
		Value: instance.InstancePublicIp(),
	})

	awscdk.NewCfnOutput(stack, jsii.String("Output2"), &awscdk.CfnOutputProps{
		Value: key.AttrKeyPairId(),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewCdkTestStack(app, "CdkTestStack", &CdkTestStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
