package main

import (
	"flag"
	"github.com/itzg/go-flagsfiller"
	"gopkg.in/yaml.v3"
	"html/template"
	"log"
	"os"
	"strings"
)

type Args struct {
	Template string `default:"{{.KubeNamespace}}@{{.KubeContext}}"`
}

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

type TemplateContext struct {
	CurrentDirectory        string
	CompactCurrentDirectory string

	KubeNamespace string
	KubeContext   string
}

func main() {
	var args Args
	filler := flagsfiller.New()
	err := filler.Fill(flag.CommandLine, &args)
	if err != nil {
		log.Fatal(err)
	}
	flag.Parse()

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

	var templateContext TemplateContext
	templateContext.KubeNamespace = "default"

	if kubeConfig.CurrentContext != "" {
		templateContext.KubeContext = kubeConfig.CurrentContext
		for _, context := range kubeConfig.Contexts {
			if context.Name == kubeConfig.CurrentContext {
				if context.Details.Namespace != "" {
					templateContext.KubeNamespace = context.Details.Namespace
				}
			}
		}
	}

	t, err := template.New("out").Parse(args.Template)
	if err != nil {
		log.Fatal("Parsing output template", err)
	}

	templateContext.CurrentDirectory, err = os.Getwd()
	if err != nil {
		log.Fatal("Getting current working directory")
	}
	templateContext.CompactCurrentDirectory = buildCompactCurrentDirectory(templateContext.CurrentDirectory)

	err = t.Execute(os.Stdout, templateContext)
	if err != nil {
		log.Fatal("Executing template", err)
	}
}

func buildCompactCurrentDirectory(directory string) (result string) {
	homeDir := getHomeDir()
	result = strings.Replace(directory, homeDir, "~", 1)
	return
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
