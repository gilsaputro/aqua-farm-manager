# aqua-farm-manager
A prototype repository for the development of an aquafarm management application.

## Getting Started
These intruction will get you a project and how to run the binary on your local machine.

### Prerequsites
The AquaFarm management system requires Go 1.19 or higher and Docker installed on the local machine in order to run the binary.

#### Docker
You need to have docker installed in your machine.
Follow this step if you don't have docker on your machine :
    a. Download the Docker CE (Community Edition) package from the Docker website (https://www.docker.com/products/docker-desktop).
    b. Install the package by following the instructions provided during the installation process.
    c. Once the installation is complete, verify that Docker has been installed correctly by running the following command in your terminal: "docker run hello-world".

#### Go Programming Language
You need to have golang 1.19 installed in your machine.
Follow this step if you don't have golang 1.19 on your machine :
    a. Download the Go 1.19 binary package from the official Go website (https://golang.org/dl/).
    b. Install the package by following the instructions provided during the installation process.
    c. Once the installation is complete, verify that Go has been installed correctly by running the following command in your terminal: "go version".

## How to run locally
### Building:
Once you have all the prerequisites properly installed, you can start by cloning this repository.
    a. Clone the Repository: Use "git clone" to clone the AquaFarm management repository to the directory you created in step 1a.
    b. Change Directory: Change to the directory of the cloned repository by using the "cd" command.
### Docker Setup:
To run the AquaFarm management system binary correctly, it is necessary to connect it with the related dependencies. This can be done simply by executing the following command: 

```azure
make deps-init
```

The deps-init command will perform the following actions:
    a. Build Vault and store secrets
    b. Build Redis and verify that it is running
    c. Build Postgres and verify that it is running
    d. Build NSQ and create a topic for the aqua_farm_tracking_event."

To stop the dependencies, run :
```azure
make deps-tear
```

### Running Binary:
Once you have cloned the repository and set up the dependencies, you can run the binary using either of the following methods:
```
go run ./cmd/aqua-farm-manager/main.go
```

or 

```
go build ./cmd/aqua-farm-manager/
./aqua-farm-manager
```

Note: The details mentioned in these steps may vary depending on your configuration.