package envtest

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Manager interface {
	SetEnv(key, value string)
	PathIfExists() []string
	Cleanup()
	New(t *testing.T)
}

type SystemEnvManager struct {
	t       *testing.T
	setEnvs map[string]string
}

func (sem *SystemEnvManager) SetEnv(key, value string) {
	require.NotEmpty(sem.t, key)
	require.NotEmpty(sem.t, value)

	if sem.setEnvs == nil {
		sem.setEnvs = make(map[string]string)
	}

	sem.setEnvs[key] = value

	os.Setenv(key, value)
}

func (sem *SystemEnvManager) PathIfExists() []string {
	return nil
}

func (sem *SystemEnvManager) Cleanup() {
	for key := range sem.setEnvs {
		require.NoError(sem.t, os.Unsetenv(key))
	}

	sem.setEnvs = make(map[string]string)
}

func (sem *SystemEnvManager) New(t *testing.T) {
	t.Cleanup(sem.Cleanup)
	sem.t = t
}

type FileEnvManager struct {
	t        *testing.T
	fileName string
}

func (fem *FileEnvManager) SetEnv(key, value string) {
	fem.t.Helper()

	if fem.fileName == "" {
		file, err := os.CreateTemp("", "cfg-env-*.env")
		fem.fileName = file.Name()
		require.NoError(fem.t, err)
	}

	content, err := os.ReadFile(fem.fileName)
	require.NoError(fem.t, err)
	content = append(content, []byte("\n"+key+"="+value)...)

	require.NoError(fem.t, os.WriteFile(fem.fileName, content, 0644))
}

func (fem *FileEnvManager) PathIfExists() []string {
	fem.t.Helper()

	require.NotEmpty(fem.t, fem.fileName)
	return []string{fem.fileName}
}

func (fem *FileEnvManager) Cleanup() {
	fem.t.Helper()
	os.Remove(fem.fileName)
	fem.fileName = ""
}

func (fem *FileEnvManager) New(t *testing.T) {
	fem.t = t
	t.Cleanup(fem.Cleanup)
}

func Test[T any](t *testing.T, loadFunc func(path ...string) (T, error), envs Inputs) {

	t.Run("returns an error when env file does not exist", func(t *testing.T) {
		t.Parallel()
		cfg, err := loadFunc(uuid.New().String())
		assert.Nil(t, cfg)
		assert.Error(t, err)
	})

	testConfig[T](loadFunc, envs)
}

func testConfig[T any](loadFunc func(path ...string) (T, error), envs Inputs) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("env file", func(t *testing.T) {
			t.Parallel()
			testEnvs[T](t, loadFunc, envs, &FileEnvManager{})
		})

		t.Run("system envs", func(t *testing.T) {
			testEnvs[T](t, loadFunc, envs, &SystemEnvManager{})
		})
	}
}

func testEnvs[T any](t *testing.T, loadFunc func(path ...string) (T, error), envs Inputs, envManager Manager) {
	for testKey, testEnv := range envs {
		if testEnv.Required() {
			t.Run("returns an error when "+testKey+" is missing", func(t *testing.T) {
				envManager.New(t)

				for key, env := range envs {
					if key != testKey {
						envManager.SetEnv(key, env.GetValue())
					}
				}

				cfg, err := loadFunc(envManager.PathIfExists()...)
				assert.Nil(t, cfg)
				assert.Error(t, err)
			})
		}

		if len(testEnv.InvalidValues()) != 0 {
			t.Run("returns an error when "+testKey+" is invalid", func(t *testing.T) {
				envManager.New(t)

				for _, invalidGetValue := range testEnv.InvalidValues() {
					for key, env := range envs {
						if key != testKey {
							envManager.SetEnv(key, env.GetValue())
						} else {
							envManager.SetEnv(key, invalidGetValue)
						}
					}

					cfg, err := loadFunc(envManager.PathIfExists()...)
					assert.Nil(t, cfg)
					assert.Error(t, err)

					envManager.Cleanup()
				}
			})
		}
	}

	t.Run("succeeds when envs are set correctly", func(t *testing.T) {
		envManager.New(t)

		for key, env := range envs {
			envManager.SetEnv(key, env.GetValue())
		}

		cfg, err := loadFunc(envManager.PathIfExists()...)
		assert.NotNil(t, cfg)
		assert.NoError(t, err)
	})
}
