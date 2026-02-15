package envtest

type Input interface {
	GetValue() string
	ExpectedValue() any
	InvalidValues() []string
	NotRequire() Input
	Required() bool
}

type Inputs map[string]Input

func (i Inputs) Join(new Inputs) {
	for k, v := range new {
		i[k] = v
	}
}

// ======================
// Simple value envs
// ======================

func String(value string) Input {
	return &stringEnv{
		Value:    value,
		required: true,
	}
}

func Int(value string) Input {
	return &intEnv{
		Value:    value,
		required: true,
	}
}

func Duration(value string) Input {
	return &durationEnv{
		Value:    value,
		required: true,
	}
}

func Durations(value string) Input {
	return &durationsEnv{
		Value:    value,
		required: true,
	}
}

func Urls(value string) Input {
	return &urlsEnv{
		Value:    value,
		required: true,
	}
}

func Url(value string) Input {
	return &urlEnv{
		Value:    value,
		required: true,
	}
}

func NatsCreds(path string) Input {
	return &natsCredsEnv{
		Value:    path,
		required: true,
	}
}

// ======================
// Structs that require multiple inputs
// (order MUST match struct field order)
// ======================

// Cert: (cert, key)
func Cert(cert string, key string) Input {
	return &certEnv{
		Cert:     cert,
		Key:      key,
		required: true,
	}
}

// CertKey: (key, cert)
func CertKey(key string, cert string) Input {
	return &certKeyEnv{
		Key:      key,
		Cert:     cert,
		required: true,
	}
}

// CaCert: (cert)
func CaCert(value string) Input {
	return &caCertEnv{
		Value:    value,
		required: true,
	}
}
