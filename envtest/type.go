package envtest

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"net/url"

	"github.com/nats-io/nkeys"
)

//
// ======================
// Base Struct
// ======================
//

type baseStruct struct {
	required bool
}

func (b baseStruct) isRequired() bool {
	return b.required
}

func valueValidation(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("value cannot be empty")
	}
	return nil
}

//
// ======================
// STRING ENV
// ======================
//

type stringEnv struct {
	baseStruct
	value string
}

func (env *stringEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env stringEnv) getValues() []string {
	return []string{env.value}
}

func (stringEnv) invalidValues() []string {
	return nil
}

func (env stringEnv) validate() error {
	return valueValidation(env.value)
}

//
// ======================
// INT ENV
// ======================
//

type intEnv struct {
	baseStruct
	value string
}

func (env *intEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env intEnv) getValues() []string {
	return []string{env.value}
}

func (intEnv) invalidValues() []string {
	return []string{"false", "true", "string", "0.1"}
}

func (env intEnv) validate() error {
	if err := valueValidation(env.value); err != nil {
		return err
	}
	_, err := strconv.Atoi(env.value)
	return err
}

//
// ======================
// DURATION ENV
// ======================
//

type durationEnv struct {
	baseStruct
	value string
}

func (env *durationEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env durationEnv) getValues() []string {
	return []string{env.value}
}

func (durationEnv) invalidValues() []string {
	return []string{"false", "true", "string", "0.1", "1"}
}

func (env durationEnv) validate() error {
	if err := valueValidation(env.value); err != nil {
		return err
	}
	_, err := time.ParseDuration(env.value)
	return err
}

//
// ======================
// DURATIONS ENV
// ======================
//

type durationsEnv struct {
	baseStruct
	value string
}

func (env *durationsEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env durationsEnv) getValues() []string {
	return []string{env.value}
}

func (durationsEnv) invalidValues() []string {
	return []string{"false", "true", "string", "0.1", "1"}
}

func (env durationsEnv) validate() error {
	if err := valueValidation(env.value); err != nil {
		return err
	}

	for _, v := range strings.Split(strings.TrimSpace(env.value), ",") {
		if _, err := time.ParseDuration(strings.TrimSpace(v)); err != nil {
			return err
		}
	}
	return nil
}

//
// ======================
// URL ENV
// ======================
//

type urlEnv struct {
	baseStruct
	value string
}

func (env *urlEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env urlEnv) getValues() []string {
	return []string{env.value}
}

func (urlEnv) invalidValues() []string {
	return []string{
		"http:/invalid#!%&)*#!",
		"http://one.com http://another-without-a.coma",
	}
}

func (env urlEnv) validate() error {
	if err := valueValidation(env.value); err != nil {
		return err
	}
	_, err := url.Parse(env.value)
	return err
}

//
// ======================
// URLS ENV
// ======================
//

type urlsEnv struct {
	baseStruct
	value string
}

func (env *urlsEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env urlsEnv) getValues() []string {
	return []string{env.value}
}

func (urlsEnv) invalidValues() []string {
	return []string{
		"http:/invalid#!%&)*#!",
		"http://one.com http://another-without-a.coma",
	}
}

func (env urlsEnv) validate() error {
	if err := valueValidation(env.value); err != nil {
		return err
	}

	for _, u := range strings.Split(strings.TrimSpace(env.value), ",") {
		if _, err := url.Parse(strings.TrimSpace(u)); err != nil {
			return err
		}
	}
	return nil
}

//
// ======================
// CERT ENV
// ======================
//

type certEnv struct {
	baseStruct
	cert string
	key  string
}

func (env *certEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env certEnv) getValues() []string {
	return []string{env.cert}
}

func (certEnv) invalidValues() []string {
	return []string{"./non-existing.pem"}
}

func (env certEnv) validate() error {
	if err := valueValidation(env.cert); err != nil {
		return err
	}
	if err := valueValidation(env.key); err != nil {
		return err
	}
	_, err := tls.LoadX509KeyPair(env.cert, env.key)
	return err
}

//
// ======================
// CERT KEY ENV
// ======================
//

type certKeyEnv struct {
	baseStruct
	key  string
	cert string
}

func (env *certKeyEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env certKeyEnv) getValues() []string {
	return []string{env.cert}
}

func (certKeyEnv) invalidValues() []string {
	return []string{"./non-existing.key"}
}

func (env certKeyEnv) validate() error {
	if err := valueValidation(env.cert); err != nil {
		return err
	}
	if err := valueValidation(env.key); err != nil {
		return err
	}
	_, err := tls.LoadX509KeyPair(env.cert, env.key)
	return err
}

//
// ======================
// CA CERT ENV
// ======================
//

type caCertEnv struct {
	baseStruct
	value string
}

func (env *caCertEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env caCertEnv) getValues() []string {
	return []string{env.value}
}

func (caCertEnv) invalidValues() []string {
	return []string{"./non-existing.pem"}
}

func (env caCertEnv) validate() error {
	if err := valueValidation(env.value); err != nil {
		return err
	}

	raw, err := os.ReadFile(env.value)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(raw)
	if block == nil || block.Type != "CERTIFICATE" {
		return fmt.Errorf("invalid certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	if !cert.IsCA {
		return fmt.Errorf("certificate is not a CA")
	}

	return nil
}

//
// ======================
// NATS CREDS ENV
// ======================
//

type natsCredsEnv struct {
	baseStruct
	value string
}

func (env *natsCredsEnv) NotRequire() Input {
	env.required = false
	return env
}

func (env natsCredsEnv) getValues() []string {
	return []string{env.value}
}

func (natsCredsEnv) invalidValues() []string {
	return []string{"./non-existing.file"}
}

func (env natsCredsEnv) validate() error {
	if err := valueValidation(env.value); err != nil {
		return err
	}

	data, err := os.ReadFile(env.value)
	if err != nil {
		return err
	}

	_, err = nkeys.ParseDecoratedJWT(data)
	return err
}

//
// ======================
// ENUM ENV
// ======================
//

type enum struct {
	baseStruct
	value   string
	allowed []string
}

func (env *enum) NotRequire() Input {
	env.required = false
	return env
}

func (env enum) getValues() []string {
	return []string{env.value}
}

func (env enum) invalidValues() []string {
	invalid := make([]string, len(env.allowed))
	for i, v := range env.allowed {
		invalid[i] = "not_" + v
	}
	return invalid
}

func (env enum) validate() error {
	if err := valueValidation(env.value); err != nil {
		return err
	}

	for _, v := range env.allowed {
		if env.value == v {
			return nil
		}
	}

	return fmt.Errorf("invalid value: %s", env.value)
}
