package main

import (
	"fmt"
	"os"
	"testing"

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
func TestMandatoryTags(t *testing.T) {

	for _, instance := range tfplan {
		if contains(taggableResources, instance.Type) {
			assert.Contains(t, instance.Attributes, "tags.qventus:stack", "All taggable resources should be tagged with 'qventus:stack'")
			assert.Contains(t, instance.Attributes, "tags.qventus:customer", "All taggable resources should be tagged with 'qventus:customer'")
			assert.Contains(t, instance.Attributes, "tags.qventus:environment", "All taggable resources should be tagged with 'qventus:environment'")
			assert.Contains(t, instance.Attributes, "tags.terraform:commitHash", "All taggable resources should be tagged with 'terraform:commitHash'")
			assert.Contains(t, instance.Attributes, "tags.role", "All taggable resources should be tagged with 'role'")

		}
	}
}

func TestEC2Tags(t *testing.T) {
	for _, instance := range tfplan {
		if contains(ec2Resources, instance.Type) {
			assert.Contains(t, instance.Attributes, "tags.Name", "All taggable resources should be tagged with 'Name'")
		}
	}
}

func TestEc2TerminationProtection(t *testing.T) {
	for _, instance := range tfplan {
		if contains([]string{"aws_instance", "aws_launch_template"}, instance.Type) {
			assert.Equal(t, "true", instance.Attributes["disable_api_termination"], "EC2 should be protected against termination")
		}
	}
}
