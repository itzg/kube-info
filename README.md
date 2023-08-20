[![GitHub release](https://img.shields.io/github/release/itzg/prompt-info.svg)](https://github.com/itzg/prompt-info/releases)

A little utility that can be used to populate command prompts. The kube and git modules avoid calling any external tools, but instead efficiently process the source of information directly.

## Installation

### Scoop

```
scoop bucket add itzg https://github.com/itzg/scoop-bucket.git
scoop install prompt-info
```

## Usage

```text
  -modules value
        Comma separated list of modules to enable. Can be
        - dir
        - git
        - kube
         (default dir,git,kube)
  -template string
        Go template that can reference
        - KubeNamespace
        - KubeContext
        - CurrentDirectory
        - CurrentDirectoryCompact
        - GitBranch
         (default "{{.KubeNamespace}}@{{.KubeContext}}")
```

## Examples

### Using as PowerShell prompt

In `$PROFILE` add:

```ps1
function prompt {
    prompt-info -template "{{.CurrentDirectoryCompact}} [{{.KubeNamespace}}@{{.KubeContext}}]{{if .GitBranch}} {{.GitBranch}}{{end}} > "
}
```
