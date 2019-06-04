// tfester project main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform/terraform"
)

type Resource struct {
	Type       string
	Name       string
	Index      string
	Attributes map[string]string
}

func main() {
	j, _, err := tfjson(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(j)
}

type output map[string]interface{}

func tfjson(planfile string) (string, []Resource, error) {
	diff := output{}
	recources := []Resource{}
	f, err := os.Open(planfile)
	if err != nil {
		return "", recources, err
	}
	defer f.Close()

	plan, err := terraform.ReadPlan(f)
	if err != nil {
		return "", recources, err
	}

	for _, v := range plan.Diff.Modules {
		convertModuleDiff(diff, &recources, v)
	}

	j, err := json.MarshalIndent(diff, "", "    ")
	if err != nil {
		return "", recources, err
	}

	return string(j), recources, nil
}

func parseResource(key string, attributes interface{}) Resource {
	parts := strings.Split(key, ".")
	rtype := parts[0]
	name := parts[1]
	index := "0"
	if len(parts) == 3 {
		index = parts[2]
	}
	a := map[string]string{}
	for k, v := range attributes.(map[string]*terraform.ResourceAttrDiff) {
		a[k] = v.New
	}
	return Resource{Type: rtype, Name: name, Index: index, Attributes: a}
}

func insert(out output, path []string, key string, value interface{}) {
	if len(path) > 0 && path[0] == "root" {
		path = path[1:]
	}
	for _, elem := range path {
		switch nested := out[elem].(type) {
		case output:
			out = nested
		default:
			new := output{}
			out[elem] = new
			out = new
		}
	}
	out[key] = value
}

func convertModuleDiff(out output, res *[]Resource, diff *terraform.ModuleDiff) {
	insert(out, diff.Path, "destroy", diff.Destroy)
	for k, v := range diff.Resources {
		*res = append(*res, parseResource(k, v.Attributes))
		convertInstanceDiff(out, append(diff.Path, k), v)
	}
}

func convertInstanceDiff(out output, path []string, diff *terraform.InstanceDiff) {
	insert(out, path, "destroy", diff.Destroy)
	insert(out, path, "destroy_tainted", diff.DestroyTainted)

	for k, v := range diff.Attributes {
		insert(out, path, k, v.New)
	}
}
