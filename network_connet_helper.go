/*
*	Description : 网络连接相关的方法，有ssh,tcp,udp,websocket,Whois查询, DNS查询等 TODO 测试
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/websocket"
)

// SSHClient 连接ssh
// addr : 主机地址, 如: 127.0.0.1:22
// user : 用户
// pass : 密码
// 返回 ssh连接
func SSHClient(user string, pass string, addr string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	sshConn, err := net.Dial("tcp", addr)
	if nil != err {
		return nil, err
	}
	clientConn, chanList, reqs, err := ssh.NewClientConn(sshConn, addr, config)
	if nil != err {
		_ = sshConn.Close()
		return nil, err
	}
	client := ssh.NewClient(clientConn, chanList, reqs)
	return client, nil
}

// LinuxSendCommand Linux Send Command Linux执行命令
func LinuxSendCommand(command string) (opStr string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command)
	stdout, stdoutErr := cmd.StdoutPipe()
	defer func() {
		_ = stdout.Close()
	}()
	if stdoutErr != nil {
		Error("ERR stdout : ", stdoutErr)
		return stdoutErr.Error()
	}
	if startErr := cmd.Start(); startErr != nil {
		Error("ERR Start : ", startErr)
		return startErr.Error()
	}
	opBytes, opBytesErr := io.ReadAll(stdout)
	if opBytesErr != nil {
		opStr = opBytesErr.Error()
	}
	opStr = string(opBytes)
	_ = cmd.Wait()
	return
}

// WindowsSendCommand Windows Send Command
func WindowsSendCommand(command []string) (opStr string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if len(command) < 1 {
		return ""
	}
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	stdout, stdoutErr := cmd.StdoutPipe()
	if stdoutErr != nil {
		Error("ERR stdout : ", stdoutErr)
		return stdoutErr.Error()
	}
	defer func() {
		_ = stdout.Close()
	}()
	if startErr := cmd.Start(); startErr != nil {
		Error("ERR Start : ", startErr)
		return startErr.Error()
	}
	opBytes, opBytesErr := io.ReadAll(stdout)
	if opBytesErr != nil {
		Error(opBytesErr)
		return opBytesErr.Error()
	}
	opStr = string(opBytes)
	_ = cmd.Wait()
	return
}

// WindowsSendPipe TODO  执行windows 管道命令
func WindowsSendPipe(command1, command2 []string) (opStr string) {
	return ""
}

// ================ TCP ================
/*
Tcp的连接 (Tcp客户端); 应用场景是模拟Tcp客户端;

// ===== Use
func main(){
	client := gt.NewTcpClient()
	client.Run("192.168.0.9:29123", f)
}

func f(client *gt.TcpClient){
	go func() {
		// 发送登录请求
		_,err := client.Send([]byte(`{
			"cmd":"Auth",
			"data":{
				"account":"a10",
				"password":"123456",
				"device":"1",
				"source":"windows"
			}
		}`))
		if err != nil {
			log.Println("err = ", err)
		}
	}()
}
*/

// TcpClient Tcp客户端
type TcpClient struct {
	Connection *net.TCPConn
	HawkServer *net.TCPAddr
	StopChan   chan struct{}
	CmdChan    chan string
	Token      string
	RConn      chan struct{}
}

func NewTcpClient() *TcpClient {
	return new(TcpClient)
}

func (c *TcpClient) Send(b []byte) (int, error) {
	if c.Connection == nil {
		return 0, fmt.Errorf("conn is null")
	}
	return c.Connection.Write(b)
}

func (c *TcpClient) Read(b []byte) (int, error) {
	if c.Connection == nil {
		return 0, fmt.Errorf("conn is null")
	}
	return c.Connection.Read(b)
}

func (c *TcpClient) Addr() string {
	if c.Connection == nil {
		return ""
	}
	return c.Connection.RemoteAddr().String()
}

func (c *TcpClient) Close() {
	if c.Connection == nil {
		return
	}
	_ = c.Connection.Close()
}

func (c *TcpClient) Stop() {
	c.StopChan <- struct{}{}
}

func (c *TcpClient) ReConn() {
	c.RConn <- struct{}{}
}

