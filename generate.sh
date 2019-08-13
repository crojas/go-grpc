#!/bin/bash

export PATH=$PATH:$GOPATH/bin
protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.