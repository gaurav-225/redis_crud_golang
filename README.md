# Implemented based on task mentioned in https://app.codecrafters.io/catalog `Build Your Own Redis`

##### Tasks completed 

1. Added support for the ECHO, PING, SET & GET commands.
2. Added support for setting a key with an expiry.
3. Propagating write commands from a master to a replica.



## Bash Commands to run 

1. Start Master Server
```bash
go run app/server.go -port 44001
```

2. Start Replica server for Master running on port 44001

```bash
go run app/server.go -port 44002 -replicaof "localhost 44001"

#-----------------------------In case of error
# to store log in file
go run app/server.go -port 44002 -replicaof "localhost 44001" | tee resultCheckWhyOccur

# to store error and log in file
(go run app/server.go -port 44002 -replicaof "localhost 44001") 2>&1 | tee resultCheckWhyOccur2
```


3. Using telnet to commuicate with Server

```bash
telnet localhost 44002
telnet localhost 44001
->Get a
->Set a 10
```


