SHELL := bash

DEFAULT_BUILD_LOCATION = /Users/$(USER)/.local/bin/
PROJECT_NAME = dbm-sandbox


build:
	go build -o $(DEFAULT_BUILD_LOCATION) . 

