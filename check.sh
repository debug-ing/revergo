# check for code
go fmt ./...

# go vet is the official Go static analysis tool that checks for common errors in code.
go vet ./...

# golint is a linter for Go source code.
golint ./...

# golangci-lint is a fast Go linters runner. It runs linters in parallel, uses caching, supports yaml config, has integrations with all major IDE and has dozens of linters included.
golangci-lint run

