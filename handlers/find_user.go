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
    "encoding/json"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/service/dynamodb/expression"

    "fmt"
    "os"
)

func main() {

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("ap-southeast-1")},
    )

    // Create DynamoDB client
    svc := dynamodb.New(sess)
    // snippet-end:[dynamodb.go.scan_items.session]

    // snippet-start:[dynamodb.go.scan_items.vars]
    tableName := "groups"
    admin := "motethansen"
    // snippet-end:[dynamodb.go.scan_items.vars]

    // snippet-start:[dynamodb.go.scan_items.expr]
    // Create the Expression to fill the input struct with.
    // Get all movies in that year; we'll pull out those with a higher rating later
    filt := expression.Name("admin").Equal(expression.Value(admin))

    // Or we could get by ratings and pull out those with the right year later
    //    filt := expression.Name("info.rating").GreaterThan(expression.Value(min_rating))

    // Get back the title, year, and rating
    // proj := expression.NamesList(expression.Name("Title"), expression.Name("Year"), expression.Name("Rating"))
    // proj := expression.NamesList(expression.Name("groupID"), expression.Name("groupname"))


    // expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
    expr, err := expression.NewBuilder().WithFilter(filt).Build()
    if err != nil {
        fmt.Println("Got error building expression:")
        fmt.Println(err.Error())
        os.Exit(1)
    }
    // snippet-end:[dynamodb.go.scan_items.expr]

    // snippet-start:[dynamodb.go.scan_items.call]
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

    for _, i := range result.Items {
        var item map[string]interface{}
        err = dynamodbattribute.UnmarshalMap(i, &item)
        av, err := json.Marshal(item)

        if err != nil {
            fmt.Println("Got error unmarshalling:")
            fmt.Println(err.Error())
            os.Exit(1)
        }
        fmt.Println(string(av))
        numItems++
    }

    fmt.Println("Found", numItems, "item(s).")
    // snippet-end:[dynamodb.go.scan_items.process]
}
// snippet-end:[dynamodb.go.scan_items]