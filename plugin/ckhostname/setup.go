package ckhostname

import (
	"fmt"
	"net"
	"strings"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("ckhostname", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	rules, err := ckParse(c)
	if err != nil {
		return plugin.Error("ckhostname", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return CKHostname{Next: next, Rules: rules}
	})

	return nil
}

func ckParse(c *caddy.Controller) (map[string]string, error) {
	rules := make(map[string]string)

	for c.Next() {
		args := c.RemainingArgs()

		if len(args) != 2 {
			return nil, fmt.Errorf("bad args: %s, NOT 2 part", fmt.Sprint(args))
		}

		if ip := net.ParseIP(args[1]); ip == nil {
			return nil, fmt.Errorf("bad args[1]: %s, MUST be legal ip.", args[1])
		}

		if !strings.HasSuffix(args[0], ".") {
			args[0] = args[0] + "."
		}

		rules[args[0]] = args[1]
	}
	fmt.Printf("after parse:\n")
	for k, v := range rules {
		fmt.Printf("key: %s, value: %s\n", k, v)
	}
	return rules, nil
}
