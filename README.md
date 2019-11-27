# OneDrive with vsftpd

forked from: [jstaf/onedriver](https://github.com/jstaf/onedriver)

## build:

```
CGO_ENABLED=0 go build
mv onedriver dockerfile
cd dockerfile
docker build -t jdkhome/onedrive-ftp:{tag} .
```

## use:

see https://www.jdkhome.com/dev-ops/deploy/onedrive-ftp.html




