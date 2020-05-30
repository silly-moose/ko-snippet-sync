#!/bin/bash

# Build App
go build -ldflags "-s -w"

# Sign app for Apple
gon gon_config.json 
