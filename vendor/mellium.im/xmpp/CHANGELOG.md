# Changelog

All notable changes to this project will be documented in this file.


## v0.22.0 — 2024-09-23

### Security

- xmpp: fix incorrect stanza handling when different stanza type had same ID.
        A CVE has been issued with the id [CVE-2024-46957].

[CVE-2024-46957]: https://mellium.im/cve/cve-2024-46957/

### Added

- file: new implementation of [XEP-0446: File Metadata Element]
- stream: add handling of `see-other-host` error

[XEP-0446: File Metadata Element]: https://xmpp.org/extensions/xep-0446.html

### Breaking

- receipts: the API for requesting read receipts has been simplified
- stream: `Error.InnerXML()` has been replaced by `Error.Content`
- pubsub: changed condition names to better reflect their type

### Fixed

- dial: fix errors dialing connections without implicit TLS
- dial: TLS dialers now respect timeouts on context
- muc: fix a goroutine leak when joining or leaving a room and respect the
       context timeout when leaving the room
- stanza: when marshaling an error, all translations are now included
- stanza: when unmarshaling an error, the condition is now unmarshaled even if
          there are unknown child elements in the stanza
- stream: whitespace is now allowed around stream start elements during session
          negotiation
- stream: fix broken error unmarshaling


## v0.21.4 — 2023-01-11

### Added

- xmpp: add `UpdateAddr` to set the JID during resource binding
- pubsub: various namespaces and generated types
- version: add handler and mux option for software version requests
- blocklist: the ability to report users for spam or abuse


### Breaking

- ibr2: the deprecated ibr2 package has been removed, no replacement has been
  provided as it implemented an experimental version of the protocol that was
  not recommended for use


### Fixed

- xmpp: authentication no longer fails if the remote server offers channel
  binding
- xmpp: better error messages are returned when performing StartTLS
- blocklist: the blocklist functions now use the correct namespace and are
  compatible with most servers


## v0.21.3 — 2022-08-18

### Added

- ibb: the `Expect` method can be used to accept only pre-negotiated IBB
  sessions
- roster: the ability to set/delete roster items on entities other than the
  account using `SetIQ` and `DeleteIQ`
- upload: new package implementing file transfer by uploading to an HTTP server
- xmpp: support SCRAM authentication with channel binding over TLS 1.3 using
  [RFC 9266]


### Breaking

- compress: the compress package has been removed; users of this package can
  find a drop in replacement at [mellium.im/legacy/compress]
- ibb: the `Listener` type is now a pointer type


### Deprecated

- ibr2: deprecate legacy version of the extensible in-band registration protocol


### Fixed

- component: ensure that the stream namespace is set on the session so that the
  `mux` package and anything else matching on namespaces can be used with
  component connections
- internal/discover: a resource leak where XRD response bodies were unlimited in
  size and could waste bandwidth until the context was canceled even though the
  XRD files should never be large
- internal/discover: a resource leak where the XRD response body was not closed
  until the context timed out
- xmpp: a bug where resource binding would respond as if it were a
  client-to-server connection even on server-to-server streams

[RFC 9266]: https://datatracker.ietf.org/doc/html/rfc9266
[mellium.im/legacy/compress]: https://mellium.im/legacy/compress/


## v0.21.2 — 2022-04-07

### Added

- bookmarks: new package implementing [XEP-0402: PEP Native Bookmarks]
- crypto: new package implementing [XEP-0300: Use of Cryptographic Hash
  Functions in XMPP] and [XEP-0414: Cryptographic Hash Function Recommendations
  for XMPP]
- disco: add support for calculating and receiving [XEP-0115: Entity
  Capabilities] hashes
- mux: add ability to advertise features and identities that do not correspond
  to a handler
- xmpp: ability to create informational-only stream features


### Deprecated

- compress: the `compress` package has been deprecated and copied to the
  [`legacy`] module as [`legacy/compress`], it will be removed entirely in a
  future release


### Fixed

- websocket: rewrote WebSocket discovery to remove TXT record lookups and fix
  broken XRD file lookups
- pubsub: the IQ passed in by the user was not respected in `DeleteIQ`
- mux: fix a possible crash if a whitespace keepalive is encountered

