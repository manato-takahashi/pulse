# pulse

複数のエンドポイントを指定して、並行でヘルスチェックを行うことができるシンプルなヘルスチェックCLIツール

## Install

コンピュータにGoがインストールされていれば、1コマンドでインストール完了

```bash
go install github.com/manato-takahashi/pulse@latest
```

`pulse` コマンドが見つからない場合は、`$GOPATH/bin` にパスを通す:

```bash
export PATH="$HOME/go/bin:$PATH"
```

## Usage

### ローカルCLI

```bash
# --config フラグで使用する設定ファイルのパスを指定可能
$ pulse --config endpoints.sample.yaml
✓ https://example.com              200  42ms
✓ https://httpbin.org/status/200   200  661ms
× https://httpbin.org/status/503   503  662ms
× https://httpbin.org/status/401   401  661ms
```

### CI/CD

```bash
# --exit-on-fail フラグを付けることで、1つでもヘルスチェックが失敗するエンドポイントがあればエラーを返す(CI上で検知可能)
pulse --config endpoints.sample.yaml --exit-on-fail
```

## Configuration

```yaml
# endpoints.sample.yaml
endpoints:
  - url: https://example.com
  - url: https://httpbin.org/status/200
  - url: https://httpbin.org/status/503
  - url: https://httpbin.org/status/401
```

## Features

- YAML設定ファイルで複数エンドポイントを一括チェック
- goroutineによる並行リクエストで高速実行
- `--exit-on-fail` フラグでCI/CDパイプラインに組み込み可能
