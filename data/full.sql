-- phpMyAdmin SQL Dump
-- version 4.9.0.1
-- https://www.phpmyadmin.net/
--
-- 主机： mysql:3306
-- 生成日期： 2019-09-25 08:39:45
-- 服务器版本： 8.0.13
-- PHP 版本： 7.2.22

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 数据库： `book_crawl`
--
CREATE DATABASE IF NOT EXISTS `book_crawl` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `book_crawl`;

-- --------------------------------------------------------

--
-- 表的结构 `book_info`
--

CREATE TABLE `book_info` (
  `book_id` bigint(20) NOT NULL COMMENT '小说ID',
  `book_name` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '小说名',
  `author` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '作者',
  `tag` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '标签',
  `status` tinyint(4) NOT NULL COMMENT '小说状态：0:完结,1:连载',
  `score` decimal(3,1) NOT NULL COMMENT '评分',
  `score_count` int(11) NOT NULL COMMENT '评分人数',
  `score_detail` json NOT NULL COMMENT '评分详情',
  `add_list_count` int(11) NOT NULL COMMENT '添加书单数',
  `last_update_time` datetime NOT NULL COMMENT '最后更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='小说详情表';


--
-- 表的结构 `book_score_daily_201909`
--

CREATE TABLE `book_score_daily_201909` (
  `log_id` bigint(20) NOT NULL COMMENT '日志id',
  `book_id` bigint(20) NOT NULL COMMENT '小说id',
  `score` decimal(3,1) NOT NULL COMMENT '评分',
  `score_count` int(11) NOT NULL COMMENT '评分人数',
  `score_detail` json NOT NULL COMMENT '评分详情',
  `date_key` date NOT NULL COMMENT '日志日期'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='小说每日积分记录';

--
-- 转储表的索引
--

--
-- 表的索引 `book_info`
--
ALTER TABLE `book_info`
  ADD PRIMARY KEY (`book_id`),
  ADD UNIQUE KEY `index_unique` (`book_name`,`author`);

--
-- 表的索引 `book_score_daily_201909`
--
ALTER TABLE `book_score_daily_201909`
  ADD PRIMARY KEY (`log_id`),
  ADD UNIQUE KEY `index_unique` (`book_id`,`date_key`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `book_info`
--
ALTER TABLE `book_info`
  MODIFY `book_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '小说ID';

--
-- 使用表AUTO_INCREMENT `book_score_daily_201909`
--
ALTER TABLE `book_score_daily_201909`
  MODIFY `log_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志id';
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