[`legacy`]: https://mellium.im/legacy/
[`legacy/compress`]: https://mellium.im/legacy/compress/
[XEP-0115: Entity Capabilities]: https://xmpp.org/extensions/xep-0115.html
[XEP-0300: Use of Cryptographic Hash Functions in XMPP]: https://xmpp.org/extensions/xep-0300.html
[XEP-0402: PEP Native Bookmarks]: https://xmpp.org/extensions/xep-0402.html
[XEP-0414: Cryptographic Hash Function Recommendations for XMPP]: https://xmpp.org/extensions/xep-0414.html


## v0.21.1 — 2022-02-11

### Security

- websocket: fix an issue where the wrong hostname was validated in connections
             made after looking up DNS TXT records, resulting in a potential
             MITM. A CVE has been issued with the id [CVE-2022-24968].

[CVE-2022-24968]: https://mellium.im/cve/cve-2022-24968/


## v0.21.0 — 2022-02-08

### Breaking

- xmpp: the deprecated methods `InSID` and `OutSID` have been dropped in favor
  of `In` and `Out`


### Added:

- carbons: add `WrapReceived` and `WrapSent` functions
- forward: add `Unwrap` function
- ibb: new package implementing [XEP-0047: In-Band Bytestreams]
- internal/integration/slixmpp: [Slixmpp] support for integration tests
- pubsub: new package implementing parts of [XEP-0060: Publish-Subscribe] and
  [XEP-0163: Personal Eventing Protocol]
- stanza: add `NSError`, `NSClient`, and `NSServer` constants


### Fixed

- xmpp: IQ error responses are now sent to the original entity, not to
  ourselves
- mux: unhandled IQ result stanzas with no payload no longer result in an error
- mux: responses are no longer sent for unhandled IQ result stanzas

[XEP-0047: In-Band Bytestreams]: https://xmpp.org/extensions/xep-0047.html
[XEP-0060: Publish-Subscribe]: https://xmpp.org/extensions/xep-0060.html
[XEP-0163: Personal Eventing Protocol]: https://xmpp.org/extensions/xep-0163.html
[Slixmpp]: https://slixmpp.readthedocs.io/


## v0.20.0 — 2021-09-26

### Security

- all: update the default TLS config to disable TLS < 1.2

### Breaking

- disco: the `Feature`, `Identity`, and `Item` types have been moved to the
  `info` and `items` packages
- mux: make namespace checking stricter by adding argument to `New`
- roster: rename `version` attribute to `ver`
- roster: the `Push` callback now takes the roster version
- roster: `FetchIQ` now takes a `roster.IQ` instead of a `stanza.IQ` so that the
  version can be set
- styling: decoding tokens now uses an iterator pattern
- stanza: added stanza namespace argument to `Is`, `AddID`, and `AddOriginID`
- xmpp: the `WebSocket` option on `StreamConfig` has been removed in favor of
  `websocket.Negotiator`
- xmpp: the `IterIQ` and `IterIQElement` methods on `Session` now return the
  start element token associated with the IQ payload
- xmpp: `Negotiator` now takes a stream config function instead of a
  `StreamConfig` struct


### Added

- blocklist: new package implementing [XEP-0191: Blocking Command]
- carbons: new package implementing [XEP-0280: Message Carbons]
- commands: new package implementing [XEP-0050: Ad-Hoc Commands]
- history: implement [XEP-0313: Message Archive Management]
- form: add `Len` and `Raw` methods to form data
- muc: new package implementing [XEP-0045: Multi-User Chat] and [XEP-0249: Direct MUC Invitations]
- mux: `mux.ServeMux` now implements `info.FeatureIter`, `info.IdentityIter`,
  `form.Iter`, and `items.Iter`
- roster: the roster `Iter` now returns the roster version being iterated over
  from the `Version` method
- roster: if a `stanza.Error` is returned from the `Push` handler it is now sent
  in response to the original roster push.
- stanza: implement [XEP-0203: Delayed Delivery]
- stanza: more general `UnmarshalError` function that doesn't focus on IQs
- stanza: add `Error` method to `Presence` and `Message`
- stanza: errors now return text if present instead of condition when converting
  to a string
- websocket: add `Negotiator` to replace the `WebSocket` option on the stream
  config
