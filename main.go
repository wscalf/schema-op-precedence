package main

import (
	"fmt"
	"os"

	"github.com/authzed/spicedb/pkg/namespace"
	corev1 "github.com/authzed/spicedb/pkg/proto/core/v1"
	"github.com/authzed/spicedb/pkg/schemadsl/compiler"
	"github.com/authzed/spicedb/pkg/schemadsl/generator"
	"github.com/authzed/spicedb/pkg/schemadsl/input"
)

func main() {
	dsl, err := os.ReadFile("schema.zed")
	if err != nil {
		fmt.Println(err)
		return
	}

	source := input.Source("schema.zed")
	compiled, err := compiler.Compile(compiler.InputSchema{Source: source, SchemaString: string(dsl)}, compiler.AllowUnprefixedObjectType())
	if err != nil {
		fmt.Println(err)
		return
	}

	role_binding := findDefinition(compiled, "role_binding")
	if role_binding == nil {
		fmt.Println("role_binding not found")
		return
	}

	role_binding.Relation = append(role_binding.Relation, namespace.MustRelation("read", namespace.Intersection(
		namespace.ComputedUserset("subject"),
		namespace.Rewrite(
			namespace.Union(
				namespace.TupleToUserset("role", "can_read_things"),
				namespace.TupleToUserset("role", "things_admin"),
				namespace.TupleToUserset("role", "global_admin"),
			),
		))))

	result, _, err := generator.GenerateSchema(compiled.OrderedDefinitions)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Generated schema: %s\n", result)
}

func findDefinition(schema *compiler.CompiledSchema, name string) *corev1.NamespaceDefinition {
	for _, d := range schema.ObjectDefinitions {
		if d.Name == name {
			return d
		}
	}

	return nil
}
