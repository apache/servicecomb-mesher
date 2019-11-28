# Profile Mesher

Mesher has a convenience way to enable go [pprof]( https://golang.org/pkg/net/http/pprof/ ), so that you can easily analyze the performance of mesher.

## Configurations

```yaml
pprof:
  enable: true
  listen: 127.0.0.0.1:6060
```

**enable**
>*(optional, bool)* Default is false


**listen**

>*(optional, string)* Listen IP and port

