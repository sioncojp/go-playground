package main

import (
	"github.com/urfave/cli"
)

// FlagSet ... flagオプションを設定
func FlagSet() *cli.App {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "awsenv",
			Usage: "aws profile name",
		},
		cli.StringFlag{
			Name:  "region",
			Value: "ap-northeast-1",
			Usage: "Region (e.g. ap-northeast-1)",
		},
		cli.Int64Flag{
			Name:  "parallel",
			Value: 1,
			Usage: "並列数",
		},
		cli.StringFlag{
			Name:     "src",
			Required: true,
			Usage:    "転送元バケット名",
		},
		cli.StringFlag{
			Name:     "dest",
			Required: true,
			Usage:    "転送先先バケット名",
		},
		cli.StringFlag{
			Name:  "dir",
			Value: "a",
			Usage: "転送元バケット内のディレクトリ。今回はs3://bucket/id/ここ",
		},
		cli.StringFlag{
			Name:  "id",
			Value: "",
			Usage: "単一のidに対して実行する",
		},
		cli.IntFlag{
			Name:  "beforeday",
			Value: 0,
			Usage: "何日前のデータを削除するか",
		},
		cli.BoolFlag{
			Name:  "delete",
			Usage: "転送元バケットの対象ディレクトリを削除する",
		},
		cli.BoolFlag{
			Name:  "check",
			Usage: "copyできたかチェックする",
		},
	}

	app.Action = func(c *cli.Context) error {
		f := Flag{
			c.String("awsenv"),
			c.String("region"),
			c.Int64("parallel"),
			c.String("src"),
			c.String("dest"),
			c.String("dir"),
			c.String("id"),
			c.Int("beforeday"),
			c.Bool("delete"),
			c.Bool("check"),
		}

		if err := S3MoveOtherBucket(f); err != nil {
			return err
		}

		return nil
	}

	return app
}
