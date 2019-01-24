# COM S 319 Homework 01

In order to compile the client and server go 1.11 or higher is needed.

## Setting up the Server

The server can be compiled with the command `go build -o <name> cmd/server/*` where `<name>` is the desired name of the compiled binary.

After compilation, the compiled binary is ran like any other executable.

You can also run the server without compilation with the `go run cmd/server/*`, making sure to have the correct relative path to the `server/*` directory.

the server uses the `-addr` flag to specify the port on which to run. If the `-addr` flag is not specified it will default to port 9000.

In order to specify the port, just pass the `-addr=":<port>"` flag when executing the program.

## Setting up the Client

Compiling/running the client is nearly identical to the server. It can be compiled with the command `go build -o <name> cmd/client/*` where `<name>` is the desired name of the compiled binary.

As with the server, you can run the client without executing with the `go run` command. Making sure to have the correct relative or absolute path to the `client/*` directory.

The client uses the `-server` flag in order to set the host:port of the server that it will connect to. By default this is `localhost:9000`. The flag should be set by passing the hostname/ip of the server, along with the desired port. `-server="<host>:<port>"`

## Limitations

The main limitations of the client and server exist only on Windows hosts. Some connectivity issues occur on Windows depending on the server port. Port 9000 has been tested on both Linux and Windows hosts and seems to be the most consistent.

Another limitation with the client is Window's CMD and PowerShells complience with ASCII escape characters. The client attemps to flush Stdout on each entered message, and prints a colored and formated version to the consol in its place. This is not possible on Windows machines which causes messages to be printed to Stdout upon hitting enter, and again as they are formatted by the client and printed to Stdout.