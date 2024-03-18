package pinningcli

import (
	"net/http"

	"gorm.io/datatypes"
)

var (
	offset = uint(0)
	limit = uint(10)
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{}
}

type Pinning struct{}

type pinCriteriaUnion struct {
	HashToPin string   `json:"hash_to_pin,omitempty"`
	Metadata  Metadata `json:"metadata,omitempty"`
}

type Metadata struct {
	Name      string         `json:"name,omitempty"`
	Type      string         `json:"type,omitempty"`
	KeyValues datatypes.JSON `gorm:"type:jsonb" json:"keyvalues"`
}

func (pinning Pinning) buildPinningMetadata(criteria pinCriteriaUnion) *Metadata {

	var metadata Metadata

	if len(criteria.Metadata.Name) == 0 && len(criteria.Metadata.KeyValues) == 0 {
		return nil
	}
	if len(criteria.Metadata.Name) != 0 {
		metadata.Name = criteria.Metadata.Name
	}
	if len(criteria.Metadata.KeyValues) != 0 {
		metadata.KeyValues = criteria.Metadata.KeyValues
	}

	return &metadata

}
