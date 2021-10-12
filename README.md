# hackathon-aks-upgrade-doctor
Hackathon project with intent to help based on heuristics for aks cluster upgrades.


## Development

```
$ make

Usage:
  make <target>

General
  help             Display this help.

Build
  build            Build the binary

Developement
  fmt              Run go fmt against code.
  vet              Run go vet against code.
  test             Run unit tests.
```

### Build binary

```
$ make build
go build -o bin/aks-doctor ./aks/upgrade
# now we can run the binary...
$ ./bin/aks-doctor
Usage:
  aks-doctor [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  demo        demo for pdb
  help        Help about any command

Flags:
  -h, --help   help for aks-doctor

Use "aks-doctor [command] --help" for more information about a command.
```