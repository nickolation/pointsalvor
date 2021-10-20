# Todoist API golang SDK for friendly point-managment 

An unofficial golang SDK for todoist-client managment 

Todoist API reference: https://developer.todoist.com/rest/v1/#overview 

### Installation
```shell
go get -u github.com/nickolation/pointsalvor
``` 

### Authorization
- To use the methods of SDK you need to create a structure Agent 
```go
    type Agent struct {
        Engine *http.Client
        Token string
    }
```
- The object __Agent__ contains the field __Token__ - token-api used for auth on the todoist-client  
- In order to make authorized calls to the REST API, your application must provide an authorization header with the appropriate _Bearer $token_. For working through the examples, you can obtain your personal API token from the __integrations settings__ for your account. - <a href="https://developer.todoist.com/rest/v1/#authorization">[More details]</a>



### Example of using:
```go
package main

import (
	"github.com/nickolation/pointsalvor"
)

func main() {
	agent, err := pointsalvor.NewAgent("<api-token>")
	if err != nil {
		fmt.Println(err.Error())
	}
}
```