set shell := ["bash", "-u", "-c"]

export scripts := ".github/workflows/scripts"
export GOBIN := `echo $PWD/.bin`

# show available commands
[private]
default:
    @just --list

# start nerfs artifact builder in demo mode
run mode: compile
    $GOBIN/nerfs {{mode}}

# compile the nerfs executable
compile: tidy
    cd cmds/nerfs && go install

# tidy up Go modules
tidy:
    go mod tidy

# vet the nerfs source tree
vet:
    go vet ./...

# run go test on the source tree
tests:
    go test -race ./...

# lint the nerfs source tree
lint: vet
    $GOBIN/golangci-lint run --config $scripts/golangci.yaml

# show host system information
@sysinfo:
    echo "{{os()/arch()}} {{num_cpus()}}c"

# locally install build dependencies
init:
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4
