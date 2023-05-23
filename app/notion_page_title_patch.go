package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const daysUntilNextThursday = 11

type NotionAPI struct {
	DatabaseURL string
	APIKey      string
}

type NotionResponse struct {
	Results []struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"Results"`
}

func (apiInstance *NotionAPI) ReadPageID() (string, string, error) {
	databaseURL := "https://api.notion.com/v1/databases/" + apiInstance.DatabaseURL + "/query"

	payload := strings.NewReader(`{
    "filter": {
        "property": "会議種別",
        "multi_select": {
            "contains": "定例"
        }
    },
    "page_size": 1
}`)

	request, err := http.NewRequest("POST", databaseURL, payload)
	if err != nil {
		return "", "", err
	}

	request.Header.Add("accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+apiInstance.APIKey)
	request.Header.Add("Notion-Version", "2022-06-28")
	request.Header.Add("content-type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	var notionResponse NotionResponse
	if err := json.NewDecoder(response.Body).Decode(&notionResponse); err != nil {
		return "", "", err
	}

	if len(notionResponse.Results) == 0 {
		return "", "", fmt.Errorf("no page found in Notion database")
	}

	pageURL := notionResponse.Results[0].URL

	return notionResponse.Results[0].ID, pageURL, nil
}

func (apiInstance *NotionAPI) PatchPageTitle(pageID string) error {
	apiEndpoint := "https://api.notion.com/v1/pages/" + pageID

	nextThursday := time.Now().AddDate(0, 0, (daysUntilNextThursday-int(time.Now().Weekday()))%7)
	nextThursdayStartStr := nextThursday.Format("2006-01-02")
	nextThursdayTitleStr := nextThursday.Format("01/02")

	payload := strings.NewReader(fmt.Sprintf(`{
    "properties": {
        "名前": {
            "title": [
                {
                    "text": {
                        "content": "定例%s"
                    }
                }
            ]
        },
		"会議開催日": {
            "date": {
   				 "start": "%s",
   				 "end": null
  			}
        }
    }
}`, nextThursdayTitleStr, nextThursdayStartStr))

	request, err := http.NewRequest("PATCH", apiEndpoint, payload)
	if err != nil {
		return err
	}

	request.Header.Add("accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+apiInstance.APIKey)
	request.Header.Add("Notion-Version", "2022-06-28")
	request.Header.Add("content-type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}
