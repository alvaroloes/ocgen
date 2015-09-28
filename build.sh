echo "Building OCGen for Darwin x64..."
GOOS=darwin GOARCH=amd64 go build -o ./bin/osx_64/ocgen ./main.go
echo "done"
echo "Building OCGen for Linux x64..."
GOOS=linux GOARCH=amd64 go build -o ./bin/linux_64/ocgen ./main.go
echo "done"
