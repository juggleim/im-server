ALTER TABLE `apps`
    ADD COLUMN `lic_conf` varchar(2048) DEFAULT NULL COMMENT '应用授权配置' AFTER `app_name`;

ALTER TABLE `appnavs`
    ADD COLUMN `admin_url` varchar(200) DEFAULT NULL AFTER `alias_no`;

INSERT INTO `globalconfs` (`conf_key`,`conf_value`)VALUES('jimdb_version','20260720') ON DUPLICATE KEY UPDATE conf_value=VALUES(conf_value);
