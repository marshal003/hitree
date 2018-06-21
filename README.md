### HiTree - Golang Implementation of Tree Command

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

## Download & Install

### Install from Source

```
go get -u github.com/marshal003/hitree
```

### Intall from Binary

```
Inprogress...
```

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
    # List only directories
    hitree -d 

    # Prune empty directories
    hitree --prune

    # Include hidden files & directories
    hitree --all

    # Include only go files
    hitree -P "*.go"

    # Exclude all md files
    hitree -I "*.md"

    # Skip reporting
    hitree --noreport

    # Follow links (if this is dir)
    hitree ---followlink

    # Show fullpath and restrict level to 2
    hitree -L 3 -f
    ```
## References
- https://linux.die.net/man/1/tree
- https://www.youtube.com/watch?v=XbKSssBftLM&t=1s (Courtesy to Francesc Campoy)

## License

HiTree is released under the Apache 2.0 license. See 