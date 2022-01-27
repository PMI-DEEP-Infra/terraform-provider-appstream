// Code is taken from aws provider 
//https://github.com/hashicorp/terraform-provider-aws/blob/9f214e53f0a1d650cb0fd4ef5951e07b1f86a2ea/internal/tags/key_value_tags.go
//https://github.com/hashicorp/terraform-provider-aws/blob/9f214e53f0a1d650cb0fd4ef5951e07b1f86a2ea/internal/service/appstream/tags_gen.go
package appstream

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appstream"
	"reflect"
  "strings"
)


const (
	AwsTagKeyPrefix                             = `aws:`
	// ElasticbeanstalkTagKeyPrefix                = `elasticbeanstalk:`
	// NameTagKey                                  = `Name`
	// RdsTagKeyPrefix                             = `rds:`
	// ServerlessApplicationRepositoryTagKeyPrefix = `serverlessrepo:`
)
// map[string]*string handling

type TagData struct {
	// Additional boolean field names and values associated with this tag.
	// Each service is responsible for properly handling this data.
	AdditionalBoolFields map[string]*bool

	// Additional string field names and values associated with this tag.
	// Each service is responsible for properly handling this data.
	AdditionalStringFields map[string]*string

	// Tag value.
	Value *string
}

type KeyValueTags map[string]*TagData

// Tags returns appstream service tags.
func Tags(tags KeyValueTags) map[string]*string {
	return aws.StringMap(tags.Map())
}

// Merge adds missing and updates existing tags.
func (tags KeyValueTags) Merge(mergeTags KeyValueTags) KeyValueTags {
	result := make(KeyValueTags)

	for k, v := range tags {
		result[k] = v
	}

	for k, v := range mergeTags {
		result[k] = v
	}

	return result
}

// Removed returns tags removed.
func (tags KeyValueTags) Removed(newTags KeyValueTags) KeyValueTags {
	result := make(KeyValueTags)

	for k, v := range tags {
		if _, ok := newTags[k]; !ok {
			result[k] = v
		}
	}

	return result
}

// Keys returns tag keys.
func (tags KeyValueTags) Keys() []string {
	result := make([]string, 0, len(tags))

	for k := range tags {
		result = append(result, k)
	}

	return result
}

// Updated returns tags added and updated.
func (tags KeyValueTags) Updated(newTags KeyValueTags) KeyValueTags {
	result := make(KeyValueTags)

	for k, newV := range newTags {
		if oldV, ok := tags[k]; !ok || !oldV.Equal(newV) {
			result[k] = newV
		}
	}

	return result
}

// Map returns tag keys mapped to their values.
func (tags KeyValueTags) Map() map[string]string {
	result := make(map[string]string, len(tags))

	for k, v := range tags {
		if v == nil || v.Value == nil {
			result[k] = ""
			continue
		}

		result[k] = *v.Value
	}

	return result
}

func (td *TagData) Equal(other *TagData) bool {
	if td == nil && other == nil {
		return true
	}

	if td == nil || other == nil {
		return false
	}

	if !reflect.DeepEqual(td.AdditionalBoolFields, other.AdditionalBoolFields) {
		return false
	}

	if !reflect.DeepEqual(td.AdditionalStringFields, other.AdditionalStringFields) {
		return false
	}

	if !reflect.DeepEqual(td.Value, other.Value) {
		return false
	}

	return true
}



func New(i interface{}) KeyValueTags {
	switch value := i.(type) {
	case KeyValueTags:
		return make(KeyValueTags).Merge(value)
	case map[string]*TagData:
		return make(KeyValueTags).Merge(KeyValueTags(value))
	case map[string]string:
		kvtm := make(KeyValueTags, len(value))

		for k, v := range value {
			str := v // Prevent referencing issues
			kvtm[k] = &TagData{Value: &str}
		}

		return kvtm
	case map[string]*string:
		kvtm := make(KeyValueTags, len(value))

		for k, v := range value {
			strPtr := v

			if strPtr == nil {
				kvtm[k] = nil
				continue
			}

			kvtm[k] = &TagData{Value: strPtr}
		}

		return kvtm
	case map[string]interface{}:
		kvtm := make(KeyValueTags, len(value))

		for k, v := range value {
			kvtm[k] = &TagData{}

			str, ok := v.(string)

			if ok {
				kvtm[k].Value = &str
			}
		}

		return kvtm
	case []string:
		kvtm := make(KeyValueTags, len(value))

		for _, v := range value {
			kvtm[v] = nil
		}

		return kvtm
	case []interface{}:
		kvtm := make(KeyValueTags, len(value))

		for _, v := range value {
			kvtm[v.(string)] = nil
		}

		return kvtm
	default:
		return make(KeyValueTags)
	}
}


// UpdateTags updates appstream service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(conn *appstream.AppStream, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := New(oldTagsMap)
	newTags := New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &appstream.UntagResourceInput{
			ResourceArn: aws.String(identifier),
			TagKeys:     aws.StringSlice(removedTags.IgnoreAWS().Keys()),
		}

		_, err := conn.UntagResource(input)

		if err != nil {
			return fmt.Errorf("error untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &appstream.TagResourceInput{
			ResourceArn: aws.String(identifier),
			Tags:        Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.TagResource(input)

		if err != nil {
			return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}


// IgnoreAWS returns non-AWS tag keys.
func (tags KeyValueTags) IgnoreAWS() KeyValueTags {
	result := make(KeyValueTags)

	for k, v := range tags {
		if !strings.HasPrefix(k, AwsTagKeyPrefix) {
			result[k] = v
		}
	}

	return result
}
