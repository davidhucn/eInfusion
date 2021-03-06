/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80011
 Source Host           : localhost:3306
 Source Schema         : transfusion

 Target Server Type    : MySQL
 Target Server Version : 80011
 File Encoding         : 65001

 Date: 01/01/2001 00:11:52
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for t_apply_dict
-- ----------------------------
DROP TABLE IF EXISTS `t_apply_dict`;
CREATE TABLE `t_apply_dict`  (
  `qcode` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '应用字典',
  `hospital` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `dept_name` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `ward_name` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `province` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `city` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  PRIMARY KEY (`qcode`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for t_device_dict
-- ----------------------------
DROP TABLE IF EXISTS `t_device_dict`;
CREATE TABLE `t_device_dict`  (
  `qcode` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '配对、二维码与设备对照用\r\n',
  `detector_id` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `remark` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `disable` tinyint(4) NULL DEFAULT NULL COMMENT '是否禁用',
  PRIMARY KEY (`detector_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of t_device_dict
-- ----------------------------
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000000', 'B0000000', NULL, 0);
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000001', 'B0000001', NULL, 0);
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000002', 'B0000002', NULL, 0);
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000003', 'B0000003', NULL, 0);
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000004', 'B0000004', NULL, 0);
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000005', 'B0000005', NULL, 0);
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000006', 'B0000006', NULL, 0);
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000007', 'B0000007', NULL, 0);

-- ----------------------------
-- Table structure for t_main
-- ----------------------------
DROP TABLE IF EXISTS `t_main`;
CREATE TABLE `t_main`  (
  `did` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `time` datetime(0) NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(0),
  `capacity` smallint(6) NULL DEFAULT NULL,
  `error` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `remark` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `alarm` tinyint(255) NULL DEFAULT NULL COMMENT '输液结束报警',
  PRIMARY KEY (`did`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci COMMENT = 'Device apply information\r\n' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for t_match_dict
-- ----------------------------
DROP TABLE IF EXISTS `t_match_dict`;
CREATE TABLE `t_match_dict`  (
  `detector_id` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `receiver_id` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  `remark` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  PRIMARY KEY (`detector_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for t_receiver_dict
-- ----------------------------
DROP TABLE IF EXISTS `t_receiver_dict`;
CREATE TABLE `t_receiver_dict`  (
  `receiver_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `detector_amount` int(11) NULL DEFAULT NULL,
  `reconn_time` int(11) NULL DEFAULT NULL,
  `remark` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `last_time` datetime(0) NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(0),
  `ip_addr` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `target_ip` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '接收器指向的局域网目标地址',
  `target_port` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `server_ip` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '接收器指向的远端服务器地址',
  `server_port` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`receiver_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of t_receiver_dict
-- ----------------------------
INSERT INTO `t_receiver_dict` VALUES ('A0000000', 1, NULL, NULL, '2018-07-20 01:23:22', '192.168.20.146', '192.168.0.107', '[30]', NULL, NULL);

SET FOREIGN_KEY_CHECKS = 1;
