package http

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"net/url"
	"syscall"
	"testing"
	"time"
)

func TestValidateURLInactivePolicyAllowsHistoricalURLs(t *testing.T) {
	p := FetchPolicy{}
	for _, raw := range []string{
		"http://127.0.0.1/a.png",
		"http://localhost/a.png",
		"https://images.example.com/a.png",
		"file:///etc/passwd",
	} {
		if err := p.ValidateURL(mustPolicyURL(t, raw)); err != nil {
			t.Fatalf("inactive policy rejected %s: %v", raw, err)
		}
	}
}

func TestValidateURL(t *testing.T) {
	p := NewFetchPolicy(true, []string{
		"CDN.Example.com.",
		"https://img.example.org/ignored",
		"[::1]",
	})

	allowed := []string{
		"https://cdn.example.com/a.png",
		"http://CDN.Example.com./b.png",
		"https://img.example.org/c.png",
	}
	for _, raw := range allowed {
		if err := p.ValidateURL(mustPolicyURL(t, raw)); err != nil {
			t.Fatalf("allowed %s: %v", raw, err)
		}
	}

	denied := []string{
		"http://evil.example/a.png",
		"http://127.0.0.1/a.png",
		"http://[::1]/a.png",
		"http://[::ffff:127.0.0.1]/a.png",
		"http://169.254.169.254/latest/meta-data/",
		"http://10.1.2.3/a.png",
		"http://100.100.100.200/a.png",
		"file:///etc/passwd",
		"gopher://cdn.example.com/a",
		"http:///a.png",
	}
	for _, raw := range denied {
		err := p.ValidateURL(mustPolicyURL(t, raw))
		if !errors.Is(err, ErrURLNotAllowed) {
			t.Fatalf("expected %s to be denied, got %v", raw, err)
		}
	}
}

func TestAllowlistDoesNotBlockListedPublicHost(t *testing.T) {
	p := NewFetchPolicy(false, []string{"images.example.com:443"})
	if err := p.ValidateURL(mustPolicyURL(t, "https://images.example.com/a.png")); err != nil {
		t.Fatal(err)
	}
	if err := p.ValidateURL(mustPolicyURL(t, "https://other.example.com/a.png")); !errors.Is(err, ErrURLNotAllowed) {
		t.Fatalf("got %v", err)
	}
	// Without the private-network block, a listed loopback literal is still a host match.
	loopback := NewFetchPolicy(false, []string{"127.0.0.1"})
	if err := loopback.ValidateURL(mustPolicyURL(t, "http://127.0.0.1/a.png")); err != nil {
		t.Fatal(err)
	}
}

func TestBlockedAddresses(t *testing.T) {
	p := NewFetchPolicy(true, nil)
	blocked := []string{
		"127.0.0.1",
		"127.255.255.255",
		"10.0.0.1",
		"10.255.255.255",
		"172.16.0.0",
		"172.31.255.255",
		"192.168.0.1",
		"169.254.169.254",
		"169.254.1.1",
		"0.0.0.0",
		"100.100.100.200",
		"::1",
		"::",
		"fc00::1",
		"fd00:ec2::254",
		"fe80::1",
		"::ffff:127.0.0.1",
		"::ffff:10.0.0.1",
		"::ffff:169.254.169.254",
		ipv4Embedded("2002", netip.MustParseAddr("127.0.0.1")).String(),
		ipv4Embedded("2002", netip.MustParseAddr("10.1.2.3")).String(),
		teredo(netip.MustParseAddr("192.168.1.1")).String(),
		nat64(netip.MustParseAddr("127.0.0.1")).String(),
		"::7f00:1",
	}
	allowed := []string{
		"8.8.8.8",
		"1.1.1.1",
		"172.15.255.255",
		"172.32.0.0",
		"11.0.0.1",
		"100.100.100.201",
		"::ffff:8.8.8.8",
		ipv4Embedded("2002", netip.MustParseAddr("8.8.8.8")).String(),
		teredo(netip.MustParseAddr("8.8.8.8")).String(),
		nat64(netip.MustParseAddr("1.1.1.1")).String(),
		"2001:db8::1",
	}

	for _, addr := range blocked {
		if err := p.control("tcp", net.JoinHostPort(addr, "80"), nil); !errors.Is(err, ErrURLNotAllowed) {
			t.Errorf("address %s: got %v, want blocked", addr, err)
		}
	}
	for _, addr := range allowed {
		if err := p.control("tcp", net.JoinHostPort(addr, "443"), nil); err != nil {
			t.Errorf("address %s: got %v, want allowed", addr, err)
		}
	}
}

