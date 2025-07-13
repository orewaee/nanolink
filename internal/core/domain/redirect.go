package domain

import (
	"errors"
	"fmt"
	"regexp"
)

type RedirectOpts struct {
	// Redirect service host
	Host string

	// Redirect service port
	Port int

	// Display an HTML template instead of an empty 404 error
	NotFoundTempl bool

	TLS      bool
	CertFile string
	KeyFile  string
}

var DefaultRedirectOpts = RedirectOpts{
	Host:          "127.0.0.1",
	Port:          2222,
	TLS:           false,
	NotFoundTempl: true,
}

type RedirectOpt func(*RedirectOpts) error

const ipv6Regex = `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
const ipv4Regex = `^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`

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
