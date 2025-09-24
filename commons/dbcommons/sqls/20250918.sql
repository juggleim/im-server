CREATE TABLE IF NOT EXISTS `accountapprels` (
  `id` int NOT NULL AUTO_INCREMENT,
  `app_key` varchar(20) DEFAULT '',
  `account_id` int DEFAULT '0',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_app` (`account_id`,`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

ALTER TABLE `accounts` 
ADD COLUMN `role_type` TINYINT NULL DEFAULT 0;

INSERT INTO `globalconfs` (`conf_key`,`conf_value`)VALUES('jimdb_version','20250918') ON DUPLICATE KEY UPDATE conf_value=VALUES(conf_value);