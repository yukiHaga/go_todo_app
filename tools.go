// go:build tools
// 上のコメントはビルドタグ。ビルド時にtoolsを渡すと、このファイルがビルド対象になる

package main

// このように該当ツールをimportしたtools.goファイルを定義しておくことで、go.modによるバージョン管理ができる
import _ "github.com/matryer/moq"
