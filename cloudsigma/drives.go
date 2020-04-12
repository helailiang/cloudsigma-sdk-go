package cloudsigma

import (
	"context"
	"fmt"
	"net/http"
)

const driveBasePath = "drives"

// DrivesService handles communication with the drives related methods of the CloudSigma API.
//
// CloudSigma API docs: http://cloudsigma-docs.readthedocs.io/en/latest/drives.html
type DrivesService service

// Drive represents a CloudSigma drive.
type Drive struct {
	Media       string `json:"media"`
	Name        string `json:"name"`
	ResourceURI string `json:"resource_uri"`
	Size        int    `json:"size"`
	Status      string `json:"status"`
	StorageType string `json:"storage_type"`
	UUID        string `json:"uuid"`
}

type DriveCloneRequest struct {
	Media       string `json:"media,omitempty"`
	Name        string `json:"name,omitempty"`
	Size        int    `json:"size,omitempty"`
	StorageType string `json:"storage_type,omitempty"`
}

type drivesRoot struct {
	Drives []Drive `json:"objects"`
}

// Get provides detailed information for drive identified by uuid.
//
// CloudSigma API docs: http://cloudsigma-docs.readthedocs.io/en/latest/drives.html#list-single-drive
func (s *DrivesService) Get(ctx context.Context, uuid string) (*Drive, *http.Response, error) {
	if uuid == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v", driveBasePath, uuid)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	drive := new(Drive)
	resp, err := s.client.Do(ctx, req, drive)
	if err != nil {
		return nil, resp, err
	}

	return drive, resp, nil
}
