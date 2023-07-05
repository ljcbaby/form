
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
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for results
-- ----------------------------
DROP TABLE IF EXISTS `results`;
CREATE TABLE `results`  (
  `id` int(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `form_id` int(0) UNSIGNED NOT NULL,
  `result` json NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `formid`(`form_id`) USING BTREE,
  CONSTRAINT `results_ibfk_1` FOREIGN KEY (`form_id`) REFERENCES `forms` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

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
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- View structure for form_bases
-- ----------------------------
DROP VIEW IF EXISTS `form_bases`;
CREATE ALGORITHM = UNDEFINED SQL SECURITY DEFINER VIEW `form_bases` AS select `f`.`id` AS `id`,`f`.`title` AS `title`,(`f`.`status` - 1) AS `isPublish`,(select count(0) from `results` `r` where (`r`.`form_id` = `f`.`id`)) AS `answerCount`,`f`.`modifiedAt` AS `modifiedAt`,`f`.`owner_id` AS `owner_id` from `forms` `f` where (`f`.`status` <> 3);

-- ----------------------------
-- Procedure structure for GetComponents
-- ----------------------------
DROP PROCEDURE IF EXISTS `GetComponents`;
delimiter ;;
CREATE PROCEDURE `GetComponents`(IN p_id INT)
BEGIN
  -- 创建临时表
  DROP TEMPORARY TABLE IF EXISTS temp_table;
  CREATE TEMPORARY TABLE temp_table LIKE forms;

  -- 复制指定ID行到临时表
  INSERT INTO temp_table
  SELECT * FROM forms WHERE id = p_id;

  -- 使用临时表进行进一步操作
  SELECT
    TRIM(BOTH '"' FROM JSON_EXTRACT(jt.component, '$.fe_id')) AS fe_id,
    TRIM(BOTH '"' FROM JSON_EXTRACT(jt.component, '$.props.title')) AS title
  FROM
    temp_table
    JOIN JSON_TABLE(temp_table.components, '$[*]' COLUMNS (component JSON PATH '$')) AS jt
  WHERE
    JSON_EXTRACT(jt.component, '$.type') IN ('questionInput', 'questionTextarea', 'questionRadio', 'questionCheckbox');

  -- 删除临时表
  DROP TEMPORARY TABLE IF EXISTS temp_table;
END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for GetComponentType
-- ----------------------------
DROP PROCEDURE IF EXISTS `GetComponentType`;
delimiter ;;
CREATE PROCEDURE `GetComponentType`(IN formId INT, IN feId VARCHAR(255))
BEGIN
  -- 创建临时表
  DROP TEMPORARY TABLE IF EXISTS temp_table;
  CREATE TEMPORARY TABLE temp_table LIKE forms;

  -- 复制指定ID行到临时表
  INSERT INTO temp_table
  SELECT * FROM forms WHERE id = formId;

  -- 使用临时表进行进一步操作
  SELECT
     (SELECT CAST(jt.value AS CHAR CHARACTER SET utf8mb4) FROM JSON_TABLE(components, '$[*]' COLUMNS (fe_id VARCHAR(255) PATH '$.fe_id', value VARCHAR(65535) PATH '$.type')) AS jt WHERE jt.fe_id COLLATE utf8mb4_general_ci = feId COLLATE utf8mb4_general_ci) AS value
  FROM
    temp_table;

  -- 删除临时表
  DROP TEMPORARY TABLE IF EXISTS temp_table;
END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for GetResults
-- ----------------------------
DROP PROCEDURE IF EXISTS `GetResults`;
delimiter ;;
CREATE PROCEDURE `GetResults`(IN fid INT, IN feid VARCHAR(255), IN L INT, IN O INT)
BEGIN
  -- 创建临时表
  DROP TEMPORARY TABLE IF EXISTS temp_table;
  CREATE TEMPORARY TABLE temp_table LIKE results;

  -- 复制指定ID行到临时表
  INSERT INTO temp_table
  SELECT * FROM results WHERE form_id = fid;

  -- 使用临时表进行进一步操作
  SELECT
    id,
     (SELECT CAST(jt.value AS CHAR CHARACTER SET utf8mb4) FROM JSON_TABLE(result, '$[*]' COLUMNS (fe_id VARCHAR(255) PATH '$.fe_id', value VARCHAR(65535) PATH '$.value')) AS jt WHERE jt.fe_id COLLATE utf8mb4_general_ci = feid COLLATE utf8mb4_general_ci) AS value
  FROM
    temp_table
  LIMIT L OFFSET O;

  -- 删除临时表
  DROP TEMPORARY TABLE IF EXISTS temp_table;
END
;;
delimiter ;

-- ----------------------------
-- Triggers structure for table forms
-- ----------------------------
DROP TRIGGER IF EXISTS `prevent_status_update`;
delimiter ;;
CREATE TRIGGER `prevent_status_update` BEFORE UPDATE ON `forms` FOR EACH ROW BEGIN
  IF NEW.status < OLD.status OR NEW.status > 3 THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid status update';
  END IF;
END
;;
delimiter ;

SET FOREIGN_KEY_CHECKS = 1;
