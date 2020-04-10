package joe

import (
	"encoding/json"
	"sort"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// A Storage provides a convenient interface to a Memory implementation. It is
// responsible for how the actual key value data is encoded and provides
// concurrent access as well as logging.
//
// The default Storage that is returned by joe.NewStorage() encodes values as
// JSON and stores them in-memory.
type Storage struct {
	logger  *zap.Logger
	mu      sync.RWMutex
	memory  Memory
	encoder MemoryEncoder
}

// The Memory interface allows the bot to persist data as key-value pairs.
// The default implementation of the Memory is to store all keys and values in
// a map (i.e. in-memory). Other implementations typically offer actual long term
// persistence into a file or to redis.
type Memory interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, bool, error)
	Delete(key string) (bool, error)
	Keys() ([]string, error)
	Close() error
}

// A MemoryEncoder is used to encode and decode any values that are stored in
// the Memory. The default implementation that is used by the Storage uses a
// JSON encoding.
type MemoryEncoder interface {
	Encode(value interface{}) ([]byte, error)
	Decode(data []byte, target interface{}) error
}

type inMemory struct {
	data map[string][]byte
}

type jsonEncoder struct{}

// NewStorage creates a new Storage instance that encodes values as JSON and
// stores them in-memory. You can change the memory and encoding via the
// provided setters.
func NewStorage(logger *zap.Logger) *Storage {
	return &Storage{
		logger:  logger,
		memory:  newInMemory(),
		encoder: new(jsonEncoder),
	}
}

// SetMemory assigns a different Memory implementation.
func (s *Storage) SetMemory(m Memory) {
	s.mu.Lock()
	s.memory = m
	s.mu.Unlock()
}

// SetMemoryEncoder assigns a different MemoryEncoder.
func (s *Storage) SetMemoryEncoder(enc MemoryEncoder) {
	s.mu.Lock()
	s.encoder = enc
	s.mu.Unlock()
}

// Keys returns all keys known to the Memory.
func (s *Storage) Keys() ([]string, error) {
	s.mu.RLock()
	keys, err := s.memory.Keys()
	s.mu.RUnlock()

	sort.Strings(keys)
	return keys, err
}

// Set encodes the given data and stores it in the Memory that is managed by the
// Storage.
func (s *Storage) Set(key string, value interface{}) error {
	data, err := s.encoder.Encode(value)
	if err != nil {
		return errors.Wrap(err, "encode data")
	}

	s.mu.Lock()
	s.logger.Debug("Writing data to memory", zap.String("key", key))
	err = s.memory.Set(key, data)
	s.mu.Unlock()

	return err
}

// Get retrieves the value under the requested key and decodes it into the
// passed "value" argument which must be a pointer. The boolean return value
// indicates if the value actually existed in the Memory and is false if it did
// not. It is legal to pass <nil> as the value if you only want to check if
// the given key exists but you do not actually care about the concrete value.
func (s *Storage) Get(key string, value interface{}) (bool, error) {
	s.mu.RLock()
	s.logger.Debug("Retrieving data from memory", zap.String("key", key))
	data, ok, err := s.memory.Get(key)
	s.mu.RUnlock()
	if err != nil {
		return false, errors.WithStack(err)
	}

	if !ok || value == nil {
		return ok, nil
	}

	err = s.encoder.Decode(data, value)
	if err != nil {
		return false, errors.Wrap(err, "decode data")
	}

	return true, nil
}

// Delete removes a key and its associated value from the memory. The boolean
// return value indicates if the key existed or not.
func (s *Storage) Delete(key string) (bool, error) {
	s.mu.Lock()
	s.logger.Debug("Deleting data from memory", zap.String("key", key))
	ok, err := s.memory.Delete(key)
	s.mu.Unlock()

	return ok, err
}

// Close closes the Memory that is managed by this Storage.
func (s *Storage) Close() error {
	s.mu.Lock()
	err := s.memory.Close()
	s.mu.Unlock()
	return err
}

func newInMemory() *inMemory {
	return &inMemory{data: map[string][]byte{}}
}

func (m *inMemory) Set(key string, value []byte) error {
	m.data[key] = value
	return nil
}

func (m *inMemory) Get(key string) ([]byte, bool, error) {
	value, ok := m.data[key]
	return value, ok, nil
}

func (m *inMemory) Delete(key string) (bool, error) {
	_, ok := m.data[key]
	delete(m.data, key)
	return ok, nil
}

func (m *inMemory) Keys() ([]string, error) {
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}

	return keys, nil
}

func (m *inMemory) Close() error {
	m.data = map[string][]byte{}
	return nil
}

func (jsonEncoder) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (jsonEncoder) Decode(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}
