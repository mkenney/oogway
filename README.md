[![Build Status](https://travis-ci.org/kcmerrill/oogway.svg?branch=master)](https://travis-ci.org/kcmerrill/oogway) [![Go Report Card](https://goreportcard.com/badge/github.com/kcmerrill/oogway)](https://goreportcard.com/report/github.com/kcmerrill/oogway)

![oogway](assets/oogway.png "oogway")

**Oogway** is simple, yet flexible, monitoring tool. At it's core, **Oogway** is a coordinated shell runner. You can determine how often checks are run, what to do when they fail, when they succeed and even try to fix problems before they go critical. What you do is completely up to you. Send stats, send notifications etc ...

**Designed to get you up and running in < 5 minutes**

## Checks

Checks are yaml files with a given extension. By default, that extension is `.oog`. Within these yaml files will contain your instructions in regards to each check.

```yaml
kcmerrill.com:
    every: 5s
    reset: 15m
    try: 5
    check:
        cmd: curl --fail https://kcmerrill.com
    fix:
        after: 5
        cmd: |
            ssh me@kcmerrill.com "cd /code/kcmerrill.com && alfred /dev staticwebserver"
    warning:
        cmd: |
            alfred /oogway/slack:http.warning "#alerts" "{{ .Name }}" "https://kcmerrill.com"
            curl -XPOST "http://localhost:8086/write?db=oogway" -d 'check,http=kcmerrill.com,status=warning warning=1'
    critical:
        cmd: |
            alfred /oogway/slack:http.critical "#alerts" "{{ .Name }}" "https://kcmerrill.com"
            curl -XPOST "http://localhost:8086/write?db=oogway" -d 'check,http=kcmerrill.com,status=critical critical=1'
    recover:
        cmd: |
            alfred /oogway/slack:http.recover "#alerts" "{{ .Name }}" "https://kcmerrill.com"
```

A few things to note:

1. I'm using pre-baked [`alfred`](https://github.com/kcmerrill/alfred) tasks here as to not repeat myself.
   1. Feel free and use these tasks too! Just remember to set `SLACK_WEBHOOK_URL` in your `env`.
1. While decent notifications, you can add whatever you want to commands.
   1. Be creative, you can use an external config like [sherlock](https://github.com/kcmerrill/sherlock) to do heartbeats.
1. `Oogway` pairs nicely with graphing databases such as [InfluxDB](https://github.com/influxdata/influxdb).
   1. Notice on the critical, sending data. The `cmd` key is just yaml, so it can accept multi-line commands.
1. The `cmd` is just run on the shell, so if you have something much more complicated than just a simple curl, feel free to call your application here.
1. You can have as many files you'd like, or many checks in one file. 100% up to you. Just note that duplicated check names will be overwritten.

## The Flow

1. **check** *(required)* is run. Depending on exit code:
   1. If succesful, go back to **check** next time around.
   1. If unsuccesful, continue on.
1. **warning** *(optional, only if after >= attempts is met)* will run
1. **fix** *(optional, only if after >= attempts is met)*
1. **critical** *(optional)* will run once **check** has failed **try** times.
1. **recover** *(optional)* will run *only* if the check is critical and then turns **ok**

## Templates

Every command that you run, will have access to the other commands and it's errors, elapsed time and results output so you can do whatever you need to. It uses [sprig](https://github.com/Masterminds/sprig) so you'll have plenty of additional functionality too.

* **{{ .Name }}** string, name of the check
* **{{ .Attempts }}** int, how many times the check has been run
* **{{ .Status }}** string, ok|warning|critical
* **{{ .Instructions.Check|Fix|Warning|Critical|Recover.Cmd }}** string, command that was run
* **{{ .Instructions.Check|Fix|Warning|Critical|Recover.Results }}** string, results of the command
* **{{ .Instructions.Check|Fix|Warning|Critical|Recover.Error }}** string, errors of the command
* **{{ .Instructions.Check|Fix|Warning|Critical|Recover.Runtime }}** int, elapsed time(in seconds) of the command

## Binaries || Installation

[![MacOSX](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/apple_logo.png "Mac OSX")](http://go-dist.kcmerrill.com/kcmerrill/oogway/mac/amd64) [![Linux](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/linux_logo.png "Linux")](http://go-dist.kcmerrill.com/kcmerrill/oogway/linux/amd64)

via go:

`$ go get -u github.com/kcmerrill/oogway`

via docker:

`$ docker run --rm -d -v $PWD:/oogway kcmerrill/oogway`

## Usage

```bash
$ oogway --check-interval 10s --check-dir /path/to/checks
```

* `--check-interval` How often the checks should be reloaded? *(default 1s)*
* `--check-dir` Directory where checks are stored *(default ".")*
* `--check-extension` Extension of check yaml files *(default "oog")*
