package envutil

import (
	"crypto/x509"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Url string

func (u *Url) UnmarshalText(text []byte) error {
	if _, err := url.Parse(string(text)); err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}

	*u = Url(text)

	return nil
}

func (s Url) Value() string {
	return string(s)
}

type Urls []string

func (us *Urls) UnmarshalText(text []byte) error {
	*us = Urls(strings.Split(strings.TrimSpace(string(text)), ", "))

	for _, u := range *us {
		if _, err := url.Parse(u); err != nil {
			return fmt.Errorf("failed to parse url: %w", err)
		}
	}

	return nil
}

func (s Urls) Value() []string {
	return []string(s)
}

func (s Urls) Join() string {
	return strings.Join(s, ", ")
}

func ToUrls(s ...string) Urls {
	return Urls(s)
}

type File string

func (f *File) UnmarshalText(text []byte) error {
	envFile := strings.TrimSpace(string(text))

	if envFile == "" {
		return nil
	}

	_, err := os.ReadFile(string(envFile))
	if err != nil {
		return fmt.Errorf("failed to read the cert: %w", err)
	}

	*f = File(string(envFile))
	return nil
}

func (f File) IsProvided() bool {
	return f != ""
}

func (f File) Path() string {
	return string(f)
}

func (f File) CertPool() (*x509.CertPool, error) {
	caCert, err := os.ReadFile(string(f))
	if err != nil {
		return nil, fmt.Errorf("failed to read the cert: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("failed to add the cert...")
	}

	return caCertPool, nil
}

type Port string

func (p *Port) UnmarshalText(text []byte) error {

	if _, err := strconv.Atoi(string(text)); err != nil {
		return errors.New("invalid port")
	}

	*p = Port(":" + string(text))

	return nil
}

func (p Port) Value() string {
	return string(p)
}

type Durations []time.Duration

func (d *Durations) UnmarshalText(text []byte) error {
	for _, value := range strings.Split(strings.TrimSpace(string(text)), ", ") {
		dur, err := time.ParseDuration(value)
		if err != nil {
			return err
		}

		*d = append(*d, dur)
	}

	return nil
}

func (d Durations) Value() []time.Duration {
	return []time.Duration(d)
}
