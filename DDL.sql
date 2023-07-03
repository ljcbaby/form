
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for forms
-- ----------------------------
DROP TABLE IF EXISTS `forms`;
CREATE TABLE `forms`  (
  `id` int(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `owner_id` int(0) UNSIGNED NOT NULL,
  `status` int(0) UNSIGNED NOT NULL DEFAULT 1,
  `modifiedAt` datetime(0) NOT NULL,
  `components` json NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `owner_id`(`owner_id`) USING BTREE,
  CONSTRAINT `forms_ibfk_1` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 28 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for results
-- ----------------------------
DROP TABLE IF EXISTS `results`;
CREATE TABLE `results`  (
  `id` int(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `form_id` int(0) UNSIGNED NOT NULL,
  `identify` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `result` json NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `formid`(`form_id`) USING BTREE,
  CONSTRAINT `results_ibfk_1` FOREIGN KEY (`form_id`) REFERENCES `forms` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `password` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `salt` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `nickname` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `username`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- View structure for form_bases
-- ----------------------------
DROP VIEW IF EXISTS `form_bases`;
CREATE ALGORITHM = UNDEFINED SQL SECURITY DEFINER VIEW `form_bases` AS select `f`.`id` AS `id`,`f`.`title` AS `title`,(`f`.`status` - 1) AS `isPublish`,(select count(0) from `results` `r` where (`r`.`form_id` = `f`.`id`)) AS `answerCount`,`f`.`modifiedAt` AS `modifiedAt`,`f`.`owner_id` AS `owner_id` from `forms` `f` where (`f`.`status` <> 3);

SET FOREIGN_KEY_CHECKS = 1;
