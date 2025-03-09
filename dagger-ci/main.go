// A generated module for GoTemplate functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/go-template/internal/dagger"
)

type GoTemplate struct{}

// Returns a container that echoes whatever string argument is provided
func (m *GoTemplate) ContainerEcho(stringArg string) *dagger.Container {
	return dag.Container().From("nexus-docker.aban.io/alpine:latest").WithExec([]string{"echo", stringArg})
}

// Returns lines that match a pattern in the files of the provided Directory
func (m *GoTemplate) GrepDir(ctx context.Context, directoryArg *dagger.Directory, pattern string) (string, error) {
	return dag.Container().
		From("nexus-docker.io/alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}

func (m *GoTemplate) GoBuilder(ctx context.Context, src *dagger.Directory) *dagger.Directory {
	gbuilder := dag.Golang(src)

	id, err := gbuilder.ID(ctx)

	if err != nil {
		return nil
	}

	container := dag.LoadContainerFromID(dagger.ContainerID(id))

	container.WithEnvVariable("GOPROXY", "https://goproxy.io,direct")

	return gbuilder.Build(
		dagger.GolangBuildOpts{
			Main: "./...",
		},
	)
}
