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
    // "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    // "github.com/aws/aws-sdk-go/service/dynamodb/expression"

    "fmt"
    "os"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

// type Item struct {
//     Email string `json:"email"`
//     Nickname string `json:"nickname"`
//     Followings []string  `json:"followings"`
//     Followers []string  `json:"followers"`
//     GroupParticipations []string  `json:"groupParticipations"`
//     Bio string `json:"bio"`
// }

// type Req struct {
//     Email string `json:"email"`
//     Nickname string `json:"nickname"`
//     Bio string `json:"bio"`
// }

func Handler(request Request) (Response, error) {

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("ap-southeast-1")},
    )
    
    svc := dynamodb.New(sess)
    tableName := "friends"

    if err != nil {
		fmt.Println(err)
		fmt.Println(err.Error())
	}

	av := &dynamodb.AttributeValue{
		S: aws.String("Hung"),
	}
	var qids []*dynamodb.AttributeValue
	qids = append(qids, av)

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String("nguyenvhung@mail3.ru"),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":qid": {
				L: qids,
			},
			":empty_list": {
				L: []*dynamodb.AttributeValue{},
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET groupParticipations = list_append(if_not_exists(groupParticipations, :empty_list), :qid)"),
		TableName:        aws.String(tableName),
    }

    _, err = svc.UpdateItem(input)
    if err != nil {
        fmt.Println("Got error calling Update:")
        fmt.Println(err.Error())
        os.Exit(1)
    }

    // THIS IS THE PART TO TRANSFER SUCCESS MESSSAGE
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
		},
	}

	return resp, nil
}


func main() {
	lambda.Start(Handler)
}