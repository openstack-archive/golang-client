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
type ImageQueryParameters struct {
	Name            string
	Status          string
	ContainerFormat string
	DiskFormat      string
	MinSize         int64
	MaxSize         int64
	SortKey         string
	SortDirection   SortDirection
	Marker          string
	Limit           int64
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

func buildQueryUrl(imageService ImageService, queryParameters *ImageQueryParameters, imagePartialUrl string) (*url.URL, error) {
	reqUrl, err := url.Parse(imageService.Url)
	if err != nil {
		return nil, err
	}

	if queryParameters != nil {
		values := url.Values{}
		if queryParameters.Name != "" {
			values.Set("name", queryParameters.Name)
		}
		if queryParameters.ContainerFormat != "" {
			values.Set("container_format", queryParameters.ContainerFormat)
		}
		if queryParameters.DiskFormat != "" {
			values.Set("disk_format", queryParameters.DiskFormat)
		}
		if queryParameters.Status != "" {
			values.Set("status", queryParameters.Status)
		}
		if queryParameters.MinSize != 0 {
			values.Set("size_min", fmt.Sprintf("%d", queryParameters.MinSize))
		}
		if queryParameters.MaxSize != 0 {
			values.Set("size_max", fmt.Sprintf("%d", queryParameters.MaxSize))
		}
		if queryParameters.Limit != 0 {
			values.Set("limit", fmt.Sprintf("%d", queryParameters.Limit))
		}
		if queryParameters.Marker != "" {
			values.Set("marker", queryParameters.Marker)
		}
		if queryParameters.SortKey != "" {
			values.Set("sort_key", queryParameters.SortKey)
		}
		if queryParameters.SortDirection != "" {
			values.Set("sort_dir", string(queryParameters.SortDirection))
		}

		if len(values) > 0 {
			reqUrl.RawQuery = values.Encode()
		}
	}
	reqUrl.Path += imagePartialUrl

	return reqUrl, nil
}

type imagesDetailResponse struct {
	Images []ImageDetailResponse `json:"images"`
}

type imagesResponse struct {
	Images []ImageResponse `json:"images"`
}
