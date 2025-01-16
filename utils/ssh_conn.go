package utils

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
)

type SShConfig struct {
	User                 string
	Password             string
	Host                 string
	Port                 string
	PrivateKey           string
	PrivateKeyPassphrase string
}

func NewSSHClient(config *SShConfig) (client *ssh.Client, err error) {
	var authMethod []ssh.AuthMethod
	var singner ssh.Signer

	if config.PrivateKey != "" {
		if config.PrivateKeyPassphrase == "" {
			singner, err = ssh.ParsePrivateKey([]byte(config.PrivateKey))
			if err != nil {
				return
			}
		} else {
			singner, err = ssh.ParsePrivateKeyWithPassphrase([]byte(config.PrivateKey), []byte(config.PrivateKeyPassphrase))
			if err != nil {
				return
			}
		}
		authMethod = append(authMethod, ssh.PublicKeys(singner))
	} else {
		if config.Password == "" {
			return nil, errors.New("password and private key is empty")
		}
		authMethod = append(authMethod, ssh.Password(config.Password))
	}

	clientConfig := &ssh.ClientConfig{
		User:            config.User,
		Auth:            authMethod,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err = ssh.Dial("tcp", config.Host+":"+config.Port, clientConfig)

	return
}

// 执行远程命令
func ExecuteSSHCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// 执行命令并获取输出
	output, err := session.CombinedOutput(command)
	if err != nil {
		return string(output), fmt.Errorf("failed to execute command: %w", err)
	}

	return string(output), nil
}

// 上传单个文件
func UploadFile(client *ssh.Client, localPath string, remotePath string) error {
	// 获取 SFTP 客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	// 打开本地文件
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer localFile.Close()

	// 创建远程文件
	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %w", err)
	}
	defer remoteFile.Close()

	// 将本地文件复制到远程文件
	_, err = remoteFile.ReadFrom(localFile)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

//func main() {
//	privateKey := `
//-----BEGIN RSA PRIVATE KEY-----
//MIIEogIBAAKCAQEA7+5e4VTY3760+TZQRRqSH4haItctguc0c41wIxsVjc1Qp9bP
//qZesf6rdhGln07qQJ1zRSdQaJafc08ChEtWAHygM5olWfSsx86E9sjGiq1TmBq0j
//Ipt1bNMdnMdAn5WQ2OBOBzaTNOe/ZnumaARG+PqNly54DTvdf9oU8aupLFXyzf2v
//xSrHpaGbK23ztx89KooRM162dziQtAWvlc0QY4fVmUcXnRXIShVfDwTiaax0hSDO
//6yhrNmX4VkHfg5SGw4Nic+8UwYqxU6/93M/klNl7sV2IcI4lWd3VIhqX0qpu9uRZ
//fkmtawexh6eGiVxlY/+50r9B09lAgSfmMPHCxQIDAQABAoIBAECKMXGRnkkJhqrm
//5k7AaAAdcImgsjhdMynGKRz4YyLi0MhlKzUmfJtW/gFpTSmSLMa52/5tFJ9+eRNo
//5KspTS6UWmwFE4PWA7jSbcMqQOSvkLTJDAN6J+sfGv8aRhLna7A7HiIolw6LLFxH
//9KpURDcjRsUdFeQRb3K92hZsI7St4KgYw51x6DA+832DKpH2+yxKbqY6g7buS7kD
//2pMTm6GH+DafgzfuCHOL6SbKzML81fo9XBu6iBzCkEIsJl7wtV4ZiAL9jcUiS2jD
//0Xvdb3zntno9fhcWahqRjAWXcC+PPbeliqTHfBs+R5dWcmEaazoMfOQQ7aXy2tR8
//AMvndT0CgYEA/HYb7Iggpw5Ez8vWZ9y5xquJpHmVE0//UByog/8aq18qe0TzsI9Z
//vRGDeo8u5L/LeKrokoUoXuTOVPBZiEIRPghgJhtp4sTiQuSu9agkj6svVnMZfYhA
//lwpnKN/biWTv9T/Nte3XG1JvG7wGq0DdHdiR34FvxHb6Xyja6dQocHMCgYEA80tM
//1YSbzE5cWRx+NH1dCdPbTv2WFjNoNJeLWsyVZjweXLeDBVxk3ZX/Av+2kqDlVzw8
//LkSIhhPWRPQytlO4v3Y4Tgo5/eNjuN6+bsagVXHyPUM0G5HZVlv1B9gSn6xH7PWs
//Jfo4cVMvTnTrdTBZzyPOwX4GYMHU3RXPH6auyecCgYBG50z0Y074XBOLYK44wU8T
//sv2XSeZKZD9KWqIhYDY3RyUBNd5TCg+kABUzCJ+c8xjMLQPgkrFB5XTlehNLJ3L8
//PxHx4eUdITqCmwNgTvbluqgy2WShUvEA+pT6b9SSg9y4vlCh9chiDgbSfT5KPo9b
//YIWnhgzD2r56l1jULxekbQKBgAGfTerakINjPmBlvT2yXE11eS/kpvyM6TP4krhP
//RuvAmN87ZgdCH3YOyv2FIP2HTyAuyaPxVwu11CbvjesDUecM7cEvdkWIH6Ea8yAf
//+O+468mWyiEo7s8Rm+eqfC1OY8hjtvsl2PyAdn9KbkuAwAiOj5FgusAoarfyrkfi
//v6WfAoGAC6WGOLA3Ep8ipzbmYqptpwiKcNUOoQREGzj6luBhkThxVhi4zZ/y6Nby
//MDzrMbmxZTmxrcJZRlwbxzyzqxRk84QbitDeBVZNa0S6olwoHh4QgAdefcWAT2vy
//TjGrz3SLG6U8tD4Kv+BPW0j62acooRCmZSuhMsHUYXZL4GzbeQw=
//-----END RSA PRIVATE KEY-----
//`
//	config := &SShConfig{
//		User:       "root",
//		PrivateKey: privateKey,
//		Host:       "192.168.128.129",
//		Port:       "22",
//	}
//	client, err := NewSSHClient(config)
//	if err != nil {
//		fmt.Println("Failed to create SSH client:", err)
//		return
//	}
//	defer client.Close()
//
//	result, err := ExecuteSSHCommand(client, "ls -l")
//	if err != nil {
//		fmt.Println("Failed to execute command:", err)
//		return
//	}
//
//	fmt.Println(result)
//}
