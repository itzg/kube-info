[![GitHub release](https://img.shields.io/github/release/itzg/kube-info.svg)](https://github.com/itzg/kube-info/releases)

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