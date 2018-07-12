/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50722
Source Host           : 127.0.0.1:3306
Source Database       : transfusion

Target Server Type    : MYSQL
Target Server Version : 50722
File Encoding         : 65001

Date: 2018-07-12 10:31:35
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for t_apply_dict
-- ----------------------------
DROP TABLE IF EXISTS `t_apply_dict`;
CREATE TABLE `t_apply_dict` (
  `qcode` varchar(255) CHARACTER SET utf8 NOT NULL COMMENT '应用字典',
  `hospital` varchar(255) DEFAULT NULL,
  `dept_name` varchar(255) DEFAULT NULL,
  `ward_name` varchar(255) DEFAULT NULL,
  `province` varchar(255) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`qcode`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
-- Records of t_apply_dict
-- ----------------------------

-- ----------------------------
-- Table structure for t_device_dict
-- ----------------------------
DROP TABLE IF EXISTS `t_device_dict`;
CREATE TABLE `t_device_dict` (
  `qcode` varchar(255) CHARACTER SET utf8 NOT NULL COMMENT '配对、二维码与设备对照用\r\n',
  `detector_id` varchar(255) NOT NULL,
  `remark` varchar(255) DEFAULT NULL,
  `disable` tinyint(4) DEFAULT NULL COMMENT '是否禁用',
  PRIMARY KEY (`detector_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
-- Records of t_device_dict
-- ----------------------------
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000000', 'B0000000', null, '0');
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000001', 'B0000001', null, '0');
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000002', 'B0000002', null, '0');
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000003', 'B0000003', null, '0');
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000004', 'B0000004', null, '0');
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000005', 'B0000005', null, '0');
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000006', 'B0000006', null, '0');
INSERT INTO `t_device_dict` VALUES ('1x0CPxx1B0000007', 'B0000007', null, '0');

-- ----------------------------
-- Table structure for t_main
-- ----------------------------
DROP TABLE IF EXISTS `t_main`;
CREATE TABLE `t_main` (
  `qcode` varchar(255) CHARACTER SET utf8 NOT NULL,
  `time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `capacity` smallint(6) DEFAULT NULL,
  `error` varchar(255) DEFAULT NULL,
  `remark` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`qcode`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='Device apply information\r\n';

-- ----------------------------
-- Records of t_main
-- ----------------------------

-- ----------------------------
-- Table structure for t_match_dict
-- ----------------------------
DROP TABLE IF EXISTS `t_match_dict`;
CREATE TABLE `t_match_dict` (
  `detector_id` varchar(255) NOT NULL,
  `receiver_id` varchar(255) DEFAULT NULL,
  `remark` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`detector_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
-- Records of t_match_dict
-- ----------------------------

-- ----------------------------
-- Table structure for t_receiver_dict
-- ----------------------------
DROP TABLE IF EXISTS `t_receiver_dict`;
CREATE TABLE `t_receiver_dict` (
  `receiver_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `detector_amount` int(11) DEFAULT NULL,
  `reconn_time` int(11) DEFAULT NULL,
  `remark` varchar(255) DEFAULT NULL,
  `last_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `ip_addr` varchar(255) DEFAULT NULL,
  `target_ip` varchar(255) DEFAULT NULL COMMENT '接收器指向的局域网目标地址',
  `target_port` varchar(255) DEFAULT NULL,
  `server_ip` varchar(255) DEFAULT NULL COMMENT '接收器指向的远端服务器地址',
  `server_port` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`receiver_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_receiver_dict
-- ----------------------------
INSERT INTO `t_receiver_dict` VALUES ('A0000000', '1', null, null, '2018-07-03 12:00:40', '192.168.137.1', '192.168.0.107', '[30]', null, null);
