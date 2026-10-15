CREATE TABLE IF NOT EXISTS `inboxmsgs_split` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `send_time` bigint NOT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(20) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `msg_type` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`,`send_time`),
  KEY `IDX_USERID_MSG` (`app_key`,`user_id`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`user_id`,`msg_id`),
  KEY `idx_appkey` (`app_key`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
PARTITION BY RANGE (`send_time`) (
  PARTITION `p_future` VALUES LESS THAN MAXVALUE
);
