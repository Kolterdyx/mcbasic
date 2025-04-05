package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"os"
	"os/exec"
	"path"
)

const (
	DefaultProjectName       = "MyProject"
	DefaultProjectNamespace  = "my_namespace"
	DefaultProjectEntrypoint = "src/main.mcb"
)

var gitInitFlag = cli.BoolFlag{
	Name:  "git-init",
	Value: false,
	Usage: "Initialize a git repository if git is installed",
}

var InitCommand = &cli.Command{
	Name:  "init",
	Usage: "Initialize a new project",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "project-name",
			Value: DefaultProjectName,
			Usage: "Name of the project",
		},
		&cli.StringFlag{
			Name:  "namespace",
			Value: DefaultProjectNamespace,
			Usage: "Namespace of the project",
		},
		&cli.StringFlag{
			Name:  "entrypoint",
			Value: DefaultProjectEntrypoint,
			Usage: "Entrypoint of the project. The compiler will look for this file to start compiling",
		},
		&cli.BoolFlag{
			Name:  "force",
			Value: false,
			Usage: "Force overwrite existing files",
		},
		&gitInitFlag,
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return Init(cmd)
	},
}

func Init(cmd *cli.Command) error {
	DebugFlagHandler(cmd)

	config := interfaces.ProjectConfig{
		Project: interfaces.Project{
			Name:        cmd.String("project-name"),
			Namespace:   cmd.String("namespace"),
			Version:     "0.0.1",
			Description: "",
		},
	}

	err := getProjectData(cmd, &config)
	if err != nil {
		return err
	}

	err = runProjectTasks(cmd, &config)
	if err != nil {
		return err
	}

	log.Debugf("Project config: %+v", config)

	return nil
}

func getProjectData(cmd *cli.Command, config *interfaces.ProjectConfig) error {
	name, err := promptStringForValue(cmd, "project-name", DefaultProjectName, "Project name")
	if err != nil {
		return err
	}
	config.Project.Name = name

	namespace, err := promptStringForValue(cmd, "namespace", DefaultProjectNamespace, "Namespace")
	if err != nil {
		return err
	}
	config.Project.Namespace = namespace

	entrypoint, err := promptStringForValue(cmd, "entrypoint", DefaultProjectEntrypoint, "Entrypoint")
	if err != nil {
		return err
	}
	config.Project.Entrypoint = entrypoint
	return nil
}

func runProjectTasks(cmd *cli.Command, config *interfaces.ProjectConfig) error {

	// Check if the project directory already exists
	projectDir := path.Join(os.Getenv("PWD"), config.Project.Name)

	forceOverwrite := cmd.Bool("force")

	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		if !forceOverwrite {
			force, err := promptConfirmation("Project directory already exists. Do you want to overwrite it?", "n")
			if err != nil {
				return err
			}
			if !force {
				log.Info("Project creation cancelled by user")
				return nil
			}
			log.Debug("Project directory already exists, overwriting...")
			forceOverwrite = true
		} else {
			log.Debug("Project directory already exists, but force flag is set, overwriting...")
		}
	}

	gitInit := cmd.Bool("git-init")
	if gitInit {
		// Check if git is installed
		if _, err := exec.LookPath("git"); err != nil {
			return fmt.Errorf("git is not installed or not found in PATH: %w", err)
		}
		log.Debug("Git init flag is set, initializing git repository")
	}
	if !gitInitFlag.IsSet() {
		git, err := promptConfirmation("Do you want to initialize a git repository?", "y")
		if err != nil {
			return err
		}
		gitInit = git
	}

	fmt.Printf(`
Project name: %s
Project path: %s
Namespace: %s
Entrypoint file: %s
Initialize git repository: %t"
`,
		config.Project.Name,
		projectDir,
		config.Project.Namespace,
		config.Project.Entrypoint,
		gitInit,
	)

	confirm, err := promptConfirmation("Create the project?", "y")
	if err != nil {
		return err
	}
	if !confirm {
		log.Info("Project creation cancelled by user")
		return nil
	}

	err = createProject(projectDir, config, forceOverwrite, gitInit)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	log.Debug("Project data collected successfully")
	return nil
}

func createProject(projectDir string, config *interfaces.ProjectConfig, forceOverwrite bool, gitInit bool) error {

	// Check if the project directory already exists
	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		if forceOverwrite {
			log.Debug("Project directory already exists, overwriting...")
			err := os.RemoveAll(projectDir)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("project directory already exists: %s", projectDir)
		}
	}

	// Create the project directory
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		return err
	}

	// Create parent directories for entrypoint
	entrypoint := config.Project.Entrypoint
	entrypointDir := path.Dir(entrypoint)
	err = os.MkdirAll(path.Join(projectDir, entrypointDir), 0755)
	if err != nil {
		return err
	}
	// Create the entrypoint file
	entrypointFile, err := os.Create(path.Join(projectDir, entrypoint))
	if err != nil {
		return err
	}
	defer func(entrypointFile *os.File) {
		err := entrypointFile.Close()
		if err != nil {
			log.Fatalf("Failed to close entrypoint file: %v", err)
		}
	}(entrypointFile)
	_, err = entrypointFile.WriteString(`
# This is the entrypoint file
func main() {
	mcb:print("Hello, world!");
}
`)
	if err != nil {
		return err
	}

	if gitInit {
		log.Debug("Initializing git repository...")
		// Change current working directory to the project directory
		err = os.Chdir(projectDir)
		if err != nil {
			return err
		}
		err := exec.Command("git", "init").Run()
		if err != nil {
			return err
		}
		log.Debug("Git repository initialized successfully")
	} else {
		log.Debug("Git init flag is not set, skipping git initialization")
	}

	// Create project.toml file
	projectFile, err := os.Create(path.Join(projectDir, "project.toml"))
	if err != nil {
		return err
	}
	defer func(projectFile *os.File) {
		err := projectFile.Close()
		if err != nil {
			log.Fatalf("Failed to close project file: %v", err)
		}
	}(projectFile)
	tomlEncoder := toml.NewEncoder(projectFile)
	err = tomlEncoder.Encode(config)
	if err != nil {
		return err
	}
	log.Debug("Project file created successfully")

	log.Debugf("Project %s created successfully", projectDir)
	return nil
}

func promptStringForValue(cmd *cli.Command, flagName, defaultValue, promptLabel string) (string, error) {
	if cmd.String(flagName) == defaultValue {
		prompt := promptui.Prompt{
			Label: fmt.Sprintf("%s (default: %s)", promptLabel, defaultValue),
		}
		value, err := prompt.Run()
		if err != nil {
			return "", fmt.Errorf("failed to read %s: %w", flagName, err)
		}
		if value == "" {
			value = defaultValue
		}
		return value, nil
	}
	return cmd.String(flagName), nil
}

func promptConfirmation(promptLabel string, defaultValue string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     promptLabel,
		IsConfirm: true,
		Default:   defaultValue,
	}
	result, err := prompt.Run()
	if err != nil && !errors.Is(err, promptui.ErrAbort) {
		return false, fmt.Errorf("failed to read confirmation: %+v", err)
	} else if errors.Is(err, promptui.ErrAbort) {
		return false, nil
	}
	return result == "y" || result == "Y" || (result == "" && (defaultValue == "y" || defaultValue == "Y")), nil
}
