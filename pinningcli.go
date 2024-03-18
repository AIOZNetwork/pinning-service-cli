package pinningcli

import (
	"fmt"
	"time"
)

type PinningCli struct {
	Pinning Pinning
}

type PinningIdentifiers struct {
	ApiKey    string
	SecretKey string
}

type PinningPin struct {
	ID           string
	IpfsPinHash  string
	Size         uint64
	UserID       string
	DatePinned   time.Time
	DateUnpinned time.Time
	Metadata     PinningPinMetadata
}

type PinningPinMetadata struct {
	Name      string
	KeyValues map[string]string
}

func (pin PinningPin) String() string {
	return fmt.Sprintf("ID: %s | IpfsPinHash: %s", pin.ID, pin.IpfsPinHash)
}

func (PinningCli *PinningCli) Noop() error {
	return nil
}
func (PinningCli *PinningCli) GetPinByID(ids PinningIdentifiers, criteria GetPinByIDCriteria) (Pin, []error) {
	pin, errs := PinningCli.Pinning.getPinByID(ids, criteria)
	if errs != nil {
		return Pin{}, errs
	}
	return pin, nil
}
func (PinningCli *PinningCli) GetPinsList(ids PinningIdentifiers, criteria ListPinsCriteria) (listPinResult, []error) {
	pins, errs := PinningCli.Pinning.getPinsList(ids, criteria)
	if errs != nil {
		return listPinResult{}, errs
	}
	return pins, nil
}

func (PinningCli *PinningCli) Unpin(ids PinningIdentifiers, criteria UnpinCriteria) (UnpinResult, []error) {
	result, errs := PinningCli.Pinning.unpin(ids, criteria)
	return result, errs
}

func (PinningCli *PinningCli) PinByHash(ids PinningIdentifiers, criteria PinByHashCriteria) (Pin, []error) {
	pin, errs := PinningCli.Pinning.pinByHash(ids, criteria)
	return pin, errs
}

func (PinningCli *PinningCli) PinFile(ids PinningIdentifiers, criteria PinFileCriteria) (Pin, []error) {
	pin, errs := PinningCli.Pinning.pinFile(ids, criteria)
	return pin, errs
}

func (PinningCli *PinningCli) TestAuthentication(ids PinningIdentifiers) (TestAuthResult, []error) {
	result, errs := PinningCli.Pinning.testAuthentication(ids)
	return result, errs
}
