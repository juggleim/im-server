CREATE TABLE IF NOT EXISTS `banuserrules` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) NOT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `end_time` bigint NOT NULL DEFAULT '0',
  `scope_key` varchar(20) NOT NULL DEFAULT 'default',
  `scope_value` varchar(1000) NOT NULL DEFAULT '',
  `scope_value_hash` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  `ext` varchar(100) DEFAULT NULL,
  `app_key` varchar(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_ban_rule` (`app_key`,`user_id`,`scope_key`,`scope_value_hash`),
  KEY `idx_ban_rule_query` (`app_key`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
