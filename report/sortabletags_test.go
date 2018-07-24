package report

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
)

func TestIfTagsAreSortedCorrectly(t *testing.T) {

	st := NewSortableTags()
	ec2Tags := []*ec2.Tag{
		{
			Key:   aws.String("BBBB"),
			Value: aws.String("SomeValue6"),
		},
		{
			Key:   aws.String("bbb"),
			Value: aws.String("SomeValue2"),
		},
		{
			Key:   aws.String("AAA"),
			Value: aws.String("SomeValue3"),
		},
		{
			Key:   aws.String("aaaa"),
			Value: aws.String("SomeValue4"),
		},
	}
	st.Add(ec2Tags)
	tableData := strings.Split(st.ToTableData(), "\n")
	assert.True(t, len(ec2Tags) == len(tableData))
	assert.EqualValues(t, "AAA", strings.Split(tableData[0], ":")[0])
	assert.EqualValues(t, "aaaa", strings.Split(tableData[1], ":")[0])
	assert.EqualValues(t, "bbb", strings.Split(tableData[2], ":")[0])
	assert.EqualValues(t, "BBBB", strings.Split(tableData[3], ":")[0])
}

func TestSorttableTagsIfSliceIsEmptyDoNotPanic(t *testing.T) {
	st := NewSortableTags()
	ec2Tags := []*ec2.Tag{}
	st.Add(ec2Tags)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did panic")
		}
	}()

	st.ToTableData()
}
