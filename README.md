# shellscript-executor

Golang program that will run all bash files provided in the given json file `test.json`:
```json
{
  "abc": "/home/tinh/abc.sh",
  "def": "/home/abc.sh"
}
```

The files must be sh/bash files:
1. The content of each file has `#!/bin/sh` or `#!/bin/bash` as the first line.
2. Executable, aka `chmod +x script`.

The json file path can be customized with env variable:
```
SETTINGS_DICTIONARY_PATH="test.json"
```
