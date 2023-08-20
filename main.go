package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type KubeConfig struct {
	Contexts       []Context
	CurrentContext string `yaml:"current-context"`
}

type ContextDetails struct {
	Namespace string
}

type Context struct {
	Name    string
	Details ContextDetails `yaml:"context"`
}

func main() {
	kubeConfigFile, err := os.Open(getHomeDir() + "/.kube/config")
	if err != nil {
		if os.IsNotExist(err) {
			println("<none>")
			os.Exit(0)
		}
		log.Fatal("Reading kube config", err)
	}

	var kubeConfig KubeConfig
	decoder := yaml.NewDecoder(kubeConfigFile)
	err = decoder.Decode(&kubeConfig)
	if err != nil {
		log.Fatal("Decoding kube config", err)
	}

	err = kubeConfigFile.Close()
	if err != nil {
		log.Fatal("Closing kube config", err)
	}

	if kubeConfig.CurrentContext != "" {
		for _, context := range kubeConfig.Contexts {
			if context.Name == kubeConfig.CurrentContext {
				ns := "default"
				if context.Details.Namespace != "" {
					ns = context.Details.Namespace
				}
				fmt.Printf("%s@%s\n", ns, kubeConfig.CurrentContext)
				return
			}
		}
	}
	fmt.Println("<unknown>")
}

func getHomeDir() string {
	for _, v := range []string{"HOME", "USERPROFILE"} {
		value := os.Getenv(v)
		if value != "" {
			return value
		}
	}

	log.Fatal("Unable to resolve home directory")
	return ""
}
