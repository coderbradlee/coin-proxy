export GOPATH=$GOPATH:`pwd`
echo $GOPATH
go build -o coin-proxy .
