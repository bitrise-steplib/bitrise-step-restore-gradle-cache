package step

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/cache"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

const stepId = "restore-gradle-cache"

// Cache key templates
// OS + Arch: to guarantee that stack-specific content (absolute paths, binaries) are stored separately
// checksum values:
// - `**/*.gradle*`: Gradle build files in any submodule, including ones written in Kotlin (*.gradle.kts)
// - `**/gradle-wrapper.properties`: contains exact Gradle version
// - `**/gradle.properties`: contains Gradle config values
// - `**/libs.versions.toml`: version catalog file, contains dependencies and their versions
var keys = []string{
	`{{ .OS }}-{{ .Arch }}-gradle-cache-{{ checksum "**/*.gradle*" "**/gradle-wrapper.properties" "**/gradle.properties" "**/libs.versions.toml" }}`,
	`{{ .OS }}-{{ .Arch }}-gradle-cache-`,
}

type Input struct {
	Verbose bool `env:"verbose,required"`
	Retries int  `env:"retries,required"`
}

type RestoreCacheStep struct {
	logger      log.Logger
	inputParser stepconf.InputParser
	envRepo     env.Repository
	cmdFactory  command.Factory
}

func New(
	logger log.Logger,
	inputParser stepconf.InputParser,
	envRepo env.Repository,
	cmdFactory command.Factory,
) RestoreCacheStep {
	return RestoreCacheStep{
		logger:      logger,
		inputParser: inputParser,
		envRepo:     envRepo,
		cmdFactory:  cmdFactory,
	}
}

func (step RestoreCacheStep) Run() error {
	var input Input
	if err := step.inputParser.Parse(&input); err != nil {
		return fmt.Errorf("failed to parse inputs: %w", err)
	}
	stepconf.Print(input)
	step.logger.Println()
	step.logger.Printf("Cache keys:")
	step.logger.Printf(strings.Join(keys, "\n"))
	step.logger.Println()

	step.logger.EnableDebugLog(input.Verbose)

	restorer := cache.NewRestorer(step.envRepo, step.logger, step.cmdFactory, nil)
	return restorer.Restore(cache.RestoreCacheInput{
		StepId:         stepId,
		Verbose:        input.Verbose,
		Keys:           keys,
		NumFullRetries: input.Retries,
	})
}
