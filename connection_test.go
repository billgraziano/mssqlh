package mssqlh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedacted(t *testing.T) {
	assert := assert.New(t)
	conn := Connection{
		User:     "u1",
		Password: "pass",
	}
	assert.Equal("sqlserver://u1:redacted@localhost", conn.Redacted(0))
	assert.Equal("sqlserver://u1:pass@localhost", conn.Redacted(5))
	assert.Equal("sqlserver://u1:pass@localhost", conn.Redacted(4))
	assert.Equal("sqlserver://u1:pas_redacted@localhost", conn.Redacted(3))
	assert.Equal("sqlserver://u1:pa_redacted@localhost", conn.Redacted(2))
	assert.Equal("sqlserver://u1:p_redacted@localhost", conn.Redacted(1))
	assert.Equal("pass", conn.Password)
}

func TestSetInstance(t *testing.T) {
	assert := assert.New(t)
	var c0 Connection
	c0.FQDN = "host"
	assert.Equal("host", c0.Computer(), "host only")
	assert.Equal("", c0.Instance(), "host only")

	var c1 Connection
	c1.FQDN = "host\\instance"
	assert.Equal("host", c1.Computer(), "host instance")
	assert.Equal("instance", c1.Instance(), "host instance")

	var c2 Connection
	c2.FQDN = "tcp:host\\instance"
	assert.Equal("host", c2.Computer(), "host instance")
	assert.Equal("instance", c2.Instance(), "host instance")
}

func TestNameSplitter(t *testing.T) {
	assert := assert.New(t)
	var tests = []struct {
		name     string
		in       string
		host     string
		inst     string
		port     int
		protocol string
	}{
		{"host", "host", "host", "", 0, ""},
		{"host", "host-name", "host-name", "", 0, ""},
		{"np:host", "np:host-name", "host-name", "", 0, "np"},
		{"host", "host_name", "host_name", "", 0, ""},
		{"host", "host_name", "host_name", "", 0, ""},
		{"fqdb", "host.domain.name", "host.domain.name", "", 0, ""},
		{"tcp:fqdb", "tcp:host.domain.name", "host.domain.name", "", 0, "tcp"},
		{"instance", "host\\instance", "host", "instance", 0, ""},
		{"lpc:instance", "lpc:host\\instance", "host", "instance", 0, "lpc"},
		{"instance", "host\\instance_name", "host", "instance_name", 0, ""},
		{"instance", "host\\instance-name", "host", "instance-name", 0, ""},
		{"fq-instance", "h.domain.name\\instance-name", "h.domain.name", "instance-name", 0, ""},
		{"port-colon", "host:1433", "host", "", 1433, ""},
		{"tcp:port-colon", "tcp:host:1433", "host", "", 1433, "tcp"},
		{"port-comma", "host,1433", "host", "", 1433, ""},
		{"tcp:port-comma", "tcp:host,1433", "host", "", 1433, "tcp"},
		{"all-three", "host\\instance,1433", "host", "instance", 1433, ""},
		{"tcp:all-three", "tcp:host\\instance,1433", "host", "instance", 1433, "tcp"},
		{"fq-port", "h.domain.com:1433", "h.domain.com", "", 1433, ""},
		{"tcp:fq-port", "tcp:h.domain.com:1433", "h.domain.com", "", 1433, "tcp"},
	}

	for _, tc := range tests {
		protocol, host, instance, port := parseFQDN(tc.in)
		assert.Equal(tc.host, host, tc.name)
		assert.Equal(tc.inst, instance, tc.name)
		assert.Equal(tc.port, port, tc.name)
		assert.Equal(tc.protocol, protocol, tc.name)
	}
}

func TestStripProtocol(t *testing.T) {
	assert := assert.New(t)
	var tests = []struct {
		got      string
		protocol string
		server   string
	}{
		{"host", "", "host"},
		{"tcp:host", "tcp", "host"},
		{"lpc:host", "lpc", "host"},
		{"np:host", "np", "host"},
		{"admin:host", "", "admin:host"},
		{"other:host", "", "other:host"},
	}
	for _, tc := range tests {
		protocol, server := StripProtocol(tc.got)
		assert.Equal(tc.protocol, protocol)
		assert.Equal(tc.server, server)
	}
}
