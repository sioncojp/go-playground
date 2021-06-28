# s3-move-other-bucket

- s3を他バケットにSTANDARD_IAにしてコピーする
- --deleteがついてる場合は、転送元のファイルを削除する
- --checkがついてる場合は、転送元/先のオブジェクト数を出力する

```shell
# bucket構成
bucket/
├── 1
│   ├── a
│   │   ├── aa
│   │   │   └── hoge.txt
│   │   └── bb
│   │       └── fuga.txt
│   ├── b
│   ├── c
│   └── d
├── 2
│   ├── a
│   │   └── bb
│   │       └── piyo.txt
│   ├── b
│   ├── c
│   └── d
.
.
.
└── 10000
│   ├── a
│   │   │   └── foo.txt
│   │   └── bb
│   │       └── bar.txt
│   ├── b
│   ├── c
│   └── d
```

# help

```
$ ./bin/s3-move-other-bucket help
NAME:
   s3-move-other-bucket - A new cli application

USAGE:
   s3-move-other-bucket [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --awsenv value     aws profile name
   --region value     Region (e.g. ap-northeast-1) (default: "ap-northeast-1")
   --parallel value   並列数 (default: 1)
   --src value        転送元バケット名
   --dest value       転送先先バケット名
   --dir value        転送元バケット内のディレクトリ。今回はs3://bucket/id/ここ (default: "a")
   --id value         単一のidに対して実行する
   --beforeday value  何日前のデータを削除するか (default: 0)
   --delete           転送元バケットの対象ディレクトリを削除する
   --check            copyできたかチェックする
   --help, -h         show help
Required flags "src, dest" not set
```

# Get Started
```shell
$ make build

# copy
$ ./bin/s3-move-other-bucket --src src_bucket --dest dest_bucket --parallel 10

# delete
$ ./bin/s3-move-other-bucket --src src_bucket --dest dest_bucket --paralell 10 --delete

# check
$ ./bin/s3-move-other-bucket --src src_bucket --dest dest_bucket --parallel 10 --check

# copy before 10day
$ ./bin/s3-move-other-bucket --src src_bucket --dest dest_bucket --parallel 10 --beforeday 10
```

# cross compile
```shell
$ make build/cross
```