- disco: new `disco.Handler` option to configure a `mux.ServeMux` to respond to
  disco#info requests
- disco/info: new package containing the former `disco.Feature` and
  `disco.Identity` types
- xmpp: add `In()` and `Out()` methods to `Session` for fetching stream info


### Fixed

- form: calling `Set` after unmarshaling into an uninitialized form no longer
  panics
- form: unmarshaling into an existing form now resets the stored values to
  prevent data leaks across forms
- paging: when unmarshaling sets, index is an attribute, not a child element
- roster: the roster version is now always included, even if empty, to signal
  that we support roster versioning
- roster: the handler now responds with an IQ result if the roster `Push`
  handler does not return an error
- stanza: unmarshaling error IQs now works even if the error is not the first
  child in the payload
- styling: pre-block start tokens with no newline had nonsensical formatting
- xmpp: empty IQ iters no longer return EOF when there is no payload
- xmpp: `UnmarshalIQ` and `UnmarshalIQElement` no longer return an XML error
  on responses with no payload


### Deprecated

- xmpp: the `InSID` and `OutSID` methods on `Session`


[XEP-0045: Multi-User Chat]: https://xmpp.org/extensions/xep-0045.html
[XEP-0050: Ad-Hoc Commands]: https://xmpp.org/extensions/xep-0050.html
[XEP-0191: Blocking Command]: https://xmpp.org/extensions/xep-0191.html
[XEP-0203: Delayed Delivery]: https://xmpp.org/extensions/xep-0203.html
[XEP-0249: Direct MUC Invitations]: https://xmpp.org/extensions/xep-0249.html
[XEP-0280: Message Carbons]: https://xmpp.org/extensions/xep-0280.html
[XEP-0313: Message Archive Management]: https://xmpp.org/extensions/xep-0313.html


## v0.19.0 — 2021-05-02

### Breaking

- color: change list of color vision deficiencies from uint8 to a new type
- disco: rename `GetItems` and `GetItemsIQ` to `FetchItems` and `FetchItemsIQ`
  to signal that they return an iterator
- xmpp: remove unused argument from `DialSession`
- jid: splitting a JID no longer trims trailing dots from the domainpart


### Added

- delay: new package implementing [XEP-0203: Delayed Delivery]
- disco: new package implementing [XEP-0030: Service Discovery]
- jid: normalization of domainparts for display purposes
- paging: new package implementing [XEP-0059: Result Set Management]
- stanza: new functions `AddID` and `AddOriginID` to support unique and stable
  stanza IDs
- stanza: ability to compare errors with `errors.Is`
- stanza: add `Is` function to check if an XMLName is a valid stanza name
- stanza: new `Wrap` method on `Error`
- stanza: new `UnmarshalIQError` function
- styling: satisfy `fmt.Stringer` for the `Style` type
- version: new package implementing [XEP-0092: Software Version]
- xmpp: satisfy `fmt.Stringer` for the `SessionState` type
- xmpp: new `UnmarshalIQ`, `UnmarshalIQElement`, `IterIQ`, and `IterIQElement`
  methods
- xtime: times can now be marshaled and unmarshaled as XML attributes


### Fixed

- form: if no field type is set the correct default (text-single) is used
- jid: JIDs created with `New` now trim trailing dots from the domainpart
- xmpp: unknown IQ error responses are now sent to the correct address
- xmpp: fixed DOS where reads/writes never timed out on `Dial*` functions


[XEP-0030: Service Discovery]: https://xmpp.org/extensions/xep-0030.html
[XEP-0059: Result Set Management]: https://xmpp.org/extensions/xep-0059.html
[XEP-0092: Software Version]: https://xmpp.org/extensions/xep-0092.html
[XEP-0203: Delayed Delivery]: https://xmpp.org/extensions/xep-0203.html


## v0.18.0 — 2021-02-14

### Breaking

- component: `NewClientSession` has been split in to `NewSession` and
  `ReceiveSession`
- mux: fixed the signature of the `IQFunc` option to differentiate it from the
  normal `IQ` option
- stream: the `SeeOtherHost` function signature has changed now that payload is
  available for all stream errors
