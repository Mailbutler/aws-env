package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func main() {
	if os.Getenv("AWS_ENV_PATH") == "" {
		log.Println("aws-env running locally, without AWS_ENV_PATH")
	} else {
		// recursivePtr := flag.Bool("recursive", false, "recursively process parameters on path")
		// flag.Parse()

		sess := CreateSession()
		client := CreateClient(sess)
		recursive := false
		if os.Getenv("AWS_ENV_RECURSIVE") != "" {
			r, err := strconv.ParseBool(os.Getenv("AWS_ENV_RECURSIVE"))
			if err != nil {
				panic(err.Error() + "\nSee: https://golang.org/pkg/strconv/#ParseBool")
			}
			recursive = r
		}
		ExportVariables(client, os.Getenv("AWS_ENV_PATH"), recursive, "")
	}
	fmt.Println("os.Args")
	fmt.Println(os.Args)
	binary, lookErr := exec.LookPath(os.Args[0])
	if lookErr != nil {
		panic(lookErr)
	}

	env := os.Environ()
	args := os.Args[1:]
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

func CreateSession() *session.Session {
	sess := session.Must(session.NewSession())
	if len(aws.StringValue(sess.Config.Region)) == 0 {
		meta := ec2metadata.New(sess)
		identity, err := meta.GetInstanceIdentityDocument()
		if err != nil {
			return session.Must(nil, err)
		}
		return session.Must(session.NewSession(&aws.Config{
			Region: aws.String(identity.Region),
		}))
	}
	return sess
}

func CreateClient(sess *session.Session) *ssm.SSM {
	return ssm.New(sess)
}

func ExportVariables(client *ssm.SSM, path string, recursive bool, nextToken string) {
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
		PrintExportParameter(path, element)
		SetExportParameter(element)
	}

	if output.NextToken != nil {
		ExportVariables(client, path, recursive, *output.NextToken)
	}
}

func PrintExportParameter(path string, parameter *ssm.Parameter) {
	name := *parameter.Name
	value := *parameter.Value

	env := strings.Replace(strings.Trim(name[len(path):], "/"), "/", "_", -1)
	value = strings.Replace(value, "\n", "\\n", -1)
	value = strings.Replace(value, "'", "\\'", -1)

	fmt.Printf("export %s=$'%s'\n", env, value)
}

func SetExportParameter(parameter *ssm.Parameter) {
	name := *parameter.Name
	value := *parameter.Value

	_, envName := path.Split(name)
	os.Setenv(envName, value)
}
