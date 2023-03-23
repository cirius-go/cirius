package awsutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func LoadEnvFromAPS(stage, appName string, nextToken *string) error {
	basePath := fmt.Sprintf("%s/%s/", appName, stage)

	svc := ssm.New(session.New())
	resp, err := svc.GetParametersByPath(&ssm.GetParametersByPathInput{
		NextToken:      nextToken,
		Path:           aws.String(basePath),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	for _, param := range resp.Parameters {
		paramName := strings.Replace(*param.Name, basePath, "", 1)
		os.Setenv(strings.ToUpper(paramName), *param.Value)
	}

	if resp.NextToken != nil {
		return LoadEnvFromAPS(appName, stage, resp.NextToken)
	}

	return nil
}
