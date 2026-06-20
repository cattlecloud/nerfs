set shell := ["bash", "-u", "-c"]

export scripts := ".github/workflows/scripts"
export GOBIN := `echo $PWD/.bin`

# show available commands
[private]
default:
    @just --list

# start nerfs artifact builder in demo mode
[group('build')]
run mode: compile
    $GOBIN/nerfs {{mode}}

# compile the nerfs executable
[group('build')]
compile: tidy
    cd cmds/nerfs && go install

# tidy up Go modules
[group('build')]
tidy:
    go mod tidy

# vet the nerfs source tree
[group('lint')]
vet:
    go vet ./...

# run go test on the source tree
[group('testing')]
tests:
    go test -race ./...

# run specific unit test
[group('testing')]
[no-cd]
test unit:
    go test -v -count=1 -race -run {{unit}} 2>/dev/null

# lint the nerfs source tree
[group('lint')]
lint: vet
    $GOBIN/golangci-lint run --config $scripts/golangci.yaml

# show host system information
[group('build')]
@sysinfo:
    echo "{{os()/arch()}} {{num_cpus()}}c"

# locally install build dependencies
[group('build')]
init:
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4

# deploy artifact
[group('deploy')]
deploy host:
    HOST={{host}} envy exec gh-cattle $scripts/deploy.sh nerfs

# create release artifact
[group('deploy')]
artifact:
    gh workflow run build.yaml
