package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
					err := db.Model(&nixgo.Project{}).Order("name").Find(&projects).Error
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
						log.Fatalln(result.Error)
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
						Name: "id",
					},
				},
				Action: func(cCtx *cli.Context) error {
					var packages []nixgo.ProjectPackage

					id := cCtx.Uint("id")
					if id == 0 {
						name := cCtx.Args().First()
						var project nixgo.Project
						db.Where(&nixgo.Project{Name: name}).Find(&project)
						id = project.ID
					}
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
						Name: "id",
					},
				},
				Action: func(cCtx *cli.Context) error {
					id := cCtx.Uint("id")
					pkgIndex := 0
					if id == 0 {
						name := cCtx.Args().First()
						var project nixgo.Project
						db.Where(&nixgo.Project{Name: name}).Find(&project)
						id = project.ID
						pkgIndex = 1
					}

					for i := pkgIndex; i < cCtx.Args().Len(); i++ {
						pkg := cCtx.Args().Get(i)
						projectPackage := nixgo.ProjectPackage{ProjectID: id, Name: pkg}

						result := db.Create(&projectPackage)
						if result.Error != nil {
							log.Fatalln(result.Error)
						}

						fmt.Printf("Added %v to %v\n", pkg, id)
					}

					return nil
				},
			},
			{
				Name:    "remove-package",
				Aliases: []string{"rp"},
				Usage:   "removes package from project/environment",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "id",
					},
				},
				Action: func(cCtx *cli.Context) error {
					id := cCtx.Uint("id")
					pkgIndex := 0
					if id == 0 {
						name := cCtx.Args().First()
						var project nixgo.Project
						db.Where(&nixgo.Project{Name: name}).Find(&project)
						id = project.ID
						pkgIndex = 1
					}

					for i := pkgIndex; i < cCtx.Args().Len(); i++ {
						pkg := cCtx.Args().Get(i)
						var projectPackage nixgo.ProjectPackage
						db.Where(&nixgo.ProjectPackage{ProjectID: id, Name: pkg}).First(&projectPackage)
						db.Delete(&projectPackage)
						fmt.Printf("Removed %v from %v\n", pkg, id)
					}

					return nil
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"del", "remove", "rm"},
				Usage:   "delete project/environment",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "id",
					},
				},
				Action: func(cCtx *cli.Context) error {
					var name string

					id := cCtx.Uint("id")
					if id == 0 {
						name = cCtx.Args().First()
						var project nixgo.Project
						db.Where(&nixgo.Project{Name: name}).Find(&project)
						id = project.ID
					}

					if id == 0 {
						log.Fatalln("Project not found")
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
						log.Fatalf("Cannot find project: %v\n", name)
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

						if project.Path != "" {
							cmd.Dir = project.Path
						}

						cmd.Run()
						cmd.Wait()
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
