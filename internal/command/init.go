package command

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Songmu/prompter"
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
	config.Project.Name = promptStringForValue(cmd, "project-name", DefaultProjectName, "Project name")
	config.Project.Namespace = promptStringForValue(cmd, "namespace", DefaultProjectNamespace, "Namespace")
	config.Project.Entrypoint = promptStringForValue(cmd, "entrypoint", DefaultProjectEntrypoint, "Entrypoint")
	return nil
}

func runProjectTasks(cmd *cli.Command, config *interfaces.ProjectConfig) error {

	// Check if the project directory already exists
	projectDir := path.Join(os.Getenv("PWD"), config.Project.Name)

	overwrite := cmd.Bool("force")

	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		if !overwrite {
			force := promptConfirmation("Project directory already exists. Overwrite", false)
			if !force {
				log.Info("Project creation cancelled by user")
				return nil
			}
			log.Debug("Project directory already exists, overwriting...")
			overwrite = true
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
		gitInit = promptConfirmation("Do you want to initialize a git repository", true)
	}

	fmt.Printf(`

This will create a new project with the following configuration:

Project name: %s
Project path: %s
DPNamespace: %s
Entrypoint file: %s

Overwrite existing files: %t
Initialize git repository: %t

`,
		config.Project.Name,
		projectDir,
		config.Project.Namespace,
		config.Project.Entrypoint,
		overwrite,
		gitInit,
	)

	if !promptConfirmation("Create the project", true) {
		log.Info("Project creation cancelled by user")
		return nil
	}

	err := createProject(projectDir, config, overwrite, gitInit)
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

func promptStringForValue(cmd *cli.Command, flagName, defaultValue, promptLabel string) string {
	if cmd.String(flagName) == defaultValue {
		return prompter.Prompt(promptLabel, defaultValue)
	}
	return cmd.String(flagName)
}

func promptConfirmation(promptLabel string, defaultValue bool) bool {
	return prompter.YN(promptLabel, defaultValue)
}
