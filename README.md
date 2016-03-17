# gwc-cli
Gateway Converter CLI sample.

## Build
```shell
go get -t -d -v ./...
go build
```

## Copy config.
```shell
cp examples/config.yml .
```
Please edit ./config.yml and put proper information.


## Run
```shell
./gwc-cli -evid {end-node vendor id} -aname {app name configured in config.yml}
```

## Respond to command.
```shell
cp examples/commandResult.json
```
Please edit ./commandResult.json and choose 3. when gwc-cli is up.
