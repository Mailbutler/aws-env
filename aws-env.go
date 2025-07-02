package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	formatExports = "exports"
	formatDotenv  = "dotenv"
)

func main() {
	// Check if AWS_ENV_PATH is set
	if os.Getenv("AWS_ENV_PATH") == "" {
		log.Println("aws-env running locally, without AWS_ENV_PATH")
		return
	}

	recursivePtr := flag.Bool("recursive", false, "recursively process parameters on path")
	format := flag.String("format", formatExports, "output format")
	flag.Parse()

	if *format != formatExports && *format != formatDotenv {
		log.Fatal("Unsupported format option. Must be 'exports' or 'dotenv'")
	}

	sess := CreateSession()
	client := CreateClient(sess)

	// if AWS_ENV_PATH contains a comma-separated list, export variables for each path
	paths := strings.Split(os.Getenv("AWS_ENV_PATH"), ",")
	for _, path := range paths {
		ExportVariables(client, strings.TrimSpace(path), *recursivePtr, *format, "")
	}
}

func CreateSession() *session.Session {
	return session.Must(session.NewSession())
}

func CreateClient(sess *session.Session) *ssm.SSM {
	return ssm.New(sess)
}

func ExportVariables(client *ssm.SSM, path string, recursive bool, format string, nextToken string) {
	input := &ssm.GetParametersByPathInput{
		Path:           &path,
		WithDecryption: aws.Bool(true),
		Recursive:      aws.Bool(recursive),
	}

	if nextToken != "" {
		input.SetNextToken(nextToken)
	}

	output, err := client.GetParametersByPath(input)

	if err != nil {
		log.Panic(err)
	}

	for _, element := range output.Parameters {
		OutputParameter(path, element, format)
	}

	if output.NextToken != nil {
		ExportVariables(client, path, recursive, format, *output.NextToken)
	}
}

func OutputParameter(path string, parameter *ssm.Parameter, format string) {
	name := *parameter.Name
	value := *parameter.Value

	env := strings.Replace(strings.Trim(name[len(path):], "/"), "/", "_", -1)
	value = strings.Replace(value, "\n", "\\n", -1)
	value = strings.Replace(value, "'", "\\'", -1)

	switch format {
	case formatExports:
		fmt.Printf("export %s=$'%s'\n", env, value)
	case formatDotenv:
		fmt.Printf("%s=\"%s\"\n", env, value)
	}
}
