package gconf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"abodemine/lib/errors"
)

type ConfigFileCtxKey string

// ResolveConfig resolves the config file path and parses it into the given value.
func ResolveConfig(v any, configPath string, envKey string) error {
	configPath = strings.TrimSpace(configPath)

	if configPath == "" {
		envKey = strings.TrimSpace(envKey)

		if envKey == "" {
			return &errors.Object{
				Id:     "13ae1af8-1cab-4792-bb27-09275c5615d1",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to resolve config file path.",
			}
		}

		configPath = os.Getenv(envKey)
	}

	// This should indicate that the config was defined elsewhere.
	if configPath == "-" {
		return nil
	}

	fileBody, err := os.ReadFile(configPath)
	if err != nil {
		return &errors.Object{
			Id:     "d57b1b31-3c58-453c-8ef1-9227765e0b45",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Failed to read config file.",
			Cause:  err.Error(),
		}
	}

	if passStr := os.Getenv(fmt.Sprintf("%s_ENC_PASS", envKey)); passStr != "" {
		iterStr := os.Getenv(fmt.Sprintf("%s_ENC_ITER", envKey))

		iter, err := strconv.Atoi(iterStr)
		if err != nil {
			return &errors.Object{
				Id:     "a85d39f6-d106-43c7-a9e7-b4cd927d84a9",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse iteration count.",
				Cause:  err.Error(),
				Meta: map[string]any{
					"iter": iterStr,
				},
			}
		}

		decryptedBody, err := DecryptConfig(bytes.NewBuffer(fileBody), passStr, iter)
		if err != nil {
			return errors.Forward(err, "3aab8a04-f594-456a-9a53-7b9afc4f4227")
		}

		fileBody = decryptedBody
	}

	expandedBody := []byte(os.ExpandEnv(string(fileBody)))

	switch strings.ToLower(filepath.Ext(configPath)) {
	case ".json":
		if err := json.Unmarshal(expandedBody, v); err != nil {
			return &errors.Object{
				Id:     "6e3cb928-8e31-4951-bfec-83b01cbf1eac",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse json config file.",
				Cause:  err.Error(),
			}
		}
	case ".toml":
		if err := toml.Unmarshal(expandedBody, v); err != nil {
			return &errors.Object{
				Id:     "7c1a4c58-a4ae-46bc-9005-44c7fcbef2d3",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse toml config file.",
				Cause:  err.Error(),
			}
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(expandedBody, v); err != nil {
			return &errors.Object{
				Id:     "0a3d3a29-ba73-49e7-9466-28e1bb1c5499",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse yaml config file.",
				Cause:  err.Error(),
			}
		}
	}

	return nil
}

// pass must be in the form of a base64 encoded string.
func DecryptConfig(reader io.Reader, pass string, iter int) ([]byte, error) {
	cmd := exec.Command(
		"openssl",
		"enc",
		"-aes-256-cbc",
		"-pbkdf2",
		"-d",
		"-iter", strconv.Itoa(iter),
		"-pass", "pass:"+pass,
	)

	out := new(bytes.Buffer)

	cmd.Stdin = reader
	cmd.Stdout = out

	if err := cmd.Run(); err != nil {
		return nil, &errors.Object{
			Id:     "845a74a2-5ab3-41e4-a5c2-583798a98af4",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to decrypt config file.",
			Cause:  err.Error(),
		}
	}

	return out.Bytes(), nil
}
