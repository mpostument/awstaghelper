package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/stretchr/testify/assert"
)

type mockedIamUser struct {
	iamiface.IAMAPI
	respListUsers    iam.ListUsersOutput
	respListUserTags iam.ListUserTagsOutput
}

func (m *mockedIamUser) ListUsersPages(input *iam.ListUsersInput, pageFunc func(*iam.ListUsersOutput, bool) bool) error {
	pageFunc(&m.respListUsers, true)
	return nil
}

func (m *mockedIamUser) ListUserTags(*iam.ListUserTagsInput) (*iam.ListUserTagsOutput, error) {
	return &m.respListUserTags, nil
}

func TestGetIamUsers(t *testing.T) {
	cases := []*mockedIamUser{
		{
			respListUsers: listUsersResponse,
		},
	}

	expectedResult := listUsersResponse.Users

	for _, c := range cases {
		t.Run("getIamUsers", func(t *testing.T) {
			result := getIamUsers(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseIamUserTags(t *testing.T) {
	cases := []*mockedIamUser{
		{
			respListUsers:    listUsersResponse,
			respListUserTags: listUserTagsResponse,
		},
	}

	expectedResult := [][]string{
		{"UserName", "Name", "Owner"},
		{"test-user1", "test-user1", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseIamUserTags", func(t *testing.T) {
			result := ParseIamUserTags("Name,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var listUsersResponse = iam.ListUsersOutput{
	Users: []*iam.User{
		{
			UserName: aws.String("test-user1"),
		},
	},
}

var listUserTagsResponse = iam.ListUserTagsOutput{
	Tags: []*iam.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String("test-user1"),
		},
		{
			Key:   aws.String("Owner"),
			Value: aws.String("mpostument"),
		},
	},
}

type mockedIamRole struct {
	iamiface.IAMAPI
	respListRoles     iam.ListRolesOutput
	respListRolesTags iam.ListRoleTagsOutput
}

func (m *mockedIamRole) ListRolesPages(input *iam.ListRolesInput, pageFunc func(*iam.ListRolesOutput, bool) bool) error {
	pageFunc(&m.respListRoles, true)
	return nil
}

func (m *mockedIamRole) ListRoleTags(*iam.ListRoleTagsInput) (*iam.ListRoleTagsOutput, error) {
	return &m.respListRolesTags, nil
}

func TestGetIamRoles(t *testing.T) {
	cases := []*mockedIamRole{
		{
			respListRoles: listRolesResponse,
		},
	}

	expectedResult := listRolesResponse.Roles

	for _, c := range cases {
		t.Run("getIamRoles", func(t *testing.T) {
			result := getIamRoles(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseIamRolesTags(t *testing.T) {
	cases := []*mockedIamRole{
		{
			respListRoles:     listRolesResponse,
			respListRolesTags: listRolesTagsResponse,
		},
	}

	expectedResult := [][]string{
		{"RoleName", "Name", "Owner"},
		{"test-role1", "test-role1", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseIamRolesTags", func(t *testing.T) {
			result := ParseIamRolesTags("Name,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var listRolesResponse = iam.ListRolesOutput{
	Roles: []*iam.Role{
		{
			RoleName: aws.String("test-role1"),
		},
	},
}

var listRolesTagsResponse = iam.ListRoleTagsOutput{
	Tags: []*iam.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String("test-role1"),
		},
		{
			Key:   aws.String("Owner"),
			Value: aws.String("mpostument"),
		},
	},
}
