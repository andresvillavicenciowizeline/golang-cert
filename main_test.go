package main

/*
response: ["alpha"]

response: [
	{domain: string}
	{domain: string}
	{domain: string}
	{domain: string}
]
*/
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"encoding/json"

	handlers "github.com/andresvillavicenciowizeline/proxy-app/api/handlers"
	server "github.com/andresvillavicenciowizeline/proxy-app/api/server"
	utils "github.com/andresvillavicenciowizeline/proxy-app/api/utils"
)

func init() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		utils.LoadEnv()
		app := server.SetUp()
		handlers.HandlerRedirection(app)
		wg.Done()
		server.RunServer(app)
	}(wg)
	wg.Wait()
	fmt.Println("Server running...")

}

type Response struct {
	Status       int            `json:"status,omitempty"`
	Response     string         `json:"result,omitempty"`
	ResponseText []ResponseText `json:"res,omitempty"`
}

type ResponseText struct {
	Domain string
}

func TestAlgorithmn(t *testing.T) {
	cases := []struct {
		Domain string
		Output string
	}{
		{Domain: "alpha", Output: `["alpha"]`},
		// {Domain: "alpha", Output: `["alpha,  alpha"]`},
		// {Domain: "omega", Output: `["alpha,  alpha, omega"]`},
		// {Domain: "beta", Output: `["alpha,  alpha, beta, omega"]`},
		{Domain: "", Output: "domain error"},
	}

	valuesToCompare := &Response{}
	client := http.Client{}
	for _, singleCase := range cases {
		req, err := http.NewRequest("GET", "http://localhost:8080/ping", nil)
		assert.Nil(t, err)
		req.Header.Add("domain", singleCase.Domain)

		response, err := client.Do(req)

		bytes, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(bytes, valuesToCompare)

		assert.Nil(t, err)
		assert.Equal(t, singleCase.Output, valuesToCompare.Response)
	}
}
