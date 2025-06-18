/*M!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19-11.6.2-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: gin-gorm
-- ------------------------------------------------------
-- Server version	11.6.2-MariaDB-ubu2404

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*M!100616 SET @OLD_NOTE_VERBOSITY=@@NOTE_VERBOSITY, NOTE_VERBOSITY=0 */;

--
-- Table structure for table `category_basic`
--

DROP TABLE IF EXISTS `category_basic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `category_basic` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL COMMENT '分类名称',
  `parent_id` int(11) DEFAULT 0 COMMENT '父级ID',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `category_basic`
--

LOCK TABLES `category_basic` WRITE;
/*!40000 ALTER TABLE `category_basic` DISABLE KEYS */;
INSERT INTO `category_basic` VALUES
(1,'Category_1','数组',0,'2025-05-18 16:15:24','2025-05-18 16:15:26',NULL),
(2,'Category2','字符串',0,'2025-05-18 16:15:52','2025-05-18 16:15:53',NULL),
(3,'Category3','图',0,'2025-05-18 16:16:23','2025-06-05 22:17:49',NULL),
(10,'1c79530d-7c1c-444d-bd44-1074289a997f','树',0,'2025-06-05 22:16:20','2025-06-05 22:16:20',NULL);
/*!40000 ALTER TABLE `category_basic` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `contest_basic`
--

DROP TABLE IF EXISTS `contest_basic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contest_basic` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) DEFAULT NULL COMMENT '唯一标识',
  `name` varchar(100) DEFAULT NULL COMMENT '名称',
  `content` text DEFAULT NULL COMMENT '竞赛描述',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `start_at` datetime DEFAULT NULL,
  `end_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `contest_basic`
--

