package images

import (
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// GetResult temporarily stores a Get response.
type GetResult struct {
	gophercloud.Result
}

// DeleteResult represents the result of an image.Delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Extract interprets a GetResult as an Image.
func (r GetResult) Extract() (*Image, error) {
	var s struct {
		Image *Image `json:"image"`
	}
	err := r.ExtractInto(&s)
	return s.Image, err
}

// Image is used for JSON (un)marshalling.
// It provides a description of an OS image.
type Image struct {
	// ID contains the image's unique identifier.
	ID string `json:"id"`

	// Name provides a human-readable moniker for the OS image.
	Name string `json:"name"`

	// Updated and Created contain ISO-8601 timestamps of when the state of the server last changed, and when it was created.
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	// MinDisk and MinRAM specify the minimum resources a server must provide to be able to install the image.
	MinDisk int `json:"min_disk"`
	MinRAM  int `json:"min_ram"`

	// The Progress and Status fields indicate image-creation status.
	// Any usable image will have 100% progress.
	Progress int    `json:"progress"`
	Status   string `json:"status"`

	Metadata map[string]interface{} `json:"metadata"`
}

// ImagePage contains a single page of results from a List operation.
// Use ExtractImages to convert it into a slice of usable structs.
type ImagePage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a page contains no Image results.
func (page ImagePage) IsEmpty() (bool, error) {
	images, err := ExtractImages(page)
	return len(images) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the next page of results.
func (page ImagePage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"images_links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractImages converts a page of List results into a slice of usable Image structs.
func ExtractImages(r pagination.Page) ([]Image, error) {
	var s struct {
		Images []Image `json:"images"`
	}
	err := (r.(ImagePage)).ExtractInto(&s)
	return s.Images, err
}