- stream: the `TokenReader` method signature has changed to make it match the
  `xmlstream.Marshaler` interface
- stream: renamed `ErrorNS` to `NSError` for consistency
- xmpp: the `Negotiator` type signature has changed
- xmpp: the `StreamConfig.Features` field is now a callback instead of a slice
- xmpp: the `StreamFeature` type's `Parse` function now takes an xml.Decoder
  instead of an `xml.TokenReader`
- xmpp: IQs sent during stream negotiation are now matched against registered
  stream features based on the namespace of their payload
- xmpp: the session negotiation functions have been renamed and revamped, see
  the main package docs for a comprehensive overview
- xtime: change the `Time` type to fix unmarshaling


### Added

- internal/integration/mcabber: [Mcabber] support for integration tests
- oob: implementations of `xmlstream.Marshaler` and `xmlstream.WriterTo` for the
  types `IQ`, `Query`, and Data
- receipts: add `Request` and `Requested` to add receipt requests to messages
  without waiting on a response
- roster: add `Set` and `Delete` functions for roster management
- roster: support multiple groups
- s2s: added a new package implementing Bidi
- stanza: add `Error` method on the `IQ` type
- stream: new `InnerXML` and `ApplicationError` methods on `Error` provide a way
  to easily construct customized stream errors
- stream: ability to compare errors with `errors.Is`
- stream: support adding human-readable text to errors
- stream: add `Info` block for use by custom `Negotiators`
- styling: add `Disable` and `Unstyled` to support disabling styling on some
  messages
- websocket: added a new package for dialing WebSocket connections
- xmpp: `SetReadDeadline` and `SetWriteDeadline` are now proxied even if the
  underlying connection is not a `net.Conn`
- xmpp: all sent stanzas are now given randomly generated IDs if no ID was
  provided (not just IQs)
- xmpp: the start token that triggers a call to `Negotiate` is now no longer
  popped from the stream before the actual call to `Negotiate` for server side
  stream features
- xmpp: a `SASLServer` stream feature for authenticating users
- xmpp: two methods for getting the stream IDs, `InSID` and `OutSID`
- xmpp: a new state bit, `s2s` to indicate whether the session is a
  server-to-server connection

[Mcabber]: https://mcabber.com/


### Changed

- dial: resolving SRV records for XMPP over implicit or opportunistic TLS is
  now done concurrently to make the initial connection faster
- dial: the fallback for dialing an XMPP connection when no SRV records exist is
  now more robust
- xmpp: stanzas sent over S2S connections now always have the "from" attribute
  set
- xmpp: the default negotiator now supports negotiating the WebSocket
  subprotocol when the `WebSocket` option is set


### Fixed

- component: the `NewSession` function (previously `NewClientSession`) now
  correctly marks the connection as having been received or initiated
- component: the wrong error is no longer returned if Prosody sends an error
  immediately after the start of a stream with no ID
- dial: dialing an XMPP connection where no xmpps SRV records exist no longer
  results in an error (fallback works correctly)
- roster: fix infinite loop marshaling lists of items
- stream: errors are now unmarshaled correctly
- xmpp: the Encode methods no longer sometimes duplicate the xmlns attribute
- xmpp: stream errors are now unmarshaled and returned from `Serve` and during
  session negotiation
- xmpp: XML tokens written directly to the session are now always flushed to the
  underlying connection when the token writer is closed
- xmpp: stream feature parse functions can no longer read beyond the end of the
  feature element
- xmpp: the resource binding feature was previously broken for server sessions.
- xmpp: the "to" and "from" attributes on incoming streams are now verified and
  an error is returned if they change between stream restarts
- xmpp: empty "from" attributes are now removed before marshaling so that a
  default value can be set (if applicable)


## v0.17.1 — 2020-11-21

### Breaking

- roster: remove workaround for a bug in Go versions prior to 1.14 which is now
  the earliest supported version
- xmpp: the `Encode` and `EncodeElement` methods now take a context and respect
  its deadline

### Added

- internal/integration: new package for writing integration tests
- internal/integration/ejabberd: [Ejabberd] support for integration tests
- internal/integration/prosody: [Prosŏdy] support for integration tests
- internal/xmpptest: new `ClientServer` for testing two connected sessions
- xmpp: new `EncodeIQ` and `EncodeIQElement` methods