func (c *TcpClient) Run(serverHost string, r func(c *TcpClient, data []byte), w func(c *TcpClient)) {
	//用于重连
Reconnection:

	hawkServer, err := net.ResolveTCPAddr("tcp", serverHost)
	if err != nil {
		log.Printf("hawk server [%s] resolve error: [%s]", serverHost, err.Error())
		time.Sleep(1 * time.Second)
		goto Reconnection
	}

	//连接服务器
	connection, err := net.DialTCP("tcp", nil, hawkServer)
	if err != nil {
		log.Printf("connect to hawk server error: [%s]", err.Error())
		time.Sleep(1 * time.Second)
		goto Reconnection
	}
	log.Println("[连接成功] 连接服务器成功")

	//创建客户端实例
	c.Connection = connection
	c.HawkServer = hawkServer
	c.StopChan = make(chan struct{})
	c.CmdChan = make(chan string)
	c.RConn = make(chan struct{})

	//启动接收
	go func(conn *TcpClient) {
		for {
			read := make([]byte, 1024)
			for {
				n, err := conn.Connection.Read(read)
				if err != nil {
					if err == io.EOF {
						log.Println(conn.Addr(), " 断开了连接!")
					}
					conn.Close()
					c.RConn <- struct{}{}
					return
				}
				if n > 0 && n < 1025 {
					log.Println(string(read[:n]))
					r(c, read[:n])
				}

			}
		}
	}(c)

	go w(c)

	for {
		select {
		case a := <-c.RConn:
			log.Println("global.RConn = ", a)
			goto Reconnection
		case <-c.StopChan:
			return

		}
	}
}

// ================ UDP ================
/*
	Description : Udp的连接 (Udp客户端); 应用场景是模拟Udp客户端;
*/

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
	defer func() {
		_ = u.Conn.Close()
	}()

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

// ================ Websocket ================
/*
	Description : websocket的连接, 模拟websocket客户端
*/

// WSClient websocket 客户端
type WSClient interface {
	Send(body []byte) error
	Read(data []byte) error
	Close()
}

// WsClient websocket 客户端
func WsClient(host, path string, isSSL bool) (WSClient, error) {
	ws := &webSocketClient{
		Host:  host,
		Path:  path,
		IsSSL: isSSL,
	}
	err := ws.conn()
	return ws, err
}

type webSocketClient struct {
	Host  string
	Path  string
	Ws    *websocket.Conn
	IsSSL bool
}

func (c *webSocketClient) conn() error {
	var err error
	u := c.Host + c.Path
	if c.IsSSL {
		c.Ws, err = websocket.Dial(u, "", "https://"+c.Host+"/")
	} else {
		c.Ws, err = websocket.Dial(u, "", "http://"+c.Host+"/")
	}
	return err
}

func (c *webSocketClient) Send(body []byte) error {
	_, err := c.Ws.Write(body)
	return err
}

func (c *webSocketClient) Close() {
	_ = c.Ws.Close()
}

func (c *webSocketClient) Read(data []byte) error {
	_, err := c.Ws.Read(data)
	return err
}

// ================ Whois ================

var RootWhoisServers = "whois.iana.org:43"

type WhoisInfo struct {
	Root string
	Rse  string
}

func Whois(host string) *WhoisInfo {
	hostList := strings.Split(host, ".")
	host = strings.Join(hostList[len(hostList)-2:], ".")
	info := &WhoisInfo{}
	rootRse := whois(RootWhoisServers, host)
	info.Root = rootRse
	referList := regFindTxt(`(?is:refer:(.*?)\n)`, rootRse)
	if len(referList) > 0 {
		refer := StrDeleteSpace(referList[0])
		rse := whois(refer+":43", host)
		info.Rse = rse
	}
	return info
}

func whois(server, host string) string {
	conn, _ := net.Dial("tcp", server)
	_, _ = conn.Write([]byte(host + " \r\n"))
	buf := make([]byte, 1024*10)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		Error(err)
		return ""
	}
	rse := string(buf[:n])
	defer func() {
		_ = conn.Close()
	}()
	return rse
}

// ================ DNS查询 ================

type DNSInfo struct {
	IPs           []string `json:"ips"`
	LookupCNAME   string   `json:"cname"`
	DnsServerIP   string   `json:"dnsServerIP"`
	DnsServerName string   `json:"dnsServerName"`
	IsCDN         bool     `json:"isCDN"`
	Ms            float64  `json:"ms"`
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
	for k := range allIPMap {
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
	"8.8.8.8":       "公共DNS|GoogleDNS",
	//"8.8.4.4":       "公共DNS|GoogleDNS",
	//"1.1.1.1":       "公共DNS|CloudflareDNS",
	//"1.0.0.1":       "公共DNS|CloudflareDNS",
	//"9.9.9.9": "公共DNS|IBM Quad9DNS",
	// "185.222.222.222": "公共DNS|DNS.SB",
	//"185.184.222.222": "公共DNS|DNS.SB",
	//"208.67.222.222": "公共DNS|OpenDNS",
	//"208.67.220.220": "公共DNS|OpenDNS",
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
