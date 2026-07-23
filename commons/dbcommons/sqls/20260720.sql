ALTER TABLE `ioscertificates`
  DROP INDEX `uniq_package`,
  ADD UNIQUE KEY `uniq_package` (`app_key`,`package`);

ALTER TABLE `androidpushconfs`
  DROP INDEX `uniq_channel`,
  ADD UNIQUE KEY `uniq_channel` (`app_key`,`push_channel`,`package`);

ALTER TABLE `androidpushconfs`
CHANGE COLUMN `push_conf` `push_conf` VARCHAR(5000) NULL DEFAULT NULL ;

ALTER TABLE `pushtokens`
  ADD KEY `idx_pushtoken` (`app_key`,`push_token`);
