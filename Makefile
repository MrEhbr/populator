SHELL := /usr/bin/env bash -o pipefail
GOPKG ?= github.com/MrEhbr/populator
DOCKER_IMAGE ?=	mrehbr/populator
GOBINS ?= cmd/populator
GO_APP ?= populator
PROTO_PATH := .
PROTOC_GEN_GO_OUT := .
PROTOC_GEN_GO_OPT := plugins=grpc

include rules.mk
