SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL,
  `nickname` varchar(255) NOT NULL DEFAULT '' COMMENT '默认为空，如果不为空，则优先用这个 nickname',
  `auth_types` bigint(20) NOT NULL DEFAULT '0' COMMENT 'bitmap',
  `password` varbinary(255) NOT NULL DEFAULT '',
  `salt` varbinary(255) NOT NULL DEFAULT '',
  `create_time` bigint(20) NOT NULL,
  `has_fixed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否已经固定，不固定的帐号视为临时帐号，不能绑定其他帐号',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for user_email
-- ----------------------------
DROP TABLE IF EXISTS `user_email`;
CREATE TABLE `user_email` (
  `user_id` bigint(20) NOT NULL,
  `nickname` varchar(255) NOT NULL DEFAULT '' COMMENT '默认为 email',
  `email` varchar(128) NOT NULL,
  `has_fixed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '对于新帐户, 为0, 只有当用户确认绑定到某个帐户或者确认新建帐户后才为1',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `email` (`email`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for user_phone
-- ----------------------------
DROP TABLE IF EXISTS `user_phone`;
CREATE TABLE `user_phone` (
  `user_id` bigint(20) NOT NULL,
  `nickname` varchar(255) NOT NULL DEFAULT '' COMMENT '默认为 phone',
  `phone` varchar(20) NOT NULL,
  `has_fixed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '对于新帐户, 为0, 只有当用户确认绑定到某个帐户或者确认新建帐户后才为1',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `phone` (`phone`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for user_qq
-- ----------------------------
DROP TABLE IF EXISTS `user_qq`;
CREATE TABLE `user_qq` (
  `user_id` bigint(20) NOT NULL,
  `nickname` varchar(255) NOT NULL DEFAULT '' COMMENT '默认为 qq nickname',
  `openid` varchar(128) NOT NULL,
  `has_fixed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '对于新帐户, 为0, 只有当用户确认绑定到某个帐户或者确认新建帐户后才为1',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `openid` (`openid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for user_wechat
-- ----------------------------
DROP TABLE IF EXISTS `user_wechat`;
CREATE TABLE `user_wechat` (
  `user_id` bigint(20) NOT NULL,
  `nickname` varchar(255) NOT NULL DEFAULT '' COMMENT '默认为 wechat nickname',
  `openid` varchar(128) NOT NULL,
  `has_fixed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '对于新帐户, 为0, 只有当用户确认绑定到某个帐户或者确认新建帐户后才为1',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `openid` (`openid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for user_weibo
-- ----------------------------
DROP TABLE IF EXISTS `user_weibo`;
CREATE TABLE `user_weibo` (
  `user_id` bigint(20) NOT NULL,
  `nickname` varchar(255) NOT NULL DEFAULT '' COMMENT '默认为 weibo nickname',
  `openid` varchar(128) NOT NULL,
  `has_fixed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '对于新帐户, 为0, 只有当用户确认绑定到某个帐户或者确认新建帐户后才为1',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `openid` (`openid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;