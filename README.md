## Commands to build:
```bash
set GOOS=windows 
set GOARCH=x64 
go build -o output/git-split.exe
```

## Commands to run:
```bash
output/git-split.exe split
  --target name-of-target-branch \
```