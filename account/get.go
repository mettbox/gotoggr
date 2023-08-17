package account

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	TogglAPI   = "https://api.track.toggl.com/api/v9"
	ReportsAPI = "https://api.track.toggl.com/reports/api/v2"
)

func Get() (account Account, err error) {
	url := TogglAPI + "/me?with_related_data=true"

	apiToken := GetToken()

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	req.SetBasicAuth(apiToken, "api_token")

	res, err := client.Do(req)
	if err != nil {
		return
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if status := res.StatusCode; status != 200 {
		fmt.Printf("Error: Could not connect to Toggl API. %s (%d)", http.StatusText(status), status)
		os.Exit(0)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	account = Account{}
	jsonErr := json.Unmarshal(body, &account)
	if jsonErr != nil {
		err = jsonErr
		return
	}

	return
}
