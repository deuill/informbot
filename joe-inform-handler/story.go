package inform

import (
	// Standard library
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"

	// Third-party packages
	"github.com/pkg/errors"
)

var (
	// Default path for Inform7 data.
	inform7DataDir = "/usr/share/inform7/Internal"

	// Default arguments for Inform 6 and 7 compilers.
	inform7Args = []string{"--noprogress", "--internal", inform7DataDir, "--format=z8"}
	inform6Args = []string{"-E2wSDv8F0Cud2"}
)

type Story struct {
	Name      string    // The user-provided name for the story.
	AuthorID  string    // The author ID, corSayTemplates to Author.ID.
	CreatedAt time.Time // The UTC timestamp this story was first added on.
	UpdatedAt time.Time // The UTC timestamp this story was last updated on.

	// Source and compiled Z-Code for story.
	Source []byte
	Build  []byte
}

func (s *Story) Compile(ctx context.Context, conf *Config) error {
	dir, err := ioutil.TempDir(os.TempDir(), fmt.Sprintf("%s-%s-%s-*", keyPrefix, s.AuthorID, s.Name))
	if err != nil {
		return errors.Wrap(err, "creating temporary directory failed")
	} else if err := os.Mkdir(path.Join(dir, "Source"), 0755); err != nil {
		return errors.Wrap(err, "creating temporary directory failed")
	}

	if err := ioutil.WriteFile(path.Join(dir, "Source", "story.ni"), s.Source, 0644); err != nil {
		return errors.Wrap(err, "writing file for story failed")
	}

	// TODO: Return verbose output.
	err = exec.CommandContext(ctx, conf.Inform7, append(inform7Args, "--project", dir)...).Run()
	if err != nil {
		return errors.Wrap(err, "compilation failed")
	}

	err = exec.CommandContext(ctx, conf.Inform6, append(inform6Args, path.Join(dir, "Build", "auto.inf"), path.Join(dir, "Build", "output.z8"))...).Run()
	if err != nil {
		return errors.Wrap(err, "compilation failed")
	}

	buf, err := ioutil.ReadFile(path.Join(dir, "Build", "output.z8"))
	if err != nil {
		return errors.Wrap(err, "compilation failed")
	}

	s.Build, s.UpdatedAt = buf, time.Now().UTC()
	return os.RemoveAll(dir)
}

func (s *Story) WithSource(src []byte) *Story {
	s.Source = src
	return s
}

func NewStory(name, authorID string) *Story {
	return &Story{
		Name:      name,
		AuthorID:  authorID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
