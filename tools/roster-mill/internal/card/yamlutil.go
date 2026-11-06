package card

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type yamlScalar = yaml.Node

func yamlFloat(f float64) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: fmt.Sprintf("%g", f)}
}

func yamlInt(n int) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%d", n)}
}
