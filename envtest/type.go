package envtest

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"net/url"

	"github.com/nats-io/nkeys"
)

// ==========================================================
// STRING ENV
// ==========================================================

type stringEnv struct {
	Value    string
	required bool
}

func (env *stringEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env stringEnv) GetValue() string {
	return env.Value
}

func (env stringEnv) ExpectedValue() any {
	return env
}

func (stringEnv) InvalidValues() []string {
	return nil
}

func (env stringEnv) Required() bool {
	return env.required
}

// ==========================================================
// INT ENV
// ==========================================================

type intEnv struct {
	Value    string
	required bool
}

func (env *intEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env intEnv) GetValue() string {
	return env.Value
}

func (env intEnv) ExpectedValue() any {
	if _, err := strconv.Atoi(env.Value); err != nil {
		log.Fatal(err)
	}
	return env.Value
}

func (intEnv) InvalidValues() []string {
	return []string{"false", "true", "string", "0.1"}
}

func (env intEnv) Required() bool {
	return env.required
}

// ==========================================================
// DURATION ENV
// ==========================================================

type durationEnv struct {
	Value    string
	required bool
}

func (env *durationEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env durationEnv) GetValue() string {
	return env.Value
}

func (env durationEnv) ExpectedValue() any {
	_, err := time.ParseDuration(env.Value)
	if err != nil {
		log.Fatal(err)
	}
	return env.Value
}

func (durationEnv) InvalidValues() []string {
	return []string{"false", "true", "string", "0.1", "1"}
}

func (env durationEnv) Required() bool {
	return env.required
}

type durationsEnv struct {
	Value    string
	required bool
}

func (env *durationsEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env durationsEnv) GetValue() string {
	return env.Value
}

func (env durationsEnv) ExpectedValue() any {
	for _, value := range strings.Split(strings.TrimSpace(env.Value), ", ") {
		_, err := time.ParseDuration(value)
		if err != nil {
			log.Fatal(err)
		}
	}
	return env.Value
}

func (durationsEnv) InvalidValues() []string {
	return []string{"false", "true", "string", "0.1", "1"}
}

func (env durationsEnv) Required() bool {
	return env.required
}

// ==========================================================
// URLS ENV (multi)
// ==========================================================

type urlsEnv struct {
	Value    string
	required bool
}

func (env *urlsEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env urlsEnv) GetValue() string {
	return env.Value
}

func (env urlsEnv) ExpectedValue() any {
	for _, envUrl := range strings.Split(strings.TrimSpace(env.Value), ",") {
		if _, err := url.Parse(envUrl); err != nil {
			log.Fatal(err)
		}
	}
	return env.Value
}

func (urlsEnv) InvalidValues() []string {
	return []string{
		"http:/invalid#!%&)*#!",
		"http://one.com http://another-without-a.coma",
	}
}

func (env urlsEnv) Required() bool {
	return env.required
}

// ==========================================================
// URL ENV (single)
// ==========================================================

type urlEnv struct {
	Value    string
	required bool
}

func (env *urlEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env urlEnv) GetValue() string {
	return env.Value
}

func (env urlEnv) ExpectedValue() any {
	if _, err := url.Parse(env.Value); err != nil {
		log.Fatal(err)
	}
	return env.Value
}

func (urlEnv) InvalidValues() []string {
	return []string{
		"http:/invalid#!%&)*#!",
		"http://one.com http://another-without-a.coma",
	}
}

func (env urlEnv) Required() bool {
	return env.required
}

// ==========================================================
// CERT ENV
// ==========================================================

type certEnv struct {
	Cert     string
	Key      string
	required bool
}

func (env *certEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env certEnv) GetValue() string {
	return env.Cert
}

func (env certEnv) ExpectedValue() any {
	_, err := tls.LoadX509KeyPair(env.Cert, env.Key)
	if err != nil {
		log.Fatal(err)
	}
	return env.Cert
}

func (certEnv) InvalidValues() []string {
	return []string{"./non-existing.pem"}
}

func (env certEnv) Required() bool {
	return env.required
}

// ==========================================================
// CERT KEY ENV
// ==========================================================

type certKeyEnv struct {
	Key      string
	Cert     string
	required bool
}

func (env *certKeyEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env certKeyEnv) GetValue() string {
	return env.Cert
}

func (env certKeyEnv) ExpectedValue() any {
	_, err := tls.LoadX509KeyPair(env.Cert, env.Key)
	if err != nil {
		log.Fatal(err)
	}
	return env.Key
}

func (certKeyEnv) InvalidValues() []string {
	return []string{"./non-existing.key"}
}

func (env certKeyEnv) Required() bool {
	return env.required
}

// ==========================================================
// CA CERT ENV
// ==========================================================

type caCertEnv struct {
	Value    string
	required bool
}

func (env *caCertEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env caCertEnv) GetValue() string {
	return env.Value
}

func (env caCertEnv) ExpectedValue() any {
	raw, err := os.ReadFile(env.Value)
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(raw)
	if block == nil || block.Type != "CERTIFICATE" {
		log.Fatal("invalid certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	if !cert.IsCA {
		log.Fatal("certificate is not a Certificate Authority")
	}

	return env.Value
}

func (caCertEnv) InvalidValues() []string {
	return []string{"./non-existing.pem"}
}

func (env caCertEnv) Required() bool {
	return env.required
}

// ==========================================================
// NATS CREDS ENV
// ==========================================================

type natsCredsEnv struct {
	Value    string
	required bool
}

func (env *natsCredsEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env natsCredsEnv) GetValue() string {
	return env.Value
}

func (env natsCredsEnv) ExpectedValue() any {
	data, err := os.ReadFile(env.Value)
	if err != nil {
		return err
	}

	_, err = nkeys.ParseDecoratedJWT(data)
	if err != nil {
		return fmt.Errorf("invalid creds format: %w", err)
	}

	return env.Value
}

func (natsCredsEnv) InvalidValues() []string {
	return []string{"./non-existing.file"}
}

func (env natsCredsEnv) Required() bool {
	return env.required
}
