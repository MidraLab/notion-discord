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

type notionResponse struct {
	Results []struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"Results"`
}

func (n *NotionAPI) ReadPageID() (string, string, error) {
	dbUrl := "https://api.notion.com/v1/databases/" + n.DatabaseURL + "/query"

	payload := strings.NewReader(`{
    "filter": {
        "property": "会議種別",
        "multi_select": {
            "contains": "定例"
        }
    },
    "page_size": 1
}`)

	req, err := http.NewRequest("POST", dbUrl, payload)
	if err != nil {
		return "", "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+n.APIKey)
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	var notionRes notionResponse
	if err := json.NewDecoder(res.Body).Decode(&notionRes); err != nil {
		return "", "", err
	}

	if len(notionRes.Results) == 0 {
		return "", "", fmt.Errorf("no page found in Notion database")
	}

	url := notionRes.Results[0].URL

	return notionRes.Results[0].ID, url, nil
}

func (n *NotionAPI) PatchPageTitle(id string) error {
	url := "https://api.notion.com/v1/pages/" + id

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

	req, err := http.NewRequest("PATCH", url, payload)
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+n.APIKey)
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}
