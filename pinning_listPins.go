package pinningcli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

type listPinResponse struct {
	Data   listPinResult `json:"data"`
	Status string        `json:"status"`
}

type listPinResult struct {
	Totals Totals `json:"totals"`
	Pins   []*Pin `json:"pins"`
}

type Totals struct {
	Files int64  `json:"files"`
	Size  uint64 `json:"size"`
}
type ListPinsCriteria struct {
	Limit          uint
	Offset         uint
	SortBy         string
	SortOrder      string
	Pinned         string
	KeyValuesQuery datatypes.JSON
}

func (pinning Pinning) getPinsList(ids PinningIdentifiers, criteria ListPinsCriteria) (listPinResult, []error) {
	if criteria.Offset != 0 {
		offset = criteria.Offset
	}
	if criteria.Limit != 0 {
		limit = criteria.Limit
	}
	query := fmt.Sprintf("%s/pinning/pins/?offset=%d&limit=%d", baseUrl, offset, limit)
	if len(criteria.SortBy) > 0 {
		query += "&sortBy=" + criteria.SortBy
	}
	if len(criteria.SortOrder) > 0 {
		query += "&sortOrder=" + criteria.SortOrder
	}
	if len(criteria.Pinned) > 0 {
		query += "&pinned=" + criteria.Pinned
	}
	if len(criteria.KeyValuesQuery) > 0 {
		query += "&metadata=" + criteria.KeyValuesQuery.String()
	}
	log.Debug("GET ", query)

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return listPinResult{}, []error{fmt.Errorf("failed to call Pinning server with GET at %s", query), err}
	}

	req.Header.Add("pinning_api_key", ids.ApiKey)
	req.Header.Add("pinning_secret_key", ids.SecretKey)

	resp, err := httpClient.Do(req)

	if err != nil {
		return listPinResult{}, []error{fmt.Errorf("failed to call Pinning server with GET at %s", query), err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return listPinResult{}, []error{fmt.Errorf("failed to read body while calling Pinning server with GET at %s", query), err}
	}

	log.Trace("Status code: ", resp.Status)
	log.Trace("Response body: ", string(body))

	if resp.StatusCode != 200 {
		return listPinResult{}, []error{fmt.Errorf("failed to call Pinning server with GET at %s", query),
			fmt.Errorf("status code: %s", resp.Status),
		}
	}

	var response listPinResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return listPinResult{}, []error{fmt.Errorf("failed to parse JSON response body while calling Pinning server with GET at %s", query), err}
	}

	return response.Data, nil
}
