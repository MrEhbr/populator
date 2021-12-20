SHELL := /usr/bin/env bash -o pipefail
GOPKG ?= github.com/MrEhbr/golang-repo-template
DOCKER_IMAGE ?=	mrehbr/golang-repo-template
GOBINS ?= cmd/golang-repo-template
GO_APP ?= golang-repo-template
PROTO_PATH := .
PROTOC_GEN_GO_OUT := .
PROTOC_GEN_GO_OPT := plugins=grpc

include rules.mk
