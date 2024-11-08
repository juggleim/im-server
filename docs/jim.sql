

DROP TABLE IF EXISTS `accounts`;
CREATE TABLE `accounts` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `account` varchar(45) DEFAULT NULL COMMENT '登录名',
  `password` varchar(45) DEFAULT NULL COMMENT '密码',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '修改时间',
  `state` tinyint DEFAULT '0' COMMENT '状态',
  `parent_account` varchar(45) DEFAULT NULL COMMENT '父账号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_account` (`account`),
  KEY `idx_parent` (`parent_account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '管理端-账号';

DROP TABLE IF EXISTS `androidpushconfs`;
CREATE TABLE `androidpushconfs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `app_key` varchar(20) DEFAULT NULL COMMENT 'appKey',
  `push_channel` varchar(10) DEFAULT NULL COMMENT '推送渠道',
  `package` varchar(100) DEFAULT NULL COMMENT 'package',
  `push_conf` varchar(500) DEFAULT NULL COMMENT '配置参数',
  `push_ext` MEDIUMBLOB DEFAULT NULL COMMENT '其他配置',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_channel` (`app_key`,`push_channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'android推送配置';

DROP TABLE IF EXISTS `appexts`;
CREATE TABLE `appexts` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `app_key` varchar(50) DEFAULT NULL COMMENT '应用key',
  `app_item_key` varchar(50) DEFAULT NULL COMMENT '参数key',
  `app_item_value` varchar(2048) DEFAULT NULL COMMENT '参数value',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_APPKEY_APPITEMKEY` (`app_key`,`app_item_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'app应用-补充配置';

DROP TABLE IF EXISTS `apps`;
CREATE TABLE `apps` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `app_key` varchar(45) NOT NULL COMMENT '应用key',
  `app_secret` varchar(45) NOT NULL COMMENT '应用密钥 16位',
  `app_secure_key` varchar(45) NOT NULL COMMENT '安全key 16位',
  `app_status` tinyint DEFAULT '0' COMMENT '状态',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `app_type` tinyint DEFAULT '0' COMMENT '应用类型',
  `app_name` varchar(100) DEFAULT NULL COMMENT '应用名称',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'app应用';

DROP TABLE IF EXISTS `banusers`;
CREATE TABLE `banusers` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) NOT NULL COMMENT '用户id',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `end_time` bigint DEFAULT '0' COMMENT '结束时间',
  `scope_key` varchar(20) NOT NULL DEFAULT 'default',
  `scope_value` varchar(1000) DEFAULT '',
  `ext` varchar(100) DEFAULT NULL COMMENT '其他配置',
  `app_key` varchar(20) NOT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey_userid` (`app_key`,`user_id`,`scope_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '禁用用户';

DROP TABLE IF EXISTS `bc_hismsgs`;
CREATE TABLE `bc_hismsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `conver_id` varchar(100) NOT NULL COMMENT '会话id',
  `sender_id` varchar(32) DEFAULT NULL COMMENT '发送者id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_type` varchar(50) DEFAULT NULL COMMENT '消息类型',
  `msg_id` varchar(20) NOT NULL COMMENT '消息id',
  `send_time` bigint DEFAULT NULL COMMENT '发送时间',
  `msg_seq_no` int DEFAULT NULL COMMENT '消息seq',
  `msg_body` mediumblob COMMENT '消息体',
  `app_key` varchar(20) NOT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  KEY `idx_msgid` (`app_key`,`conver_id`,`msg_id`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '广播-历史消息表';

DROP TABLE IF EXISTS `blocks`;
CREATE TABLE `blocks` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) DEFAULT NULL COMMENT '操作人',
  `block_user_id` varchar(32) DEFAULT NULL   COMMENT '被锁定者id',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey_userid` (`app_key`,`user_id`,`block_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '锁定记录';

