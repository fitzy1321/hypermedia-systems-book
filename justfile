default:
    just --list

run:
    air

install_deps:
    go install github.com/air-verse/air@latest
    go get .
    go mod tidy

update:
    go get -u ./...
    go mod tidy
