## Commands to build:
```bash
go build -o output/git-split
```

## Commands to run:
```bash
output/git-split split \
  --base main \
  --target feature-big \
  --size 3 \
  --prefix feature-part
```