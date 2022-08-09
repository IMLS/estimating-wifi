package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

//***Define Success and Error structures***

type AuthSuccess struct {
	/* variables */
}
type AuthError struct {
	/* variables */
}

//***Configuration for POST/GET***
var timeOut int = 15

//***Begin POST function***
func newPost(uri string, data []map[string]interface{}, key string) {

	//***for testing***
	fmt.Println("Running Post")
	myJson, _ := json.Marshal(data)
	fmt.Println(string(myJson))

	//***create client and conditions needed for the client***
	client := resty.New()
	client.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() == http.StatusTooManyRequests
		},
	)
	client.SetTimeout(time.Duration(timeOut) * time.Second)

	//***Client post data***
	resp, err := client.R().
		SetBody(data).
		SetAuthToken(key).
		//SetResult(&AuthSuccess{}). Could be incorperated once we have defined response
		//SetError(&AuthError{}). Could be incorperated once we have defined response
		Post(uri)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("  Body       :\n", resp)
	}

	//***for testing response code***
	//fmt.Println("Response Info:")
	//fmt.Println("Status Code:", resp.StatusCode())
	//fmt.Println("Status:", resp.Status())

}
