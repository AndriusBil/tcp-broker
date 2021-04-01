# tcp-broker

## Built with
Project built with go1.16, using native go libraries.

## Getting started
***

Test

### Prerequisites

Install [golang](https://golang.org/doc/install)

### Installation
1. Clone the repo
```
git clone git@github.com:AndriusBil/tcp-broker.git
```

2. Install dependencies
```
go get ./...
```

3. Build the project
```
go build
```

## Usage
***

Start the broker, inject consumer and publisher ports via environment variables
```bash
PUBLISHERS_PORT=:3003 CONSUMERS_PORT=:3004 ./tcp-broker 
```

Fire up multiple consumers using netcat

```
nc localhost 3004
```

Feed the messages

```
echo "Message" | nc localhost 3003
```