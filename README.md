[![GitHub release](https://img.shields.io/github/release/itzg/kube-info.svg)](https://github.com/itzg/kube-info/releases)

## Installation

### Scoop

```
scoop bucket add itzg https://github.com/itzg/scoop-bucket.git
scoop install kube-info
```

## Usage

```text
  -template string
        Go template that can reference
        - KubeNamespace
        - KubeContext
        - CurrentDirectory
        - CompactCurrentDirectory
        - GitBranch
         (default "{{.KubeNamespace}}@{{.KubeContext}}")
```

## Examples

### Using as PowerShell prompt

In `$PROFILE` add:

```ps1
function prompt {
    kube-info -template "{{.CompactCurrentDirectory}} [{{.KubeNamespace}}@{{.KubeContext}}]{{if .GitBranch}} {{.GitBranch}}{{end}} > "
}
```
