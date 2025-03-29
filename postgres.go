package ctr

import (
	"fmt"
	"os"
	"time"

	"lesiw.io/ctrctl"
	"lesiw.io/defers"
)

func Postgres() error {
	id, err := ctrctl.ContainerCreate(&ctrctl.ContainerCreateOpts{
		Env:     []string{"POSTGRES_PASSWORD=postgres"},
		Publish: []string{"5432"},
	}, "postgres:17", "")
	if err != nil {
		return fmt.Errorf("could not create postgres container: %w", err)
	}
	defers.Add(func() {
		_, _ = ctrctl.ContainerRm(&ctrctl.ContainerRmOpts{Force: true}, id)
	})
	_, err = ctrctl.ContainerStart(nil, id)
	if err != nil {
		return fmt.Errorf("could not start postgres container: %w", err)
	}
	port, err := ctrctl.ContainerInspect(&ctrctl.ContainerInspectOpts{
		// text/template was a mistake.
		Format: `{{ with (index .NetworkSettings.Ports "5432/tcp" 0) }}` +
			`{{ .HostPort }}{{ end }}`,
	}, id)
	if err != nil {
		return fmt.Errorf("could not get port of postgres container: %w", err)
	}
	os.Setenv("PGHOST", "localhost")
	os.Setenv("PGPORT", port)
	os.Setenv("PGUSER", "postgres")
	os.Setenv("PGDATABASE", "postgres")
	os.Setenv("PGPASSWORD", "postgres")
	for {
		_, err := ctrctl.ContainerExec(nil, id,
			"psql", "-U", "postgres", "-c", "SELECT VERSION();")
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}
