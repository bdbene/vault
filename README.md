# Vault

This project is an excercise in Go to develop a restful service that can encrypt and store secrets. The goal for the project was to practice coding in Go. As such, this service is not meant to go into production. Furthermore, certain improvements (with respect to security, performance, best-practices) could obviously be improved, but that is not the goal. The learnings of this project will be used elsewhere instead of perfecting this one.

## Getting Started

Instructions for running locally:

### Prerequisites
Go needs to be installed to compile the code.

To download the dependancies:
```
go get ./...
```

External dependancies:
- github.com/BurntSushi/toml
- github.com/golang/mock/gomock
- github.com/gorilla/mux

### Installation
Either clone this repository or use "go get github.com/bdbene/vault"

cd application/
go build
go install

## Running
By default, running the application will have it listen on port 8080.
```
$ vault
```

## Configurations
Edit the config.tml file to configure the application at start-up. 

## Testing
Generate mocked classes, then run tests:
```
go generate ./...
go test ./...

```

## Areas of Improvement
Areas to improve if I were to perfect this project instead of using it as an exercise. 
### Design
Different parts of the project could obviously be designed using better software engineering principles. For example, the module that handles encryption should have been written as an interface. This would have allow for different encryption implementations, and decrease coupling with the other modules that use it. Furthermore, with lower coupling it would be easier to unit test the different modules. This could be fixed by dependancy-injecting the different interfaces where they are needed. This solution is however implemented with the storage drivers and handler. 

Mutual exlusion to the data store is handled by the Handler, instead of at the data store. This simple solution causes all reads and writes to happen sequentially whitch is horrible for performance. Instead, the service should take advantage of the datastore for concurrent accesses. For example, if there was a storage driver for a SQL database, concurrent reads and writes would be handled by the database itself. Furthermore, independant transactions could be interleaved by the database to improve performance.

Storage drivers could be written as a plugin that can be picked up. This would allow anyone to write their own storage drivers for their datastores, making the service compatible with any database. The only existing storage driver is garbage both for reads and writes. No values are cached in memory, requiring the entire file to be reread for every request. For writes, every action is done individually. They should be buffered in memory for some time or until the buffer is full, so that slow writes to disk can be done at once. 

### Usage
Error handling could be improved, and http statuses should be used to indicate errors. 

There is no graceful shutdown implemented. Therefore any pending requests could be dropped. If this happens when a user posts a secret, the secret could be written but the user will see an error and think the request failed. 

### Security
There are probably mulptiple different attack vectors. One security problem is that plaintext secrets could be copied into swap-space. Anyone who has access to the harddrive could have access to secrets as they are being received or sent, if the service's data is swapped. 

Vault is currently configurable to use 1-way SSL. This means that secrets are encrypted in transit. However, there is no mutual SSL, so clients cannot be authenticated. 
