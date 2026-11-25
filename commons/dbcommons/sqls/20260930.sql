CREATE TABLE IF NOT EXISTS `sendboxmsgs_split` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `send_time` bigint NOT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(20) DEFAULT NULL,
  `target_id` varchar(45) DEFAULT NULL,
  `msg_type` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`,`send_time`),
  KEY `idx_user_id_send_time` (`app_key`,`user_id`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`user_id`,`msg_id`),
  KEY `idx_appkey` (`app_key`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
PARTITION BY RANGE (`send_time`) (
  PARTITION `p_future` VALUES LESS THAN MAXVALUE
);

CREATE TABLE IF NOT EXISTS `cmdinboxmsgs_split` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `send_time` bigint NOT NULL,
  `partition_day` bigint NOT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(20) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `msg_type` varchar(50) DEFAULT NULL,
  `uniq_tag` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`,`partition_day`),
  UNIQUE KEY `uniq_tag` (`app_key`,`user_id`,`uniq_tag`,`partition_day`),
  KEY `idx_appkey_time` (`app_key`,`user_id`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`user_id`,`msg_id`),
  KEY `idx_appkey` (`app_key`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
PARTITION BY RANGE (`partition_day`) (
  PARTITION `p_future` VALUES LESS THAN MAXVALUE
);

CREATE TABLE IF NOT EXISTS `cmdsendboxmsgs_split` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `send_time` bigint NOT NULL,
  `partition_day` bigint NOT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(20) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `msg_type` varchar(50) DEFAULT NULL,
  `uniq_tag` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`,`partition_day`),
  UNIQUE KEY `uniq_tag` (`app_key`,`user_id`,`uniq_tag`,`partition_day`),
  KEY `idx_appkey_userid_time` (`app_key`,`user_id`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`user_id`,`msg_id`),
  KEY `idx_appkey` (`app_key`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
PARTITION BY RANGE (`partition_day`) (
  PARTITION `p_future` VALUES LESS THAN MAXVALUE
);
