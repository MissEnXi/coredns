package ckhostname

import (
	"context"
	"fmt"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type CKHostname struct {
	Next  plugin.Handler
	Rules map[string]string
}

func (hn CKHostname) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	reply, ok := hn.Rules[qname]
	if !ok {
		return plugin.NextOrFailure(hn.Name(), hn.Next, ctx, w, r)
	}

	fmt.Printf("YES! match %s for request: %s\n", reply, qname)

	resp := hn.toMsg(r, qname, reply)
	w.WriteMsg(resp)
	return dns.RcodeSuccess, nil
}

func (hn CKHostname) toMsg(r *dns.Msg, request, reply string) *dns.Msg {

	rr := new(dns.A)
	rr.Hdr = dns.RR_Header{Name: request, Rrtype: dns.TypeA, Class: dns.ClassINET}
	rr.A = net.ParseIP(reply).To4()

	var answers []dns.RR
	answers = append(answers, rr)

	m1 := new(dns.Msg)
	m1.SetReply(r)
	m1.Authoritative = true
	m1.Answer = answers

	return m1
}

func (hn CKHostname) Name() string { return "ckhostname" }
