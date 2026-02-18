# Golang project practice: GYM backend manager

## Using gRPC 
For using gRPC we need .proto files to create contracts, as an example for rest apis we use json as response, in gRPC uses .proto files that should be "translated" to GO languaje.

Steps to install compiler for proto files to go language:

1. install compiler (protoc) (mac steps)

brew install protobuf

2. install go plugins

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

critical step: env vars
to ensure computer can find the plugins the BIN folder of GO must be in the PATH 
check command: export PATH="$PATH:$(go env GOPATH)/bin"
validate with echo $PATH, should contain the go path

3. generate the code (lets start or validate the Magic)
go to root of the project then run: 
protoc --go_out=. --go-grpc_out=. api/proto/my-file.proto

if everything goes well, a proto folder should be created (if not exists), that contains 2 files with .pb.go extensions.

NOTE: NEVER touch this files directly, if we need to update the contract, update the .proto file then run the generate comand again.

## Migrations
Lets handle DB work with migrations, versioning, rollback and team work with escential tools.

Steps to migrate:

1. tool: golang-migrate

install (for mac) brew install golang-migrate

2. structurate mirgations:

under folder: migrations/ use files:
 up -> for applying changes
 down -> for reverting

to generate a new migration use:

migrate create -ext sql -dir migrations -seq migration_name

## Goroutines
1. What is a Goroutine?
Unlike traditional threads in other languages, which are heavyweight and managed by the operating system, goroutines are extremely lightweight threads managed by the Go runtime itself.

Size: A goroutine starts with a mere 2KB of memory (a Java or C++ thread can occupy 1MB). You can have thousands or millions of them on a typical laptop.

Concurrency vs. Parallelism: Goroutines allow concurrency (handling many things at once). If your CPU has multiple cores, Go will automatically execute them in parallel.


2. Where it works? (Reales escenario)
They excel at tasks that are I/O bound (waiting for something to finish without blocking the program):

External calls: Querying 3 different APIs simultaneously.

Heavy background processes: Sending emails, generating PDFs, or processing images without making the user wait.

Data streaming: Processing information as it arrives.

For this project, email is the example of using goroutines in real world task.

IMPORTANT:

The 3 Pillars
To master concurrency in Go, simply knowing how to type the word "go" isn't enough. You must study these three concepts:

- Channels: "Don't communicate by sharing memory, share memory by communicating." Channels are the "pipes" for passing data between goroutines safely.

- WaitGroups: Used when you need your program to wait for several goroutines to finish before closing (very common in batch processes).

- Context: Used to cancel goroutines if a request takes too long (vital in gRPC).

Goroutines are powerful, but you must avoid Goroutine Leaks (launching a goroutine that never ends and consumes memory eternally).
