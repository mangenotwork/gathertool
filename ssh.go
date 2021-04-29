package gathertool

import (
	"net"

	"golang.org/x/crypto/ssh"
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

	clientConn, chans, reqs, err := ssh.NewClientConn(sshConn, addr, config)
	if nil != err {
		sshConn.Close()
		return nil, err
	}

	client := ssh.NewClient(clientConn, chans, reqs)
	return client, nil
}
