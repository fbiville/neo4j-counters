package main

import (
	"context"
	"fmt"
	container "github.com/fbiville/neo4j-counters/pkg/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"os"
)

func main() {
	neo4jVersion := "5"
	if len(os.Args) > 1 {
		neo4jVersion = os.Args[1]
	}
	ctx := context.Background()

	db, driver, err := container.StartSingleInstance(ctx, container.ContainerConfiguration{
		Neo4jVersion: neo4jVersion,
		Username:     "neo4j",
		Password:     "letmein!",
	})
	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)
	defer db.Stop(ctx, nil)

	result, err := neo4j.ExecuteQuery(ctx, driver, "CREATE (:Node {prop: 'hello'})", nil, neo4j.EagerResultTransformer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("PropertiesSet counter: %d\n", result.Summary.Counters().PropertiesSet())

	result, err = neo4j.ExecuteQuery(ctx, driver, "MATCH (n:Node) SET n.prop = 'bonjour'", nil, neo4j.EagerResultTransformer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("PropertiesSet counter: %d\n", result.Summary.Counters().PropertiesSet())

	result, err = neo4j.ExecuteQuery(ctx, driver, "MATCH (n:Node) REMOVE n.prop", nil, neo4j.EagerResultTransformer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("PropertiesSet counter: %d\n", result.Summary.Counters().PropertiesSet())
}
