package http

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"syscall"
	"time"
)

// ErrURLNotAllowed is returned when a user-supplied URL is rejected by the
// fetch policy. It is wrapped around the specific reason.
var ErrURLNotAllowed = errors.New("url is not allowed")

// alibabaMetadata is the Alibaba Cloud metadata service. It is not inside a
// link-local or RFC1918 range, so the class checks below do not cover it.
var alibabaMetadata = netip.MustParseAddr("100.100.100.200")

// nat64WellKnown is the 96-bit prefix 64:ff9b::/96 from RFC 6052.
var nat64WellKnown = [12]byte{0x00, 0x64, 0xff, 0x9b}

// FetchPolicy restricts outbound requests made for a user-supplied url.
// A zero policy allows every URL, which is the historical behavior.
type FetchPolicy struct {
	// BlockPrivateNetworks refuses connections to loopback, private,
	// link-local, unspecified, and well-known cloud metadata addresses.
	// The decision is made in the dialer after DNS resolution.
	BlockPrivateNetworks bool

	allowed map[string]struct{}
}

// NewFetchPolicy builds a policy. hosts is the optional allow list; an empty
// list does not restrict hosts. Entries are hostnames or IP addresses.
// Schemes, paths, ports, and a trailing dot are ignored. Matching is
// case-insensitive.
func NewFetchPolicy(blockPrivateNetworks bool, hosts []string) FetchPolicy {
	p := FetchPolicy{BlockPrivateNetworks: blockPrivateNetworks}
	if len(hosts) == 0 {
		return p
	}
	allowed := make(map[string]struct{}, len(hosts))
	for _, host := range hosts {
		host = normalizeHost(host)
		if host == "" {
			continue
		}
		allowed[host] = struct{}{}
	}
	if len(allowed) > 0 {
		p.allowed = allowed
	}
	return p
}

// Active reports whether the policy rejects any URL that the historical
// client would have fetched.
func (p FetchPolicy) Active() bool {
	return p.BlockPrivateNetworks || len(p.allowed) > 0
}

// ValidateURL checks the scheme, the allow list, and IP literals.
// Hostnames are not resolved here; the dialer checks those after DNS so a
// rebinding lookup cannot skip the address check.
func (p FetchPolicy) ValidateURL(u *url.URL) error {
	if !p.Active() {
		return nil
	}
	if u == nil {
		return fmt.Errorf("%w: empty url", ErrURLNotAllowed)
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
	default:
		return fmt.Errorf("%w: scheme %q is not permitted", ErrURLNotAllowed, u.Scheme)
	}
	host := normalizeHost(u.Hostname())
	if host == "" {
		return fmt.Errorf("%w: missing host", ErrURLNotAllowed)
	}
	if !p.hostAllowed(host) {
		return fmt.Errorf("%w: host %q is not allowed", ErrURLNotAllowed, host)
	}
	if p.BlockPrivateNetworks {
		if ip, err := netip.ParseAddr(host); err == nil && isBlockedAddr(ip) {
			return fmt.Errorf("%w: host %q is a blocked address", ErrURLNotAllowed, host)
		}
	}
	return nil
}

func (p FetchPolicy) hostAllowed(host string) bool {
	if len(p.allowed) == 0 {
		return true
	}
	_, ok := p.allowed[normalizeHost(host)]
	return ok
}

func normalizeHost(host string) string {
	host = strings.TrimSpace(host)
	if host == "" {
		return ""
	}
	if strings.Contains(host, "://") {
		if u, err := url.Parse(host); err == nil {
			if name := u.Hostname(); name != "" {
				host = name
			}
		}
	}
	if strings.HasPrefix(host, "[") {
		if name, _, err := net.SplitHostPort(host); err == nil {
			host = name
		} else {
			host = strings.Trim(host, "[]")
		}
	} else if name, _, err := net.SplitHostPort(host); err == nil {
		host = name
	}
	host = strings.TrimSuffix(host, ".")
	return strings.ToLower(host)
}

