// Copyright (c) 2014 Hewlett-Packard Development Company, L.P.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

/*
Package image implements a client library for accessing OpenStack Image V1 service

Images and ImageDetails can be retrieved using the api.

In addition more complex filtering and sort queries can by using the ImageQueryParameters.

*/
package image

import (
	"fmt"
	"git.openstack.org/stackforge/golang-client.git/misc"
	"net/http"
	"net/url"
)

// ImageService is a client service that can make
// requests against a OpenStack version 1 image service.
// Below is an example on creating an image service and getting images:
// 	imageService := image.ImageService{Client: *http.DefaultClient, TokenId: tokenId, Url: "http://imageservicelocation"}
//  images:= imageService.Images()
type ImageService struct {
	Client  http.Client
	TokenId string
	Url     string
}

// ImageResponse is a structure for all properties of
// an image for a non detailed query
type ImageResponse struct {
	CheckSum        string `json:"checksum"`
	ContainerFormat string `json:"container_format"`
	DiskFormat      string `json:"disk_format"`
	Id              string `json:"id"`
	Name            string `json:"name"`
	Size            int64  `json:"size"`
}

// ImageDetailResponse is a structure for all properties of
// an image for a detailed query
type ImageDetailResponse struct {
	CheckSum        string                `json:"checksum"`
	ContainerFormat string                `json:"container_format"`
	CreatedAt       misc.RFC8601DateTime  `json:"created_at"`
	Deleted         bool                  `json:"deleted"`
	DeletedAt       *misc.RFC8601DateTime `json:"deleted_at"`
	DiskFormat      string                `json:"disk_format"`
	Id              string                `json:"id"`
	IsPublic        bool                  `json:"is_public"`
	MinDisk         int64                 `json:"min_disk"`
	MinRam          int64                 `json:"min_ram"`
	Name            string                `json:"name"`
	Owner           *string               `json:"owner"`
	UpdatedAt       misc.RFC8601DateTime  `json:"updated_at"`
	Properties      map[string]string     `json:"properties"`
	Protected       bool                  `json:"protected"`
	Status          string                `json:"status"`
	Size            int64                 `json:"size"`
	VirtualSize     *int64                `json:"virtual_size"` // Note: Property exists in OpenStack dev stack payloads but not Helion public cloud.
}

// ImageQueryParameters is a structure that
// contains the filter, sort, and paging parameters for
// an image or imagedetail query.
type ImageQueryParameters struct{ url.Values }

// NewImageQueryParameters creates an initialized value. This value can
// then be used to add the supported parameters. Multiple values for the
// same query parameter overwrite previous values.
func NewImageQueryParameters() ImageQueryParameters {
	return ImageQueryParameters{url.Values{}}
}

// NameFilter will add a "name" query parameter with the value that is specified.
// All the images with the specified Name will be retrieved.
func (i *ImageQueryParameters) NameFilter(value string) *ImageQueryParameters {
	i.Set("name", value)
	return i
}

// StatusFiler will add a "status" query parameter with the value that is specified.
// All the images with the specified Status will be retrieved.
func (i *ImageQueryParameters) StatusFilter(value string) *ImageQueryParameters {
	i.Set("status", value)
	return i
}

// ContainerFormat will add a "container_format" query parameter with the value that is specified.
// All the images with the specified ContainerFormat will be retrieved.
func (i *ImageQueryParameters) ContainerFormatFilter(value string) *ImageQueryParameters {
	i.Set("container_format", value)
	return i
}

// DiskFormatFilter will add a "disk_format" query parameter with the value that is specified.
// All the images with the specified DiskFormat will be retrieved.
func (i *ImageQueryParameters) DiskFormatFilter(value string) *ImageQueryParameters {
	i.Set("disk_format", value)
	return i
}

// MinSizeFilter will add a "size_min" query parameter with the value that is specified.
// All the images with at least the min size bytes will be retrieved. If MaxSizeFilter
// is also specified then all images will be within the range of min and max specified.
func (i *ImageQueryParameters) MinSizeFilter(value int64) *ImageQueryParameters {
	i.Set("size_min", fmt.Sprintf("%d", value))
	return i
}