[Ejabberd]: https://www.ejabberd.im/
[Prosŏdy]: https://prosody.im/


### Fixed

- stanza: converting stanzas with empty to/from attributes no longer fails
- xmpp: fixed data race that could result in invalid session state and lead to
  writes on a closed session and other state related issues
- xmpp: the `Send` and `SendElement` methods now respect the context deadline


## v0.17.0 — 2020-11-11

### Breaking

- sasl2: removed experimental package
- xmpp: removed option to make STARTTLS feature optional


### Added

- styling: new package implementing [XEP-0393: Message Styling]
- xmpp: `ConnectionState` method


### Fixed

- roster: iters that have errored no longer panic when closed
- xmpp: using TeeIn/TeeOut no longer breaks SCRAM based SASL mechanisms
- xmpp: stream negotiation no longer fails when the only required features
  cannot yet be negotiated because they depend on optional features


[XEP-0393: Message Styling]: https://xmpp.org/extensions/xep-0393.html


## v0.16.0 — 2020-03-08

### Breaking

- xmpp: the end element is now included in the token stream passed to handlers
- xmpp: SendElement now wraps the entire stream, not just the first element


### Added

- receipts: new package implementing [XEP-0333: Chat Markers]
- roster: add handler and mux option for roster pushes


[XEP-0333: Chat Markers]: https://xmpp.org/extensions/xep-0333.html


### Fixed

- mux: fix broken `Decode` and possible infinite loop due to cutting off the
  last token in a buffered XML token stream
- roster: work around a bug in Go 1.13 where `io.EOF` may be returned from the
  XML decoder


## v0.15.0 — 2020-02-28

### Breaking

- all: dropped support for versions of Go before 1.13
- mux: move `Wrap{IQ,Presence,Message}` functions to methods on the stanza types


### Added

- mux: ability to select handlers by stanza payload
- mux: new handler types and API
- ping: a function for easily encoding pings and handling errors
- ping: a handler and mux option for responding to pings
- stanza: ability to convert stanzas to/from `xml.StartElement`
- stanza: API to simplify replying to IQs
- uri: new package for parsing XMPP URI's and IRI's
- xtime: new package for handling [XEP-0202: Entity Time] and [XEP-0082: XMPP Date and Time Profiles]


[XEP-0202: Entity Time]: https://xmpp.org/extensions/xep-0202.html
[XEP-0082: XMPP Date and Time Profiles]: https://xmpp.org/extensions/xep-0082.html


### Fixed

- dial: if a port number is present in a JID it was previously ignored


## v0.14.0 — 2019-08-18

### Breaking

- ping: remove `IQ` function and replace with struct based API


### Added

- ping: add `IQ` struct based encoding API


### Changed

- stanza: a zero value `IQType` now marshals as "get"
- xmpp: read timeouts are now returned instead of ignored


### Fixed

- dial: fix broken fallback to domainpart
- xmpp: allow whitespace keepalives
- roster: the iterator now correctly closes the underlying TokenReadCloser
- xmpp: fix bug where stream processing could stop after an IQ was received


## v0.13.0 — 2019-07-27

### Breaking

- xmpp: change `Handler` to take an `xmlstream.TokenReadEncoder`
- xmpp: replace `EncodeToken` and `Flush` with `TokenWriter`
- xmpp: replace `Token` with `TokenReader`


### Added

- examples/echobot: add graceful shutdown on SIGINT
- xmpp: `Encode` and `EncodeElement` methods


### Changed

- xmpp: calls to `Serve` no longer return `io.EOF` on success


### Fixed

- examples/echobot: calling `Send` from within the handler resulted in deadlock
- xmpp: closing the input stream was racy, resulting in invalid XML


## v0.12.0

### Breaking

- dial: moved network dialing types and functions into new package.
- dial: use underlying net.Dialer's DNS Resolver in Dialer.
- stanza: change API of `WrapIQ` and `WrapPresence` to not abuse pointers
- xmpp: add new `SendIQ` API and remove response from `Send` and `SendElement`
- xmpp: new API for writing custom tokens to a session

### Fixed

- xmpp: let `Session.Close` operate concurrently with `SendElement` et al.
