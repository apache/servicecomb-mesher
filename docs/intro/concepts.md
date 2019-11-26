# Concepts

### Sidecar
Mesher leverages 
[distributed design pattern, sidecar](https://kubernetes.io/blog/2015/06/the-distributed-system-toolkit-patterns/) 
to work along with service. 

### go chassis 
Mesher is a light weight sidecar proxy developed on top of go-chassis,
so it has the same [concepts](http://go-chassis.readthedocs.io/en/latest/intro/concepts.html) with it 
and it has all features of go chassis

### Destination Resolver

Destination Resolver parses request into a service name

### Source Resolver

Source resolver gets remote IP and based on remote IP, it provides a standard way for the applications to create media sources.

### Admin API

Admin API listens on isolated port, it gives a way to interact with mesher

