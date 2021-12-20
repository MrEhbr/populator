// +build tools

package tools

import (
	_ "github.com/golang/protobuf/proto"                    // required by rules.mk
	_ "github.com/golang/protobuf/protoc-gen-go"            // required by rules.mk
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint" // required by rules.mk
	_ "github.com/tailscale/depaware"                       // required by rules.mk
	_ "mvdan.cc/gofumpt"
)
