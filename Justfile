set shell := ["bash", "-u", "-c"]

export scripts := ".github/workflows/scripts"
export GOBIN := `echo $PWD/.bin`

# print available commands
default:
    @just --list

# start nerfs-builds artifact builder in demo mode
run mode: compile
    $GOBIN/nerfs-builds {{mode}}

# compile the nerfs-builds executable
compile: tidy
    cd cmds/nerfs-builds && go install

# tidy up Go modules
tidy:
    go mod tidy

# vet the nerfs-compile source tree
vet:
    go vet ./...

# run go test on the source tree
test:
    go test -race ./...

# lint the nerfs-compile source tree
lint: vet
    $GOBIN/golangci-lint --config $scripts/golangci.yaml run

# show host system information
@sysinfo:
    echo "{{os()/arch()}} {{num_cpus()}}c"

# locally install build dependencies
init:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2
