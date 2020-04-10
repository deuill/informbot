package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/go-joe/joe"
	"go.uber.org/zap"
)

// memory is an implementation of a joe.Memory which stores all values as a JSON
// encoded file. Note that there is no need for a joe.Memory to handle
// synchronization for concurrent access (e.g. via locks) because this is
// automatically handled by the joe.Brain.
type memory struct {
	path   string
	logger *zap.Logger
	data   map[string][]byte
}

// Memory is a joe.Option which is supposed to be passed to joe.New(…) to
// configure a new bot. The path indicates the destination file at which the
// memory will store its values encoded as JSON object. If there is already a
// JSON encoded file at the given path it will be loaded and decoded into memory
// to serve future requests. If the file exists but cannot be opened or does not
// contain a valid JSON object its error will be deferred until the bot is
// actually started via its Run() function.
//
// Example usage:
//     b := joe.New("example",
//         file.Memory("/tmp/joe.json"),
//         …
//     )
func Memory(path string) joe.Module {
	return joe.ModuleFunc(func(conf *joe.Config) error {
		memory, err := NewMemory(path, WithLogger(conf.Logger("memory")))
		if err != nil {
			return err
		}

		conf.SetMemory(memory)
		return nil
	})
}

// NewMemory creates a new Memory instance that persists all values to the given
// path. If there is already a JSON encoded file at the given path it is loaded
// and decoded into memory to serve future requests. An error is returned if the
// file exists but cannot be opened or does not contain a valid JSON object.
func NewMemory(path string, opts ...Option) (joe.Memory, error) {
	memory := &memory{
		path: path,
		data: map[string][]byte{},
	}

	for _, opt := range opts {
		err := opt(memory)
		if err != nil {
			return nil, err
		}
	}

	if memory.logger == nil {
		memory.logger = zap.NewNop()
	}

	memory.logger.Debug("Opening memory file", zap.String("path", path))
	f, err := os.Open(path)
	switch {
	case os.IsNotExist(err):
		memory.logger.Debug("File does not exist. Continuing with empty memory", zap.String("path", path))
	case err != nil:
		return nil, fmt.Errorf("failed to open file: %w", err)
	default:
		memory.logger.Debug("Decoding JSON from memory file", zap.String("path", path))
		err := json.NewDecoder(f).Decode(&memory.data)
		_ = f.Close()
		if err != nil {
			return nil, fmt.Errorf("failed decode data as JSON: %w", err)
		}
	}

	memory.logger.Info("Memory initialized successfully",
		zap.String("path", path),
		zap.Int("num_memories", len(memory.data)),
	)

	return memory, nil
}

// Set assign the key to the value and then saves the updated memory to its JSON
// file. An error is returned if this function is called after the memory was
// closed already or if the file could not be written or updated.
func (m *memory) Set(key string, value []byte) error {
	if m.data == nil {
		return errors.New("brain was already shut down")
	}

	m.data[key] = value
	return m.persist()
}

// Get returns the value that is associated with the given key. The second
// return value indicates if the key actually existed in the memory.
//
// An error is only returned if this function is called after the memory was
// closed already.
func (m *memory) Get(key string) ([]byte, bool, error) {
	if m.data == nil {
		return nil, false, errors.New("brain was already shut down")
	}

	value, ok := m.data[key]
	return value, ok, nil
}

// Delete removes any value that might have been assigned to the key earlier.
// The boolean return value indicates if the memory contained the key. If it did
// not contain the key the function does nothing and returns without an error.
// If the key existed it is removed and the corresponding JSON file is updated.
//
// An error is returned if this function is called after the memory was closed
// already or if the file could not be written or updated.
func (m *memory) Delete(key string) (bool, error) {
	if m.data == nil {
		return false, errors.New("brain was already shut down")
	}

	_, ok := m.data[key]
	if !ok {
		return false, nil
	}

	delete(m.data, key)
	return ok, m.persist()
}

// Keys returns a list of all keys known to this memory.
// An error is only returned if this function is called after the memory was
// closed already.
func (m *memory) Keys() ([]string, error) {
	if m.data == nil {
		return nil, errors.New("brain was already shut down")
	}

	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}

	// provide a stable result
	sort.Strings(keys)

	return keys, nil
}

// Close removes all data from the memory. Note that all calls to the memory
// will fail after this function has been called.
func (m *memory) Close() error {
	if m.data == nil {
		return errors.New("brain was already closed")
	}

	m.data = nil
	return nil
}

func (m *memory) persist() error {
	f, err := os.OpenFile(m.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return fmt.Errorf("failed to open file to persist data: %w", err)
	}

	err = json.NewEncoder(f).Encode(m.data)
	if err != nil {
		_ = f.Close()
		return fmt.Errorf("failed to encode data as JSON: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("failed to close file; data might not have been fully persisted to disk: %w", err)
	}

	return nil
}
