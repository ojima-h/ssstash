package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:   "bucket",
		EnvVar: "SSSTASH_S3_BUCKET",
	},
	cli.StringFlag{
		Name:   "prefix",
		EnvVar: "SSSTASH_S3_PREFIX",
	},
	cli.StringFlag{
		Name: "profile",
	},
}

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Flags:   flags,
			Action: func(c *cli.Context) error {
				app, err := buildApp(c)
				if err != nil {
					return err
				}

				err = app.ListIter(func(name string) bool {
					fmt.Println(name)
					return true
				})
				return wrapError(err)
			},
		},
		{
			Name: "put",
			Flags: append(
				flags,
				cli.StringFlag{
					Name:   "key",
					EnvVar: "SSSTASH_KEY_ID",
				},
			),
			Action: func(c *cli.Context) error {
				if err := validateArgsLength(c, 2, 2); err != nil {
					return err
				}
				name := c.Args().Get(0)
				val := c.Args().Get(1)

				keyID := c.String("key")
				if keyID == "" {
					return fmt.Errorf("key ID is not specified")
				}

				app, err := buildApp(c)
				if err != nil {
					return err
				}

				return wrapError(app.Put(name, val, keyID))
			},
		},
		{
			Name:  "get",
			Flags: flags,
			Action: func(c *cli.Context) error {
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
			Name:    "delete",
			Aliases: []string{"rm"},
			Flags:   flags,
			Action: func(c *cli.Context) error {
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
	bucket := c.String("bucket")
	prefix := c.String("prefix")
	profile := c.String("profile")

	if bucket == "" {
		return nil, cli.NewExitError("S3 Bucket is not specified", 1)
	}

	return NewApp(bucket, prefix, profile), nil
}

func wrapError(err error) error {
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}
