// app.post("/find_user", async (request, response) => {
//     let client = await MongoClient.connect(CONNECTION_URL, { useNewUrlParser: true });
//     database = client.db(DATABASE_NAME);
//     collection = database.collection("users_friends");

//     collection.find({ "email": request.body.email_to_find }).toArray((error, result) => {
//         if (error) {
//             return response.status(500).send(error);
//         }
//         response.send(result[0]);
//     });
// });

// snippet-comment:[These are tags for the AWS doc team's sample catalog. Do not remove.]
// snippet-sourceauthor:[Doug-AWS]
// snippet-sourcedescription:[DynamoDBScanItems.go gets items from and Amazon DymanoDB table using the Expression Builder package.]
// snippet-keyword:[Amazon DynamoDB]
// snippet-keyword:[Scan function]
// snippet-keyword:[Expression Builder]
// snippet-keyword:[Go]
// snippet-sourcesyntax:[go]
// snippet-service:[dynamodb]
// snippet-keyword:[Code Sample]
// snippet-sourcetype:[full-example]
// snippet-sourcedate:[2019-03-19]
/*
   Copyright 2010-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.

   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at

    http://aws.amazon.com/apache2.0/

   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
*/
// snippet-start:[dynamodb.go.scan_items]
package main

// snippet-start:[dynamodb.go.scan_items.imports]
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
	Email string `json:"email"`
}

// func Handler(ctx context.Context) (Response, error) {
func Handler(request Request) (Response, error) {

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
    expr, err := expression.NewBuilder().WithFilter(filt).Build()
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
        // ProjectionExpression:      expr.Projection(),
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

    var toReturn []map[string]interface{}

    for _, i := range result.Items {
        var item map[string]interface{}
        err = dynamodbattribute.UnmarshalMap(i, &item)

        av, err := json.Marshal(item) // This av is just for testing purpose:

        if err != nil {
            fmt.Println("Got error unmarshalling:")
            fmt.Println(err.Error())
            os.Exit(1)
        }
        fmt.Println(string(av)) // Printing av is just for testing purpose:
        toReturn = append(toReturn, item)
        numItems++
    }

    fmt.Println("Found", numItems, "item(s).")

    var buf bytes.Buffer
    // body, err := json.Marshal(map[string]interface{}{
	// 	"message": "This is handler 2!",
    // })
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
		},
	}

	return resp, nil
}


func main() {
	lambda.Start(Handler)
}