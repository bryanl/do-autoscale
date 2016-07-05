package ctl

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	errInvalidFormat = errors.New("invalid format")
)

// CaddyConfig is configuration for caddy.
type CaddyConfig struct{}

// Generate generates a Caddyfile based on a Config.
func (cc *CaddyConfig) Generate(c *Config) string {
	var buf bytes.Buffer

	if c.AutoConfigureTLS {
		fmt.Fprintln(&buf, c.Hostname)
		fmt.Fprintf(&buf, "tls %s\n", c.TLSEmail)
	} else {
		fmt.Fprintln(&buf, c.Hostname+":443")
		fmt.Fprintln(&buf, "tls /etc/autoscale/ssl/autoscale.crt /etc/autoscale/ssl/autoscale.key")
	}

	fmt.Fprintln(&buf, "proxy / autoscale:8888")
	return buf.String()
}

// Config is configuration for autoscale setup.
type Config struct {
	Token             string
	Hostname          string
	AutoConfigureTLS  bool
	TLSEmail          string
	BasicAuthPassword string
	AllowedHosts      string
}

// Env converts the Config to env style statements.
func (c *Config) Env() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%s=%s\n", "AUTOSCALE_ACCESS_TOKEN", c.Token)
	fmt.Fprintln(&buf, "AUTSCALE_USE_MEMORY_RESOURCES=0")
	fmt.Fprintln(&buf, "AUTSCALE_REGISTER_DEFAULT_METRICS=1")
	fmt.Fprintln(&buf, "AUTSCALE_REGISTER_OFFLINE_METRICS=0")
	fmt.Fprintln(&buf, "AUTSCALE_HTTP_ADDR=:8888")
	fmt.Fprintln(&buf, "AUTSCALE_ENV=production")
	fmt.Fprintln(&buf, "AUTOSCALE_HTTP_ADDR=:8888")
	fmt.Fprintf(&buf, "%s=%s\n", "AUTOSCALE_WEB_PASSWORD", c.BasicAuthPassword)

	return buf.String()
}

func setup() error {
	if _, err := os.Stat("/etc/autoscale/.configured"); err == nil {
		fmt.Println("* setup has already run")
		return nil
	}

	if err := setupDocker(); err != nil {
		return err
	}

	if err := setupEnv(); err != nil {
		return err
	}

	if err := installServices(); err != nil {
		return err
	}

	return nil
}

func setupDocker() error {
	fmt.Println("* installing docker")
	return runScript(installDocker)
}

func runScript(script string) error {
	tmpfile, err := ioutil.TempFile("", "script")
	if err != nil {
		return fmt.Errorf("creating script: %s", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(script)); err != nil {
		return fmt.Errorf("writing script: %s", err)
	}

	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("closing script: %s", err)
	}

	if err := os.Chmod(tmpfile.Name(), 0700); err != nil {
		return fmt.Errorf("making script executable: %s", err)
	}

	cmd := exec.Command(tmpfile.Name())

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running  script: %s", err)
	}

	return nil

}

func installServices() error {
	if err := pullImages(
		"bryanl/postgres-autoscale",
		"prom/prometheus",
		"bryanl/do-autoscale",
		"abiosoft/caddy",
	); err != nil {
		return err
	}

	fmt.Println("* creating docker bridge network")
	if err := runScript(createDockerNetwork); err != nil {
		return err
	}

	fmt.Println("* install postgres service")
	if err := ioutil.WriteFile("/etc/systemd/system/postgres.service", []byte(postgresService), 0644); err != nil {
		return err
	}

	fmt.Println("* install prometheus service")
	if err := ioutil.WriteFile("/etc/systemd/system/prometheus.service", []byte(prometheusService), 0644); err != nil {
		return err
	}

	fmt.Println("* install autoscale service")
	if err := ioutil.WriteFile("/etc/systemd/system/autoscale.service", []byte(autoscaleService), 0644); err != nil {
		return err
	}

	fmt.Println("* install caddy service")
	if err := ioutil.WriteFile("/etc/systemd/system/caddy.service", []byte(caddyService), 0644); err != nil {
		return err
	}

	fmt.Println("* reload systemctl")
	if err := exec.Command("/bin/systemctl", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("unable to reload daemons: %s", err)
	}

	fmt.Println("* start postgres")
	if err := exec.Command("/bin/systemctl", "start", "postgres").Run(); err != nil {
		return fmt.Errorf("unable to start postgres: %s", err)
	}

	fmt.Println("* start prometheus")
	if err := exec.Command("/bin/systemctl", "start", "prometheus").Run(); err != nil {
		return fmt.Errorf("unable to start prometheus: %s", err)
	}

	fmt.Println("* start autoscale")
	if err := exec.Command("/bin/systemctl", "start", "autoscale").Run(); err != nil {
		return fmt.Errorf("unable to start autoscale: %s", err)
	}

	fmt.Println("* start caddy")
	if err := exec.Command("/bin/systemctl", "start", "caddy").Run(); err != nil {
		return fmt.Errorf("unable to start caddy: %s", err)
	}

	if err := ioutil.WriteFile("/etc/autoscale/.configured", []byte(""), 0644); err != nil {
		return fmt.Errorf("unable to mark configuration status: %s", err)
	}

	return nil
}

func pullImages(images ...string) error {
	for _, image := range images {
		fmt.Println("* pulling", image, "image")
		if err := exec.Command("/usr/bin/docker", "pull", image).Run(); err != nil {
			return err
		}
	}

	return nil
}

