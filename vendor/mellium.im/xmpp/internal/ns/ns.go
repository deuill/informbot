// Copyright 2016 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

// Package ns provides namespace constants that are used by the xmpp package and
// other internal packages.
package ns // import "mellium.im/xmpp/internal/ns"

// List of commonly used namespaces.
const (
	Bind     = "urn:ietf:params:xml:ns:xmpp-bind"
	Client   = "jabber:client"
	SASL     = "urn:ietf:params:xml:ns:xmpp-sasl"
	Server   = "jabber:server"
	Stanza   = "urn:ietf:params:xml:ns:xmpp-stanzas"
	StartTLS = "urn:ietf:params:xml:ns:xmpp-tls"
	WS       = "urn:ietf:params:xml:ns:xmpp-framing"
	XML      = "http://www.w3.org/XML/1998/namespace"
)
