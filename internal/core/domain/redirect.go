package domain

import (
	"errors"
	"fmt"
	"regexp"
)

type HostType uint

const (
	IPv4 HostType = iota
	IPv6
)

type RedirectOpts struct {
	// Host specifies the IP address for the redirect service to bind to.
	//
	// Example: "127.0.0.1" or "::1"
	Host string

	// HostType defines the type of host (e.g., IPv4 or IPv6).
	HostType HostType

	// Port defines the TCP port number on which the redirect service will listen.
	// Must be a valid port number (1-65535).
	//
	// Example: 8080, 443
	Port int

	// TLS enables HTTPS (TLS) for secure encrypted connections.
	// When enabled, CertFile and KeyFile must be provided.
	//
	// Default: false (plain HTTP)
	TLS bool

	// CertFile specifies the path to the TLS certificate file (PEM-encoded).
	// Required when TLS is enabled.
	//
	// Example: "/etc/ssl/certs/server.crt" or "cert.pem"
	CertFile string

	// KeyFile specifies the path to the private key file for the TLS certificate.
	// Required when TLS is enabled.
	//
	// Example: "/etc/ssl/private/server.key" or "key.pem"
	KeyFile string

	// NotFoundTempl enables rendering a custom HTML template for 404 responses
	// when a requested link is not found, instead of a plain text error.
	//
	// Default: true
	NotFoundTempl bool
}

var DefaultRedirectOpts = RedirectOpts{
	Host:          "127.0.0.1",
	HostType:      IPv4,
	Port:          2000,
	NotFoundTempl: true,
}

type RedirectOpt func(*RedirectOpts) error

const ipv4Regex = `^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`

const ipv6Regex = `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|` +
	`([0-9a-fA-F]{1,4}:){1,7}:|` +
	`([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|` +
	`([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|` +
	`([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|` +
	`([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|` +
	`([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|` +
	`[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|` +
	`:((:[0-9a-fA-F]{1,4}){1,7}|:)|` +
	`fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|` +
	`::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|` +
	`([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`

func WithHost(host string) RedirectOpt {
	return func(opts *RedirectOpts) error {
		isIPv4, err := regexp.MatchString(ipv4Regex, host)
		if err != nil {
			return err
		}

		isIPv6, err := regexp.MatchString(ipv6Regex, host)
		if err != nil {
			return err
		}

		if !isIPv4 && !isIPv6 {
			return errors.New("host must be IPv4 or IPv6")
		}

		if isIPv6 {
			opts.HostType = IPv6
		}

		opts.Host = host
		return nil
	}
}

func WithPort(port int) RedirectOpt {
	minPort := 0
	maxPort := 65535

	return func(opts *RedirectOpts) error {
		if minPort > port && maxPort < port {
			return fmt.Errorf("port must be >%d and <%d", minPort, maxPort)
		}

		opts.Port = port
		return nil
	}
}

func WithNotFoundTempl(use bool) RedirectOpt {
	return func(opts *RedirectOpts) error {
		opts.NotFoundTempl = use
		return nil
	}
}

func WithTLS(certFile, keyFile string) RedirectOpt {
	return func(opts *RedirectOpts) error {
		opts.TLS = true
		opts.CertFile = certFile
		opts.KeyFile = keyFile
		return nil
	}
}
