//go:build tools

package main

// 上のコメントはビルドタグ。ビルド時にtoolsを渡すと、このファイルがビルド対象になる
// スペースを入れちゃダメ。

// このように該当ツールをimportしたtools.goファイルを定義しておくことで、go.modによるバージョン管理ができる
import _ "github.com/matryer/moq"
