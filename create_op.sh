#!/bin/sh

# Check if Go is installed
if ! command -v go >/dev/null 2>&1; then
    echo "Go is not installed. Please install Go from https://go.dev/dl/"
    exit 1
fi

echo "Go found: $(go version)"

# Check if op is installed, install if not
if ! command -v op >/dev/null 2>&1; then
    echo "Installing op..."
    go install lesiw.io/op@latest
else
    echo "op is already installed."
fi

# Create .ops directory structure
mkdir -p .ops/commands
mkdir -p .ops/composite

# Create go.mod if it doesn't exist
if [ ! -f .ops/go.mod ]; then
    cat > .ops/go.mod << 'EOF'
module ops

go 1.24
EOF
    echo "Created .ops/go.mod"
else
    echo ".ops/go.mod already exists, skipping."
fi

# Create main.go if it doesn't exist
if [ ! -f .ops/main.go ]; then
    cat > .ops/main.go << 'EOF'
package main

import (
	"lesiw.io/ops"
	"ops/composite"
)

func main() {
	ops.Handle(composite.Ops{})
}
EOF
    echo "Created .ops/main.go"
else
    echo ".ops/main.go already exists, skipping."
fi

# Create commands/ops.go if it doesn't exist
if [ ! -f .ops/commands/ops.go ]; then
    cat > .ops/commands/ops.go << 'EOF'
package commands

type Ops struct{}
EOF
    echo "Created .ops/commands/ops.go"
else
    echo ".ops/commands/ops.go already exists, skipping."
fi

# Create commands/hello.go if it doesn't exist
if [ ! -f .ops/commands/hello.go ]; then
    cat > .ops/commands/hello.go << 'EOF'
package commands

import "fmt"

func (Ops) Hello() error {
	fmt.Println("Hello from op!")
	return nil
}
EOF
    echo "Created .ops/commands/hello.go"
else
    echo ".ops/commands/hello.go already exists, skipping."
fi

# Create composite/ops.go if it doesn't exist
if [ ! -f .ops/composite/ops.go ]; then
    cat > .ops/composite/ops.go << 'EOF'
package composite

import "ops/commands"

type Ops struct {
	commands.Ops
}
EOF
    echo "Created .ops/composite/ops.go"
else
    echo ".ops/composite/ops.go already exists, skipping."
fi

# Run go mod tidy to fetch dependencies
echo "Running go mod tidy..."
cd .ops && go mod tidy

echo "Setup complete! Run 'op hello' to test."
