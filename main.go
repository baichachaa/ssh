package main

import (
	"fmt"
	"main/app"
)

//func main()  {
//	// 建立SSH客户端连接
//	client, err := ssh.Dial("tcp", "baichacha.cn:22", &ssh.ClientConfig{
//		User:            "root",
//		Auth:            []ssh.AuthMethod{ssh.Password("Linux@11")},
//		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//	})
//	if err != nil {
//		log.Fatalf("SSH dial error: %s", err.Error())
//	}
//	// 建立新会话
//	session, err := client.NewSession()
//	if err != nil {
//		log.Fatalf("new session error: %s", err.Error())
//	}
//
//	defer session.Close()
//	session.Stdout = os.Stdout // 会话输出关联到系统标准输出设备
//	session.Stderr = os.Stderr // 会话错误输出关联到系统标准错误输出设备
//	session.Stdin = os.Stdin   // 会话输入关联到系统标准输入设备
//	modes := ssh.TerminalModes{
//		ssh.ECHO:          0,  // 禁用回显（0禁用，1启动）
//		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
//		ssh.TTY_OP_OSPEED: 14400, //output speed = 14.4kbaud
//	}
//	if err = session.RequestPty("linux", 32, 160, modes); err != nil {
//		log.Fatalf("request pty error: %s", err.Error())
//	}
//	if err = session.Shell(); err != nil {
//		log.Fatalf("start shell error: %s", err.Error())
//	}
//	if err = session.Wait(); err != nil {
//		log.Fatalf("return error: %s", err.Error())
//	}
//}
//

func main() {
	s := app.NewSSHService("root", "Linux@11", "baichacha.cn", 22)
	b := s.GetMuxShell([]string{"id", "exit"})
	fmt.Println(b)
}
