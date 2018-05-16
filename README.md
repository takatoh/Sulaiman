# Sulaimān

A simple photo uploader.

シンプルな画像アップローダーです。スライマーンと読んでください。

## Install
Github からクローンして
``` clone https://github.com/takatoh/sulaiman.git```
依存ライブラリをインストールします。
``` dep ensure```
go build します。
``` go build```

## Usage
ビルドしてできた実行ファイルと config.json.example、それから static をディレクトリごと
適当なディレクトリにコピーします。
confing.json.exmaple を config.json にリネームして、適当に編集します。
```JSON
{
    "site_name" : "Sulaimān",
    "host_name" : "localhost:1323",
    "port" : 1323,
    "Photo_dir" : "photos"
}
```
画像保存用のディレクトリを作ります。上の例では photos/img と photos/thumb が必要です。
あとは実行するだけです。
``` sulaiman```

## License
MIT License
