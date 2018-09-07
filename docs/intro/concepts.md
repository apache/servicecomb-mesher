# Concepts

### Sidecar
Mesher leverage 
[distributed design pattern, sidecar](https://kubernetes.io/blog/2015/06/the-distributed-system-toolkit-patterns/) 
to work along with service. 

### go chassis 
Mesher is a light weight sidecar proxy developed on top of go-chassis,
so it has the same [concepts](http://go-chassis.readthedocs.io/en/latest/intro/concepts.html) with it 
and it has all features of go chassis

### DestinationResolver

Destination Resolver parse request into a service name

### Source Resolver

source resolver get remote IP and based on remote IP, it 

### Admin API

Listen on isolated port, it gives a way to interact with mesher

