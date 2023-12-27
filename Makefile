# MIT License
# 
# Copyright (c) 2023 Juan Medina
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTOOL=$(GOCMD) tool
GOFORMAT=$(GOCMD) fmt
GORUN=$(GOCMD) run
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod
BUILD_DIR=bin

ifeq ($(OS),Windows_NT)
	BINARY_NAME=$(BUILD_DIR)/twitch-clips.exe
else
	BINARY_NAME=$(BUILD_DIR)/twitch-clips
endif	

APP_PATH="./cmd/twitch-clips"

#ARGS

default: build

build: clean
	$(GOBUILD) -o $(BINARY_NAME) -v $(APP_PATH)
vet:
	$(GOVET) "./..."
clean:
	$(GOCLEAN) $(APP_PATH)
format:
	$(GOFORMAT) "./..."
run: build
	./$(BINARY_NAME) -client-id $(CLIENT_ID) -channel $(CHANNEL)
update:
	$(GOGET) -u all
	$(GOMOD) tidy