func setupEnv() error {
	token, err := collectResponse("Enter DigitalOcean token")
	if err != nil {
		return err
	}

	hostname, err := collectResponse("Enter hostname for TLS configuration")
	if err != nil {
		return err
	}

	autoConfigureTLS, err := collectYesNo(fmt.Sprintf("Configure TLS for %s automatically", hostname))
	if err != nil {
		return err
	}

	basicAuthPassword, err := collectResponse("Basic auto password")
	if err != nil {
		return err
	}

	config := Config{
		Token:             token,
		Hostname:          hostname,
		AutoConfigureTLS:  autoConfigureTLS,
		BasicAuthPassword: basicAuthPassword,
	}

	if config.AutoConfigureTLS {
		tlsEmail, err := collectResponse("Lets Encrypt email address")
		if err != nil {
			return err
		}

		config.TLSEmail = tlsEmail
	} else {
		fmt.Println("\n** Automatic TLS will not be configured.",
			"Place your PEM encoded certificate and key in",
			"/etc/autoscale/ssl/autoscale.crt and /etc/autoscale/ssl/autoscale.key.",
			"Afterwards, restart caddy with 'systemctl restart caddy'.")
		fmt.Println()
	}

	if err := os.MkdirAll("/etc/autoscale", 0755); err != nil {
		return fmt.Errorf("creating autoscale config dir: %s", err)
	}

	if err := os.MkdirAll("/etc/autoscale/ssl", 0700); err != nil {
		return fmt.Errorf("creating autoscale ssl dir: %s", err)
	}

	if err := ioutil.WriteFile("/etc/autoscale/autoscale.env", []byte(config.Env()), 0600); err != nil {
		return fmt.Errorf("creating autoscale environment: %s", err)
	}

	cc := &CaddyConfig{}
	if err := ioutil.WriteFile("/etc/autoscale/Caddyfile", []byte(cc.Generate(&config)), 0644); err != nil {
		return fmt.Errorf("creating Caddyfile: %s", err)
	}

	return nil
}

func collectResponse(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("unable to read input: %s", err)
	}

	return strings.TrimSpace(text), nil
}

func collectFormattedResponse(prompt string, format *regexp.Regexp) (string, error) {
	resp, err := collectResponse(prompt)
	if err != nil {
		return "", fmt.Errorf("unable to read input: %s", err)
	}

	if !format.Match([]byte(resp)) {
		return "", errInvalidFormat
	}

	return resp, nil
}

func collectYesNo(prompt string) (bool, error) {
	for {
		resp, err := collectResponse(fmt.Sprintf("%s (y or n)", prompt))
		if err != nil {
			return false, err
		}

		l := strings.ToLower(resp)
		switch l {
		case "y":
			return true, nil
		case "n":
			return false, nil
		}
	}
}

var installDocker = `#!/bin/bash
echo "* checking for docker"
systemctl status docker > /dev/null
if [[ $? == 0 ]]; then
  echo "* docker is installed"
  exit 0
fi

echo "* installing docker"
apt-get update
apt-get install -q -y apt-transport-https ca-certificates
apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
echo "deb https://apt.dockerproject.org/repo ubuntu-xenial main" > /etc/apt/sources.list.d/docker.list
apt-get update

n=0
until [ $n -ge 5 ]
do
  apt-get install -q -y docker-engine && break # sometimes the key server is wonky
  n=$[$n+1]
  sleep 15
done

systemctl start docker
`

var createDockerNetwork = `#!/bin/bash
echo "* creating docker network"
docker network create -d bridge autoscale > /dev/null
`

var postgresService = `[Unit]
Description=Postgres
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
ExecStartPre=-/usr/bin/docker kill db
ExecStartPre=-/usr/bin/docker rm db
ExecStart=/usr/bin/docker run --net=autoscale --name db bryanl/postgres-autoscale
ExecStop=/usr/bin/docker stop -t 2 db

[Install]
WantedBy=multi-user.target`

var prometheusService = `[Unit]
Description=Promethus
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
ExecStartPre=-/usr/bin/docker kill prometheus
ExecStartPre=-/usr/bin/docker rm prometheus
ExecStart=/usr/bin/docker run --net=autoscale --name prometheus prom/prometheus
ExecStop=/usr/bin/docker stop -t 2 prometheus

[Install]
WantedBy=multi-user.target`

var autoscaleService = `[Unit]
Description=Autoscale service
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
ExecStartPre=-/usr/bin/docker kill autoscale
ExecStartPre=-/usr/bin/docker rm autoscale
ExecStart=/usr/bin/docker run --net=autoscale --env-file=/etc/autoscale/autoscale.env --name autoscale bryanl/do-autoscale
ExecStop=/usr/bin/docker stop -t 2 autoscale

[Install]
WantedBy=multi-user.target`

var caddyService = `[Unit]
Description=Caddy service
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
ExecStartPre=-/usr/bin/docker kill caddy
ExecStartPre=-/usr/bin/docker rm caddy
ExecStart=/usr/bin/docker run --net=autoscale -v /etc/autoscale/Caddyfile:/etc/Caddyfile -v /etc/autoscale/ssl:/etc/autoscale/ssl -p 80:80 -p 443:443 --name caddy abiosoft/caddy
ExecStop=/usr/bin/docker stop -t 2 caddy

[Install]
WantedBy=multi-user.target`
