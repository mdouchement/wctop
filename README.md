# wctop

wctop is simple **web-ui** for monitoring local running Docker containers like the CLI tool [ctop](http://ctop.sh).

![wctop](https://user-images.githubusercontent.com/6150317/39636085-7cd6f0b8-4fbf-11e8-8a6d-58642ebfb201.png)

## Installation

```sh
docker run -v /var/run/docker.sock:/var/run/docker.sock -p 8080:8080 -d mdouchement/wctop
```

## Technologies

- [Golang](https://golang.org/)
- [Vue.js](https://vuejs.org/)
- [Bulma](https://bulma.io/)

## Development

- Golang 1.8 or greater

- Initialization

```sh
go get -u github.com/mjibson/esc
go get github.com/Masterminds/glide
glide install
```

- Run the server with assets loaded from the filesystem:

```sh
LOCAL_ASSETS=1 go run -race wctop.go -b localhost
```

- When a new asset is added, run `go generate`


## License

**MIT**

## Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request
