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
		value:    value,
		required: true,
	}
}

func Int(value string) Input {
	return &intEnv{
		value:    value,
		required: true,
	}
}

func Duration(value string) Input {
	return &durationEnv{
		value:    value,
		required: true,
	}
}

func Durations(value string) Input {
	return &durationsEnv{
		value:    value,
		required: true,
	}
}

func Urls(value string) Input {
	return &urlsEnv{
		value:    value,
		required: true,
	}
}

func Url(value string) Input {
	return &urlEnv{
		value:    value,
		required: true,
	}
}

func NatsCreds(path string) Input {
	return &natsCredsEnv{
		value:    path,
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
		cert:     cert,
		key:      key,
		required: true,
	}
}

// CertKey: (key, cert)
func CertKey(key string, cert string) Input {
	return &certKeyEnv{
		key:      key,
		cert:     cert,
		required: true,
	}
}

// CaCert: (cert)
func CaCert(value string) Input {
	return &caCertEnv{
		value:    value,
		required: true,
	}
}

func Enum(testvalue string, allowedvalues ...string) Input {
	return &enum{
		value:   testvalue,
		allowed: allowedvalues,
	}
}
