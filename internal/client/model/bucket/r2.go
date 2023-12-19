package bucket

import (
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type STSResponse struct {
	Success bool        `json:"success"`
	Result  string      `json:"result"`
	Errors  interface{} `json:"errors"`
}

func GetR2STS(ctx context.Context) {
	url := "https://api.cloudflare.com/client/v4/user/tokens/verify"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 设置请求头
	//d71097891fd7fd38ccfe49720b37d20d
	//2cede312fab5e244110f09fb03f9a021421665bf76122aed8463bc8e563454eb
	req.Header.Add("Authorization", "Bearer ReWS_FyQQQqbKaWSv-qUjwewxA21S58HKCora9GP")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var stsResp STSResponse
	var res1 map[string]interface{}
	err = json.Unmarshal(body, &res1)
	fmt.Println("Result:", res1)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	if stsResp.Success {
		fmt.Println("STS verification successful")
		fmt.Println("Result:", stsResp.Result)
	} else {
		fmt.Println("STS verification failed", stsResp)
	}
}
