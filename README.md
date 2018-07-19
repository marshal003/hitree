### HiTree - Golang Implementation of Tree Command
 
[![Circle Ci](https://circleci.com/gh/marshal003/hitree.svg)](https://circleci.com/gh/marshal003/hitree.svg)
[![Coverage Status](https://coveralls.io/repos/github/marshal003/hitree/badge.svg)](https://coveralls.io/github/marshal003/hitree)

HiTree is a golang based implementation of popular linux tree command. This cameout as a learning exercise of golang. 

### Features

- Listing only directories
- Pruning empty directories
- Colored output
- Following Symlinks(if it is link to dir)
- Including / Excluding files based on regex
- Printing fullpath of the files
- Printing report - count of dir and files 
- Including / Excluding hidden files in result
- Controlling Max Level in the output
- Output Tree Structure as JSON on Console
- Output Tree Structure in file
- Filter Dirs based on filelimit
- Output with UserId, GroupId, Permission & Modification Time
- Sort in reverse alphabatic order
- Sort by Modification time
- Include Stats in JSON structure
- Redirect output in File

## Demo (using termtosvg)

![Alt text](./hitree.svg)

## Download & Install

### Install from Source

```
go get -u github.com/marshal003/hitree
```

### Intall from Binary

- Download from [Release](https://github.com/marshal003/hitree/releases)
- Extract tar and put the hitree binary in your path
- Open terminal and type `hitree`
- You can create alias if `hitree` looks too big to type :)

### Libraries Used

- Cobra & Viper (building CLI)
- Aurora (Coloring output)

### Usage Examples

- Get Help
    ```
    hitree --help
    ```

- Defaults to current directory
    ```
    hitree 
    ```
- With flags
    ```
    // List only directories
    hitree -d 

    // Prune empty directories
    hitree --prune

    // Include hidden files & directories
    hitree --all

    // Include only go files
    hitree -P "*.go"

    // Exclude all md files
    hitree -I "*.md"

    // Skip reporting
    hitree --noreport

    // Follow links (if this is dir)
    hitree ---followlink

    // Show fullpath and restrict level to 2
    hitree -L 3 -f

    // Output tree structure as JSON on console
    hitree --json > tree.json
    
    // Output tree structure as JSON on console
    hitree --json -o output.json 
    ```
### TODO
- Add filtering support based on modification time

## References
- https://linux.die.net/man/1/tree
- https://www.youtube.com/watch?v=XbKSssBftLM&t=1s (Courtesy to Francesc Campoy)

## License

HiTree is released under the Apache 2.0 license. See 