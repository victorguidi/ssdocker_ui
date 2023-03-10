package main

import (
	"log"
	"os/user"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type SshServer struct {
	Hosts         []string
	User          *User
	SshConnection []*SshConnection
}

type User struct {
	Username string
	Password string
}

type Container struct {
	ID      string `json:"id"`
	Created string `json:"created"`
	Status  string `json:"status"`
	Ports   string `json:"ports"`
	Names   string `json:"names"`
	Server  string `json:"server"`
}

type Images struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	ImageID    string `json:"image_id"`
	Created    string `json:"created"`
	Size       string `json:"size"`
	Server     string `json:"server"`
}

type SshConnection struct {
	Host    string
	Channel chan string
}

func (s *SshServer) Connect() error {

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	username := user.Username

	channels := make([]chan string, len(s.Hosts))
	var hostKeyCallback ssh.HostKeyCallback
	hostKeyCallback, err = knownhosts.New("/home/" + username + "/.ssh/known_hosts")
	if err != nil {
		// Range over the servers and users and ssh into each server
		for i, server := range s.Hosts {
			config := &ssh.ClientConfig{
				User: s.User.Username,
				Auth: []ssh.AuthMethod{
					ssh.Password(s.User.Password),
				},
				HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
				HostKeyAlgorithms: []string{ssh.KeyAlgoED25519, ssh.KeyAlgoRSA, ssh.KeyAlgoDSA, ssh.KeyAlgoECDSA256, ssh.KeyAlgoECDSA384, ssh.KeyAlgoECDSA521},
			}
			c := make(chan string)
			channels[i] = c
			s.SshConnection = append(s.SshConnection, &SshConnection{
				Host:    server,
				Channel: c,
			})
			go OpenConnection(server, c, config)
		}
	} else {
		// Range over the servers and users and ssh into each server
		for i, server := range s.Hosts {
			config := &ssh.ClientConfig{
				User: s.User.Username,
				Auth: []ssh.AuthMethod{
					ssh.Password(s.User.Password),
				},
				HostKeyCallback:   hostKeyCallback,
				HostKeyAlgorithms: []string{ssh.KeyAlgoED25519, ssh.KeyAlgoRSA, ssh.KeyAlgoDSA, ssh.KeyAlgoECDSA256, ssh.KeyAlgoECDSA384, ssh.KeyAlgoECDSA521},
			}
			c := make(chan string)
			channels[i] = c
			s.SshConnection = append(s.SshConnection, &SshConnection{
				Host:    server,
				Channel: c,
			})
			go OpenConnection(server, c, config)
		}
	}
	return nil
}

func OpenConnection(host string, c chan string, config *ssh.ClientConfig) {
	for {
		// Reading the command from the channel
		cmd := <-c
		if cmd == "exit" {
			break
		}
		// Connecting to the remote server and perform the SSH handshake.
		client, err := ssh.Dial("tcp", host+":"+"22", config)
		if err != nil {
			log.Fatal(err)
		}
		session, err := client.NewSession()
		// Create a new session.
		if err != nil {
			log.Fatal(err)
		}
		// Once a Session is created, we can execute a single command on
		// the remote side using the Run method.
		var b []byte
		b, err = session.Output(cmd)
		if err != nil {
			log.Fatal(err)
		}
		c <- string(b)
		// Once we're done with a session, close it.
		session.Close()
	}
}
