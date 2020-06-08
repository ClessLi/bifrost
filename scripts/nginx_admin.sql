SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
-- ----------------------------
-- Database structure for ng_conf_admin
-- ----------------------------
CREATE DATABASE `ng_conf_admin` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

-- ----------------------------
-- Table structure for ng_conf_admin.user
-- ----------------------------
DROP TABLE IF EXISTS `ng_conf_admin`.`user`;
CREATE TABLE `ng_conf_admin`.`user` (
    `user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    PRIMARY KEY (`user_name`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FOMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `ng_conf_admin`.`user` VALUES ('ngadmin', '696e76697369626c655f636c6f616bdb06c91d894ad7bfe7675deb6a880dbfbc62f6c7');

SET FOREIGN_KEY_CHECKS = 1;