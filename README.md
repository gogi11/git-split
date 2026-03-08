If you wish to use it with your own repo, set the API key as an environment variable:
```bash
set GITHUB_TOKEN=your_github_token
set GITLAB_TOKEN=your_gitlab_token
set BITBUCKET_TOKEN=your_bitbucket_token
```

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
  --mode directory \
  --dry-run
```