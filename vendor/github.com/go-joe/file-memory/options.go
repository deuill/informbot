// Package file implements file based memory for the Joe bot library.
// https://github.com/go-joe/joe
package file

import "go.uber.org/zap"

// Option corresponds to a configuration setting of the file memory.
// All available options are the exported functions of this package that share
// the prefix "With" in their names.
type Option func(*memory) error

// WithLogger is a memory option that allows the caller to set a different
// logger. By default this option is not required because the file.Memory(…)
// function automatically uses the logger of the given joe.Config.
func WithLogger(logger *zap.Logger) Option {
	return func(memory *memory) error {
		memory.logger = logger
		return nil
	}
}

// IDEA: encrypted brain?
// IDEA: only decrypt keys on demand?
