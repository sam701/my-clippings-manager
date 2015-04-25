# my-clippings-manager
Amazon Kindle My Clippings Manager

## Build

    export GOPATH=somedir
    go get github.com/sam701/my-clippings-manager
    cd somedir/bin
    ./my-clippings-manager

## Contributing

### Working on web assets.

The web assets reside in the folder `web`.
All files from this folder are embedded into the output binary by [go-bindata](https://github.com/jteeuwen/go-bindata).
When you are working on the web resources, it is helpful to call first

    go-bindata -debug -prefix web web

The server will serve the assets from the original files on disk.

**Do not forget to call `go generate` before you commit!**
