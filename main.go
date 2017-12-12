package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
	"fmt"
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name: "bucket",
		EnvVar: "SSSTASH_BUCKET",
	},
	cli.StringFlag{
		Name: "prefix",
		EnvVar: "SSSTASH_PREFIX",
	},
	cli.StringFlag{
		Name: "key",
		EnvVar: "SSSTASH_KEY_ID",
	},
}

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name: "list",
			Aliases: []string{"ls"},
			Flags: flags,
			Action: func (c *cli.Context) error {
				app, err := buildApp(c)
				if err != nil {
					return err
				}

				err = app.ListIter(func (name string) bool {
					fmt.Println(name)
					return true
				})
				return wrapError(err)
			},
		},
		{
			Name: "put",
			Flags: flags,
			Action: func (c *cli.Context) error {
				if err := validateArgsLength(c, 2, 2); err != nil {
					return err
				}
				key := c.Args().Get(0)
				val := c.Args().Get(1)

				app, err := buildApp(c)
				if err != nil {
					return err
				}

				return wrapError(app.Put(key, val))
			},
		},
		{
			Name: "get",
			Flags: flags,
			Action: func (c *cli.Context) error {
				if err := validateArgsLength(c, 1, 1); err != nil {
					return err
				}
				name := c.Args().Get(0)

				app, err := buildApp(c)
				if err != nil {
					return err
				}

				return wrapError(app.Get(name))
			},
		},
		{
			Name: "delete",
			Aliases: []string{"rm"},
			Flags: flags,
			Action: func (c *cli.Context) error {
				if err := validateArgsLength(c, 1, 1); err != nil {
					return err
				}
				name := c.Args().Get(0)

				app, err := buildApp(c)
				if err != nil {
					return err
				}

				return wrapError(app.Delete(name))
			},
		},
	}
	app.Run(os.Args)
}

func validateArgsLength(c *cli.Context, min int, max int) error {
	l := len(c.Args())
	if max < 0 {
		if l < min {
			return cli.NewExitError(fmt.Sprintf("at least %d args required", min), 1)
		}
	} else if min == max {
		if l != min {
			return cli.NewExitError(fmt.Sprintf("%d args required", min), 1)
		}
	} else {
		if !(min <= l && l <= max) {
			msg := fmt.Sprintf("%d..%d args required", min, max)
			return cli.NewExitError(msg, 1)
		}
	}

	return nil
}

func buildApp(c *cli.Context) (*App, error) {
	b := c.String("bucket")
	p := c.String("prefix")
	k := c.String("key")

	if b == "" {
		return nil, cli.NewExitError("S3 Bucket is not specified", 1)
	}

	return NewApp(b, p, k), nil
}

func wrapError(err error) error {
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}