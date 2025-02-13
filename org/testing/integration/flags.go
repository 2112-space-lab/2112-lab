package integration

import (
	"flag"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/alecthomas/kong"
	"github.com/cucumber/godog"
	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	models_db "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db/models"
)

var CLI struct {
	FlagWrapper struct { //! wrapper is used to prevent collision with other flag sources
		ShowHelp bool `default:"false" help:"display app testing options usage"`
		Log      struct {
			Level  string `default:"debug" env:"LOG_LEVEL"`
			Format string `default:"json" enum:"json,text"`
		} `embed:"" prefix:"log."`
		DatabaseServer struct {
			Host         string        `default:"localhost" env:"TEST_DB_HOST"`
			Port         int           `default:"5432" env:"TEST_DB_PORT"`
			HostDocker   string        `default:"2112-database.db" env:"TEST_DB_HOST_DOCKER"`
			PortDocker   int           `default:"5432" env:"TEST_DB_PORT_DOCKER"`
			UserName     string        `default:"postgres" env:"TEST_DB_USER"`
			Password     string        `default:"postgres" env:"TEST_DB_PASS"`
			DatabaseName string        `default:"postgres" env:"TEST_DB_NAME"`
			Timeout      time.Duration `default:"90s"`
			MaxConn      int32         `default:"100"`
			MaxIdleConns int           `default:"10"`
			SSLMode      string        `default:"disable" enum:"require,disable"`
		} `embed:"" prefix:"db."`
		Options struct {
			TestRunningEnv        string `default:"host" env:"TEST_IT_RUNNING_ENV" enum:"host,docker" help:"whether tests are being executed on host or in docker container"`
			DockerNetwork         string `default:"2112_net" env:"TEST_IT_DOCKER_NETWORK" help:"Existing Docker network on which containers should be started"`
			AppWorkspacePath      string `default:"../.." aliases:"root-path" type:"path" help:"Root folder of App codebase"`
			ArtifactsPath         string `default:"./_artifacts" type:"path" help:"Folder where test artifacts per scenario are located"`
			AppMigrationsPath     string `default:"../../assets/migrations" type:"path" help:"Folder containing App database SQL migrations"`
			PreserveDatabases     bool   `default:"false"`
			PreserveContainers    bool   `default:"false"`
			AppServiceDockerImage string `default:"app-service:latest"`
			PropagatorDockerImage string `default:"propagator-service:latest"`
		} `embed:"" prefix:"options."`
	} `embed:"" prefix:"tapp."`
}

func GetDatabaseConnectionInfo() models_db.DatabaseConnectionInfo {
	return models_db.NewDatabaseConnectionInfo(
		models_db.DatabaseHostName(CLI.FlagWrapper.DatabaseServer.Host),
		models_db.DatabaseHostPort(CLI.FlagWrapper.DatabaseServer.Port),
		models_db.DatabaseHostNameDocker(CLI.FlagWrapper.DatabaseServer.HostDocker),
		models_db.DatabaseHostPortDocker(CLI.FlagWrapper.DatabaseServer.PortDocker),
		models_db.DatabaseName(CLI.FlagWrapper.DatabaseServer.DatabaseName),
		CLI.FlagWrapper.DatabaseServer.MaxConn,
		CLI.FlagWrapper.DatabaseServer.UserName,
		CLI.FlagWrapper.DatabaseServer.Password,
		models_db.NewDatabaseTLSConfig(
			"", //TODO add certs config
			"",
			"",
			CLI.FlagWrapper.DatabaseServer.SSLMode,
		),
	)
}

var opts = godog.Options{
	// Output:      os.Stdout,
	Concurrency: 1,
	Paths:       []string{"specs"},
}

func keepKongFlags(item string) bool {
	return strings.HasPrefix(item, "--tapp")
}

func godogFlags(item string) bool {
	return fx.NegateFilter(keepKongFlags)(item)
}

func parseAppTestKongFlags(logger *slog.Logger, kongAppTestArgs []string) {
	kk, err := kong.New(&CLI,
		kong.UsageOnError(),
		DefaultEnvarsApp(""),
	)
	if err != nil {
		panic(err)
	}
	kctx, err := kk.Parse(kongAppTestArgs)
	if err != nil {
		panic(err)
	}
	logger.Info("kong parsed app test args",
		slog.Any("kong.AppTestOptions", CLI),
	)
	kk.FatalIfErrorf(err)
	if CLI.FlagWrapper.ShowHelp {
		err = kctx.PrintUsage(false)
		if err != nil {
			panic(err)
		}
	}
}

func parseGodogFlags(logger *slog.Logger, artifactsFolderPath string, godogArgs []string) {
	hasFormat := false
	for _, v := range godogArgs {
		if strings.HasPrefix(v, "--godog.format=") {
			hasFormat = true
			break
		}
	}
	if !hasFormat {
		godogArgs = append(godogArgs, fmt.Sprintf(
			"--godog.format=cucumber:%s/results/godog-cucumber.log.json,pretty:%s/results/godog-pretty.log",
			artifactsFolderPath,
			artifactsFolderPath,
		))
	}
	godog.BindFlags("godog.", flag.CommandLine, &opts)
	err := flag.CommandLine.Parse(godogArgs[1:])
	if err != nil {
		panic(err)
	}
	logger.Info("godog parsed args",
		slog.Any("godog.Options", opts),
	)
	if CLI.FlagWrapper.ShowHelp {
		flag.CommandLine.Usage()
	}
}

// DefaultEnvarsApp option inits environment names for flags.
// It is copied from kong library
// and slightly modified to allow env precedence with prefix
func DefaultEnvarsApp(prefix string) kong.Option {
	processFlag := func(flag *kong.Flag) {
		switch env := flag.Envs; {
		case flag.Name == "help":
			return
		case len(env) == 1 && env[0] == "-":
			flag.Envs = nil
			return
			// case len(env) > 0:
			// 	return
		}
		replacer := strings.NewReplacer("-", "_", ".", "_")
		names := append([]string{prefix}, camelCase(replacer.Replace(flag.Name))...)
		names = siftStrings(names, func(s string) bool { return !(s == "_" || strings.TrimSpace(s) == "") })
		name := strings.ToUpper(strings.Join(names, "_"))
		flag.Envs = append([]string{name}, flag.Envs...)
		flag.Value.Tag.Envs = append([]string{name}, flag.Value.Tag.Envs...)
	}

	var processNode func(node *kong.Node)
	processNode = func(node *kong.Node) {
		for _, flag := range node.Flags {
			processFlag(flag)
		}
		for _, node := range node.Children {
			processNode(node)
		}
	}

	return kong.PostBuild(func(k *kong.Kong) error {
		processNode(k.Model.Node)
		return nil
	})
}

func siftStrings(ss []string, filter func(s string) bool) []string {
	i := 0
	ss = append([]string(nil), ss...)
	for _, s := range ss {
		if filter(s) {
			ss[i] = s
			i++
		}
	}
	return ss[0:i]
}

func camelCase(src string) (entries []string) {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return []string{src}
	}
	entries = []string{}
	var runes [][]rune
	lastClass := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		var class int
		switch {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}
	return entries
}
