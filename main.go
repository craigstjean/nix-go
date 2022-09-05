package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"craigstjean.com/nix-go/nixgo"

	"github.com/urfave/cli/v2"
)

func main() {
	db := nixgo.Start()

	app := &cli.App{
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls", "l"},
				Usage:   "list existing projects/environments",
				Action: func(cCtx *cli.Context) error {
					var projects []nixgo.Project
					err := db.Model(&nixgo.Project{}).Find(&projects).Error
					if err == nil {
						for _, p := range projects {
							if p.Path != "" {
								fmt.Printf("%v: %v (%v)\n", p.ID, p.Name, p.Path)
							} else {
								fmt.Printf("%v: %v\n", p.ID, p.Name)
							}
						}
					}
					return nil
				},
			},
			{
				Name:  "new",
				Usage: "create new project/environment",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "path",
						Value: "",
					},
				},
				Action: func(cCtx *cli.Context) error {
					newName := cCtx.Args().First()
					project := nixgo.Project{Name: newName}

					path := cCtx.String("path")
					if path != "" {
						project.Path = path
					}

					result := db.Create(&project)
					if result.Error != nil {
						log.Fatal(result.Error)
					}

					fmt.Printf("Created %v\n", project.ID)

					return nil
				},
			},
			{
				Name:    "list-packages",
				Aliases: []string{"lp"},
				Usage:   "list package in project/environment",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					id := cCtx.Uint("id")
					var packages []nixgo.ProjectPackage
					db.Where(&nixgo.ProjectPackage{ProjectID: id}).Find(&packages)

					for _, pp := range packages {
						fmt.Printf("%v\n", pp.Name)
					}

					return nil
				},
			},
			{
				Name:    "add-package",
				Aliases: []string{"ap"},
				Usage:   "add package to project/environment",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					id := cCtx.Uint64("id")
					pkg := cCtx.Args().First()
					projectPackage := nixgo.ProjectPackage{ProjectID: uint(id), Name: pkg}

					result := db.Create(&projectPackage)
					if result.Error != nil {
						log.Fatal(result.Error)
					}

					fmt.Printf("Added to %v\n", id)

					return nil
				},
			},
			{
				Name:    "remove-package",
				Aliases: []string{"rp"},
				Usage:   "removes package from project/environment",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					id := cCtx.Uint("id")
					pkg := cCtx.Args().First()

					if pkg == "" {
						fmt.Println("Package name required")
					} else {
						var projectPackage nixgo.ProjectPackage
						db.Where(&nixgo.ProjectPackage{ProjectID: id, Name: pkg}).First(&projectPackage)
						db.Delete(&projectPackage)
						fmt.Printf("Deleted from %v\n", id)
					}

					return nil
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"del", "remove", "rm"},
				Usage:   "delete project/environment",
				Action: func(cCtx *cli.Context) error {
					id, err := strconv.Atoi(cCtx.Args().First())
					if err != nil {
						log.Fatalf("%v not a number", id)
					} else {
						db.Delete(&nixgo.Project{}, id)
						fmt.Printf("Deleted %v\n", id)
					}

					return nil
				},
			},
			{
				Name:    "shell",
				Aliases: []string{"run", "go"},
				Usage:   "start project/environment",
				Action: func(cCtx *cli.Context) error {
					name := cCtx.Args().First()

					var project nixgo.Project
					err := db.Where(&nixgo.Project{Name: name}).Preload("Packages").First(&project).Error
					if err != nil {
						log.Fatalf("Cannot find project: %v", name)
					} else {
						var sb strings.Builder
						sb.WriteString("-p ")
						for _, pp := range project.Packages {
							sb.WriteString(pp.Name)
							sb.WriteString(" ")
						}

						if len(project.Packages) == 0 {
							sb.WriteString("hello ")
						}

						sb.WriteString("--run zsh")

						cmd := exec.Command("nix-shell", strings.Split(sb.String(), " ")...)
						cmd.Env = os.Environ()
						cmd.Env = append(cmd.Env, fmt.Sprintf("NIX_ENV=%v", name))
						cmd.Stdin = os.Stdin
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						cmd.Run()
						cmd.Wait()
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
