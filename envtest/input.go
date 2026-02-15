package envtest

type Input interface {
	getValues() []string
	invalidValues() []string
	validate() error
	NotRequire() Input
	isRequired() bool
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
		baseStruct: baseStruct{required: true},
		value:      value,
	}
}

func Int(value string) Input {
	return &intEnv{
		baseStruct: baseStruct{required: true},
		value:      value,
	}
}

func Duration(value string) Input {
	return &durationEnv{
		baseStruct: baseStruct{required: true},
		value:      value,
	}
}

func Durations(value string) Input {
	return &durationsEnv{
		baseStruct: baseStruct{required: true},
		value:      value,
	}
}

func Urls(value string) Input {
	return &urlsEnv{
		baseStruct: baseStruct{required: true},
		value:      value,
	}
}

func Url(value string) Input {
	return &urlEnv{
		baseStruct: baseStruct{required: true},
		value:      value,
	}
}

func NatsCreds(path string) Input {
	return &natsCredsEnv{
		baseStruct: baseStruct{required: true},
		value:      path,
	}
}

// ======================
// Structs that require multiple inputs
// ======================

// Cert: (cert, key)
func Cert(cert string, key string) Input {
	return &certEnv{
		baseStruct: baseStruct{required: true},
		cert:       cert,
		key:        key,
	}
}

// CertKey: (key, cert)
func CertKey(key string, cert string) Input {
	return &certKeyEnv{
		baseStruct: baseStruct{required: true},
		key:        key,
		cert:       cert,
	}
}

// CaCert: (cert)
func CaCert(value string) Input {
	return &caCertEnv{
		baseStruct: baseStruct{required: true},
		value:      value,
	}
}

// Enum: (value, allowed...)
func Enum(value string, allowed ...string) Input {
	return &enum{
		baseStruct: baseStruct{required: true},
		value:      value,
		allowed:    allowed,
	}
}
