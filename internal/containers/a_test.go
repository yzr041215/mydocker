package containers

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	err := SaveContainerInfo(&ContainerInfo{
		ContainerID: "123456",
		Image:       "test/image",
		Command:     "test command",
		Created:     "2021-01-01 00:00:00",
		Status:      "running",
		Pid:         12345,
	})
	if err != nil {
		fmt.Println(err)
	}
}
