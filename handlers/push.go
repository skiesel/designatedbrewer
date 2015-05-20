package handlers

import (
	"fmt"
	"net/http"
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
	// "github.com/skiesel/designatedbrewer/utils"
	"net"
)

func Push(w http.ResponseWriter, r *http.Request) {
	config := ssh.Config{}
	config.SetDefaults()

	clientConfig := &ssh.ClientConfig{
		Config : config,
		User : "designatedbrewer@kieselnet.com",
		Auth : []ssh.AuthMethod{ ssh.Password("F7bJ8g,pf(utT2Z?Hk") },
		HostKeyCallback : nil,
		ClientVersion : "",
	}

	tcpCon, err := net.Dial("tcp", "ftp.longsincehere.com:22")
	if err != nil {
		panic(err)
	}

	sshCon, sshChan, sshReq, err := ssh.NewClientConn(tcpCon, "ftp.longsincehere.com:22", clientConfig)

	sshClient := ssh.NewClient(sshCon, sshChan, sshReq)

	// open an SFTP sesison over an existing ssh connection.
	sftp, err := sftp.NewClient(sshClient)
	if err != nil {
		panic(err)
	}
	defer sftp.Close()

	// walk a directory
	walk := sftp.Walk("/home1/longsinc/designatedbrewer")
	for walk.Step() {
		if walk.Err() != nil {
			continue
		}
		// log.Println(w.Path())
	}

	// leave your mark
	f, err := sftp.Create("hello.txt")
	if err != nil {
		panic(err)
	}
	if _, err := f.Write([]byte("Hello world!")); err != nil {
		panic(err)
	}

	// check it's there
	fi, err := sftp.Lstat("hello.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(fi)


	fmt.Fprint(w, "success")
}