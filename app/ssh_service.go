package app

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	ek "main/app/error"
	"strconv"
	"strings"
	"sync"
)

type SSHService struct {
	client *ssh.Client
}

func NewSSHService(user, password, ipaddr string, port int) *SSHService {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	host := ipaddr + ":" + strconv.Itoa(port)
	c, err := ssh.Dial("tcp", host, config)
	ek.CheckError(err, "ssh登录失败")
	return &SSHService{client: c}
}

func (this *SSHService) GetMuxShell(cmd []string) string {
	session, err := this.client.NewSession()
	ek.CheckError(err, "创建shell")
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // 禁用回显（0禁用，1启动）
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("vt100", 80, 4000, modes); err != nil {
		log.Fatal(err)
	}

	w, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}
	r, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}
	e, err := session.StderrPipe()
	if err != nil {
		panic(err)
	}

	in, out := muxShell(w, r, e)
	if err := session.Shell(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(<-out) //ignore the shell output
	fmt.Println("--------------------------------------------")
	in <- "show arp"
	in <- "show int status"

	in <- "exit"
	in <- "exit"

	fmt.Printf("%s\n%s\n", <-out, <-out)

	_, _ = <-out, <-out
	session.Wait()
	return ""
}

func muxShell(w io.Writer, r, e io.Reader) (chan<- string, <-chan string) {
	in := make(chan string, 3)
	out := make(chan string, 5)
	var wg sync.WaitGroup
	wg.Add(1) //for the shell itself
	go func() {
		for cmd := range in {
			wg.Add(1)
			w.Write([]byte(cmd + "\n"))
			wg.Wait()
		}
	}()

	go func() {
		var (
			buf [65 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:])
			if err != nil {
				fmt.Println(err.Error())
				close(in)
				close(out)
				return
			}
			t += n
			result := string(buf[:t])
			if strings.Contains(result, "Username:") ||
				strings.Contains(result, "Password:") ||
				strings.Contains(result, "#") {
				out <- string(buf[:t])
				t = 0
				wg.Done()
			}
		}
	}()
	return in, out
}
