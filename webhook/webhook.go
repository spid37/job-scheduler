package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Data -
type Data struct {
	Method                 string `json:"method"`
	URL                    string `json:"url"`
	Timeout                int    `json:"timeout"`
	ExpectedResponseStatus int    `json:"expectedResponseStatus"`
	ExpectedContentType    string `json:"expectedResponseContentType"`
	ExpectedResponse       string `json:"expectedResponse"`
}

// LoadData load the job data fro webhook
func (d *Data) LoadData(data []byte) error {
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	return nil
}

// Send Send a Job request
func (d *Data) Send() error {
	var err error
	time.Sleep(10 * time.Second)

	client := http.Client{}
	if d.Timeout > 0 {
		client.Timeout = time.Duration(time.Duration(d.Timeout) * time.Millisecond)
	}

	resp, err := client.Get(d.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// check status matches expected if required
	if d.ExpectedResponseStatus > 0 && resp.StatusCode != d.ExpectedResponseStatus {
		err = fmt.Errorf("Status code %d does not match expected %d",
			resp.StatusCode,
			d.ExpectedResponseStatus,
		)
		return err
	}

	// check context type matched expected if required
	contentType := strings.ToLower(resp.Header.Get("Content-type"))
	expectedType := strings.ToLower(d.ExpectedContentType)
	if expectedType != "" && !strings.Contains(contentType, expectedType) {
		err = fmt.Errorf("Context-type %s does not match expected %s",
			contentType,
			d.ExpectedContentType,
		)
		return err
	}

	// check the body response
	bodyStr := string(body)
	if d.ExpectedResponse != "" && d.ExpectedResponse != bodyStr {
		err = fmt.Errorf("Body %s does not match expected %s",
			bodyStr,
			d.ExpectedResponse,
		)
		return err
	}

	return err
}
