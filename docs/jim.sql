

DROP TABLE IF EXISTS `accounts`;
CREATE TABLE `accounts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `account` varchar(45) DEFAULT NULL,
  `password` varchar(45) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `state` tinyint DEFAULT '0',
  `parent_account` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_account` (`account`),
  KEY `idx_parent` (`parent_account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `androidpushconfs`;
CREATE TABLE `androidpushconfs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `app_key` varchar(20) DEFAULT NULL,
  `push_channel` varchar(10) DEFAULT NULL,
  `package` varchar(100) DEFAULT NULL,
  `push_conf` varchar(500) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_channel` (`app_key`,`push_channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `appexts`;
CREATE TABLE `appexts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `app_key` varchar(50) DEFAULT NULL,
  `app_item_key` varchar(50) DEFAULT NULL,
  `app_item_value` varchar(2048) DEFAULT NULL,
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_APPKEY_APPITEMKEY` (`app_key`,`app_item_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `apps`;
CREATE TABLE `apps` (
  `id` int NOT NULL AUTO_INCREMENT,
  `app_key` varchar(45) NOT NULL,
  `app_secret` varchar(45) NOT NULL,
  `app_secure_key` varchar(45) NOT NULL,
  `app_status` tinyint DEFAULT '0',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `app_type` tinyint DEFAULT '0',
  `app_name` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `banusers`;
CREATE TABLE `banusers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) NOT NULL,
  `ban_type` tinyint DEFAULT '0',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `end_time` bigint DEFAULT '0',
  `scope_key` varchar(20) NOT NULL DEFAULT 'default',
  `scope_value` varchar(1000) DEFAULT '',
  `ext` varchar(100) DEFAULT NULL,
  `app_key` varchar(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey_userid` (`app_key`,`user_id`,`scope_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `bc_hismsgs`;
CREATE TABLE `bc_hismsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `conver_id` varchar(100) NOT NULL,
  `sender_id` varchar(32) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_type` varchar(20) DEFAULT NULL,
  `msg_id` varchar(20) NOT NULL,
  `send_time` bigint DEFAULT NULL,
  `msg_seq_no` int DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_msgid` (`app_key`,`conver_id`,`msg_id`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `blocks`;
CREATE TABLE `blocks` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `block_user_id` varchar(32) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey_userid` (`app_key`,`user_id`,`block_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `brdinboxmsgs`;
CREATE TABLE `brdinboxmsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `sender_id` varchar(32) DEFAULT NULL,
  `send_time` bigint DEFAULT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_sendtime` (`app_key`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `chatroominfos`;
CREATE TABLE `chatroominfos` (
  `id` int NOT NULL AUTO_INCREMENT,
  `chat_id` varchar(32) DEFAULT NULL,
  `chat_name` varchar(45) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_chatid` (`app_key`,`chat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `cmdinboxmsgs`;
CREATE TABLE `cmdinboxmsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `send_time` bigint DEFAULT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(20) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `msg_type` varchar(20) DEFAULT NULL,
  `uniq_tag` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_tag` (`app_key`,`user_id`,`uniq_tag`),
  KEY `idx_appkey_time` (`app_key`,`user_id`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`user_id`,`msg_id`),
  KEY `idx_appkey` (`app_key`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `cmdsendboxmsgs`;
CREATE TABLE `cmdsendboxmsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `send_time` bigint DEFAULT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(20) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `msg_type` varchar(20) DEFAULT NULL,
  `uniq_tag` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_tag` (`app_key`,`user_id`,`uniq_tag`),
  KEY `idx_appkey_userid_time` (`app_key`,`user_id`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`user_id`,`msg_id`),
  KEY `idx_appkey` (`app_key`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `convercleantimes`;
CREATE TABLE `convercleantimes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `conver_id` varchar(100) DEFAULT NULL,
  `channel_type` tinyint DEFAULT '0',
  `clean_time` bigint DEFAULT '0',
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_destroy` (`app_key`,`conver_id`,`channel_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `conversations`;
CREATE TABLE `conversations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `channel_type` tinyint DEFAULT '0',
  `latest_msg_id` varchar(20) DEFAULT NULL,
  `latest_msg` mediumblob,
  `latest_unread_msg_index` int DEFAULT '0',
  `latest_read_msg_index` int DEFAULT '0',
  `latest_read_msg_id` varchar(20) DEFAULT NULL,
  `latest_read_msg_time` bigint DEFAULT '0',
  `sort_time` bigint DEFAULT '0',
  `is_deleted` tinyint DEFAULT '0',
  `is_top` tinyint DEFAULT '0',
  `top_updated_time` bigint DEFAULT '0',
  `undisturb_type` tinyint DEFAULT '0',
  `sync_time` bigint DEFAULT '0',
  `unread_tag` tinyint DEFAULT '0',
  `group` varchar(20) DEFAULT NULL,
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_app_key_user_id_target_id` (`app_key`,`user_id`,`target_id`,`channel_type`),
  KEY `idx_sync_time` (`app_key`,`user_id`,`sync_time`),
  KEY `idx_update_time` (`app_key`,`user_id`,`sort_time`),
  KEY `idx_undisturb` (`app_key`,`user_id`,`target_id`,`channel_type`),
  KEY `idx_group` (`app_key`,`user_id`,`group`,`sort_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `msgstats`;
CREATE TABLE `msgstats` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `stat_type` TINYINT NULL DEFAULT 0,
  `channel_type` TINYINT NULL,
  `time_mark` BIGINT NULL,
  `count` INT NULL,
  `app_key` VARCHAR(20) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uniq_mark` (`app_key` ASC, `stat_type` ASC, `channel_type` ASC, `time_mark` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `useractivities`;
CREATE TABLE `useractivities` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `user_id` VARCHAR(32) NULL,
  `time_mark` BIGINT NULL,
  `created_time` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
  `count` INT NULL,
  `app_key` VARCHAR(20) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uniq_userid` (`app_key` ASC, `time_mark` ASC, `user_id` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `fileconfs`;
CREATE TABLE `fileconfs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `app_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `channel` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `conf` json DEFAULT NULL,
  `enable` tinyint(1) DEFAULT '0',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `app_key` (`app_key`,`channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `g_delhismsgs`;
CREATE TABLE `g_delhismsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `msg_time` bigint DEFAULT NULL,
  `msg_seq` int DEFAULT NULL,
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_msgid` (`app_key`,`user_id`,`target_id`,`msg_id`),
  KEY `idx_target` (`app_key`,`user_id`,`target_id`,`msg_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `g_hismsgs`;
CREATE TABLE `g_hismsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `conver_id` varchar(100) DEFAULT NULL,
  `sender_id` varchar(32) DEFAULT NULL,
  `receiver_id` varchar(32) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_type` varchar(45) DEFAULT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `send_time` bigint DEFAULT NULL,
  `msg_seq_no` int DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(20) DEFAULT NULL,
  `member_count` int DEFAULT '0',
  `read_count` int DEFAULT '0',
  `is_delete` tinyint DEFAULT '0',
  `is_ext` tinyint DEFAULT '0',
  `is_reaction` tinyint DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_appkey_converid` (`app_key`,`conver_id`,`send_time`),
  KEY `idx_msgid` (`app_key`,`conver_id`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `gc_hismsgs`;
CREATE TABLE `gc_hismsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `conver_id` varchar(100) DEFAULT NULL,
  `sender_id` varchar(32) DEFAULT NULL,
  `receiver_id` varchar(32) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_type` varchar(45) DEFAULT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `send_time` bigint DEFAULT NULL,
  `msg_seq_no` int DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_msg` (`app_key`,`conver_id`,`channel_type`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `globalconfs`;
CREATE TABLE `globalconfs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `conf_key` varchar(50) DEFAULT NULL,
  `conf_value` varchar(2000) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_key` (`conf_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `globalconvers`;
CREATE TABLE `globalconvers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `conver_id` varchar(100) DEFAULT NULL,
  `sender_id` varchar(32) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `updated_time` bigint DEFAULT NULL,
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_conver` (`app_key`,`conver_id`,`channel_type`),
  KEY `idx_time` (`app_key`,`channel_type`,`updated_time`),
  KEY `idx_targetid` (`app_key`,`channel_type`,`target_id`,`updated_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `groupinfoexts`;
CREATE TABLE `groupinfoexts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `group_id` varchar(32) DEFAULT NULL,
  `item_key` varchar(50) DEFAULT NULL,
  `item_value` varchar(100) DEFAULT NULL,
  `item_type` tinyint DEFAULT '0',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey_groupid` (`app_key`,`group_id`,`item_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `groupinfos`;
CREATE TABLE `groupinfos` (
  `id` int NOT NULL AUTO_INCREMENT,
  `group_id` varchar(64) DEFAULT NULL,
  `group_name` varchar(64) DEFAULT NULL,
  `group_portrait` varchar(200) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  `is_mute` tinyint DEFAULT '0',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey_groupid` (`app_key`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `groupmemberexts`;
CREATE TABLE `groupmemberexts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `group_id` varchar(32) DEFAULT NULL,
  `member_id` varchar(32) DEFAULT NULL,
  `item_key` varchar(50) DEFAULT NULL,
  `item_value` varchar(100) DEFAULT NULL,
  `item_type` tinyint DEFAULT '0',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_item_key` (`app_key`,`group_id`,`member_id`,`item_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `groupmembers`;
CREATE TABLE `groupmembers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `group_id` varchar(64) DEFAULT NULL,
  `member_id` varchar(64) DEFAULT NULL,
  `member_type` tinyint DEFAULT '0',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(45) DEFAULT NULL,
  `is_mute` tinyint DEFAULT '0',
  `is_allow` tinyint DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey_grpid_memid` (`app_key`,`group_id`,`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `grpassistantrels`;
CREATE TABLE `grpassistantrels` (
  `id` int NOT NULL,
  `assistant_id` varchar(32) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_target` (`assistant_id`,`target_id`,`channel_type`,`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `grpsnapshots`;
CREATE TABLE `grpsnapshots` (
  `id` int NOT NULL AUTO_INCREMENT,
  `app_key` varchar(20) NOT NULL,
  `group_id` varchar(32) NOT NULL,
  `created_time` bigint DEFAULT '0',
  `snapshot` mediumblob,
  PRIMARY KEY (`id`),
  KEY `idx_group_id` (`app_key`,`group_id`,`created_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `ic_conditions`;
CREATE TABLE `ic_conditions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `channel_type` varchar(100) DEFAULT NULL,
  `msg_type` varchar(1000) DEFAULT NULL,
  `sender_id` varchar(1000) DEFAULT NULL,
  `receiver_id` varchar(1000) DEFAULT NULL,
  `interceptor_id` int DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_icid` (`app_key`,`interceptor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `inboxmsgs`;
CREATE TABLE `inboxmsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `send_time` bigint DEFAULT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(20) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `msg_type` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_USERID_MSG` (`app_key`,`user_id`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`user_id`,`msg_id`),
  KEY `idx_appkey` (`app_key`,`send_time`)
) ENGINE=InnoDB AUTO_INCREMENT=45612 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `interceptors`;
CREATE TABLE `interceptors` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) DEFAULT NULL,
  `sort` int NOT NULL DEFAULT '0',
  `request_url` varchar(500) DEFAULT NULL,
  `request_template` text,
  `succ_template` varchar(200) DEFAULT NULL,
  `is_async` tinyint DEFAULT '0',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_sort` (`app_key`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `ioscertificates`;
CREATE TABLE `ioscertificates` (
  `id` int NOT NULL AUTO_INCREMENT,
  `package` varchar(100) DEFAULT NULL,
  `certificate` mediumblob,
  `cert_pwd` varchar(50) DEFAULT NULL,
  `app_key` varchar(20) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `is_product` tinyint DEFAULT '0',
  `cert_path` varchar(255) DEFAULT NULL,
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_package` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `mentionmsgs`;
CREATE TABLE `mentionmsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `sender_id` varchar(32) DEFAULT NULL,
  `mention_type` tinyint DEFAULT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `msg_time` bigint DEFAULT NULL,
  `msg_index` int DEFAULT NULL,
  `is_read` tinyint DEFAULT NULL,
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_uid_tid_type` (`app_key`,`user_id`,`target_id`,`channel_type`,`msg_index`,`msg_time`),
  KEY `idx_read` (`app_key`,`user_id`,`target_id`,`channel_type`,`is_read`,`msg_time`),
  KEY `idx_user_msgid` (`app_key`,`user_id`,`target_id`,`channel_type`,`msg_id`),
  KEY `idx_target_msgid` (`app_key`, `msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `mergedmsgs`;
CREATE TABLE `mergedmsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `parent_msg_id` varchar(20) DEFAULT NULL,
  `from_id` varchar(32) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `msg_time` bigint DEFAULT '0',
  `msg_body` mediumblob,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_appkey_pmsg` (`app_key`,`parent_msg_id`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `msgexts`;
CREATE TABLE `msgexts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `msg_id` varchar(20) DEFAULT NULL,
  `key` varchar(50) DEFAULT NULL,
  `value` varchar(1000) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_msgid` (`app_key`,`msg_id`,`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `p_delhismsgs`;
CREATE TABLE `p_delhismsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `msg_time` bigint DEFAULT NULL,
  `msg_seq` int DEFAULT NULL,
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_msgid` (`app_key`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `p_hismsgs`;
CREATE TABLE `p_hismsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `conver_id` varchar(100) DEFAULT NULL,
  `sender_id` varchar(32) DEFAULT NULL,
  `receiver_id` varchar(32) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_type` varchar(45) DEFAULT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `send_time` bigint DEFAULT NULL,
  `msg_seq_no` int DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(20) DEFAULT NULL,
  `is_read` tinyint DEFAULT '0',
  `is_delete` tinyint DEFAULT '0',
  `is_ext` tinyint DEFAULT '0',
  `is_reaction` tinyint DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_app_key_conver_id` (`app_key`,`conver_id`,`send_time`),
  KEY `idx_msgid` (`app_key`,`conver_id`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `pushtokens`;
CREATE TABLE `pushtokens` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `device_id` varchar(200) DEFAULT NULL,
  `platform` varchar(10) DEFAULT NULL,
  `push_channel` varchar(10) DEFAULT NULL,
  `package` varchar(200) DEFAULT NULL,
  `push_token` varchar(200) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`app_key`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `readinfos`;
CREATE TABLE `readinfos` (
  `id` int NOT NULL AUTO_INCREMENT,
  `app_key` varchar(20) NOT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `group_id` varchar(32) DEFAULT NULL,
  `member_id` varchar(32) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_member` (`app_key`,`channel_type`,`group_id`,`msg_id`,`member_id`),
  KEY `idx_memberid` (`app_key`,`channel_type`,`group_id`,`member_id`,`msg_id`),
  KEY `idx_msgid` (`app_key`,`channel_type`,`group_id`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `s_hismsgs`;
CREATE TABLE `s_hismsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `conver_id` varchar(100) DEFAULT NULL,
  `sender_id` varchar(32) DEFAULT NULL,
  `receiver_id` varchar(32) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_type` varchar(45) DEFAULT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `send_time` bigint DEFAULT NULL,
  `msg_seq_no` int DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(20) DEFAULT NULL,
  `is_read` tinyint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_appkey_converid` (`app_key`,`conver_id`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `sendboxmsgs`;
CREATE TABLE `sendboxmsgs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `send_time` bigint DEFAULT NULL,
  `msg_id` varchar(20) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `msg_body` mediumblob,
  `app_key` varchar(20) DEFAULT NULL,
  `target_id` varchar(45) DEFAULT NULL,
  `msg_type` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id_send_time` (`app_key`,`user_id`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`user_id`,`msg_id`),
  KEY `idx_appkey` (`app_key`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `sensitivewords`;
CREATE TABLE `sensitivewords` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `app_key` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `word` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `word_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '12',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_word` (`app_key`,`word`),
  KEY `idx_appkey` (`app_key`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `subrelations`;
CREATE TABLE `subrelations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `subscriber` varchar(32) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_sub` (`app_key`,`user_id`,`subscriber`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `usercleantimes`;
CREATE TABLE `usercleantimes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `target_id` varchar(32) DEFAULT NULL,
  `channel_type` tinyint DEFAULT NULL,
  `clean_time` bigint DEFAULT NULL,
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_app_key_user_id_target_id` (`app_key`,`user_id`,`target_id`,`channel_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `userexts`;
CREATE TABLE `userexts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `item_key` varchar(50) DEFAULT NULL,
  `item_value` varchar(100) DEFAULT NULL,
  `item_type` tinyint DEFAULT '0',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_item_key` (`app_key`,`user_id`,`item_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_type` tinyint DEFAULT '0',
  `user_id` varchar(32) NOT NULL,
  `nickname` varchar(45) DEFAULT NULL,
  `user_portrait` varchar(200) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(45) DEFAULT NULL,
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_userid` (`app_key`,`user_id`),
  KEY `idx_userid` (`app_key`,`user_type`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `clientlogs`;
CREATE TABLE `clientlogs` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `app_key` VARCHAR(20) NULL,
  `user_id` VARCHAR(32) NULL,
  `created_time` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
  `start` BIGINT NULL,
  `end` BIGINT NULL,
  `log` MEDIUMBLOB NULL,
  `state` TINYINT NULL DEFAULT 0,
   `platform` VARCHAR(20) NULL,
  `device_id` VARCHAR(100) NULL,
  `log_url` VARCHAR(200) NULL,
  `trace_id` VARCHAR(50) NULL,
   `msg_id` VARCHAR(20) NULL,
  `fail_reason` VARCHAR(100) NULL,
  `description` VARCHAR(100) NULL,
  
  PRIMARY KEY (`id`),
  INDEX `idx_userid` (`app_key` ASC, `user_id` ASC),
  UNIQUE KEY `uniq_msgid` (`app_key`,`msg_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT IGNORE INTO `accounts`(`account`,`password`)VALUES('admin1','7c4a8d09ca3762af61e59520943dc26494f8941b');