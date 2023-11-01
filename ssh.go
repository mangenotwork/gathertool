/*
	Description : ssh 连接等相关的方法
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"context"
	"io"
	"net"
	"os/exec"

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
		_ = sshConn.Close()
		return nil, err
	}
	client := ssh.NewClient(clientConn, chans, reqs)
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
	cmd := exec.CommandContext(ctx, command[0], command[1:len(command)]...)
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
