# dbm-sandbox

The goal of this project is to speed up deploying a sandbox that already has a configured Datadog Agent for DBM and a DBMS that can be used for replication or just general use.

## Usage

### Example

The example below uses this tool to create a directory called `sandbox-demo` that will contain a complete Docker Compose manifest that will deploy a Datadog Agent of the latest version, and a MySQL instance running version `8.0.37`. 

This project can be deployed using `docker compose up -d`.

<img alt="dbm-sandbox demo" src="assets/dbm-sandbox.gif" width="600" />

### Requirements

To run this application you will need Go v1.20 or higher installed, this is for building or straight up running the application.

To make use of the DBMS Providers, you may need specific tools installed.

For example to make use of the default Docker DBMS Provider, you will need to have Docker and Docker Compose installed.

You will also need to have your Datadog API Key set in the environment variable `DD_API_KEY`. This is required because when the provider is creating the configuration files for the project, it will need the API Key to inject it into the template.

Below is an example of the error you would get if we can't find the `DD_API_KEY` environment variable:

<img alt="missing api key error" src="assets/missingapikey.gif" width="600" />

### Building

To use the tool as a binary, you will need to build the tool, the simplest way to do this is to run the command:

``` bash
make build
```

This will build the tool and move the `dbm-sandbox` binary into the directory `/Users/$(USER)/.local/bin/`. Note that this location assumes that you are using a Mac, when building feel free to change this default location in the [Makefile](./Makefile) as you see fit.
