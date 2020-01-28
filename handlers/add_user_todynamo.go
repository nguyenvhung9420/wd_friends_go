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
    // "github.com/aws/aws-sdk-go/service/dynamodb/expression"

    "fmt"
    "os"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

type Item struct {
    Email string `json:"email"`
    Nickname string `json:"nickname"`
    Followings []string  `json:"followings"`
    Followers []string  `json:"followers"`
    GroupParticipations []string  `json:"groupParticipations"`
    Bio string `json:"bio"`
}

type Req struct {
    Email string `json:"email"`
    Nickname string `json:"nickname"`
    Bio string `json:"bio"`
}


// var newvalues = {
//     "email": request.body.email,
//     "nickname": request.body.nickname,
//     "followings": [],
//     "followers": [],
//     "groupParticipations": []
// };

// func Handler(ctx context.Context) (Response, error) {
func Handler(request Request) (Response, error) {

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("ap-southeast-1")},
    )
    // var requestBody map[string]interface{}
    var req Req
    // requestBody := Req{
	// 	email : "",
    //     nickname : "",
    //     bio : "",
	// }
    err = json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return Response{Body: err.Error(), StatusCode: 404}, nil
    }

    // var emailToFind string
    // emailToFind = fmt.Sprintf("%v", requestBody["email"])
    // // emailToFind = requestBody.email
    // fmt.Println(request.Body)
    fmt.Println(req)
    
    // var emptyStringSlive []string

    item := Item{
        Email: "nguyenvhung@mail3.ru",
        Nickname: "ru_hung",
        Followings: []string{"me", "you"} ,
        Followers: []string{"me", "you"}  ,
        GroupParticipations: []string{"me", "you"} ,
        Bio: "C'est le bio",
    }
    
    // Create DynamoDB client
    svc := dynamodb.New(sess)
    tableName := "friends"

    av, err := dynamodbattribute.MarshalMap(item)
    if err != nil {
        fmt.Println("Got error marshalling new item:")
        fmt.Println(err.Error())
        os.Exit(1)
    }

    fmt.Println(av)

    input := &dynamodb.PutItemInput{
        Item:      av,
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
		},
	}

	return resp, nil
}


func main() {
	lambda.Start(Handler)
}