package main

var (
	taggableresources = []string{"aws_instance", "aws_ebs_volume"}
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func isTaggable(resourceType string) bool {
	return contains(taggableresources, resourceType)
}

func isCamelCase(word string) bool {
	return true
}
