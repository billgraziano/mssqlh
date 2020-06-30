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
	assert.Equal("sqlserver://u1:redacted@localhost", conn.Redacted())
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
}

func TestNameSplitter(t *testing.T) {
	assert := assert.New(t)
	var tests = []struct {
		name string
		in   string
		host string
		inst string
		port int
	}{
		{"host", "host", "host", "", 0},
		{"host", "host-name", "host-name", "", 0},
		{"host", "host_name", "host_name", "", 0},
		{"fqdb", "host.domain.name", "host.domain.name", "", 0},
		{"instance", "host\\instance", "host", "instance", 0},
		{"instance", "host\\instance_name", "host", "instance_name", 0},
		{"instance", "host\\instance-name", "host", "instance-name", 0},
		{"fq-instance", "h.domain.name\\instance-name", "h.domain.name", "instance-name", 0},
		{"port-colon", "host:1433", "host", "", 1433},
		{"port-comma", "host,1433", "host", "", 1433},
		{"all-three", "host\\instance,1433", "host", "instance", 1433},
		{"fq-port", "h.domain.com:1433", "h.domain.com", "", 1433},
	}

	for _, v := range tests {
		host, instance, port := parseFQDN(v.in)
		assert.Equal(v.host, host, v.name)
		assert.Equal(v.inst, instance, v.name)
		assert.Equal(v.port, port, v.name)
	}
}
