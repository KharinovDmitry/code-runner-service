package dockerCLI

import "fmt"

type DockerCommandBuilder struct {
	command []string
}

func NewDockerCommandBuilder() *DockerCommandBuilder {
	return &DockerCommandBuilder{
		command: []string{"docker", "run", "-rm", "-d"},
	}
}

func (builder *DockerCommandBuilder) SetName(name string) {
	builder.command = append(builder.command, "--name", name)
}

func (builder *DockerCommandBuilder) Mount(from, to string) {
	builder.command = append(builder.command, "--mount", fmt.Sprintf("type=bind,source=%s,target=%s", from, to))
}

func (builder *DockerCommandBuilder) SetMemoryLimit(memLimitKb int) {
	builder.command = append(builder.command, "--memory-limit", fmt.Sprintf("%d", memLimitKb))
}

func (builder *DockerCommandBuilder) Args() []string {
	return builder.command
}
