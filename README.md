# chibi

A simple, self-hosted link catalog / shortener / redirect service. 

## Building 

```bash 
go build -v .
```

## Deploying
`chibi` does not require much configuration. You need to set up a single 
environment variable, which is required for the frontend to show the full
URLs, which can be copied. 
```bash 
FRONTEND_URL=chibi.srev.in chibi
```

Alternatively, using the super-lightweight docker container for deploying 
chibi.

```bash
docker run -e FRONTEND_URL=chibi.srev.in ghcr.io/srevinsaju/chibi 
```

## License 
This project is licensed under the [MIT License](./LICENSE)


