package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/itzg/go-flagsfiller"
	"github.com/itzg/zapconfigs"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Args struct {
	Template string   `usage:"Go template that can reference\n- KubeNamespace\n- KubeContext\n- CurrentDirectory\n- CurrentDirectoryCompact\n- GitBranch\n" default:"{{.KubeNamespace}}@{{.KubeContext}}"`
	Modules  []string `usage:"Comma separated list of modules to enable. Can be\n- dir\n- git\n- kube\n" default:"dir,git,kube"`
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
	CurrentDirectoryCompact string

	KubeNamespace string
	KubeContext   string

	GitBranch string
}

func main() {
	var args Args
	err := flagsfiller.Parse(&args)
	if err != nil {
		log.Fatal(err)
	}
	flag.Parse()

	logger := zapconfigs.NewDefaultLogger()
	//goland:noinspection GoUnhandledErrorResult
	defer logger.Sync()

	var templateContext TemplateContext

	t, err := template.New("out").Parse(args.Template)
	if err != nil {
		logger.Fatal("Invalid output template", zap.Error(err), zap.String("template", args.Template))
	}

	if moduleEnabled(args.Modules, "kube") {
		err = loadKubeInfo(&templateContext)
		if err != nil {
			logger.Fatal("Failed loading kube config", zap.Error(err))
		}
	}

	if moduleEnabled(args.Modules, "dir") {
		err = loadCurrentDirectory(&templateContext)
		if err != nil {
			logger.Fatal("Failed loading directory info", zap.Error(err))
		}
	}

	if moduleEnabled(args.Modules, "git") {
		err = loadGitInfo(templateContext.CurrentDirectory, &templateContext)
		if err != nil {
			logger.Fatal("Failed loading git info", zap.Error(err))
		}
	}

	err = t.Execute(os.Stdout, templateContext)
	if err != nil {
		logger.Fatal("Executing template", zap.Error(err))
	}
}

func moduleEnabled(enabledModules []string, module string) bool {
	for _, s := range enabledModules {
		if s == module {
			return true
		}
	}
	return false
}

func loadCurrentDirectory(templateContext *TemplateContext) error {
	var err error
	templateContext.CurrentDirectory, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	templateContext.CurrentDirectoryCompact = buildCompactCurrentDirectory(templateContext.CurrentDirectory)

	return nil
}

func loadKubeInfo(templateContext *TemplateContext) error {
	kubeConfigFile, err := os.Open(getHomeDir() + "/.kube/config")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open kube config file: %w", err)
	}

	var kubeConfig KubeConfig
	decoder := yaml.NewDecoder(kubeConfigFile)
	err = decoder.Decode(&kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to decode kube config: %w", err)
	}

	err = kubeConfigFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close kube config: %w", err)
	}

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

	return nil
}

func loadGitInfo(currentDirectory string, templateContext *TemplateContext) error {

	repo, err := findRepoDir(currentDirectory)
	if err != nil {
		return fmt.Errorf("failed to open repo: %w", err)
	} else if repo != nil {
		head, err := repo.Head()
		if err != nil {
			return fmt.Errorf("failed to get repo head: %w", err)
		}
		templateContext.GitBranch = head.Name().Short()
	}
	return nil
}

func findRepoDir(currentDirectory string) (*git.Repository, error) {
	dir := currentDirectory
	for {
		repo, err := git.PlainOpen(dir)
		if err != nil {
			if !errors.Is(err, git.ErrRepositoryNotExists) {
				return nil, fmt.Errorf("failed to open repo: %w", err)
			}
		} else {
			return repo, nil
		}

		dir = filepath.Dir(dir)
		// only root path ends with separator, see https://pkg.go.dev/path/filepath@go1.21.0#Dir
		if strings.HasSuffix(dir, string(filepath.Separator)) {
			return nil, nil
		}
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
