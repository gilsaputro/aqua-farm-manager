# aqua-farm-manager
This repository contains the prototype for an Aquafarm Management Application.

## High-Level Design (HLD)
For a detailed understanding of the system architecture and design, please refer to the [High-Level Design (HLD)](https://github.com/gilsaputro/aqua-farm-manager/wiki) document.

## API Documentation
The API documentation can be found in the [API Documentation](https://github.com/gilsaputro/aqua-farm-manager/wiki/Farm-Create) file in the repository. This file contains information on the endpoints, request and response formats, and any necessary details.

## Getting Started
These intruction will get you a project and how to run the binary on your local machine.

### Prerequsites
The AquaFarm management system requires Go 1.19 or higher and Docker installed on the local machine in order to run the binary.

#### Docker
You need to have docker installed in your machine.
Follow this step if you don't have docker on your machine :
- Download the Docker CE (Community Edition) package from the Docker website (https://www.docker.com/products/docker-desktop).
- Install the package by following the instructions provided during the installation process.
- Once the installation is complete, verify that Docker has been installed correctly by running the following command in your terminal: "docker run hello-world".

#### Go Programming Language
You need to have golang 1.19 installed in your machine.
Follow this step if you don't have golang 1.19 on your machine :
- Download the Go 1.19 binary package from the official Go website (https://golang.org/dl/).
- Install the package by following the instructions provided during the installation process.
- Once the installation is complete, verify that Go has been installed correctly by running the following command in your terminal: "go version".

## How to run locally
### Clone Repository:
Once you have all the prerequisites properly installed, you can start by cloning this repository.
- To clone the repository, run the following command in your terminal:
```
git clone https://github.com/gilsaputro/aqua-farm-manager.git
```
- To navigate to the repository directory, run the following command in your terminal:
```
cd aqua-farm-manager
```
Note: These instructions assume that you have Git installed on your machine. If you don't have Git installed, you can follow the instructions on the Git website to install it.

### Docker Setup:
To run the AquaFarm management system binary correctly, it is necessary to connect it with the related dependencies. This can be done simply by executing the following command: 

```azure
make deps-init
```

The deps-init command will perform the following actions:
- Build Vault and store secrets
- Build Redis and verify that it is running
- Build Postgres and verify that it is running
- Build NSQ and create a topic for the aqua_farm_tracking_event

To stop the dependencies, run :
```azure
make deps-tear
```

### Running Binary:
Once you have cloned the repository and set up the docker dependencies, you can run the binary using either of the following methods:

Run vendor to download package dependencies

```
go mod vendor
```

And run using :

```
go run ./cmd/aqua-farm-manager/main.go
```

or 

```
go build ./cmd/aqua-farm-manager/
./aqua-farm-manager
```

### Postman Collection
You can import postman collection in Repo File with Name : 
```
Aquafarm Management System.postman_collection.json
```
Or import from [this](https://github.com/gilsaputro/aqua-farm-manager/wiki/Postman-Collection)

Note: The details mentioned in these steps may vary depending on your configuration.
