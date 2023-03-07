# Intro to SSDocker

SSDocker is a application where you can manage multiple dockers running in different ssh servers.

All you have to do is provide a config.yml file and you will be able to manage from a single endpoint all Dockers running on each server.

The application, at the current stage can run basic Docker commands and also give you visibility of the logs for each Container.

## Config file example:

```
  hosts:
    - <ip>
    - 10.1.111.1 (for example)

  username: "username"
  password: "password"
```

In order to run I provided a Makefile and also a Dockerfile for the convenience.

  - Make run - will build and run the code in your app in your machine.

## Dockerfile

You can build the project using: 

```
  docker buildx build -t ssdocker .
```

And in order to run it you can use the command: 

```
  docker run --name=ssdocker -d -p 80:8080 ssdocker
```

In case you require a VPN in order to access you should and is using docker you must use the docker --network host flag instead of passing the -p 80:8080 flag. This is because docker will use the same network that the host machine is using, in this case the VPN.

## Possible aware errors:

In case you get an error related to the ssh Key Algorithms, there is a block in the code that we are passing the config, it should look like this: 

``` golang
  config := &ssh.ClientConfig{
    User: s.User.Username,
    Auth: []ssh.AuthMethod{
      ssh.Password(s.User.Password),
    },
    HostKeyCallback:   hostKeyCallback,
    HostKeyAlgorithms: []string{ssh.KeyAlgoED25519},
  }

For now you will have to play aroung with different algorithms, like the ones bellow: 

	- KeyAlgoRSA        = "ssh-rsa"
	- KeyAlgoDSA        = "ssh-dss"
	- KeyAlgoECDSA256   = "ecdsa-sha2-nistp256"
	- KeyAlgoSKECDSA256 = "sk-ecdsa-sha2-nistp256@openssh.com"
	- KeyAlgoECDSA384   = "ecdsa-sha2-nistp384"
	- KeyAlgoECDSA521   = "ecdsa-sha2-nistp521"
	- KeyAlgoED25519    = "ssh-ed25519"
	- KeyAlgoSKED25519  = "sk-ssh-ed25519@openssh.com"

```