// MaxSizeFilter will add a "size_max" query parameter with the value that is specified.
// All the images with no more than max size bytes will be retrieved. If MinSizeFilter
// is also specified then all images will be within the range of min and max specified.
func (i *ImageQueryParameters) MaxSizeFilter(value int64) *ImageQueryParameters {
	i.Set("size_max", fmt.Sprintf("%d", value))
	return i
}

// SortKey will add a "sort_key" query parameter with the value that is specified.
// The value of the SortKey can only be "name", "status", "container_format", "disk_format",
// "size", "id", "created_at", or "updated_at" when querying for Images or ImagesDetails
func (i *ImageQueryParameters) SortKey(value string) *ImageQueryParameters {
	i.Set("sort_key", value)
	return i
}

// SortDirection will add a "sort_dir" query parameter with the specified value. This will
// ensure that the sort will be ordered as ascending or descending.
// "asc" and "desc" are the only allowed values that can be specified for the sort direction.
func (i *ImageQueryParameters) SortDirection(value SortDirection) *ImageQueryParameters {
	i.Set("sort_dir", string(value))
	return i
}

// MarkerSort will add a "marker" query parameter with the value that is specified.
// The value specified must be an image id value. All the images that are after the
// specified image will be returned. Marker and Limit query parameters can be used
// in combination to get a specific page of image results.
func (i *ImageQueryParameters) Marker(value string) *ImageQueryParameters {
	i.Set("marker", value)
	return i
}

// Limit will add a "limit" query parameter with the value that is specified.
// The number of Images returned will not be larger than the number specified.
func (i *ImageQueryParameters) Limit(value int64) *ImageQueryParameters {
	i.Set("limit", fmt.Sprintf("%d", value))
	return i
}

// Direction of the sort, ascending or descending.
type SortDirection string

const (
	Desc SortDirection = "desc"
	Asc  SortDirection = "asc"
)

// Images will issue a get request to OpenStack to retrieve the list of images.
func (imageService ImageService) Images() (image []ImageResponse, err error) {
	return imageService.QueryImages(nil)
}

// ImagesDetail will issue a get request to OpenStack to retrieve the list of images complete with
// additional details.
func (imageService ImageService) ImagesDetail() (image []ImageDetailResponse, err error) {
	return imageService.QueryImagesDetail(nil)
}

// QueryImages will issue a get request with the specified ImageQueryParameters to retrieve the list of
// images.
func (imageService ImageService) QueryImages(queryParameters *ImageQueryParameters) ([]ImageResponse, error) {
	imagesContainer := imagesResponse{}
	err := imageService.queryImages(false /*includeDetails*/, &imagesContainer, queryParameters)
	if err != nil {
		return nil, err
	}

	return imagesContainer.Images, nil
}

// QueryImagesDetails will issue a get request with the specified ImageQueryParameters to retrieve the list of
// images with additional details.
func (imageService ImageService) QueryImagesDetail(queryParameters *ImageQueryParameters) ([]ImageDetailResponse, error) {
	imagesDetailContainer := imagesDetailResponse{}
	err := imageService.queryImages(true /*includeDetails*/, &imagesDetailContainer, queryParameters)
	if err != nil {
		return nil, err
	}

	return imagesDetailContainer.Images, nil
}

func (imageService ImageService) queryImages(includeDetails bool, imagesResponseContainer interface{}, queryParameters *ImageQueryParameters) error {
	urlPostFix := "/images"
	if includeDetails {
		urlPostFix = urlPostFix + "/detail"
	}

	reqUrl, err := buildQueryUrl(imageService, queryParameters, urlPostFix)
	if err != nil {
		return err
	}

	err = misc.GetJSON(reqUrl.String(), imageService.TokenId, imageService.Client, &imagesResponseContainer)
	if err != nil {
		return err
	}

	return nil
}

func buildQueryUrl(imageService ImageService, queryParameters *ImageQueryParameters, imagePartialUrl string) (url *url.URL, err error) {
	url, err = url.Parse(imageService.Url)
	if err != nil {
		return nil, err
	}

	if queryParameters != nil {
		url.RawQuery = queryParameters.Encode()
	}
	url.Path += imagePartialUrl

	return url, nil
}

type imagesDetailResponse struct {
	Images []ImageDetailResponse `json:"images"`
}

type imagesResponse struct {
	Images []ImageResponse `json:"images"`
}
