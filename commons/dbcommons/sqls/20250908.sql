ALTER TABLE `g_hismsgs` 
ADD COLUMN `is_portion` TINYINT NULL DEFAULT 0;

CREATE TABLE `g_portionrels` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `conver_id` VARCHAR(100) NULL,
  `sub_channel` VARCHAR(32) NULL DEFAULT '',
  `channel_type` TINYINT NULL DEFAULT 0,
  `user_id` VARCHAR(32) NULL,
  `msg_id` VARCHAR(32) NULL,
  `msg_time` BIGINT NULL DEFAULT 0,
  `app_key` VARCHAR(20) NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uniq_msgid` (`app_key`, `conver_id`, `sub_channel`, `user_id`, `msg_id`),
  INDEX `idx_msg_time` (`app_key`, `conver_id`, `sub_channel`, `user_id`, `msg_time`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `globalconfs` (`conf_key`,`conf_value`)VALUES('jimdb_version','20250908') ON DUPLICATE KEY UPDATE conf_value=VALUES(conf_value);