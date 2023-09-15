/*
	Description : Udp的连接 (Udp客户端); 应用场景是模拟Udp客户端;
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// UdpClient Udp客户端
type UdpClient struct {
	SrcAddr *net.UDPAddr
	DstAddr *net.UDPAddr
	Conn    *net.UDPConn
	Token   string
}

func (u *UdpClient) Run(hostServer string, port int, r func(u *UdpClient, data []byte), w func(u *UdpClient)) {
	var err error
	sip := net.ParseIP(hostServer)
	u.SrcAddr = &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	u.DstAddr = &net.UDPAddr{IP: sip, Port: port}
	u.Conn, err = net.DialUDP("udp", u.SrcAddr, u.DstAddr)
	if err != nil {
		log.Println(err)
	}
	log.Println("连接成功; c = ", u.Conn)
	defer u.Conn.Close()

	go func() {
		data := make([]byte, 1024)
		for {
			n, remoteAddr, err := u.Conn.ReadFromUDP(data)
			if err != nil {
				log.Printf("error during read: %s", err)
			}
			log.Printf("<%s> %s\n", remoteAddr, data[:n])
			r(u, data[:n])
		}
	}()

	go w(u)

	select {}
}

func NewUdpClient() *UdpClient {
	return new(UdpClient)
}

func (u *UdpClient) Send(b []byte) (int, error) {
	if u.Conn == nil {
		return 0, fmt.Errorf("conn is null")
	}
	return u.Conn.Write(b)
}

func (u *UdpClient) Read(b []byte) (int, *net.UDPAddr, error) {
	if u.Conn == nil {
		return 0, u.DstAddr, fmt.Errorf("conn is null")
	}
	return u.Conn.ReadFromUDP(b)
}

func (u *UdpClient) Addr() string {
	return u.DstAddr.String()
}

func (u *UdpClient) Close() {
	if u.Conn == nil {
		return
	}
	_ = u.Conn.Close()
}

// DNSInfo
type DNSInfo struct {
	IPs           []string `json:"ips"`
	LookupCNAME   string   `json:"cname"`
	DnsServerIP   string   `json:"dnsServerIP"`
	DnsServerName string   `json:"dnsServerName"`
	IsCDN         bool     `json:"isCDN"`
	Ms            float64  `json:"ms"` // ms
}

// NsLookUp DNS查询
func NsLookUp(host string) *DNSInfo {
	dnsInfo := &DNSInfo{
		DnsServerIP:   "Local",
		DnsServerName: "Local",
		IsCDN:         false,
	}
	start := time.Now().UnixNano()
	ips, err := net.LookupHost(host)
	if err == nil {
		dnsInfo.IPs = ips
	}
	dnsInfo.Ms = float64(time.Now().UnixNano()-start) / 100000
	cname, err := net.LookupCNAME(host)
	if err == nil {
		dnsInfo.LookupCNAME = cname
	}
	if len(dnsInfo.LookupCNAME) > 0 && strings.Index(dnsInfo.LookupCNAME, host) == -1 {
		dnsInfo.IsCDN = true
	}
	return dnsInfo
}

// NsLookUpFromDNSServer 指定DNS服务器查询
func NsLookUpFromDNSServer(host, dnsServer string) *DNSInfo {
	return nsLookUpFromServer(host, dnsServer, "")
}

func nsLookUpFromServer(host, dnsServer, name string) *DNSInfo {
	r := &net.Resolver{
		PreferGo:     true,
		StrictErrors: false,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: time.Millisecond * 200}
			c, e := d.DialContext(ctx, "udp", dnsServer)
			return c, e
		},
	}
	dnsInfo := &DNSInfo{
		DnsServerIP:   dnsServer,
		DnsServerName: name,
		IsCDN:         false,
	}
	start := time.Now().UnixNano()
	ip, err := r.LookupHost(context.Background(), host)
	if err == nil {
		dnsInfo.IPs = ip
	} else {
		Error(dnsServer, " | ", err)
	}
	dnsInfo.Ms = float64(time.Now().UnixNano()-start) / 100000
	cname, err := r.LookupCNAME(context.Background(), host)
	if err == nil {
		dnsInfo.LookupCNAME = cname
	} else {
		Error(dnsServer, " | ", err)
	}
	if len(dnsInfo.LookupCNAME) > 0 && strings.Index(dnsInfo.LookupCNAME, host) == -1 {
		dnsInfo.IsCDN = true
	}
	return dnsInfo
}

// NsLookUpAll 在多个DNS服务器查询
func NsLookUpAll(host string) ([]*DNSInfo, []string) {
	all := make([]*DNSInfo, 0)
	allIPMap := make(map[string]struct{})
	allIP := make([]string, 0)
	results := make(chan *DNSInfo, len(DNSServer))
	for k, v := range DNSServer {
		go func(k, v string) {
			results <- nsLookUpFromServer(host, k+":53", v)
		}(k, v)
	}
	for i := 0; i < len(DNSServer); i++ {
		res := <-results
		all = append(all, res)
	}
	for _, v := range all {
		for _, i := range v.IPs {
			allIPMap[i] = struct{}{}
		}
	}
	for k, _ := range allIPMap {
		allIP = append(allIP, k)
	}
	return all, allIP
}

// DNSServer DNS服务器地址大全  ip:服务商
var DNSServer = map[string]string{
	"114.114.114.114": "公共DNS|114DNS",
	"114.114.115.115": "公共DNS|114DNS",
	//"119.29.29.29":    "公共DNS|DNSPod DNS+",
	// "182.254.116.116": "公共DNS|DNSPod DNS+",
	"101.226.4.6":   "公共DNS|DNS 派 电信/移动/铁通",
	"218.30.118.6":  "公共DNS|DNS 派 电信/移动/铁通",
	"123.125.81.6":  "公共DNS|DNS DNS 派 联通",
	"140.207.198.6": "公共DNS|DNS DNS 派 联通",
	//"1.2.4.8":       "公共DNS|cnnicDNS",
	//"210.2.4.8":     "公共DNS|cnnicDNS",
	"8.8.8.8": "公共DNS|GoogleDNS",
	//"8.8.4.4":       "公共DNS|GoogleDNS",
	//"1.1.1.1":       "公共DNS|CloudflareDNS",
	//"1.0.0.1":       "公共DNS|CloudflareDNS",
	//"9.9.9.9": "公共DNS|IBM Quad9DNS",
	// "185.222.222.222": "公共DNS|DNS.SB",
	//"185.184.222.222": "公共DNS|DNS.SB",
	//"208.67.222.222": "公共DNS|OpenDNS",
	//"208.67.220.220": "公共DNS|OpenDNS",
	// "199.91.73.222":   "公共DNS|V2EXDNS",
	// "178.79.131.110":  "公共DNS|V2EXDNS",
	"223.5.5.5": "公共DNS|阿里云DNS",
	"223.6.6.6": "公共DNS|阿里云DNS",
	//"183.60.83.19": "公共DNS|腾讯云DNS",
	//"183.60.82.98":    "公共DNS|腾讯云DNS",
	"180.76.76.76": "公共DNS|百度云DNS",
	"4.2.2.1":      "公共DNS|微软云DNS",
	//"4.2.2.2":         "公共DNS|微软云DNS",
	"122.112.208.1":   "公共DNS|华为云DNS",
	"139.9.23.90":     "公共DNS|华为云DNS",
	"114.115.192.11":  "公共DNS|华为云DNS",
	"116.205.5.1":     "公共DNS|华为云DNS",
	"116.205.5.30":    "公共DNS|华为云DNS",
	"122.112.208.175": "公共DNS|华为云DNS",
	"139.159.208.206": "公共DNS|华为云DNS",
	"180.184.1.1":     "公共DNS|字节跳动",
	"180.184.2.2":     "公共DNS|字节跳动",
	"168.95.192.1":    "公共DNS|中華電信DNS",
	//"168.95.1.1":      "公共DNS|中華電信DNS",
	//"203.80.96.10":    "公共DNS|香港宽频DNS",
	// "203.80.96.9":     "公共DNS|香港宽频DNS",
	//"199.85.126.10": "公共DNS|赛门铁克诺顿DNS",
	//"199.85.127.10": "公共DNS|赛门铁克诺顿DNS",
	//"216.146.35.35": "公共DNS|oracle+dynDNS",
	//"216.146.36.36":  "公共DNS|oracle+dynDNS",
	//"64.6.64.6":      "公共DNS|瑞士瑞信银行DNS",
	//"64.6.65.6":      "公共DNS|瑞士瑞信银行DNS",
	"61.132.163.68":  "电信|安徽",
	"202.102.213.68": "电信|安徽",
	"202.98.192.67":  "电信|贵州",
	"202.98.198.167": "电信|贵州",
	"202.101.224.69": "电信|江西",
	"61.139.2.69":    "电信|四川",
	"218.6.200.139":  "电信|四川",
	"202.98.96.68":   "电信|四川成都",
	"202.102.192.68": "电信|安徽合肥",
	"221.7.1.21":     "联通|新疆维吾尔自治区乌鲁木齐",
	"202.106.46.151": "联通|北京",
}
