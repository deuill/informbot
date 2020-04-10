<h1 align="center">Joe Bot - File Storage</h1>
<p align="center">Basic file storage memory adapater. https://github.com/go-joe/joe</p>
<p align="center">
	<a href="https://circleci.com/gh/go-joe/file-memory/tree/master"><img src="https://circleci.com/gh/go-joe/file-memory/tree/master.svg?style=shield"></a>
	<a href="https://goreportcard.com/report/github.com/go-joe/file-memory"><img src="https://goreportcard.com/badge/github.com/go-joe/file-memory"></a>
	<a href="https://codecov.io/gh/go-joe/file-memory"><img src="https://codecov.io/gh/go-joe/file-memory/branch/master/graph/badge.svg"/></a>
	<a href="https://pkg.go.dev/github.com/go-joe/file-memory?tab=doc"><img src="https://img.shields.io/badge/godoc-reference-blue.svg?color=blue"></a>
	<a href="https://github.com/go-joe/file-memory/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-BSD--3--Clause-blue.svg"></a>
</p>

---

This repository contains a module for the [Joe Bot library][joe].

## Getting Started

This library is packaged as [Go module][go-modules]. You can get it via:

```
go get github.com/go-joe/file-memory
```

## Example usage

```go
b := &ExampleBot{
	Bot: joe.New("example", file.Memory("foobar.json")),
}
```

## Built With

* [testify](https://github.com/stretchr/testify) - A simple unit test library
* [zap](https://github.com/uber-go/zap) - Blazing fast, structured, leveled logging in Go

## Contributing

If you want to hack on this repository, please read the short [CONTRIBUTING.md](CONTRIBUTING.md)
guide first.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available,
see the [tags on this repository][tags]. 

## Authors

- **Friedrich Große** - *Initial work* - [fgrosse](https://github.com/fgrosse)
- **Stefan Warman** - *Unit tests* - [warmans](https://github.com/warmans)

See also the list of [contributors][contributors] who participated in this project.

## License

This project is licensed under the BSD-3-Clause License - see the [LICENSE](LICENSE) file for details.

[joe]: https://github.com/go-joe/joe
[go-modules]: https://github.com/golang/go/wiki/Modules
[tags]: https://github.com/go-joe/file-memory/tags
[contributors]: https://github.com/go-joe/file-memory/contributors
