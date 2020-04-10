package inform

import (
	// Standard library
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	// Third-party packages
	"github.com/pkg/errors"
)

// Default arguments for executable dependencies.
var (
	// The prefix used for Frotz meta-commands.
	frotzMetaPrefix = "\\"
	frotzArgs       = []string{"-r", "lt", "-r", "cm", "-r", "ch1", "-p", "-m", "-R"}
)

type Session struct {
	path string
	name string

	proc *os.Process

	in  io.WriteCloser
	out io.ReadCloser
	err io.ReadCloser
}

func (s *Session) Run(cmd string) error {
	switch strings.ToLower(cmd) {
	case "restore":
		return errors.New("TODO: Implement restoring")
	case "save":
		return errors.New("TODO: Implement saving")
	case "\\x", "quit":
		return errors.New("TODO: Implement session stopping")
	case "script", "unscript":
		return errors.New("transcripts are disabled")
	case "\\<", "\\>", "\\^", "\\.": // Cursor motion
	case "\\1", "\\2", "\\3", "\\4", "\\5", "\\6", "\\7", "\\8", "\\9", "\\0": // F1 - F10
	case "\\n", "\\u": // Frotz hot-keys
	default:
		if strings.HasPrefix(cmd, frotzMetaPrefix) {
			return errors.New("meta-commands are disabled")
		}
	}

	if _, err := s.in.Write(append([]byte(cmd), '\n')); err != nil {
		return errors.New("failed writing command")
	}

	return s.Error()
}

func (s *Session) Output() string {
	var buf = bytes.TrimSuffix(readPipe(s.out), []byte{'\n', '>'})
	return string(buf)
}

func (s *Session) Error() error {
	var buf = bytes.ReplaceAll(readPipe(s.err), []byte{'\n'}, []byte{':', ' '})
	if len(buf) > 0 {
		return errors.New(string(buf))
	}

	return nil
}

func (s *Session) Start(ctx context.Context, conf *Config) error {
	var err error
	var cmd = exec.CommandContext(ctx, conf.DumbFrotz, append(frotzArgs, s.path, s.name)...)

	if s.in, err = cmd.StdinPipe(); err != nil {
		return errors.Wrap(err, "starting session failed")
	} else if s.out, err = cmd.StdoutPipe(); err != nil {
		return errors.Wrap(err, "starting session failed")
	} else if s.err, err = cmd.StderrPipe(); err != nil {
		return errors.Wrap(err, "starting session failed")
	} else if err = cmd.Start(); err != nil {
		return errors.Wrap(err, "starting session failed")
	} else if err := s.Error(); err != nil {
		return err
	}

	s.proc = cmd.Process
	return nil
}

func (s *Session) Close() error {
	var err error
	if s.proc != nil {
		err = s.proc.Kill()
		s.proc = nil
	}

	return err
}

func NewSession(story *Story) (*Session, error) {
	dir, err := ioutil.TempDir(os.TempDir(), fmt.Sprintf("%s-%s-%s-*", keyPrefix, story.AuthorID, story.Name))
	if err != nil {
		return nil, errors.Wrap(err, "creating temporary directory failed")
	}

	f, err := os.Create(path.Join(dir, "output.z8"))
	if err != nil {
		return nil, errors.Wrap(err, "writing temporary story file failed")
	}

	defer f.Close()
	if _, err = f.Write(story.Build); err != nil {
		return nil, errors.Wrap(err, "writing temporary story file failed")
	}

	return &Session{
		path: dir,
		name: f.Name(),
	}, nil
}

func readPipe(r io.Reader) []byte {
	var chunks = make(chan []byte)
	go func() {
		var chunk = make([]byte, 1024)
		for {
			n, err := r.Read(chunk)
			if err != nil || n == 0 {
				close(chunks)
				return
			}

			chunks <- chunk[:n]
			if n < 1024 {
				close(chunks)
				return
			}
		}
	}()

	var buf []byte
	for {
		select {
		case chunk, ok := <-chunks:
			if !ok {
				return buf
			}
			buf = append(buf, chunk...)
		case <-time.After(10 * time.Millisecond):
			return buf
		}
	}
}
