package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// Container contains the information about the conainer.
type Container struct {
	ID   string
	Port string
}

// StartNeo4j runs a mongo container to execute commands.
func StartNeo4j() (*Container, error) {
	cmd := exec.Command("docker", "run", "--publish=7687:7687", "-d", "--env=NEO4J_AUTH=none", "neo4j")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("starting container: %v", err)
	}

	id := out.String()[:12]
	log.Println("DB ContainerID:", id)

	cmd = exec.Command("docker", "inspect", id)
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("inspect container: %v", err)
	}

	var doc []struct {
		NetworkSettings struct {
			Ports struct {
				TCP7687 []struct {
					HostPort string `json:"HostPort"`
				} `json:"7687/tcp"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		return nil, fmt.Errorf("decoding json: %v", err)
	}

	c := Container{
		ID:   id,
		Port: doc[0].NetworkSettings.Ports.TCP7687[0].HostPort,
	}

	log.Println("DB Port:", c.Port)

	return &c, nil
}

// StopNeo4j stops and removes the specified container.
func StopNeo4j(c *Container) error {
	if err := exec.Command("docker", "stop", c.ID).Run(); err != nil {
		return err
	}
	log.Println("Stopped:", c.ID)

	if err := exec.Command("docker", "rm", c.ID, "-v").Run(); err != nil {
		return err
	}
	log.Println("Removed:", c.ID)

	return nil
}
