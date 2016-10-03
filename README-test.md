### Install [fswatch](https://github.com/emcrisostomo/fswatch)
```
brew install fswatch
```

### Monitor changes on path
fswatch . | (while read; do go test ././...; done)
