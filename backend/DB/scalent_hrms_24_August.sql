/*M!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19  Distrib 10.6.23-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: scalent_hrms
-- ------------------------------------------------------
-- Server version	10.6.23-MariaDB-0ubuntu0.22.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `attendance_logs`
--

DROP TABLE IF EXISTS `attendance_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `attendance_logs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `emp_id` varchar(64) NOT NULL,
  `timestamp` datetime NOT NULL,
  `status` int(11) NOT NULL,
  `punch` int(11) NOT NULL,
  `attendance_state` varchar(32) NOT NULL,
  `device_name` varchar(128) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_attendance_logs` (`uid`,`emp_id`,`timestamp`,`status`,`punch`,`device_name`),
  KEY `idx_attendance_logs_emp_id` (`emp_id`),
  KEY `idx_attendance_logs_timestamp` (`timestamp`),
  KEY `idx_attendance_logs_uid` (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=10317 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `attendance_logs`
--

LOCK TABLES `attendance_logs` WRITE;
/*!40000 ALTER TABLE `attendance_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `attendance_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary table structure for view `casbin_policy`
--

DROP TABLE IF EXISTS `casbin_policy`;
/*!50001 DROP VIEW IF EXISTS `casbin_policy`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8mb4;
/*!50001 CREATE VIEW `casbin_policy` AS SELECT
 1 AS `p_type`,
  1 AS `v0`,
  1 AS `v1`,
  1 AS `v2`,
  1 AS `v3`,
  1 AS `v4`,
  1 AS `v5` */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `employees`
--

DROP TABLE IF EXISTS `employees`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `employees` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `emp_id` varchar(64) NOT NULL,
  `emp_name` varchar(255) NOT NULL,
  `privilege` int(11) DEFAULT 0,
  `password` varchar(64) DEFAULT '',
  `group_id` varchar(64) DEFAULT '',
  `card` bigint(20) DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_employees_emp_id` (`emp_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1361 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `employees`
--

