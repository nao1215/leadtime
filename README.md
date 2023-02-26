![Coverage](https://raw.githubusercontent.com/nao1215/octocovs-central-repo/main/badges/nao1215/leadtime/coverage.svg)
![Test Execution Time](https://raw.githubusercontent.com/nao1215/octocovs-central-repo/main/badges/nao1215/leadtime/time.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/nao1215/leadtime)](https://goreportcard.com/report/github.com/nao1215/leadtime)
[![reviewdog](https://github.com/nao1215/leadtime/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/nao1215/leadtime/actions/workflows/reviewdog.yml)
[![LinuxUnitTest](https://github.com/nao1215/leadtime/actions/workflows/linux_test.yml/badge.svg)](https://github.com/nao1215/leadtime/actions/workflows/linux_test.yml)
[![MacUnitTest](https://github.com/nao1215/leadtime/actions/workflows/mac_test.yml/badge.svg)](https://github.com/nao1215/leadtime/actions/workflows/mac_test.yml)
[![WindowsUnitTest](https://github.com/nao1215/leadtime/actions/workflows/windows_test.yml/badge.svg)](https://github.com/nao1215/leadtime/actions/workflows/windows_test.yml)
# leadtime - calculate PR lead time statistics on GitHub
leedtime is a command that outputs statistics about the time it takes for a GitHub Pull Request to be merged. The leadtime command was developed under the influence of the following books.
- Eng: [Accelerate: The Science of Lean Software and DevOps: Building and Scaling High Performing Technology Organizations](https://www.amazon.com/dp/1942788339/ref=cm_sw_r_cp_ep_dp_sBN8BbGC11MBS)
- JP: [LeanとDevOpsの科学［Accelerate］](https://www.amazon.co.jp/Lean%E3%81%A8DevOps%E3%81%AE%E7%A7%91%E5%AD%A6%EF%BC%BBAccelerate%EF%BC%BD-%E3%83%86%E3%82%AF%E3%83%8E%E3%83%AD%E3%82%B8%E3%83%BC%E3%81%AE%E6%88%A6%E7%95%A5%E7%9A%84%E6%B4%BB%E7%94%A8%E3%81%8C%E7%B5%84%E7%B9%94%E5%A4%89%E9%9D%A9%E3%82%92%E5%8A%A0%E9%80%9F%E3%81%99%E3%82%8B-impress-top-gear%E3%82%B7%E3%83%AA%E3%83%BC%E3%82%BA-ebook/dp/B07L2R3LTN)

The motivation for developing the leadtime command is to measure lead time for changes. I used unit test coverage as a measure of software quality. However, as the number of unit tests increased but the code was not rewritten, I questioned whether the quality was improving.  
  
Therefore, I considered measuring lead time, one of the indicators presented in the above book.

## How to install
### Use "go install"
If you does not have the golang development environment installed on your system, please install golang from the [golang official website](https://go.dev/doc/install).
```
 go install github.com/nao1215/leadtime@latest
```

### Use homebrew (aarch64)
```
$ brew tap nao1215/tap
$ brew install nao1215/tap/leadtime
```

## How to use
You need to set GitHub access token in environment variable "LT_GITHUB_ACCESS_TOKEN". If you want to check github.com/nao1215/sqly repository, you execute bellow.
```
$ leadtime stat --owner=nao1215 --repo=sqly
PR      Author  Bot     LeadTime[min]   Title
#29     dependabot[bot] yes     21144   Bump github.com/fatih/color from 1.13.0 to 1.14.1
#28     nao1215 no      12      Change golden pacakge import path
#27     nao1215 no      17      add unit test for infra package
#26     nao1215 no      686     Add basic unit test for shell
#25     dependabot[bot] yes     1850    Bump github.com/google/go-cmp from 0.2.0 to 0.5.9
#24     nao1215 no      6458    Add unit test for model package
#23     nao1215 no      187     Change golden test package from goldie to golden and more
#22     nao1215 no      1       Add sqlite3 syntax completion
#21     nao1215 no      1769    Add unit test for argument paser
#20     nao1215 no      53      Feat dump tsv ltsv json
#19     nao1215 no      6       Add featuer thar print date by markdown table format
#18     nao1215 no      10      Feat import ltsv
#17     nao1215 no      117     Feat import tsv
#15     nao1215 no      57      Fix panic bug when import file that is without extension
#14     nao1215 no      42      Feat import json
#13     nao1215 no      139     Fix input delays when increasing records
#12     nao1215 no      18      Add header command
#11     nao1215 no      1552    Fixed a display collapse problem when multiple lines are entered
#10     nao1215 no      4       Fixed a bug that caused SQL to fail if there was a trailing semicolon
#9      nao1215 no      29      Add move cursor function in intaractive shell
#8      nao1215 no      3       Fixed a bug in which the wrong arguments were used
#7      nao1215 no      76      Added CSV output mode
#6      nao1215 no      222     Improve execute query
#5      nao1215 no      498     Add history usecase, repository, infra. sqly manage history by sqlite3
#4      nao1215 no      139     Add function that execute select query
#3      nao1215 no      37      Add import command
#2      nao1215 no      57      Add .tables command
#1      nao1215 no      127     Add .exit/.help command and history manager

[statistics]
 Total PR       = 28
 Lead Time(Max) = 21144[min]
 Lead Time(Min) = 1[min]
 Lead Time(Sum) = 35310[min]
 Lead Time(Ave) = 1261.07[min]
 Lead Time(Median) = 66.50[min]
```

### markdown format output
If you change output format to markdown, you use --markdown option. Markdown output sample is [here](doc/sample_leadtime.md).
```
$ leadtime stat --owner=nao1215 --repo=gup --markdown
```

## Features to be added
- [ ] CSV output format
- [ ] JSON output format
- [ ] Markdown file output
- [ ] Output to file
- [ ] Supports GitHub Actions
- [ ] Exclude the bot's PR
- [ ] faster by goroutine

## Contributing / Contact
First off, thanks for taking the time to contribute! heart Contributions are not only related to development. For example, GitHub Star motivates me to develop!
  
If you would like to send comments such as "find a bug" or "request for additional features" to the developer, please use one of the following contacts.
- [GitHub Issue](https://github.com/nao1215/leadtime/issues)

## LICENSE
The leadtime project is licensed under the terms of [MIT LICENSE](./LICENSE).