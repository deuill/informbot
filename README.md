# InformBot - A Chat Bot for Inform 7 Stories, Built in Go

[![API Documentation][godoc-svg]][godoc-url] [![MIT License][license-svg]][license-url]

This package contains chat-bot for executing and running Inform 7 stories, and built around the [Joe
Bot][joe-url] framework. Any supported chat adapter can be used, and this repository contains an
example integration against a built-in XMPP adapter (which is not currently part of main-line
support).

## Building and Running

You can build the default test deployment of `informbot` (which currently works against an XMPP
server) by running `go get`, e.g.:

```go
go get github.com/deuill/informbot
```

Depending on your server setup, you may need to set up a number of options, as environment
variables, for instance:

``` go
INFORMBOT_JID="informbot@test.com" INFORMBOT_PASSWORD="123" INFORMBOT_USE_STARTTLS=true INFORMBOT_NO_TLS=true informbot
```

## Status

This package is still in early development, and is neither feature-complete nor bug-free. A large
amount of work remains on improving integration with Inform and Frotz, and optimizing against larger
stories.

## License

All code in this repository is covered by the terms of the MIT License, the full text of which can
be found in the LICENSE file.


[joe-url]: https://github.com/go-joe/joe

[godoc-url]: https://godoc.org/github.com/deuill/informbot
[godoc-svg]: https://godoc.org/github.com/deuill/informbot?status.svg

[license-url]: https://github.com/deuill/informbot/blob/master/LICENSE
[license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
