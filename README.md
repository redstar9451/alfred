# alfred
Generate alfred items from configuration

## Quick start
1. cp -r conf_example conf
2. create script workflow like this
```shell
#!/usr/bin/env sh

/Users/redstar/go/src/alfred/alfred /Users/redstar/go/src/alfred/conf/example.json $1
```
3. add step to open url
4. go to enjoy it.

## Feature
1. generate alfre items from configuration
2. fuzzy matching, this is better, huaman kind usage, especially if too many items
