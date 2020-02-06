package main

import (
    "bytes"
	// "context"
    "encoding/json"

    "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"

    "fmt"
    "os"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

type Req struct {
    Email string `json:"email"`
    Fullname string `json:"fullname"`
    Nickname string `json:"nickname"`
    Bio string `json:"bio"`
}

func AddUserToDynamo(request Request) (Response, error) {

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("ap-southeast-1")},
    )
    
    var req Req
    err = json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return Response{Body: err.Error(), StatusCode: 404}, nil
    }
    fmt.Println(req)

    svc := dynamodb.New(sess)
    tableName := "friends"

    input := &dynamodb.PutItemInput{
        Item: map[string]*dynamodb.AttributeValue{
            "email": {   S: aws.String(req.Email),   },
            "nickname": {   S: aws.String(req.Nickname),   },
            "fullname": {   S: aws.String(req.Fullname),   },
            "bio": {   S: aws.String(req.Bio),   },
            "groupParticipations": {
                L: []*dynamodb.AttributeValue{},
            },
            "followers": {
                L: []*dynamodb.AttributeValue{},
            },
            "followings": {
                L: []*dynamodb.AttributeValue{},
            },
        },
        TableName: aws.String(tableName),
    }

    _, err = svc.PutItem(input)
    if err != nil {
        fmt.Println("Got error calling PutItem:")
        fmt.Println(err.Error())
        os.Exit(1)
    }


    var buf bytes.Buffer
    body, err := json.Marshal(map[string]interface{}{
		"message": "Done",
    })
    // body, err := json.Marshal(toReturn)
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

    resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
            "X-MyCompany-Func-Reply": "hello-handler",
            "Access-Control-Allow-Origin": "*",
		},
	}

	return resp, nil
}


func main() {
	lambda.Start(AddUserToDynamo)
}