DROP TABLE IF EXISTS `brdinboxmsgs`;
CREATE TABLE `brdinboxmsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `sender_id` varchar(32) DEFAULT NULL COMMENT '发送者id',
  `send_time` bigint DEFAULT NULL COMMENT '发送时间',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_body` mediumblob COMMENT '消息体',
  `app_key` varchar(32) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  KEY `idx_sendtime` (`app_key`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '广播-收件箱';

DROP TABLE IF EXISTS `chatroominfos`;
CREATE TABLE `chatroominfos` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `chat_id` varchar(32) DEFAULT NULL COMMENT '聊天室id',
  `chat_name` varchar(45) DEFAULT NULL COMMENT '聊天室名称',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_chatid` (`app_key`,`chat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '聊天室信息';

DROP TABLE IF EXISTS `cmdinboxmsgs`;
CREATE TABLE `cmdinboxmsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `send_time` bigint DEFAULT NULL COMMENT '发送时间',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_body` mediumblob COMMENT '消息体',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  `target_id` varchar(32) DEFAULT NULL COMMENT '接收者id',
  `msg_type` varchar(50) DEFAULT NULL COMMENT '消息类型',
  `uniq_tag` varchar(100) DEFAULT NULL COMMENT '唯一标签',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_tag` (`app_key`,`user_id`,`uniq_tag`),
  KEY `idx_appkey_time` (`app_key`,`user_id`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`user_id`,`msg_id`),
  KEY `idx_appkey` (`app_key`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'cmd收件箱';

DROP TABLE IF EXISTS `cmdsendboxmsgs`;
CREATE TABLE `cmdsendboxmsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `send_time` bigint DEFAULT NULL COMMENT '发送时间',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_body` mediumblob COMMENT '消息体',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  `target_id` varchar(32) DEFAULT NULL COMMENT '接收者id',
  `msg_type` varchar(50) DEFAULT NULL COMMENT '消息类型',
  `uniq_tag` varchar(100) DEFAULT NULL COMMENT '唯一标签',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_tag` (`app_key`,`user_id`,`uniq_tag`),
  KEY `idx_appkey_userid_time` (`app_key`,`user_id`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`user_id`,`msg_id`),
  KEY `idx_appkey` (`app_key`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'cmd发件箱';

DROP TABLE IF EXISTS `convercleantimes`;
CREATE TABLE `convercleantimes` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `conver_id` varchar(100) DEFAULT NULL COMMENT '会话id',
  `channel_type` tinyint DEFAULT '0' COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `clean_time` bigint DEFAULT '0' COMMENT '清除时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_destroy` (`app_key`,`conver_id`,`channel_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '会话清理记录';

DROP TABLE IF EXISTS `conversations`;
CREATE TABLE `conversations` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `target_id` varchar(32) DEFAULT NULL COMMENT '接收者id',
  `channel_type` tinyint DEFAULT '0' COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `latest_msg_id` varchar(20) DEFAULT NULL COMMENT '最新消息id',
  `latest_msg` mediumblob COMMENT '最新消息体',
  `latest_unread_msg_index` int DEFAULT '0' COMMENT '最新未读index',
  `latest_read_msg_index` int DEFAULT '0' COMMENT '最新已读消息index',
  `latest_read_msg_id` varchar(20) DEFAULT NULL COMMENT '最新已读消息id',
  `latest_read_msg_time` bigint DEFAULT '0' COMMENT '最新已读时间',
  `sort_time` bigint DEFAULT '0' COMMENT 'sort time',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除 0 未删除，1已删除',
  `is_top` tinyint DEFAULT '0' COMMENT '是否置顶 0未指定，1置顶',
  `top_updated_time` bigint DEFAULT '0' COMMENT '置顶更新时间',
  `undisturb_type` tinyint DEFAULT '0',
  `sync_time` bigint DEFAULT '0' COMMENT '同步消息位点',
  `unread_tag` tinyint DEFAULT '0' COMMENT '未读tag',
  `group` varchar(20) DEFAULT NULL,
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_app_key_user_id_target_id` (`app_key`,`user_id`,`target_id`,`channel_type`),
  KEY `idx_sync_time` (`app_key`,`user_id`,`sync_time`),
  KEY `idx_update_time` (`app_key`,`user_id`,`sort_time`),
  KEY `idx_undisturb` (`app_key`,`user_id`,`target_id`,`channel_type`),
  KEY `idx_group` (`app_key`,`user_id`,`group`,`sort_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '会话';

DROP TABLE IF EXISTS `userconvertags`;
CREATE TABLE `userconvertags` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` VARCHAR(32) NULL COMMENT '用户id',
  `tag` VARCHAR(50) NULL COMMENT 'tag',
  `tag_name` VARCHAR(50) NULL COMMENT '分组名称',
  `created_time` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_time` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `app_key` VARCHAR(20) NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uniq_tag` (`app_key`, `user_id`, `tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '会话-分组';

DROP TABLE IF EXISTS `convertagrels`;
CREATE TABLE `convertagrels` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` VARCHAR(32) NULL COMMENT '用户id',
  `tag` VARCHAR(50) NULL COMMENT 'tag',
  `target_id` VARCHAR(32) NULL COMMENT '目标id',
  `channel_type` TINYINT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `created_time` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `app_key` VARCHAR(20) NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uniq_tag_target` (`app_key`, `user_id`, `tag`, `target_id`, `channel_type`),
  KEY `idx_target` (`app_key`, `user_id`, `target_id`, `channel_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '会话-分组绑定关系表';


DROP TABLE IF EXISTS `msgstats`;
CREATE TABLE `msgstats` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `stat_type` TINYINT NULL DEFAULT 0 COMMENT '统计类型 1上行消息，2分发，3下行消息',
  `channel_type` TINYINT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `time_mark` BIGINT NULL COMMENT '标记时间',
  `count` INT NULL COMMENT '数量',
  `app_key` VARCHAR(20) NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uniq_mark` (`app_key` ASC, `stat_type` ASC, `channel_type` ASC, `time_mark` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '统计-消息统计';

DROP TABLE IF EXISTS `useractivities`;
CREATE TABLE `useractivities` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` VARCHAR(32) NULL COMMENT '用户id',
  `time_mark` BIGINT NULL COMMENT '标记点',
  `created_time` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `count` INT NULL COMMENT '数量',
  `app_key` VARCHAR(20) NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uniq_userid` (`app_key` ASC, `time_mark` ASC, `user_id` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '统计-用户活跃统计';

DROP TABLE IF EXISTS `fileconfs`;
CREATE TABLE `fileconfs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `app_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '应用key',
  `channel` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'oss供应商',
  `conf` json DEFAULT NULL COMMENT '配置',
  `enable` tinyint(1) DEFAULT '0' COMMENT '是否可用 0不可用，1可用',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_time` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `app_key` (`app_key`,`channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT = '文件配置';

DROP TABLE IF EXISTS `g_delhismsgs`;
CREATE TABLE `g_delhismsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `target_id` varchar(32) DEFAULT NULL COMMENT '目标id',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `msg_time` bigint DEFAULT NULL COMMENT '消息时间',
  `msg_seq` int DEFAULT NULL COMMENT '消息seq',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_msgid` (`app_key`,`user_id`,`target_id`,`msg_id`),
  KEY `idx_target` (`app_key`,`user_id`,`target_id`,`msg_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群聊-历史消息已删记录表';

DROP TABLE IF EXISTS `g_hismsgs`;
CREATE TABLE `g_hismsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `conver_id` varchar(100) DEFAULT NULL COMMENT '会话id',
  `sender_id` varchar(32) DEFAULT NULL COMMENT '发送者id',
  `receiver_id` varchar(32) DEFAULT NULL COMMENT '接收者id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_type` varchar(50) DEFAULT NULL COMMENT '消息类型',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `send_time` bigint DEFAULT NULL COMMENT '发送时间',
  `msg_seq_no` int DEFAULT NULL COMMENT '消息seq',
  `msg_body` mediumblob COMMENT '消息体',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  `member_count` int DEFAULT '0' COMMENT '成员数',
  `read_count` int DEFAULT '0' COMMENT '已读数',
  `is_delete` tinyint DEFAULT '0' COMMENT '是否删除 0未删除， 1已删除',
  `is_ext` tinyint DEFAULT '0',
  `is_exset` tinyint DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_appkey_converid` (`app_key`,`conver_id`,`msg_id`,`send_time`),
  KEY `idx_conver_time` (`app_key`,`conver_id`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群聊-历史消息表';

DROP TABLE IF EXISTS `gc_hismsgs`;
CREATE TABLE `gc_hismsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `conver_id` varchar(100) DEFAULT NULL COMMENT '会话id',
  `sender_id` varchar(32) DEFAULT NULL COMMENT '发送者id',
  `receiver_id` varchar(32) DEFAULT NULL COMMENT '接收者id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_type` varchar(50) DEFAULT NULL COMMENT '消息类型',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `send_time` bigint DEFAULT NULL COMMENT '发送时间',
  `msg_seq_no` int DEFAULT NULL COMMENT '消息seq',
  `msg_body` mediumblob COMMENT '消息体',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  KEY `idx_msg` (`app_key`,`conver_id`,`channel_type`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群聊-公告消息表';

DROP TABLE IF EXISTS `globalconfs`;
CREATE TABLE `globalconfs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `conf_key` varchar(50) DEFAULT NULL COMMENT '配置key',
  `conf_value` varchar(2000) DEFAULT NULL COMMENT '配置value',
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建者',
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_key` (`conf_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '全局配置';

DROP TABLE IF EXISTS `globalconvers`;
CREATE TABLE `globalconvers` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `conver_id` varchar(100) DEFAULT NULL COMMENT '会话id',
  `sender_id` varchar(32) DEFAULT NULL COMMENT '发送者',
  `target_id` varchar(32) DEFAULT NULL COMMENT '接收者',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `updated_time` bigint DEFAULT NULL COMMENT '更新时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_conver` (`app_key`,`conver_id`,`channel_type`),
  KEY `idx_time` (`app_key`,`channel_type`,`updated_time`),
  KEY `idx_targetid` (`app_key`,`channel_type`,`target_id`,`updated_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '全局会话表';

DROP TABLE IF EXISTS `groupinfoexts`;
CREATE TABLE `groupinfoexts` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `group_id` varchar(32) DEFAULT NULL COMMENT '群id',
  `item_key` varchar(50) DEFAULT NULL COMMENT '参数key',
  `item_value` varchar(100) DEFAULT NULL COMMENT '参数value',
  `item_type` tinyint DEFAULT '0' COMMENT '参数类型',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey_groupid` (`app_key`,`group_id`,`item_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群补充信息';

DROP TABLE IF EXISTS `groupinfos`;
CREATE TABLE `groupinfos` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `group_id` varchar(64) DEFAULT NULL COMMENT '群id',
  `group_name` varchar(64) DEFAULT NULL COMMENT '群名称',
  `group_portrait` varchar(200) DEFAULT NULL COMMENT '群头像',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  `is_mute` tinyint DEFAULT '0',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey_groupid` (`app_key`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群基本信息';

DROP TABLE IF EXISTS `groupmemberexts`;
CREATE TABLE `groupmemberexts` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `group_id` varchar(32) DEFAULT NULL COMMENT '群id',
  `member_id` varchar(32) DEFAULT NULL COMMENT '群成员id',
  `item_key` varchar(50) DEFAULT NULL COMMENT '参数key',
  `item_value` varchar(100) DEFAULT NULL COMMENT '参数value',
  `item_type` tinyint DEFAULT '0' COMMENT '参数类型',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_item_key` (`app_key`,`group_id`,`member_id`,`item_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群成员-补充信息';

DROP TABLE IF EXISTS `groupmembers`;
CREATE TABLE `groupmembers` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `group_id` varchar(64) DEFAULT NULL COMMENT '群id',
  `member_id` varchar(64) DEFAULT NULL COMMENT '成员id',
  `member_type` tinyint DEFAULT '0' COMMENT '成员类型',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `app_key` varchar(45) DEFAULT NULL COMMENT '应用key',
  `is_mute` tinyint DEFAULT '0' COMMENT '是否静音',
  `is_allow` tinyint DEFAULT '0' COMMENT '是否允许发言',
  `mute_end_at` bigint DEFAULT '0' COMMENT '禁言结束时间戳',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey_grpid_memid` (`app_key`,`group_id`,`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群成员';

DROP TABLE IF EXISTS `grpassistantrels`;
CREATE TABLE `grpassistantrels` (
  `id` int NOT NULL COMMENT '主键id',
  `assistant_id` varchar(32) DEFAULT NULL COMMENT '助手id',
  `target_id` varchar(32) DEFAULT NULL COMMENT '接收者id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_target` (`assistant_id`,`target_id`,`channel_type`,`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `grpsnapshots`;
CREATE TABLE `grpsnapshots` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `app_key` varchar(20) NOT NULL COMMENT '应用key',
  `group_id` varchar(32) NOT NULL COMMENT '群id',
  `created_time` bigint DEFAULT '0' COMMENT '创建时间',
  `snapshot` mediumblob COMMENT '详情',
  PRIMARY KEY (`id`),
  KEY `idx_group_id` (`app_key`,`group_id`,`created_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群快照';

DROP TABLE IF EXISTS `ic_conditions`;
CREATE TABLE `ic_conditions` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `channel_type` varchar(100) DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_type` varchar(1000) DEFAULT NULL COMMENT '消息类型',
  `sender_id` varchar(1000) DEFAULT NULL COMMENT '发送者id',
  `receiver_id` varchar(1000) DEFAULT NULL COMMENT '接收者id',
  `interceptor_id` int DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  KEY `idx_icid` (`app_key`,`interceptor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `inboxmsgs`;
CREATE TABLE `inboxmsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `send_time` bigint DEFAULT NULL COMMENT '发送时间',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_body` mediumblob COMMENT '消息体',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  `target_id` varchar(32) DEFAULT NULL COMMENT '接收者id',
  `msg_type` varchar(50) DEFAULT NULL COMMENT '消息类型',
  PRIMARY KEY (`id`),
  KEY `IDX_USERID_MSG` (`app_key`,`user_id`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`user_id`,`msg_id`),
  KEY `idx_appkey` (`app_key`,`send_time`)
) ENGINE=InnoDB AUTO_INCREMENT=45612 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '收件箱';

DROP TABLE IF EXISTS `interceptors`;
CREATE TABLE `interceptors` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `name` varchar(50) DEFAULT NULL COMMENT '名称',
  `sort` int NOT NULL DEFAULT '0',
  `request_url` varchar(500) DEFAULT NULL COMMENT '请求地址',
  `request_template` text COMMENT '请求模板',
  `succ_template` varchar(200) DEFAULT NULL COMMENT '成功模板',
  `is_async` tinyint DEFAULT '0' COMMENT '是否异步',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `app_key` varchar(20) NOT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  KEY `idx_sort` (`app_key`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `ioscertificates`;
CREATE TABLE `ioscertificates` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `package` varchar(100) DEFAULT NULL,
  `certificate` mediumblob,
  `cert_pwd` varchar(50) DEFAULT NULL COMMENT '认证密码',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `is_product` tinyint DEFAULT '0',
  `cert_path` varchar(255) DEFAULT NULL COMMENT 'cert存入路径',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_package` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `mentionmsgs`;
CREATE TABLE `mentionmsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `target_id` varchar(32) DEFAULT NULL COMMENT '接收者id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `sender_id` varchar(32) DEFAULT NULL COMMENT '发送者id',
  `mention_type` tinyint DEFAULT NULL COMMENT '类型',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `msg_time` bigint DEFAULT NULL COMMENT '消息时间',
  `msg_index` int DEFAULT NULL COMMENT '消息index',
  `is_read` tinyint DEFAULT NULL COMMENT '是否已读',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  KEY `idx_uid_tid_type` (`app_key`,`user_id`,`target_id`,`channel_type`,`msg_index`,`msg_time`),
  KEY `idx_read` (`app_key`,`user_id`,`target_id`,`channel_type`,`is_read`,`msg_time`),
  KEY `idx_user_msgid` (`app_key`,`user_id`,`target_id`,`channel_type`,`msg_id`),
  KEY `idx_target_msgid` (`app_key`, `msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '@消息记录';

DROP TABLE IF EXISTS `mergedmsgs`;
CREATE TABLE `mergedmsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `parent_msg_id` varchar(20) DEFAULT NULL COMMENT '父消息id',
  `from_id` varchar(32) DEFAULT NULL COMMENT '发送者',
  `target_id` varchar(32) DEFAULT NULL COMMENT '接收者',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `msg_time` bigint DEFAULT '0' COMMENT '消息时间',
  `msg_body` mediumblob COMMENT '消息体',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_appkey_pmsg` (`app_key`,`parent_msg_id`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '聚合消息';

DROP TABLE IF EXISTS `msgexts`;
CREATE TABLE `msgexts` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `key` varchar(50) DEFAULT NULL COMMENT '参数key',
  `value` varchar(1000) DEFAULT NULL COMMENT '参数value',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `app_key` varchar(45) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_msgid` (`app_key`,`msg_id`,`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '消息-补充内容';

DROP TABLE IF EXISTS `msgexsets`;
CREATE TABLE `msgexsets` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `msg_id` VARCHAR(20) NULL COMMENT '消息id',
  `key` VARCHAR(50) NULL COMMENT 'key',
  `item` VARCHAR(100) NULL,
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `created_time` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `app_key` VARCHAR(20) NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uniq_item` (`app_key`, `msg_id`, `key`, `item`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `p_delhismsgs`;
CREATE TABLE `p_delhismsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `target_id` varchar(32) DEFAULT NULL COMMENT '接收者',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `msg_time` bigint DEFAULT NULL COMMENT '消息时间',
  `msg_seq` int DEFAULT NULL COMMENT '消息seq',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_msgid` (`app_key`,`user_id`,`target_id`,`msg_id`),
  KEY `idx_target` (`app_key`,`user_id`,`target_id`,`msg_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '单聊-历史消息已删记录表';

DROP TABLE IF EXISTS `p_hismsgs`;
CREATE TABLE `p_hismsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `conver_id` varchar(100) DEFAULT NULL COMMENT '会话id',
  `sender_id` varchar(32) DEFAULT NULL COMMENT '发送者',
  `receiver_id` varchar(32) DEFAULT NULL COMMENT '接收者',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_type` varchar(50) DEFAULT NULL COMMENT '消息类型',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `send_time` bigint DEFAULT NULL COMMENT '发送时间',
  `msg_seq_no` int DEFAULT NULL COMMENT '消息seq',
  `msg_body` mediumblob COMMENT '消息体',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  `is_read` tinyint DEFAULT '0' COMMENT '是否已读 0未读，1已读',
  `is_delete` tinyint DEFAULT '0' COMMENT '是否删除 0未删除，1已删除',
  `is_ext` tinyint DEFAULT '0',
  `is_exset` tinyint DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_app_key_conver_id` (`app_key`,`conver_id`,`msg_id`,`send_time`),
  KEY `idx_conver_time` (`app_key`,`conver_id`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '单聊-历史消息';

DROP TABLE IF EXISTS `pushtokens`;
CREATE TABLE `pushtokens` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `device_id` varchar(200) DEFAULT NULL COMMENT '设备id',
  `platform` varchar(10) DEFAULT NULL COMMENT '平台',
  `push_channel` varchar(10) DEFAULT NULL COMMENT '推送渠道',
  `package` varchar(200) DEFAULT NULL,
  `push_token` varchar(200) DEFAULT NULL COMMENT '推送token',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`app_key`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '推送token';

DROP TABLE IF EXISTS `readinfos`;
CREATE TABLE `readinfos` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `app_key` varchar(20) NOT NULL COMMENT '应用key',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `group_id` varchar(32) DEFAULT NULL COMMENT '群id',
  `member_id` varchar(32) DEFAULT NULL COMMENT '群成员id',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_member` (`app_key`,`channel_type`,`group_id`,`msg_id`,`member_id`),
  KEY `idx_memberid` (`app_key`,`channel_type`,`group_id`,`member_id`,`msg_id`),
  KEY `idx_msgid` (`app_key`,`channel_type`,`group_id`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '群已读表';

DROP TABLE IF EXISTS `s_hismsgs`;
CREATE TABLE `s_hismsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `conver_id` varchar(100) DEFAULT NULL COMMENT '会话id',
  `sender_id` varchar(32) DEFAULT NULL COMMENT '发送者id',
  `receiver_id` varchar(32) DEFAULT NULL COMMENT '接收者id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_type` varchar(50) DEFAULT NULL COMMENT '消息类型',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `send_time` bigint DEFAULT NULL COMMENT '发送时间',
  `msg_seq_no` int DEFAULT NULL COMMENT '消息seq',
  `msg_body` mediumblob COMMENT '消息体',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  `is_read` tinyint DEFAULT NULL COMMENT '是否已读',
  PRIMARY KEY (`id`),
  KEY `idx_appkey_converid` (`app_key`,`conver_id`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `sendboxmsgs`;
CREATE TABLE `sendboxmsgs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `send_time` bigint DEFAULT NULL COMMENT '发送时间',
  `msg_id` varchar(20) DEFAULT NULL COMMENT '消息id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `msg_body` mediumblob COMMENT '消息体',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  `target_id` varchar(45) DEFAULT NULL COMMENT '接收者id',
  `msg_type` varchar(50) DEFAULT NULL COMMENT '消息类型',
  PRIMARY KEY (`id`),
  KEY `idx_user_id_send_time` (`app_key`,`user_id`,`send_time`),
  KEY `idx_msg_id` (`app_key`,`user_id`,`msg_id`),
  KEY `idx_appkey` (`app_key`,`send_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '发件箱';

DROP TABLE IF EXISTS `sensitivewords`;
CREATE TABLE `sensitivewords` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `app_key` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '应用key',
  `word` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '敏感词',
  `word_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '12' COMMENT '词类型',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_word` (`app_key`,`word`),
  KEY `idx_appkey` (`app_key`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT = '敏感词';

DROP TABLE IF EXISTS `subrelations`;
CREATE TABLE `subrelations` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `subscriber` varchar(32) DEFAULT NULL COMMENT '订阅者',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_sub` (`app_key`,`user_id`,`subscriber`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `usercleantimes`;
CREATE TABLE `usercleantimes` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `target_id` varchar(32) DEFAULT NULL COMMENT '接收者id',
  `channel_type` tinyint DEFAULT NULL COMMENT '会话类型 1单聊, 2群聊，3聊天室，4系统，5群公告，6广播',
  `clean_time` bigint DEFAULT NULL COMMENT '清理时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_app_key_user_id_target_id` (`app_key`,`user_id`,`target_id`,`channel_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `userexts`;
CREATE TABLE `userexts` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id',
  `item_key` varchar(50) DEFAULT NULL COMMENT '参数key',
  `item_value` varchar(100) DEFAULT NULL COMMENT '参数value',
  `item_type` tinyint DEFAULT '0' COMMENT '参数类型',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_item_key` (`app_key`,`user_id`,`item_key`),
  KEY `idx_item_key` (`app_key`,`item_key`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '用户-补充信息表';

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_type` tinyint DEFAULT '0' COMMENT '用户类型 0用户, 1机器人',
  `user_id` varchar(32) NOT NULL COMMENT '用户id',
  `nickname` varchar(45) DEFAULT NULL COMMENT '昵称',
  `user_portrait` varchar(200) DEFAULT NULL COMMENT '用户头像',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `app_key` varchar(45) DEFAULT NULL COMMENT '应用key',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_userid` (`app_key`,`user_id`),
  KEY `idx_userid` (`app_key`,`user_type`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '用户-基本信息表';

DROP TABLE IF EXISTS `clientlogs`;
CREATE TABLE `clientlogs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `app_key` VARCHAR(20) NULL COMMENT '应用key',
  `user_id` VARCHAR(32) NULL COMMENT '用户id',
  `created_time` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `start` BIGINT NULL COMMENT '开始时间',
  `end` BIGINT NULL COMMENT '结束时间',
  `log` MEDIUMBLOB NULL COMMENT '详情',
  `state` TINYINT NULL DEFAULT 0 COMMENT '状态',
  `platform` VARCHAR(20) NULL COMMENT '客户端类型 ios, android, web',
  `device_id` VARCHAR(100) NULL COMMENT '设备id',
  `log_url` VARCHAR(200) NULL COMMENT '日志链接',
  `trace_id` VARCHAR(50) NULL COMMENT '跟踪id',
  `msg_id` VARCHAR(20) NULL COMMENT '消息id',
  `fail_reason` VARCHAR(100) NULL COMMENT '失败原因',
  `description` VARCHAR(100) NULL COMMENT '描述',
  PRIMARY KEY (`id`),
  INDEX `idx_userid` (`app_key` ASC, `user_id` ASC),
  UNIQUE KEY `uniq_msgid` (`app_key`,`msg_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '客户端日志表';

DROP TABLE IF EXISTS `rtcrooms`;
CREATE TABLE `rtcrooms` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `room_id` varchar(50) DEFAULT NULL COMMENT '房间id',
  `room_type` tinyint DEFAULT '0' COMMENT '房间类型',
  `owner_id` varchar(32) DEFAULT NULL COMMENT '创建者',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_roomid` (`app_key`,`room_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'rtc-房间表';

DROP TABLE IF EXISTS `rtcmembers`;
CREATE TABLE `rtcmembers` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `room_id` varchar(50) DEFAULT NULL COMMENT '房间id',
  `member_id` varchar(32) DEFAULT NULL COMMENT '成员id',
  `device_id` varchar(50) DEFAULT NULL COMMENT '设备id',
  `rtc_state` tinyint DEFAULT '0'  COMMENT '状态',
  `inviter_id` varchar(32) DEFAULT NULL COMMENT '邀请人id',
  `latest_ping_time` bigint DEFAULT '0' COMMENT '最新ping时间',
  `call_time` bigint DEFAULT '0' COMMENT '拨号时间',
  `connect_time` bigint DEFAULT '0' COMMENT '通话开始时间',
  `hangup_time` bigint DEFAULT '0' COMMENT '通话挂断时间',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `app_key` varchar(20) DEFAULT NULL COMMENT '应用key',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_member` (`app_key`,`room_id`,`member_id`),
  KEY `idx_room` (`app_key`,`member_id`,`room_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'rtc-成员表';

INSERT IGNORE INTO `accounts`(`account`,`password`)VALUES('admin1','7c4a8d09ca3762af61e59520943dc26494f8941b');