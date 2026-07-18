ALTER TABLE `ioscertificates`
  DROP INDEX `uniq_package`,
  ADD UNIQUE KEY `uniq_package` (`app_key`,`package`);

ALTER TABLE `androidpushconfs`
  DROP INDEX `uniq_channel`,
  ADD UNIQUE KEY `uniq_channel` (`app_key`,`push_channel`,`package`);