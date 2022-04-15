SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
-- ----------------------------
-- Database structure for bifrost
-- ----------------------------
CREATE DATABASE `bifrost` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

-- ----------------------------
-- Table structure for bifrost.user
-- ----------------------------
DROP TABLE IF EXISTS `bifrost`.`user`;
CREATE TABLE `bifrost`.`user`
(
    `user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `password`  varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    PRIMARY KEY (`user_name`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `bifrost`.`user`
VALUES ('heimdall', '696e76697369626c655f636c6f616b876f5380561c2cc8af7a77f7cc0d72c98bafa4c6');

SET FOREIGN_KEY_CHECKS = 1;