CREATE TABLE users (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT "ユーザーの識別子",
  name       VARCHAR(20) NOT NULL COMMENT "ユーザー名",
  password   VARCHAR(80) NOT NULL COMMENT "パスワード",
  role       VARCHAR(80) NOT NULL COMMENT "役割",
  created_at DATETIME(6) NOT NULL COMMENT "レコード作成日時",
  updated_at DATETIME(6) NOT NULL COMMENT "レコード修正日時",
  PRIMARY KEY(id),
  UNIQUE KEY uix_name (name) USING BTREE
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT="ユーザー";

CREATE TABLE tasks (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT "タスクの識別子",
  title      VARCHAR(128) NOT NULL COMMENT "タスクのタイトル",
  status     VARCHAR(20) NOT NULL COMMENT "タスクの状態",
  created_at DATETIME(6) NOT NULL COMMENT "レコード作成日時",
  updated_at DATETIME(6) NOT NULL COMMENT "レコード修正日時",
  PRIMARY KEY(id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT="タスク";