// control is net.Dialer.Control. The standard library calls it with the
// resolved remote address after DNS resolution and before connect, including
// every address of a multi-address dial and every redirect.
func (p FetchPolicy) control(network, address string, _ syscall.RawConn) error {
	if !p.BlockPrivateNetworks {
		return nil
	}
	ip, err := parseDialAddr(address)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrURLNotAllowed, err.Error())
	}
	if isBlockedAddr(ip) {
		return fmt.Errorf("%w: refusing to connect to %s", ErrURLNotAllowed, ip)
	}
	return nil
}

func (p FetchPolicy) dialer() *net.Dialer {
	return &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
		Control:   p.control,
	}
}

func (p FetchPolicy) transport() http.RoundTripper {
	var tr *http.Transport
	if base, ok := http.DefaultTransport.(*http.Transport); ok {
		tr = base.Clone()
	} else {
		tr = &http.Transport{}
	}
	tr.DialContext = p.dialer().DialContext
	// A proxy would connect to the target itself. The dialer would only see
	// the proxy address, so environment proxies are not used for these fetches.
	tr.Proxy = nil
	return tr
}

func (p FetchPolicy) checkRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return fmt.Errorf("%w: stopped after 10 redirects", ErrURLNotAllowed)
	}
	return p.ValidateURL(req.URL)
}

func parseDialAddr(address string) (netip.Addr, error) {
	host := address
	if h, _, err := net.SplitHostPort(address); err == nil {
		host = h
	}
	host = strings.TrimPrefix(host, "[")
	host = strings.TrimSuffix(host, "]")
	ip, err := netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("unresolved dial address %q", address)
	}
	return ip, nil
}

func isBlockedAddr(ip netip.Addr) bool {
	if !ip.IsValid() {
		return true
	}
	ip = ip.Unmap()
	if blockedClass(ip) {
		return true
	}
	if v4, ok := embeddedIPv4(ip); ok && blockedClass(v4) {
		return true
	}
	return false
}

func blockedClass(ip netip.Addr) bool {
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
		return true
	}
	return ip == alibabaMetadata
}

// embeddedIPv4 returns an IPv4 address carried inside an IPv6 address.
// IPv4-mapped addresses are handled by Addr.Unmap before this is called.
func embeddedIPv4(ip netip.Addr) (netip.Addr, bool) {
	if !ip.Is6() {
		return netip.Addr{}, false
	}
	b := ip.As16()
	switch {
	case b[0] == 0x20 && b[1] == 0x02:
		// 6to4 (RFC 3056): 2002:AABB:CCDD::/48.
		return netip.AddrFrom4(fourBytes(b[2:6])), true
	case b[0] == 0x20 && b[1] == 0x01 && b[2] == 0 && b[3] == 0:
		// Teredo (RFC 4380): the client IPv4 is the last 4 bytes XOR 0xff.
		return netip.AddrFrom4([4]byte{b[12] ^ 0xff, b[13] ^ 0xff, b[14] ^ 0xff, b[15] ^ 0xff}), true
	case [12]byte(b[:12]) == nat64WellKnown:
		// Well-known NAT64 prefix 64:ff9b::/96 (RFC 6052).
		return netip.AddrFrom4(fourBytes(b[12:16])), true
	case b[0]|b[1]|b[2]|b[3]|b[4]|b[5]|b[6]|b[7]|b[8]|b[9]|b[10]|b[11] == 0:
		// Deprecated IPv4-compatible form ::AABB:CCDD. :: and ::1 are addresses
		// in their own right and are classified above.
		if b[15] == 1 && b[12]|b[13]|b[14] == 0 {
			return netip.Addr{}, false
		}
		v4 := netip.AddrFrom4(fourBytes(b[12:16]))
		if v4.IsUnspecified() {
			return netip.Addr{}, false
		}
		return v4, true
	default:
		return netip.Addr{}, false
	}
}

func fourBytes(b []byte) [4]byte {
	var out [4]byte
	copy(out[:], b)
	return out
}
