package pinningcli

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

type pinFileResult struct {
	Status string `json:"status"`
	Data   Pin    `json:"data"`
}
type Pin struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id,omitempty"`
	CID          string    `gorm:"not null" json:"cid,omitempty"`
	Size         uint64    `gorm:"not null" json:"size,omitempty"`
	UserID       uuid.UUID `gorm:"not null" json:"user_id,omitempty"`
	DatePinned   time.Time `gorm:"not null" json:"date_pinned,omitempty"`
	DateUnpinned time.Time `gorm:"null" json:"date_unpinned"`
	Pinned       bool      `gorm:"not null" json:"pinned"`
	IsDir        bool      `gorm:"null" json:"is_dir"`
	Metadata     Metadata  `gorm:"embedded" json:"metadata,omitempty"`
}

type PinFileCriteria struct {
	File              string
	CIDVersion        string
	WrapWithDirectory bool
	Name              string
	KeyValues         datatypes.JSON
}

func (result pinFileResult) String() string {
	return fmt.Sprintf("IpfsHash: %s | PinSize: %d | Timestamp: %s",
		result.Data.CID, result.Data.Size, result.Data.DatePinned)
}

func (result Pin) String() string {
	return fmt.Sprintf("IpfsHash: %s | PinSize: %d | Timestamp: %s",
		result.CID, result.Size, result.DatePinned)
}

func (pinning Pinning) writeAllContent(input string, metadata *Metadata) (*io.PipeReader, *string, []error) {

	metadataBytes, err := json.Marshal(metadata)
	metadataStr := string(metadataBytes)

	stats, err := os.Stat(input)
	if os.IsNotExist(err) {
		return nil, nil, []error{fmt.Errorf("file or directory is missing: %s", input)}
	}

	var files []string
	fileIsASingleFile := !stats.IsDir()
	if fileIsASingleFile {
		files = append(files, stats.Name())
	} else {
		err = filepath.Walk(input,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					relPath, _ := filepath.Rel(input, path)
					files = append(files, relPath)
				}
				return nil
			})

		if err != nil {
			return nil, nil, []error{fmt.Errorf("fatal error while exploring directory '%s'", input), err}
		}
	}

	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()

		for _, f := range files {
			var part io.Writer
			var err error
			if fileIsASingleFile {
				part, err = m.CreateFormFile("file", f)
			} else {
				part, err = m.CreateFormFile("file", filepath.Join(stats.Name(), f))
			}
			if err != nil {
				return
			}
			var file *os.File
			if fileIsASingleFile {
				file, err = os.Open(input)
			} else {
				file, err = os.Open(filepath.Join(input, f))
			}
			if err != nil {
				return
			}
			defer file.Close()
			if _, err = io.Copy(part, file); err != nil {
				return
			}
		}

		if metadata != nil {
			err = m.WriteField("metadata", metadataStr)
		}

	}()

	boundary := m.Boundary()
	return r, &boundary, nil

}

func (pinning Pinning) pinFile(ids PinningIdentifiers, criteria PinFileCriteria) (Pin, []error) {

	var result Pin

	unionCriteria := pinCriteriaUnion{
		Metadata: Metadata{
			Name:      criteria.Name,
			KeyValues: criteria.KeyValues,
		},
	}

	pinningMetadata := pinning.buildPinningMetadata(unionCriteria)

	reader, boundary, errs := pinning.writeAllContent(criteria.File, pinningMetadata)
	if errs != nil {
		return result, append([]error{fmt.Errorf("failed to prepare content")}, errs...)
	}

	query := fmt.Sprintf("%s/pinning/", baseUrl)
	log.Debug("POST ", query)

	req, err := http.NewRequest("POST", query, reader)
	if err != nil {
		return result, []error{fmt.Errorf("failed to prepare request for Pinning server with POST at %s", query), err}
	}

	req.Header.Add("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", *boundary))
	req.Header.Add("pinning_api_key", ids.ApiKey)
	req.Header.Add("pinning_secret_key", ids.SecretKey)

	resp, err := httpClient.Do(req)

	if err != nil {
		return result, []error{fmt.Errorf("failed to call Pinning server with POST at %s", query), err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, []error{fmt.Errorf("failed to read body while calling Pinning server with POST at %s", query), err}
	}

	log.Trace("Status code: ", resp.Status)
	log.Trace("Response body: ", string(body))

	if resp.StatusCode != 201 {
		return result, []error{fmt.Errorf("error while calling Pinning server with POST at %s", query),
			fmt.Errorf("status code: %s", resp.Status),
			fmt.Errorf("response body: %s", string(body)),
		}
	}

	var response pinFileResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return result, []error{fmt.Errorf("failed to parse JSON response body while calling Pinning server with POST at %s", query), err}
	}

	log.Trace("Response unmarshalled: ", response.String())

	result = Pin{
		ID:         response.Data.ID,
		CID:        response.Data.CID,
		Size:       response.Data.Size,
		UserID:     response.Data.UserID,
		DatePinned: response.Data.DatePinned,
		Metadata:   response.Data.Metadata,
	}

	log.Trace("Return object: ", result.String())

	return result, nil
}