LOCK TABLES `contest_basic` WRITE;
/*!40000 ALTER TABLE `contest_basic` DISABLE KEYS */;
/*!40000 ALTER TABLE `contest_basic` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `contest_problem`
--

DROP TABLE IF EXISTS `contest_problem`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contest_problem` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `contest_id` int(11) DEFAULT NULL COMMENT '竞赛id',
  `problem_id` int(11) DEFAULT NULL COMMENT '问题id',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `contest_problem`
--

LOCK TABLES `contest_problem` WRITE;
/*!40000 ALTER TABLE `contest_problem` DISABLE KEYS */;
/*!40000 ALTER TABLE `contest_problem` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `contest_user`
--

DROP TABLE IF EXISTS `contest_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contest_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `contest_id` int(11) DEFAULT NULL COMMENT '竞赛id',
  `user_identity` varchar(36) DEFAULT NULL COMMENT '用户identity',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `contest_user`
--

LOCK TABLES `contest_user` WRITE;
/*!40000 ALTER TABLE `contest_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `contest_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `problem_basic`
--

DROP TABLE IF EXISTS `problem_basic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `problem_basic` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL COMMENT '问题的题目',
  `content` text DEFAULT NULL COMMENT '问题的正文描述',
  `max_runtime` int(11) DEFAULT NULL COMMENT '最大的运行时间',
  `max_mem` int(11) DEFAULT NULL COMMENT '最大的运行内存',
  `pass_num` int(11) DEFAULT 0 COMMENT '通过的问题个数',
  `submit_num` int(11) DEFAULT 0 COMMENT '提交次数',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `problem_basic`
--

LOCK TABLES `problem_basic` WRITE;
/*!40000 ALTER TABLE `problem_basic` DISABLE KEYS */;
INSERT INTO `problem_basic` VALUES
(1,'problem_1','文章标题1','文章正文1',1000,1000,1,1,'2025-05-18 15:36:20','2025-05-18 15:36:16',NULL),
(2,'problem_2','文章标题2','文章正文2',2000,2000,2,2,'2025-05-18 15:36:20','2025-05-18 15:36:16',NULL),
(3,'problem_3','文章标题3','文章正文3',3000,3000,4,4,'2025-05-18 15:36:20','2025-05-18 15:36:16',NULL),
(28,'164da810-83ea-4b1e-b03b-d05ab23ae9e6','title','content',3000,250,0,1,'2025-06-05 21:32:08','2025-06-05 21:32:08',NULL);
/*!40000 ALTER TABLE `problem_basic` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `problem_category`
--

DROP TABLE IF EXISTS `problem_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `problem_category` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `problem_id` int(11) DEFAULT NULL COMMENT '问题ID',
  `category_id` int(11) DEFAULT NULL COMMENT '分类ID',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `problem_category`
--

LOCK TABLES `problem_category` WRITE;
/*!40000 ALTER TABLE `problem_category` DISABLE KEYS */;
INSERT INTO `problem_category` VALUES
(25,1,1,'2025-05-19 16:16:46','2025-05-18 16:16:49',NULL),
(26,1,2,'2025-05-18 16:16:55','2025-05-18 16:16:58',NULL),
(27,2,1,'2025-05-18 16:17:20','2025-05-18 16:17:22',NULL),
(28,28,1,'2025-06-05 21:32:08','2025-06-05 21:32:08',NULL);
/*!40000 ALTER TABLE `problem_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `submit_basic`
--

DROP TABLE IF EXISTS `submit_basic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `submit_basic` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) DEFAULT NULL,
  `problem_identity` varchar(36) DEFAULT NULL COMMENT '问题的唯一标识',
  `user_identity` varchar(36) DEFAULT NULL COMMENT '用户的唯一标识',
  `path` varchar(255) DEFAULT NULL COMMENT '代码路径',
  `status` tinyint(1) DEFAULT -1 COMMENT '【-1-待判断，1-答案正确，2-答案错误，3-运行超时，4-运行超内存】',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `submit_basic`
--

LOCK TABLES `submit_basic` WRITE;
/*!40000 ALTER TABLE `submit_basic` DISABLE KEYS */;
INSERT INTO `submit_basic` VALUES
(1,'submit_1','problem_1','user_1','/code/x1.go',0,'2025-05-18 20:05:16','2025-05-18 20:05:16',NULL),
(2,'submit_2','problem_2','user_2','/code/x2.go',0,'2025-05-18 20:06:40','2025-05-18 20:06:41',NULL),
(18,'35c66d26-6600-4c0c-a073-e309e54aa0bb','problem_3','user_1','code/0527dafd-1cdc-49b3-a26a-726e5b9ca0d7/main.go',1,'2025-06-06 00:06:15','2025-06-06 00:06:15',NULL),
(19,'8a37f9c0-dc74-4843-8917-aee1e51b9a39','164da810-83ea-4b1e-b03b-d05ab23ae9e6','user_1','code/4f955d95-cdd2-42bc-814f-88f54b839018/main.go',2,'2025-06-06 00:06:45','2025-06-06 00:06:45',NULL);
/*!40000 ALTER TABLE `submit_basic` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `test_case`
--

DROP TABLE IF EXISTS `test_case`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `test_case` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) NOT NULL,
  `problem_identity` varchar(36) DEFAULT NULL,
  `input` text DEFAULT NULL,
  `output` text DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `test_case`
--

LOCK TABLES `test_case` WRITE;
/*!40000 ALTER TABLE `test_case` DISABLE KEYS */;
INSERT INTO `test_case` VALUES
(6,'09e8d897-b00f-41d9-aa50-723604bd8ce1','164da810-83ea-4b1e-b03b-d05ab23ae9e6','1 2 \n','3\n','2025-06-05 21:32:08','2025-06-05 21:32:08',NULL);
/*!40000 ALTER TABLE `test_case` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_basic`
--

DROP TABLE IF EXISTS `user_basic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_basic` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) DEFAULT NULL COMMENT '唯一标识',
  `name` varchar(100) DEFAULT NULL COMMENT '名称',
  `password` varchar(32) DEFAULT NULL COMMENT '密码',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `mail` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `pass_num` int(11) DEFAULT 0 COMMENT '完成问题的个数',
  `submit_num` int(11) DEFAULT 0 COMMENT '总提交次数',
  `is_admin` tinyint(1) DEFAULT 0 COMMENT '是否是管理员【0-否，1-是】',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_basic`
--

LOCK TABLES `user_basic` WRITE;
/*!40000 ALTER TABLE `user_basic` DISABLE KEYS */;
INSERT INTO `user_basic` VALUES
(1,'user_1','xieyang','e10adc3949ba59abbe56e057f20f883e','15472341778','xie@qq.com',101,102,1,'2025-05-18 19:35:57','2025-05-18 19:35:59',NULL),
(5,'f4d7359c-aeac-4f78-b57b-fc5e22c6b684','xieyang1','e10adc3949ba59abbe56e057f20f883e','','xie59199@gmail.com',200,300,0,'2025-05-19 10:46:00','2025-05-19 10:46:00',NULL);
/*!40000 ALTER TABLE `user_basic` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'gin-gorm'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*M!100616 SET NOTE_VERBOSITY=@OLD_NOTE_VERBOSITY */;

-- Dump completed on 2025-06-06  1:35:18
