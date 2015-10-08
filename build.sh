echo "Building OCGen for Darwin x64..."
GOOS=darwin GOARCH=amd64 go build -o ./bin/osx_64/ocgen ./main.go
echo "done"
