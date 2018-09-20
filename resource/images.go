package resource

import (
	"sort"

	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/tyrsession"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Images []*ec2.Image

func (im Images) SortByDate() {
	sort.SliceStable(im, func(i, j int) bool {
		return *(im)[i].CreationDate < *(im)[j].CreationDate
	})
}

func (im *Images) FindByTags(tags map[string]string) Images {
	found := Images{}
	n := 0
	for _, image := range *im {
		for _, tag := range image.Tags {
			if v, ok := tags[*tag.Key]; ok && v == *tag.Value {
				n++
				if len(tags) == n {
					found = append(found, image)
				}
			}
		}
		n = 0
	}
	return found
}

func (im *Images) LoadFromAWS(config *configuration.Config, region string) error {
	ec2API, err := config.ClientFactory.GetEc2Client(tyrsession.SessionConfig{Profile: config.Profile, Region: region})
	if err != nil {
		return err
	}
	result, err := ec2API.DescribeImages(&ec2.DescribeImagesInput{
		Owners: []*string{aws.String("self")},
	})
	if err != nil {
		return err
	}
	*im = result.Images
	return err
}
