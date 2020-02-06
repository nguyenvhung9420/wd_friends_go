package main

// snippet-start:[dynamodb.go.scan_items.imports]
import (
    "bytes"
	"context"
    "encoding/json"

    "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/service/dynamodb/expression"

    "fmt"
    // "os"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

type BodyRequest struct {
	Email string `json:"email"`
}

func FindUser(ctx context.Context, request Request) (Response, error) {
// func FindUser(request Request) (Response, error) {

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("ap-southeast-1")},
    )

    var requestBody BodyRequest
    err = json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		return Response{Body: err.Error(), StatusCode: 404}, nil
	}

    var emailToFind string
    emailToFind = requestBody.Email
    fmt.Println(requestBody)

    // Create DynamoDB client
    svc := dynamodb.New(sess)
    tableName := "friends"
    email := emailToFind
    // email := "nguyenvhung@live.fr"
    filt := expression.Name("email").Equal(expression.Value(email))
    expr, err := expression.NewBuilder().WithFilter(filt).Build()
    if err != nil {
        fmt.Println("Got error building expression:")
        fmt.Println(err.Error())
        // os.Exit(1)
    }

    // Build the query input parameters
    params := &dynamodb.ScanInput{
        ExpressionAttributeNames:  expr.Names(),
        ExpressionAttributeValues: expr.Values(),
        FilterExpression:          expr.Filter(),
        // ProjectionExpression:      expr.Projection(),
        TableName:                 aws.String(tableName),
    }

    // Make the DynamoDB Query API call
    result, err := svc.Scan(params)
    if err != nil {
        fmt.Println("Query API call failed:")
        fmt.Println((err.Error()))
        // os.Exit(1)
    }
    numItems := 0

    var toReturn []map[string]interface{}

    for _, i := range result.Items {
        var item map[string]interface{}
        err = dynamodbattribute.UnmarshalMap(i, &item)

        av, err := json.Marshal(item) // This av is just for testing purpose:

        if err != nil {
            fmt.Println("Got error unmarshalling:")
            fmt.Println(err.Error())
            // os.Exit(1)
        }
        fmt.Println(string(av)) // Printing av is just for testing purpose:
        toReturn = append(toReturn, item)
        numItems++
    }

    fmt.Println("Found", numItems, "item(s).")

    var buf bytes.Buffer

    // body, err := json.Marshal(map[string]map[string]interface{}   {
	// 	"message": toReturn,
    // } )
    body, err := json.Marshal(&toReturn)
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)
    
    resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
        // Body:            string(body),
		Headers: map[string]string{
			"Content-Type":           "application/json",
            "X-MyCompany-Func-Reply": "hello-handler",
            "Access-Control-Allow-Origin": "*",
		},
	}

    fmt.Println(resp)

	return resp, nil
}


func main() {
	lambda.Start(FindUser)
}