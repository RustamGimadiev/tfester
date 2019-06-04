package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/stretchr/testify/assert"
)

var tfplan = []Resource{}

func init() {
	_, tf, err := tfjson(os.Args[1])
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
	tfplan = tf

}
func TestValidateTags(t *testing.T) {

	for _, instance := range tfplan {
		if isTaggable(instance.Type) {
			assert.Contains(t, instance.Attributes, "tags.Name", "All taggable resources should be tagged with Name")
			assert.Contains(t, instance.Attributes, "tags.Owner", "All taggable resources should be tagged with Owner")
		}
	}
}

func TestValidateTagNameStyle(t *testing.T) {
	for _, instance := range tfplan {
		if isTaggable(instance.Type) {
			for k, v := range instance.Attributes {
				if strings.HasPrefix(k, "tags.") {
					assert.Equal(t, strcase.ToLowerCamel(v), v, "All tags value should be formated with camel case")
				}
			}
		}
	}
}

func TestSecurityGroupRules(t *testing.T) {
	for _, instance := range tfplan {
		if instance.Type == "aws_security_group" {
			for k, v := range instance.Attributes {
				if strings.HasPrefix(k, "ingress") && strings.HasSuffix(k, "cidr_blocks.0") {
					assert.NotEqual(t, "0.0.0.0/0", v, "Security groups should be not world open")
				}
			}
		}
	}
}

func TestEc2TerminationProtection(t *testing.T) {
	for _, instance := range tfplan {
		if instance.Type == "aws_instance" {
			assert.Equal(t, "true", instance.Attributes["disable_api_termination"], "EC2 should be protected against termination")
		}
	}
}
