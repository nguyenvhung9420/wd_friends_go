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
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/service/dynamodb/expression"

    "fmt"
    "os"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

type BodyRequest struct {
	Email string `json:"email_to_find"`
}

// func Handler(ctx context.Context) (Response, error) {
func FetchFollower(request Request) (Response, error) {

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("ap-southeast-1")},
    )
    // var requestBody map[string]interface{}
    var requestBody BodyRequest
    err = json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		return Response{Body: err.Error(), StatusCode: 404}, nil
	}

    var emailToFind string
    // emailToFind = fmt.Sprintf("%v", requestBody["email"])
    emailToFind = requestBody.Email
    fmt.Println(requestBody)

    // Create DynamoDB client
    svc := dynamodb.New(sess)
    tableName := "friends"
    email := emailToFind
    filt := expression.Name("email").Equal(expression.Value(email))
    proj := expression.NamesList(expression.Name("email"), expression.Name("followers"))
    expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()


    if err != nil {
        fmt.Println("Got error building expression:")
        fmt.Println(err.Error())
        os.Exit(1)
    }

    // Build the query input parameters
    params := &dynamodb.ScanInput{
        ExpressionAttributeNames:  expr.Names(),
        ExpressionAttributeValues: expr.Values(),
        FilterExpression:          expr.Filter(),
        ProjectionExpression:      expr.Projection(),
        TableName:                 aws.String(tableName),
    }

    // Make the DynamoDB Query API call
    result, err := svc.Scan(params)
    if err != nil {
        fmt.Println("Query API call failed:")
        fmt.Println((err.Error()))
        os.Exit(1)
    }
    numItems := 0
    
    var allFollowers map[string][]interface{}
    var toReturn []map[string]interface{}
    // var toReturn []interface{}
    // var toReturn []string
    err = dynamodbattribute.UnmarshalMap(result.Items[0], &allFollowers)

    for _, i := range allFollowers["followers"] {
        emailOfFollower := fmt.Sprintf("%v", i)
        resultat, err := svc.GetItem(&dynamodb.GetItemInput{
            TableName: aws.String(tableName),
            Key: map[string]*dynamodb.AttributeValue{
                "email": {
                    S: aws.String(emailOfFollower),
                },
            },
        })
        if err != nil {
            fmt.Println(err.Error())
        }
        var aFollower map[string]interface{}
        err = dynamodbattribute.UnmarshalMap(resultat.Item, &aFollower)
        if err != nil {
            panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
        }
        toReturn = append(toReturn, aFollower)
        numItems++
    }
    fmt.Println(toReturn)
    var buf bytes.Buffer
    body, err := json.Marshal(toReturn)
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
	lambda.Start(FetchFollower)
}