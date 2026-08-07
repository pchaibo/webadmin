/*
SQLyog Ultimate v11.24 (32 bit)
MySQL - 5.7.40 : Database - webadmin
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`webadmin` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `webadmin`;

/*Table structure for table `tu_admin` */

DROP TABLE IF EXISTS `tu_admin`;

CREATE TABLE `tu_admin` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '用户头像，相对于upload/avatar目录',
  `email` varchar(255) DEFAULT NULL,
  `email_code` varchar(60) DEFAULT NULL COMMENT '激活码',
  `phone` bigint(20) unsigned DEFAULT NULL COMMENT '手机号',
  `status` bigint(20) DEFAULT NULL,
  `register_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '注册时间',
  `last_login_ip` varchar(16) NOT NULL DEFAULT '' COMMENT '最后登录ip',
  `last_login_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最后登录时间',
  `remark` varchar(100) DEFAULT NULL,
  `group_id` int(10) DEFAULT '1' COMMENT '组id',
  `auth_group` int(10) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `user_login_key` (`username`)
) ENGINE=MyISAM AUTO_INCREMENT=24 DEFAULT CHARSET=utf8;

/*Data for the table `tu_admin` */

insert  into `tu_admin`(`id`,`username`,`password`,`avatar`,`email`,`email_code`,`phone`,`status`,`register_time`,`last_login_ip`,`last_login_time`,`remark`,`group_id`,`auth_group`) values (1,'admin','a9bf03cd0485dfde54b0225bd4144e41','/uploads/file/20241114/6cac6c24901ddaf5e7e80090f3d120f9.png','admin++@qq.com',NULL,13912345600,1,1716813295,'',0,'',1,1),(20,'test','cc03e747a6afbbcbf8be7668acfebee5','','cqhaibo@gmail.com',NULL,NULL,0,0,'',0,NULL,1,2);

/*Table structure for table `tu_auth_group` */

DROP TABLE IF EXISTS `tu_auth_group`;

CREATE TABLE `tu_auth_group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(200) DEFAULT NULL,
  `status` bigint(20) DEFAULT NULL,
  `rules` text,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COMMENT='用户组表';

/*Data for the table `tu_auth_group` */

insert  into `tu_auth_group`(`id`,`title`,`status`,`rules`) values (1,'超级管理员',1,'1,2,10,5,3,7,8,6,4,11'),(2,'一般管理员',1,'7,8,6,2,1');

/*Table structure for table `tu_auth_group_access` */

DROP TABLE IF EXISTS `tu_auth_group_access`;

CREATE TABLE `tu_auth_group_access` (
  `admin_id` int(10) unsigned NOT NULL COMMENT '用户id',
  `group_id` int(10) unsigned NOT NULL COMMENT '用户组id',
  UNIQUE KEY `uid_group_id` (`admin_id`,`group_id`),
  KEY `uid` (`admin_id`),
  KEY `group_id` (`group_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='用户组明细表';

/*Data for the table `tu_auth_group_access` */

insert  into `tu_auth_group_access`(`admin_id`,`group_id`) values (1,1),(20,2);

/*Table structure for table `tu_auth_rule` */

DROP TABLE IF EXISTS `tu_auth_rule`;

CREATE TABLE `tu_auth_rule` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `pid` bigint(20) DEFAULT NULL,
  `name` longtext,
  `title` longtext,
  `status` bigint(20) DEFAULT NULL,
  `is_menu` tinyint(3) unsigned DEFAULT '0' COMMENT '菜单显示',
  `condition` longtext,
  `type` bigint(20) DEFAULT NULL,
  `sort` int(11) DEFAULT '999',
  `icon` longtext,
  `z_index` int(11) NOT NULL DEFAULT '0' COMMENT '菜单位置 0 左侧  1 顶部(只有一级菜单 才能开启顶部展示)',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=12 DEFAULT CHARSET=utf8 COMMENT='规则表';

/*Data for the table `tu_auth_rule` */

insert  into `tu_auth_rule`(`id`,`pid`,`name`,`title`,`status`,`is_menu`,`condition`,`type`,`sort`,`icon`,`z_index`) values (1,1,'index','首页',1,1,'',1,1,'',0),(2,1,'user','用户管理',1,1,'',1,999,'',0),(3,1,'admin','后台管理员',1,1,'',1,3,'',0),(4,0,'admingroup','管理组',1,1,'',1,1,'layui-icon-home',0),(5,0,'authrule','权限规则',1,1,'',1,999,'layui-icon-set',0),(6,0,'coin','coin价格',1,0,'',1,2,'layui-icon-picture',0),(7,0,'heyue','合约列表',1,0,'',1,99,'',0),(8,0,'heyueorder','合约日志',1,0,'',1,999,'',0),(10,0,'shellgroup','分组管理',1,1,'',1,999,'',0),(11,0,'shell','shell管理',1,1,'',1,999,'',0);

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
