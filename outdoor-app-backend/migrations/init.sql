-- 插入测试用户
INSERT IGNORE INTO `users` (`id`, `email`, `password_hash`, `nickname`, `fitness_level`, `created_at`, `updated_at`) VALUES
(1, 'admin@outdoor.com', 'hashed_pwd_123', '超级领队', 5, NOW(), NOW()),
(2, 'test@outdoor.com', 'hashed_pwd_456', '户外萌新', 1, NOW(), NOW());

-- 插入一条测试的急救百科
INSERT IGNORE INTO `articles` (`id`, `title`, `category`, `content`, `created_at`, `updated_at`) VALUES
(1, '户外突发失温的黄金抢救时间', '急救', '失温（Hypothermia）是指人体核心温度降至35℃以下...', NOW(), NOW());