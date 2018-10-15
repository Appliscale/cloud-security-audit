package report

import (
	"bytes"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type SortableTags struct {
	Tags map[string]string
	Keys []string
}

func NewSortableTags() *SortableTags {
	return &SortableTags{Tags: make(map[string]string)}
}

func (st *SortableTags) Add(tags []*ec2.Tag) {

	for _, tag := range tags {
		st.Keys = append(st.Keys, *tag.Key)
		st.Tags[*tag.Key] = *tag.Value
	}
	less := func(i, j int) bool {
		return strings.ToLower(st.Keys[i]) < strings.ToLower(st.Keys[j])
	}
	sort.Slice(st.Keys, less)
}

func (st *SortableTags) ToTableData() string {
	n := len(st.Keys)
	if n == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for _, key := range st.Keys[:n-1] {
		maxWidth := 50
		if len(st.Tags[key]+key) > maxWidth-1 {
			i := maxWidth - 1 - len(key)
			for i < len(st.Tags[key]) {
				st.Tags[key] = st.Tags[key][:i] + "\n  " + st.Tags[key][i:]
				i += maxWidth - 2
			}
		}
		buffer.WriteString(key + ":" + st.Tags[key] + "\n")
	}
	buffer.WriteString(st.Keys[n-1] + ":" + st.Tags[st.Keys[n-1]])
	return buffer.String()
}
