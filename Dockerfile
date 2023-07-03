# デプロイ用コンテナに含めるバイナリを作成するコンテナ
FROM golang:1.18.2-bullseye AS deploy-builder

# WORKDIR命令は、RUN, CMD, ENTRYPOINT, ADD, COPYの際の、コンテナ内の作業ディレクトリを指定する
# WORKDIR が存在していなければ作成される
# 基本、命令で指定したコマンドを実行する場所は、イメージレイヤのルートだが、WORKDIR命令を使うと、次のイメージレイヤでもWORKDIRが引き継がれる。
WORKDIR /app

# COPY 命令では、追加したいファイル、ディレクトリを <コピー元> で指定すると、これらをイメージのファイルシステム上のパス <コピー先> に追加する
# COPY命令は、Dockerfileを置いたディレクトリ内のファイルやフォルダをイメージレイヤーの特定のディレクトリにコピーしたいときに使う。
# ADD命令と似ているがADD命令は挙動がわかりづらいためCOPY命令のが推奨されているそう
# 今回の場合、ホストのDockerfileがあるディレクトリにあるファイルを、appディレクトリにコピーする
COPY go.mod god.sum ./

# RUN 命令は、現在のイメージよりも上にある新しいレイヤでコマンドを実行し、その結果を コミット（確定）commit する。
# RUN命令はdocker buildを実行するタイミング(つまり、イメージを生成するとき)に実行する。
# イメージの時点で実行しておきたいコマンド（パッケージやライブラリのインストール、ファイルコピー、変更など）を実行するために書く。
# go mod downloadは、引数にモジュールの指定がない場合は、go.modファイルに記載されたすべてのモジュールをダウンロードする
RUN go mod download

# 自分のホストのルートディレクトリにあるファイルを、WORKDIRにコピーしていると思われる
COPY . .

# go build で ldflags オプションと-w, -sを指定するとシンボルやデバッグ情報を削ることができる。
# その結果、ビルドされたシングルバイナリファイルのビルドサイズを小さくすることができる
# trimpathは、panicしたときのスタックトレースのパス表示に、ビルド時の実パスがでなくなる。
# -oオプションは、実行結果を指定したファイルやディレクトリ配下に書き込むためのオプション。スラッシュで終わっていたらディレクトリを表す
RUN go build -trimpath -ldflags "-w -s" -o app

# ---------------------------------------

# デプロイ用のコンテナ
FROM debian:bullseye-slim as deploy

# apt-get updateは、インストール可能なパッケージの「一覧」を更新する。
# 実際のパッケージのインストール、アップグレードなどはおこなわない。
RUN apt-get update

# COPY --from命令では別のイメージからコピーすることができる。
# 以下では、イメージのappディレクトリ配下にあるappファイルをイメージのカレントディレクトリにコピーしている
COPY --from=deploy-builder /app/app .

# CMD命令は、コンテナ実行時に、コンテナが見えるファイルシステム上であるプログラムを自動で実行しなさいという命令
# 以下の場合、コンテナを立ち上げる際に、シングルバイナリファイルが実行される
CMD ["./apps"]

# ---------------------------------------

# ローカル開発環境で利用するホットリロード環境
FROM golang:1.18.2 AS dev
# go installは指定したパッケージをダウンロード後にビルドし、実行可能なファイルを$GOBINへ格納する
RUN go install github.com/cosmtrek/air@latest

WORKDIR /app

# air -c [tomlファイル名] // 設定ファイルを指定してair実行(WORKDIRに.air.tomlを配置しておくこと)
# -c .ari.toml のオプションは省略可能。その場合、air はカレントディレクトリから .air.toml ファイルを探し起動します。
CMD ["air"]