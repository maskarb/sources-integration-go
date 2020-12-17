package sources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var URL, _ = url.Parse("http://localhost:3000/api/v1.0/")

// url = "sources" source_id

type SourceDetails struct {
	SourceID     uint64 `json:"source_id"`
	SourceTypeID uint64 `json:"source_type_id"`
	Name         string `json:"name"`
	UID          string `json:"uid"`
	EndpointID   uint64 `json:"endpoint_id"`
	SourceType   string
}

func GetResponse(url string, header string) *http.Response {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-rh-identity", "value")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	return resp
}

func GetSourceData(resp *http.Response) map[string]interface{} {
	var data map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

// func GetSourceDetails(sourceID uint64) SourceDetails {

// }

func main() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("x-rh-identity", "value")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var data SourceDetails
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
	}

}
