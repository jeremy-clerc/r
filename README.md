# R (Redirector)

```sh
go get github.com/jeremy-clerc/r
```

```
r -help
Usage of r:
  -listen string
    	Address and port to listen on. (default "127.0.0.1:8008")
  -path string
    	Path to a file of shortcut!link. (default "links")
```

If you already have a webserver listening on port 80, like Apache, create a
vhost, and edit your hosts file to have an host `r` points to your webserver.
In your browser just type `r/shorcut`.

An other way is to set a custom search engine in your webbrowser. Then you use
it by typing `r<space>shorcut`.
