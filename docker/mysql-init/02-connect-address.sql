-- Seed the WebSocket long-connection address delivered to client SDKs by the
-- navigator service. For local development the host reaches the container on
-- 127.0.0.1:9003. In production, replace this with your public IP / domain, e.g.
--   {"default":["im.example.com:9003"]}
INSERT INTO `globalconfs` (`conf_key`, `conf_value`)
VALUES ('connect_address', '{"default":["127.0.0.1:9003"]}');