LOCK TABLES `employees` WRITE;
/*!40000 ALTER TABLE `employees` DISABLE KEYS */;
/*!40000 ALTER TABLE `employees` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_accesses`
--

DROP TABLE IF EXISTS `role_accesses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `role_accesses` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `role_id` int(10) unsigned NOT NULL,
  `uri` varchar(255) NOT NULL,
  `method` varchar(10) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `role_accesses_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=514 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_accesses`
--

LOCK TABLES `role_accesses` WRITE;
/*!40000 ALTER TABLE `role_accesses` DISABLE KEYS */;
INSERT INTO `role_accesses` VALUES (1,1,'/scalent-hrms/login','POST'),(2,1,'/scalent-hrms/email-list/list','GET'),(3,1,'/scalent-hrms/email-list/{id}','GET'),(4,1,'/scalent-hrms/email-list/{id}','PUT'),(5,1,'/scalent-hrms/email-list/{id}','PATCH'),(6,1,'/scalent-hrms/email-list/{id}','DELETE'),(7,1,'/scalent-hrms/email-list','POST'),(8,1,'/scalent-hrms/user/list','GET'),(9,1,'/scalent-hrms/user','POST'),(10,1,'/scalent-hrms/user/{id}','PUT'),(11,1,'/scalent-hrms/user/{id}','DELETE'),(12,1,'/scalent-hrms/user/{id}','PATCH'),(13,1,'/scalent-hrms/user/{id}','GET'),(14,1,'/scalent-hrms/lead/list','GET'),(15,1,'/scalent-hrms/lead/{id}','GET'),(16,1,'/scalent-hrms/lead/{id}','PUT'),(17,1,'/scalent-hrms/lead/{id}','PATCH'),(18,1,'/scalent-hrms/lead/{id}','DELETE'),(19,1,'/scalent-hrms/lead','POST'),(20,1,'/scalent-hrms/logout','POST'),(507,1,'/scalent-hrms/lead/reverify/{id}','PUT'),(508,1,'/scalent-hrms/user-setting/me','GET'),(509,1,'/scalent-hrms/user-setting','POST'),(510,1,'/scalent-hrms/user-setting/{id}','PUT'),(511,1,'/scalent-hrms/user-setting/{id}','PATCH'),(512,1,'/scalent-hrms/user-setting/list','GET'),(513,1,'/scalent-hrms/user/change-password','PATCH');
/*!40000 ALTER TABLE `role_accesses` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `roles` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `code` varchar(50) NOT NULL,
  `status` enum('ACTIVE','INACTIVE') DEFAULT 'ACTIVE',
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES (1,'ADMIN','ADMIN','ACTIVE','2026-06-30 11:50:10',NULL,NULL);
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_roles`
--

DROP TABLE IF EXISTS `user_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_roles` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `role_id` int(10) unsigned NOT NULL,
  `created_at` datetime DEFAULT current_timestamp(),
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_user_role` (`user_id`,`role_id`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `user_roles_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `user_roles_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_roles`
--

LOCK TABLES `user_roles` WRITE;
/*!40000 ALTER TABLE `user_roles` DISABLE KEYS */;
INSERT INTO `user_roles` VALUES (1,1,1,'2026-06-30 17:21:22',NULL),(3,68,1,'2026-07-24 14:47:02',NULL),(4,69,1,'2026-08-17 13:30:19',NULL),(5,70,1,'2026-08-17 13:40:44',NULL),(6,71,1,'2026-08-17 13:47:28',NULL),(7,77,1,'2026-08-18 05:46:32',NULL),(8,78,1,'2026-08-18 06:49:17',NULL),(9,79,1,'2026-08-18 06:50:48',NULL),(10,80,1,'2026-08-18 07:41:16',NULL),(11,81,1,'2026-08-18 08:04:19',NULL),(12,82,1,'2026-08-18 08:09:12',NULL),(13,83,1,'2026-08-18 08:09:31',NULL);
/*!40000 ALTER TABLE `user_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_settings`
--

DROP TABLE IF EXISTS `user_settings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_settings` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `verification_interval` int(10) unsigned NOT NULL DEFAULT 15,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_settings_user_id` (`user_id`),
  CONSTRAINT `fk_user_settings_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_settings`
--

LOCK TABLES `user_settings` WRITE;
/*!40000 ALTER TABLE `user_settings` DISABLE KEYS */;
INSERT INTO `user_settings` VALUES (8,1,60,'2026-08-18 07:43:57','2026-08-19 08:03:34',NULL);
/*!40000 ALTER TABLE `user_settings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(100) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `is_password_set` enum('YES','NO') NOT NULL DEFAULT 'NO',
  `status` enum('ACTIVE','INACTIVE') DEFAULT 'ACTIVE',
  `last_login_at` datetime DEFAULT NULL,
  `session_token` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deleted_at` datetime DEFAULT NULL,
  `active_email` varchar(255) GENERATED ALWAYS AS (case when `deleted_at` is null then `email` else NULL end) STORED,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_active_user_email` (`active_email`)
) ENGINE=InnoDB AUTO_INCREMENT=84 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'admin@scalent.io','$2a$10$mxfrVZCZsPr4oCRONJ26kuszlbe7WUg3XisFC7pIt3voAxGLtbFiS','YES','ACTIVE','2026-08-24 12:21:11','2fa85908-f826-466a-8977-ef67e3fccd6b39166168-52ac-47e3-9ab7-e444b9b9ca85','2026-06-30 16:03:47','2026-08-24 12:21:11',NULL,'admin@scalent.io'),(68,'pankaj@scalent.io','$2a$10$2mYmYYBvOteZ0.RgfqgAkuizVjl1cGXj/aPTb5MAqXoDrvgKvirSq','YES','ACTIVE','2026-08-24 06:50:14',NULL,'2026-07-24 14:44:16','2026-08-24 12:20:14',NULL,'pankaj@scalent.io'),(69,'namita@scalent.io','$2a$10$v70YpQl/hd0UfbwkWyrc8uZp7TThB9cFl78dqAsR2GjB7GMn4Wy16','NO','ACTIVE',NULL,'','2026-08-17 13:30:19','2026-08-18 06:48:39','2026-08-18 04:48:39',NULL),(70,'gokul@scalent.io','$2a$10$ITcxz6MsiUq.twWEfUNJj.chF0sRtfDjaRbvo1Hdxiwj2oKB2TARW','NO','ACTIVE',NULL,'','2026-08-17 13:40:44','2026-08-17 13:56:38','2026-08-17 11:56:38',NULL),(71,'gstarle92@gmail.com','$2a$10$XA50QzXZ5gV3aqU9HNgycO/ymmcYmnK5WP6yQG5yJO2RsDhkxlv1S','NO','ACTIVE',NULL,'','2026-08-17 13:47:28','2026-08-17 13:56:35','2026-08-17 11:56:35',NULL),(77,'gokul@scalent.io','$2a$10$GKNynURkXFTRhU/Leh55gOy97Bd6GjM7ZWk/RNwJM0ZVcIED7zSiO','NO','ACTIVE',NULL,'','2026-08-18 05:46:32','2026-08-18 06:50:35','2026-08-18 04:50:36',NULL),(78,'namita@scalent.io','$2a$10$xcA0X/Svn2eN/Ta5mrEvROp8Qrj5x/e5OU1NE8fDxaH425SiIyEm.','NO','ACTIVE',NULL,'','2026-08-18 06:49:17','2026-08-18 07:41:04','2026-08-18 05:41:04',NULL),(79,'gokul@scalent.io','$2a$10$dkOWX1kzHLwcNNrOhP/14uU/t4XHzbW3YYoqmvseDnOvOs78LTIsW','NO','ACTIVE',NULL,'','2026-08-18 06:50:48','2026-08-18 06:50:48',NULL,'gokul@scalent.io'),(80,'namita@scalent.io','$2a$10$PjyBJbqz3tDht3bGD92/fe6yc/CKC5fSEJNEox27laaYEA8h3kywe','NO','ACTIVE',NULL,'','2026-08-18 07:41:16','2026-08-18 08:04:05','2026-08-18 06:04:06',NULL),(81,'namita@scalent.io','$2a$10$68QEwxYut3d/h6xcbrCXLux4P9haMehzNyuKs7rgiL7gCf0NpXsK6','NO','ACTIVE',NULL,'','2026-08-18 08:04:19','2026-08-18 08:08:47','2026-08-18 06:08:47',NULL),(82,'testingmailora@yopmail.com','$2a$10$WEj7iJ/voy6NQG05IBL1Y.2TuD.UPWKGRDpJ.zeYFJXeZDVxtyaMi','NO','ACTIVE',NULL,'','2026-08-18 08:09:12','2026-08-18 08:10:44','2026-08-18 06:10:45',NULL),(83,'namita@scalent.io','$2a$10$SXegnuB6sQ2nuZ2toxQa9.srLRusKdGi6MmuZnRJdzOcNXuRBSJY6','YES','ACTIVE','2026-08-18 07:50:08',NULL,'2026-08-18 08:09:31','2026-08-18 09:50:07',NULL,'namita@scalent.io');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Final view structure for view `casbin_policy`
--

/*!50001 DROP VIEW IF EXISTS `casbin_policy`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `casbin_policy` AS select 'p' AS `p_type`,`r`.`code` AS `v0`,`ra`.`uri` AS `v1`,`ra`.`method` AS `v2`,'' AS `v3`,'' AS `v4`,'' AS `v5` from (`role_accesses` `ra` join `roles` `r` on(`ra`.`role_id` = `r`.`id`)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-08-24 15:30:49
