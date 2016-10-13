# GitReleaseUpdater

## Basic usage

### Install

```bash
go get github.com/geeksteam/GitReleaseUpdater
```

### Import

```go
import updater "github.com/geeksteam/GitReleaseUpdater"
```

### Start

```go
go updater.Updater{
		CheckInterval:  1 * time.Hour,
		NeedUpdate:     updater.VCSimpleDiff("0.1"),
		UpdateMethod:   updater.UMUntarGz("/path/to/my/program"),
		RestartCommand: updater.RCService("my-service-name"),
		Source: updater.SourceGitReleases(
			"repo-owner",
			"repo-name",
			"github-username",
			"github-password"),
	}.Watch()
```
You can leave github credentials just "" if use public repository 

---

Or if you want to start update manually after some event (e.g. got update command from rabbitMQ)

```go
var (
	versionURL = "http://my.site/updates/version"
	latestURL = "http://my.site/updates/latest.tar.gz"
)

updater.Updater{
		NeedUpdate:     updater.VCAlways(),
		UpdateMethod:   updater.UMReplaceBin(),
		RestartCommand: updater.RCService("my-service-name"),

		Source: updater.SourceDirectLinks(versionURL, latestURL),
	}.Update()
```

---

### Custom commands
You can use your own functions with additional logic as parameters for updater.
While **NeedUpdate**, **UpdateMethod** and **RestartCommand** is simle funcions, 
Source is interface with 2 methods: **LastVersion() (string, error)** and **Download() (string, error)**. 