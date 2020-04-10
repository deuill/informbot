// Copyright 2017 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

// Package reader contains small, reusable APIs that build on the io.Reader
// interface.
package reader // import "mellium.im/reader"

import (
	"io"
	"net"
	"sync"
)

// Error returns a reader that always returns the given error for all calls to
// Read.
func Error(err error) io.Reader {
	return Func(func(p []byte) (int, error) {
		return 0, err
	})
}

// Before returns an io.Reader that proxies calls to Read and executes the given
// function exactly once before the first call.
// If the function errors, the error is returned and the call to Read is never
// proxied to the inner io.Reader (subsequent calls to Read will still be
// proxied).
// Because no call to Read returns until the one call to f returns, if f causes
// Read to be called, it will deadlock.
// If f panics, future calls of Read return without calling f.
// For more information see the documentation for sync.Once.
func Before(r io.Reader, f func() error) io.Reader {
	return &beforeReader{
		r:    r,
		f:    f,
		once: &sync.Once{},
	}
}

type beforeReader struct {
	r    io.Reader
	f    func() error
	once *sync.Once
}

func (br *beforeReader) Read(p []byte) (n int, err error) {
	br.once.Do(func() {
		err = br.f()
	})
	if err != nil {
		return
	}
	return br.r.Read(p)
}

func (br *beforeReader) Reset(r io.Reader) {
	br.r = r
	br.once = &sync.Once{}
}

// After returns an io.Reader that proxies to another Reader and calls f after
// each Read.
// The return value of f is returned from the call to r.Read instead of the
// original return value.
func After(r io.Reader, f func(n int, err error) (int, error)) io.Reader {
	return Func(func(p []byte) (n int, err error) {
		n, err = r.Read(p)
		if f == nil {
			return n, err
		}
		return f(n, err)
	})
}

// Conn replaces the Read method of c with r.Read.
// Generally, r wraps the Read method of c.
func Conn(c net.Conn, r io.Reader) net.Conn {
	return &conn{
		Conn: c,
		r:    r,
	}
}

type conn struct {
	net.Conn
	r io.Reader
}

func (c *conn) Read(p []byte) (n int, err error) {
	return c.r.Read(p)
}

// Func is an adapter to allow the use of ordinary functions as io.Readers.
// If f is a function with the appropriate signature, Func(f) is an io.Reader
// that calls f.
type Func func(p []byte) (n int, err error)

// Read calls f.
func (f Func) Read(p []byte) (n int, err error) {
	return f(p)
}
