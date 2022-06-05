package main

import (
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
)

var IpMap map[string]bool

func dnsquery(domain string, ip string, DnsServer string, OnlyIp bool, repeat int) {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(domain, dns.TypeA)
	o := new(dns.OPT)
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	e := new(dns.EDNS0_SUBNET) //EDNS
	e.Code = dns.EDNS0SUBNET
	e.Family = 1         // 1 IPv4 2 IPv6
	e.SourceNetmask = 24 //  地址掩码 一般为24 (谷歌不支持大于24)
	e.SourceScope = 0
	e.Address = net.ParseIP(ip).To4()
	o.Option = append(o.Option, e)
	m.Extra = append(m.Extra, o)
	for i := 0; i < repeat; i++ {
		in, _, err := c.Exchange(m, DnsServer) //注意:要选择支持自定义EDNS的DNS 或者是 目标NS服务器  国内DNS大部分不支持自定义EDNS数据

		if err != nil {
			log.Fatal(err)
		}
		for _, answer := range in.Answer {

			if answer.Header().Rrtype == dns.TypeA {
				if OnlyIp {
					IpMap[answer.(*dns.A).A.String()] = true
				} else {
					println(answer.(*dns.A).A.String())
				}
			}
		}
	}

}

func main() {
	Initlist()
	var domain = flag.String("d", "www.taobao.com", "domain")
	var DnsServer = flag.String("s", "8.8.8.8:53", "dns server addr")
	var ip = flag.String("ip", "", "client ip")
	var OnlyIp = flag.Bool("i", false, "Only output ip addr")
	var repeat = flag.Int("r", 1, "repeat query rounds")
	flag.Parse()
	IpMap = make(map[string]bool)
	if *ip != "" {
		*OnlyIp = true
		dnsquery(*domain, *ip, *DnsServer, *OnlyIp, *repeat)
	} else {
		for city, ip := range CityMap {
			if !*OnlyIp {
				fmt.Println(city)
			}
			dnsquery(*domain, ip, *DnsServer, *OnlyIp, *repeat)

		}
	}

	if *OnlyIp {
		for ip, _ := range IpMap {
			println(ip)
		}
	}
	//log.Println(in)

}
