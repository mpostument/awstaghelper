/*
Copyright © 2024 Jaemok Hong jaemokhong@lguplus.co.kr
Copyright © 2020 Maksym Postument 777rip777@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/aws/aws-sdk-go/service/wafv2/wafv2iface"
	"github.com/stretchr/testify/assert"
)

type mockedWebACL struct {
	wafv2iface.WAFV2API
	respListWebACLs         wafv2.ListWebACLsOutput
	respListTagsForResource wafv2.ListTagsForResourceOutput
}

func (m *mockedWebACL) ListWebACLs(*wafv2.ListWebACLsInput) (*wafv2.ListWebACLsOutput, error) {
	return &m.respListWebACLs, nil
}

func (m *mockedWebACL) ListTagsForResource(*wafv2.ListTagsForResourceInput) (*wafv2.ListTagsForResourceOutput, error) {
	return &m.respListTagsForResource, nil
}

func TestGetWebACLs(t *testing.T) {
	cases := []*mockedWebACL{
		{
			respListWebACLs: listWebACLsResponse,
		},
	}

	expectedResult := &listWebACLsResponse

	for _, c := range cases {
		t.Run("getWebACL", func(t *testing.T) {
			result := getWebACL("REGIONAL", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseWebACLTags(t *testing.T) {
	cases := []*mockedWebACL{
		{
			respListWebACLs:         listWebACLsResponse,
			respListTagsForResource: listWebACLsTags,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:wafv2:us-east-1:084888172679:regional/webacl/test-webacl/12345678-abcd-1234-abcd-12345678abcd", "test-webacl", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseWebACLTags", func(t *testing.T) {
			result := ParseWebACLTags("Name,Owner", "REGIONAL", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var listWebACLsResponse = wafv2.ListWebACLsOutput{
	WebACLs: []*wafv2.WebACLSummary{
		{
			ARN: aws.String("arn:aws:wafv2:us-east-1:084888172679:regional/webacl/test-webacl/12345678-abcd-1234-abcd-12345678abcd"),
		},
	},
}

var listWebACLsTags = wafv2.ListTagsForResourceOutput{
	TagInfoForResource: &wafv2.TagInfoForResource{
		TagList: []*wafv2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("test-webacl"),
			},
			{
				Key:   aws.String("Owner"),
				Value: aws.String("mpostument"),
			},
		},
	},
}

type mockedCFWebACL struct {
	wafv2iface.WAFV2API
	respListWebACLs         wafv2.ListWebACLsOutput
	respListTagsForResource wafv2.ListTagsForResourceOutput
}

func (m *mockedCFWebACL) ListWebACLs(*wafv2.ListWebACLsInput) (*wafv2.ListWebACLsOutput, error) {
	return &m.respListWebACLs, nil
}

func (m *mockedCFWebACL) ListTagsForResource(*wafv2.ListTagsForResourceInput) (*wafv2.ListTagsForResourceOutput, error) {
	return &m.respListTagsForResource, nil
}

func TestGetCFWebACLs(t *testing.T) {
	cases := []*mockedCFWebACL{
		{
			respListWebACLs: listCFWebACLsResponse,
		},
	}

	expectedResult := &listCFWebACLsResponse

	for _, c := range cases {
		t.Run("getWebACLs", func(t *testing.T) {
			result := getWebACL("REGIONAL", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseCFWebACLTags(t *testing.T) {
	cases := []*mockedCFWebACL{
		{
			respListWebACLs:         listCFWebACLsResponse,
			respListTagsForResource: listCFWebACLsTags,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:wafv2:us-east-1:084888172679:global/webacl/test-webacl/12345678-abcd-1234-abcd-12345678abcd", "test-webacl", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseWebACLTags", func(t *testing.T) {
			result := ParseWebACLTags("Name,Owner", "CLOUDFRONT", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var listCFWebACLsResponse = wafv2.ListWebACLsOutput{
	WebACLs: []*wafv2.WebACLSummary{
		{
			ARN: aws.String("arn:aws:wafv2:us-east-1:084888172679:global/webacl/test-webacl/12345678-abcd-1234-abcd-12345678abcd"),
		},
	},
}

var listCFWebACLsTags = wafv2.ListTagsForResourceOutput{
	TagInfoForResource: &wafv2.TagInfoForResource{
		TagList: []*wafv2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("test-webacl"),
			},
			{
				Key:   aws.String("Owner"),
				Value: aws.String("mpostument"),
			},
		},
	},
}
