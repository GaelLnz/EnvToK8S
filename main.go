package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Secret struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   map[string]string `yaml:"metadata"`
	Type       string            `yaml:"type"`
	Data       map[string]string `yaml:"data"`
}

func main() {
	envFile := flag.String("env", ".env", "Path to the .env file")
	outputFile := flag.String("output", "secrets.yaml", "Path to the output YAML file")
	namespace := flag.String("namespace", "", "Kubernetes namespace (optional)")
	name := flag.String("name", "env-secrets", "Name of the Kubernetes Secret")
	flag.Parse()

	secrets, err := readEnvFile(*envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading .env file: %v\n", err)
		os.Exit(1)
	}

	yamlContent, err := createSecretYAML(secrets, *name, *namespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating YAML: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(*outputFile, yamlContent, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing YAML file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Secrets YAML file created: %s\n", *outputFile)
}

func readEnvFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	secrets := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		secrets[key] = base64.StdEncoding.EncodeToString([]byte(value))
	}

	return secrets, scanner.Err()
}

func createSecretYAML(secrets map[string]string, name, namespace string) ([]byte, error) {
	secret := Secret{
		ApiVersion: "v1",
		Kind:       "Secret",
		Metadata: map[string]string{
			"name": name,
		},
		Type: "Opaque",
		Data: secrets,
	}

	if namespace != "" {
		secret.Metadata["namespace"] = namespace
	}

	return yaml.Marshal(secret)
}
