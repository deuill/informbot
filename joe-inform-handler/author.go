package inform

import (
	// Standard library
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	// Third-party packages
	"github.com/pkg/errors"
)

// Options represent user-configurable values, which are used in processing commands
// and formatting output.
type Options struct {
	Prefix string
}

// Default values for options, as assigned to newly created Author instances.
var defaultOptions = Options{
	Prefix: "?",
}

type Author struct {
	ID      string
	Options Options
	Stories []*Story
}

func (a *Author) GetStory(name string) (*Story, error) {
	if name == "" {
		return nil, errors.New("story name is empty")
	}

	for i := range a.Stories {
		if a.Stories[i].Name == name {
			return a.Stories[i], nil
		}
	}

	return nil, errors.New("no story found with name '" + name + "'")
}

func (a *Author) AddStory(name, path string) (*Story, error) {
	var story *Story
	if name == "" {
		return nil, errors.New("story name is empty")
	} else if s, _ := a.GetStory(name); s != nil {
		story = s
	} else if u, err := url.Parse(path); err != nil || (u.Scheme != "https" && u.Scheme != "http") {
		return nil, errors.New("location given is not a valid HTTP URL")
	}

	resp, err := http.Get(path)
	if err != nil {
		return nil, errors.New("could not fetch story file from URL given")
	}

	defer resp.Body.Close()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, errors.New("could not fetch story file from URL given")
	} else if story == nil {
		story = NewStory(name, a.ID).WithSource(body)
		a.Stories = append(a.Stories, story)
	} else {
		story.WithSource(body)
	}

	return story, nil
}

func (a *Author) RemoveStory(name string) error {
	if name == "" {
		return errors.New("story name is empty")
	}

	// Find story with given name, and remove from list assigned to author.
	for i := range a.Stories {
		if a.Stories[i].Name == name {
			a.Stories = append(a.Stories[:i], a.Stories[i+1:]...)
			return nil
		}
	}

	return errors.New("no story found with name '" + name + "'")
}

func (a *Author) SetOption(name, value string) error {
	switch strings.ToLower(name) {
	case "prefix":
		if value == "" {
			return errors.New("cannot set empty prefix value")
		}
		a.Options.Prefix = value
	default:
		return errors.New("option name '" + name + "' is unknown")
	}

	return nil
}

func NewAuthor(id string) *Author {
	return &Author{ID: id, Options: defaultOptions}
}