func TestControlRunsAfterDNSResolution(t *testing.T) {
	p := NewFetchPolicy(true, nil)
	d := p.dialer()
	d.Timeout = time.Second

	var seen []string
	orig := d.Control
	d.Control = func(network, address string, c syscall.RawConn) error {
		seen = append(seen, address)
		return orig(network, address, c)
	}

	_, err := d.Dial("tcp", "localhost:80")
	if !errors.Is(err, ErrURLNotAllowed) {
		t.Fatalf("localhost dial error = %v, want blocked address", err)
	}
	if len(seen) == 0 {
		t.Fatal("dialer Control was not called")
	}
	for _, address := range seen {
		host, _, splitErr := net.SplitHostPort(address)
		if splitErr != nil {
			t.Fatalf("control address %q is not host:port", address)
		}
		ip, parseErr := netip.ParseAddr(host)
		if parseErr != nil || !ip.Unmap().IsLoopback() {
			t.Fatalf("control saw %q before the name was resolved to a loopback address", address)
		}
	}
}

func TestClientBlocksLoopbackAndAllowsItWhenDisabled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}))
	t.Cleanup(srv.Close)

	plain := NewClient()
	resp, err := plain.Do(mustPolicyRequest(t, srv.URL))
	if err != nil {
		t.Fatalf("default client: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("default client status %d", resp.StatusCode)
	}

	blocked := NewClient(WithFetchPolicy(NewFetchPolicy(true, nil)))
	_, err = blocked.Do(mustPolicyRequest(t, srv.URL))
	if !errors.Is(err, ErrURLNotAllowed) {
		t.Fatalf("blocked client error = %v", err)
	}

	// A listed loopback host is still refused when the address block is on.
	both := NewClient(WithFetchPolicy(NewFetchPolicy(true, []string{"127.0.0.1"})))
	_, err = both.Do(mustPolicyRequest(t, srv.URL))
	if !errors.Is(err, ErrURLNotAllowed) {
		t.Fatalf("allow list did not keep the dialer block, error = %v", err)
	}

	listed := NewClient(WithFetchPolicy(NewFetchPolicy(false, []string{"127.0.0.1"})))
	resp, err = listed.Do(mustPolicyRequest(t, srv.URL))
	if err != nil {
		t.Fatalf("allow list without the address block: %v", err)
	}
	resp.Body.Close()
}

func TestClientRejectsRedirectOutsideAllowlist(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://example.invalid/secret", http.StatusFound)
	}))
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	client := NewClient(WithFetchPolicy(NewFetchPolicy(false, []string{u.Hostname()})))
	_, err = client.Do(mustPolicyRequest(t, srv.URL+"/a.png"))
	if !errors.Is(err, ErrURLNotAllowed) {
		t.Fatalf("redirect error = %v, want allow list rejection", err)
	}
}

func TestClientBlocksLocalhostName(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}))
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	u.Host = net.JoinHostPort("localhost", u.Port())

	client := NewClient(WithFetchPolicy(NewFetchPolicy(true, nil)))
	_, err = client.Do(mustPolicyRequest(t, u.String()))
	if !errors.Is(err, ErrURLNotAllowed) {
		t.Fatalf("localhost fetch error = %v", err)
	}
}

func mustPolicyURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func mustPolicyRequest(t *testing.T, raw string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, raw, nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

// ipv4Embedded builds 2002:AABB:CCDD::/48 (prefix "2002") from an IPv4 address.
func ipv4Embedded(prefix string, v4 netip.Addr) netip.Addr {
	b := [16]byte{0x20, 0x02}
	if prefix != "2002" {
		panic(prefix)
	}
	raw := v4.As4()
	copy(b[2:6], raw[:])
	return netip.AddrFrom16(b)
}

func teredo(v4 netip.Addr) netip.Addr {
	raw := v4.As4()
	b := [16]byte{0x20, 0x01}
	b[12] = raw[0] ^ 0xff
	b[13] = raw[1] ^ 0xff
	b[14] = raw[2] ^ 0xff
	b[15] = raw[3] ^ 0xff
	return netip.AddrFrom16(b)
}

func nat64(v4 netip.Addr) netip.Addr {
	b := [16]byte{}
	copy(b[:12], nat64WellKnown[:])
	raw := v4.As4()
	copy(b[12:], raw[:])
	return netip.AddrFrom16(b)
}
