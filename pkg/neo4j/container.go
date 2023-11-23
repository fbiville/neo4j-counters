package container

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ContainerConfiguration struct {
	Neo4jVersion string
	Username     string
	Password     string
}

func (config ContainerConfiguration) neo4jAuthEnvVar() string {
	return fmt.Sprintf("%s/%s", config.Username, config.Password)
}

func (config ContainerConfiguration) neo4jAuthToken() neo4j.AuthToken {
	return neo4j.BasicAuth(config.Username, config.Password, "")
}

func StartSingleInstance(ctx context.Context, config ContainerConfiguration) (testcontainers.Container, neo4j.DriverWithContext, error) {
	container, err := startContainer(ctx, config)
	if err != nil {
		return nil, nil, err
	}
	driver, err := newNeo4jDriver(ctx, container, config.neo4jAuthToken())
	if err != nil {
		_ = container.Stop(ctx, nil)
		return nil, nil, err
	}
	if err := driver.VerifyConnectivity(ctx); err != nil {
		_ = container.Stop(ctx, nil)
		_ = driver.Close(ctx)
		return nil, nil, err
	}
	return container, driver, err
}

func startContainer(ctx context.Context, config ContainerConfiguration) (testcontainers.Container, error) {
	version := config.Neo4jVersion
	request := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("neo4j:%s", version),
		ExposedPorts: []string{"7687/tcp"},
		Env: map[string]string{
			"NEO4J_AUTH":                     config.neo4jAuthEnvVar(),
			"NEO4J_ACCEPT_LICENSE_AGREEMENT": "yes",
		},
		WaitingFor: boltReadyStrategy(),
	}
	return testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: request,
			Started:          true,
		})
}

func boltReadyStrategy() *wait.LogStrategy {
	return wait.ForLog("Bolt enabled")
}

func newNeo4jDriver(ctx context.Context, container testcontainers.Container, auth neo4j.AuthToken) (neo4j.DriverWithContext, error) {
	port, err := container.MappedPort(ctx, "7687")
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("neo4j://localhost:%d", port.Int())
	return neo4j.NewDriverWithContext(url, auth)
}
