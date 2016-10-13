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

		// Bool func which compare current version with latest
		NeedUpdate:     updater.VCSimpleDiff("0.1"),

		// This func takes path to downloaded file and should replace current program files
		UpdateMethod:   updater.UMUntarGz("/path/to/my/program"),

		// And this one restarts your program when update finished
		RestartCommand: updater.RCService("my-service-name"),

		// Source is entity wich know where to find latest version and how download it
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

err := updater.Updater{
		NeedUpdate:     updater.VCAlways(),
		UpdateMethod:   updater.UMReplaceBin(),
		RestartCommand: updater.RCService("my-service-name"),

		Source: updater.SourceDirectLinks(versionURL, latestURL),
	}.Update()
```

---

## Predefined functions
### Version checkers
All starts with **VC** and returns **func(string) bool**: 

* **VCSimpleDiff(currentVer string)** - simply compare new version are equal with current
* **VCAlways()** - only returns true

### Update methods
All starts with **UM** and return **func(string) error**
* **UMUntarGz(destDir string)** - call tar system utility, uses for tar.gz archives
* **UMReplaceBin()** - replaces own executable with new one

### Restart command
All starts with **RC** and returns **func() error**
* **RCService(service string)** - it calls "service name restart"
* **RCNothing()** - for cases when no restart needed

### Custom commands
You can use your own functions with additional logic as parameters for updater.
While **NeedUpdate**, **UpdateMethod** and **RestartCommand** is simle funcions, 
Source is interface with 2 methods: **LastVersion() (string, error)** and **Download() (string, error)**. 