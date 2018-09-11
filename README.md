# config

Manage env vars with a config.json file


# Quick start

    go get github.com/mozey/config

Set `APP_DIR`, change command below to use your project path

    export APP_DIR=${GOPATH}/src/github.com/mozey/config

Compile

    cd ${APP_DIR}

    go build \
    -ldflags "-X main.AppDir=${APP_DIR}" \
    -o ${APP_DIR}/config ./cmd/config
    
Create `config.dev.json` and set a key
                        
    cd ${APP_DIR}
    
    cp config.dev.sample.json config.dev.json
    
    ./config -key APP_FOO -value xxx
    
    cat config.dev.json
    
Set env from config

    export APP_NOT_IN_CONFIG_FILE=undefined
    
    # Print commands
    ./config

    # Set env    
    eval "$(./config)"
    
    # Print env
    printenv | sort | grep -E "APP_"
    
    
# Testing

    cd ${GOPATH}/src/github.com/mozey/config

    export APP_DEBUG=true
    gotest -v ./...
    
Debug

    go run -ldflags "-X main.AppDir=${APP_DIR}" cmd/config/main.go
    
    
# Prod env

Create `config.prod.json` and set a key

    cp config.prod.sample.json config.prod.json
    
    ./config -env prod -key APP_BEER -value pilsner
    
    cat config.prod.json
    
All config files must have the same keys,
if a key does not apply set the value to an empty string.
Compare config files and print un-matched keys

    ./config -env dev -compare prod
    
    # Config exits with error code if the keys don't match
    echo $?


# Generate config helper

    mkdir -p internal/config
    
    ./config -env prod -generate internal/config
    
    go fmt ./internal/config/config.go


# Dry run

For commands that update files,
use the `-dry-run` flag to print the result and skip the update


# Aliases

Use aliases to quickly toggle between env and print the current env.
For example

    alias dev='eval "$(./config)" && printenv | sort | grep -E "APP_|AWS_"'
    alias prod='eval "$(./config -env prod)" && printenv | sort | grep -E "APP_|AWS_"'

