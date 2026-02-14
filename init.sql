-- MySQL dump 10.13  Distrib 8.0.31, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: dare2lead
-- ------------------------------------------------------
-- Server version	8.0.31

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `academic_year`
--

CREATE DATABASE IF NOT EXISTS `dare2lead` CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci;
USE `dare2lead`;

DROP TABLE IF EXISTS `academic_year`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `academic_year` (
  `Academic_YearID` bigint NOT NULL AUTO_INCREMENT,
  `YearName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`Academic_YearID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `academic_year`
--

LOCK TABLES `academic_year` WRITE;
/*!40000 ALTER TABLE `academic_year` DISABLE KEYS */;
/*!40000 ALTER TABLE `academic_year` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `accesslevel`
--

DROP TABLE IF EXISTS `accesslevel`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `accesslevel` (
  `AccessLevelID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `LevelOfAccess` int DEFAULT '0',
  UNIQUE KEY `AccessLevelID` (`AccessLevelID`),
  KEY `WDIDX_accesslevel_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_accesslevel_levelOfAccess` (`LevelOfAccess`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `accesslevel`
--

LOCK TABLES `accesslevel` WRITE;
/*!40000 ALTER TABLE `accesslevel` DISABLE KEYS */;
/*!40000 ALTER TABLE `accesslevel` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `action`
--

DROP TABLE IF EXISTS `action`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `action` (
  `WhenAct` bigint DEFAULT '0',
  `Lesson_outlineID` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  `Act` tinyint DEFAULT '0',
  KEY `WDIDX_action_Lesson_outlineID` (`Lesson_outlineID`),
  KEY `WDIDX_action_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_action_IndividualsIDPosAct` (`IndividualsID`,`Act`),
  KEY `WDIDX_action_whenActIndividualsID` (`WhenAct`,`IndividualsID`),
  KEY `WDIDX_actions_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_actions_Lesson_outlineID` (`Lesson_outlineID`),
  KEY `WDIDX_actions_WDIDX_action_IndividualsIDPosAct` (`IndividualsID`,`Act`),
  KEY `WDIDX_actions_WDIDX_action_whenActIndividualsID` (`WhenAct`,`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `action`
--

LOCK TABLES `action` WRITE;
/*!40000 ALTER TABLE `action` DISABLE KEYS */;
/*!40000 ALTER TABLE `action` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `actions`
--

DROP TABLE IF EXISTS `actions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `actions` (
  `Act` tinyint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  `Lesson_outlineID` bigint DEFAULT '0',
  `WhenAct` bigint DEFAULT '0'
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `actions`
--

LOCK TABLES `actions` WRITE;
/*!40000 ALTER TABLE `actions` DISABLE KEYS */;
/*!40000 ALTER TABLE `actions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `applicants`
--

DROP TABLE IF EXISTS `applicants`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `applicants` (
  `ApplicantsID` bigint NOT NULL,
  `PupilForname` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PupilSurname` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PupilMiddleNames` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PupilDOB` bigint DEFAULT NULL,
  `UPN` varchar(13) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Gender` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `CurrentLA` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `CurrentSchool` varchar(9) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `LookedAfter` tinyint DEFAULT '0',
  `Statement` tinyint DEFAULT '0',
  `PupilAddressLine1` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PupilAddressLine2` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PupilAddressLine3` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PupilAddressLine4` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PupilAddressLine5` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ApplicantForename` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ApplicantSurname` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ApplicantMiddleName` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `RelationshipToStudent` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ApplicantAddress1` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ApplicantAddress2` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ApplicantAddress3` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ApplicantAddress4` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ApplicantAddress5` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PhoneContactID` int DEFAULT '0',
  `Email_Address` varchar(260) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PrefContact` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ApplicationOutcome` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `YearEntryGroup` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ADTFileStatus` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `CouncilTaxRef` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ApplicantTitle` varchar(9) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Sibling` tinyint DEFAULT '0',
  `AptitudeCode` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `AptitudeText` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `ADTFileText` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `PupiilVerifiedAddress` tinyint DEFAULT '0',
  `PreferenceRank` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PreferenceReason` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PreferenceText` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `FaithText` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `DistanceFromSchool` int DEFAULT '0',
  `CrownService` tinyint DEFAULT '0',
  `InYearEntry` tinyint DEFAULT '0',
  `DatePlaceRequired` bigint DEFAULT NULL,
  `InYearText` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `ApplAddressAsPupil` tinyint DEFAULT '0',
  `OfferStatus` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`ApplicantsID`),
  KEY `WDIDX_applicants_pupilForname` (`PupilForname`),
  KEY `WDIDX_applicants_PupilSurname` (`PupilSurname`),
  KEY `WDIDX_applicants_gender` (`Gender`),
  KEY `WDIDX_applicants_LookedAfter` (`LookedAfter`),
  KEY `WDIDX_applicants_Statement` (`Statement`),
  KEY `WDIDX_applicants_PhoneContactID` (`PhoneContactID`),
  KEY `WDIDX_applicants_Email_Address` (`Email_Address`),
  KEY `WDIDX_applicants_PrefContact` (`PrefContact`),
  KEY `WDIDX_applicants_ApplicationOutcome` (`ApplicationOutcome`),
  KEY `WDIDX_applicants_YearEntryGroup` (`YearEntryGroup`),
  KEY `WDIDX_applicants_ADTFileStatus` (`ADTFileStatus`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `applicants`
--

LOCK TABLES `applicants` WRITE;
/*!40000 ALTER TABLE `applicants` DISABLE KEYS */;
/*!40000 ALTER TABLE `applicants` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `applicationreferences`
--

DROP TABLE IF EXISTS `applicationreferences`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `applicationreferences` (
  `ApplicationReferencesID` bigint NOT NULL,
  `UPN` varchar(13) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ApplicationRef` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`ApplicationReferencesID`),
  KEY `WDIDX_applicationreferences_UPN` (`UPN`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `applicationreferences`
--

LOCK TABLES `applicationreferences` WRITE;
/*!40000 ALTER TABLE `applicationreferences` DISABLE KEYS */;
/*!40000 ALTER TABLE `applicationreferences` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `appointment_appliesto`
--

DROP TABLE IF EXISTS `appointment_appliesto`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `appointment_appliesto` (
  `Groupname` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Appointment_appliesToID` bigint NOT NULL AUTO_INCREMENT,
  `IndivOrGroupID` bigint DEFAULT '0',
  `AppointmentsID` bigint DEFAULT '0',
  `Staff` tinyint DEFAULT '0',
  `Student` tinyint DEFAULT '0',
  `Parents` tinyint DEFAULT '0',
  PRIMARY KEY (`Appointment_appliesToID`),
  KEY `WDIDX_appointment_appliesto_groupname` (`Groupname`),
  KEY `WDIDX_appointment_appliesto_indivOrGroupID` (`IndivOrGroupID`),
  KEY `WDIDX_appointment_appliesto_AppointmentsID` (`AppointmentsID`),
  KEY `WDIDX_appointment_appliesto_Staff` (`Staff`),
  KEY `WDIDX_appointment_appliesto_Student` (`Student`),
  KEY `WDIDX_appointment_appliesto_Parents` (`Parents`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `appointment_appliesto`
--

LOCK TABLES `appointment_appliesto` WRITE;
/*!40000 ALTER TABLE `appointment_appliesto` DISABLE KEYS */;
/*!40000 ALTER TABLE `appointment_appliesto` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `appointments`
--

DROP TABLE IF EXISTS `appointments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `appointments` (
  `AppointmentsID` bigint NOT NULL AUTO_INCREMENT,
  `AppName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `location` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AppStart` bigint DEFAULT NULL,
  `EndTime` time DEFAULT NULL,
  `ReplaceTimetable` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `MemoID` bigint DEFAULT '0',
  `AppointRepeatID` bigint DEFAULT '0',
  `UserID` bigint DEFAULT NULL,
  PRIMARY KEY (`AppointmentsID`),
  KEY `WDIDX_appointments_AppName` (`AppName`),
  KEY `WDIDX_appointments_location` (`location`),
  KEY `WDIDX_appointments_AppStart` (`AppStart`),
  KEY `WDIDX_appointments_EndTime` (`EndTime`),
  KEY `WDIDX_appointments_memoID` (`MemoID`),
  KEY `WDIDX_appointments_appointRepeatID` (`AppointRepeatID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `appointments`
--

LOCK TABLES `appointments` WRITE;
/*!40000 ALTER TABLE `appointments` DISABLE KEYS */;
/*!40000 ALTER TABLE `appointments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `appointrepeat`
--

DROP TABLE IF EXISTS `appointrepeat`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `appointrepeat` (
  `AppointRepeatID` bigint NOT NULL AUTO_INCREMENT,
  `RepWeeks` int DEFAULT '0',
  `AppointmentsID` bigint DEFAULT '0',
  UNIQUE KEY `AppointRepeatID` (`AppointRepeatID`),
  KEY `WDIDX_appointrepeat_AppointmentsID` (`AppointmentsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `appointrepeat`
--

LOCK TABLES `appointrepeat` WRITE;
/*!40000 ALTER TABLE `appointrepeat` DISABLE KEYS */;
/*!40000 ALTER TABLE `appointrepeat` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `assessment_description`
--

DROP TABLE IF EXISTS `assessment_description`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `assessment_description` (
  `Assessment_DescriptionID` bigint NOT NULL AUTO_INCREMENT,
  `AssessmentCode` varchar(6) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Assessment_description` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `MadatoryCTF` tinyint DEFAULT '0',
  `ToBeDetermined` tinyint DEFAULT '0',
  `Year7_2007equiv` tinyint DEFAULT '0',
  PRIMARY KEY (`Assessment_DescriptionID`),
  KEY `WDIDX_assessment_description_MadatoryCTF` (`MadatoryCTF`),
  KEY `WDIDX_assessment_description_ToBeDetermined` (`ToBeDetermined`),
  KEY `WDIDX_assessment_description_year7_2007equiv` (`Year7_2007equiv`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `assessment_description`
--

LOCK TABLES `assessment_description` WRITE;
/*!40000 ALTER TABLE `assessment_description` DISABLE KEYS */;
/*!40000 ALTER TABLE `assessment_description` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `assessment_focuses`
--

DROP TABLE IF EXISTS `assessment_focuses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `assessment_focuses` (
  `Assessment_FocusesID` bigint NOT NULL AUTO_INCREMENT,
  `Subject_code` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Component` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Component_Description` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `AssessmentFocus` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `FocusID` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Assessment_Focus_Description` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`Assessment_FocusesID`),
  KEY `WDIDX_assessment_focuses_subjectcode` (`Subject_code`),
  KEY `WDIDX_assessment_focuses_Component` (`Component`),
  KEY `WDIDX_assessment_focuses_Component_Description` (`Component_Description`),
  KEY `WDIDX_assessment_focuses_AssessmentFocus` (`AssessmentFocus`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `assessment_focuses`
--

LOCK TABLES `assessment_focuses` WRITE;
/*!40000 ALTER TABLE `assessment_focuses` DISABLE KEYS */;
/*!40000 ALTER TABLE `assessment_focuses` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `attendance`
--

DROP TABLE IF EXISTS `attendance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `attendance` (
  `Attendance_Code` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `SessionsID` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  KEY `WDIDX_attendance_SessionsID` (`SessionsID`),
  KEY `WDIDX_attendance_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_attendance_WDIDX_attendance_SessionsIDAttendance_Code` (`SessionsID`,`Attendance_Code`),
  KEY `WDIDX_attendance_WDIDX_Attendance_sessionsIDIndividualsID` (`SessionsID`,`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `attendance`
--

LOCK TABLES `attendance` WRITE;
/*!40000 ALTER TABLE `attendance` DISABLE KEYS */;
/*!40000 ALTER TABLE `attendance` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `attendance_key`
--

DROP TABLE IF EXISTS `attendance_key`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `attendance_key` (
  `Attendance_keyID` bigint NOT NULL AUTO_INCREMENT,
  `Description` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`Attendance_keyID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `attendance_key`
--

LOCK TABLES `attendance_key` WRITE;
/*!40000 ALTER TABLE `attendance_key` DISABLE KEYS */;
/*!40000 ALTER TABLE `attendance_key` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `audio_assignments`
--

DROP TABLE IF EXISTS `audio_assignments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `audio_assignments` (
  `Audio_AssignmentsID` bigint NOT NULL AUTO_INCREMENT,
  `Major_Incident_audioID` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  `Known` tinyint DEFAULT '0',
  UNIQUE KEY `Audio_AssignmentsID` (`Audio_AssignmentsID`),
  KEY `WDIDX_audio_assignments_Major_Incident_audioID` (`Major_Incident_audioID`),
  KEY `WDIDX_audio_assignments_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_audio_assignments_WDIDX_Audio_Assignments_Major_Incid00000` (`Major_Incident_audioID`,`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `audio_assignments`
--

LOCK TABLES `audio_assignments` WRITE;
/*!40000 ALTER TABLE `audio_assignments` DISABLE KEYS */;
/*!40000 ALTER TABLE `audio_assignments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `audio_file`
--

DROP TABLE IF EXISTS `audio_file`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `audio_file` (
  `Audio_fileID` bigint NOT NULL AUTO_INCREMENT,
  `TargetID` bigint DEFAULT '0',
  `Types_keyID` bigint DEFAULT '0',
  `Filename` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `CloudID` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  UNIQUE KEY `Audio_fileID` (`Audio_fileID`),
  KEY `WDIDX_audio_file_TargetID` (`TargetID`),
  KEY `WDIDX_audio_file_Types_keyID` (`Types_keyID`),
  KEY `WDIDX_audio_file_filename` (`Filename`),
  KEY `WDIDX_audio_file_CloudID` (`CloudID`),
  KEY `WDIDX_audio_file_WDIDX_audio_file_WDIDX_audio_file_WDIDX_Au00001` (`TargetID`,`Types_keyID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `audio_file`
--

LOCK TABLES `audio_file` WRITE;
/*!40000 ALTER TABLE `audio_file` DISABLE KEYS */;
/*!40000 ALTER TABLE `audio_file` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_group`
--

DROP TABLE IF EXISTS `auth_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `auth_group` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `Name` varchar(150) NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `Name` (`Name`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_group`
--

LOCK TABLES `auth_group` WRITE;
/*!40000 ALTER TABLE `auth_group` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_group_permissions`
--

DROP TABLE IF EXISTS `auth_group_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `auth_group_permissions` (
  `ID` bigint NOT NULL AUTO_INCREMENT,
  `Group_id` int NOT NULL,
  `Permission_id` int NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `auth_group_permissions_group_id_permission_id_0cd325b0_uniq` (`Group_id`,`Permission_id`),
  KEY `auth_group_permissions_group_id_b120cbf9` (`Group_id`),
  KEY `auth_group_permissions_permission_id_84c5c92e` (`Permission_id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_group_permissions`
--

LOCK TABLES `auth_group_permissions` WRITE;
/*!40000 ALTER TABLE `auth_group_permissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth_group_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_permission`
--

DROP TABLE IF EXISTS `auth_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `auth_permission` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `Name` varchar(255) NOT NULL,
  `Content_type_id` int NOT NULL,
  `Codename` varchar(100) NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `auth_permission_content_type_id_codename_01ab375a_uniq` (`Content_type_id`,`Codename`),
  KEY `auth_permission_content_type_id_2f476e4b` (`Content_type_id`)
) ENGINE=MyISAM AUTO_INCREMENT=25 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_permission`
--

LOCK TABLES `auth_permission` WRITE;
/*!40000 ALTER TABLE `auth_permission` DISABLE KEYS */;
INSERT INTO `auth_permission` VALUES (1,'Can add log Entry',1,'add_logentry'),(2,'Can change log Entry',1,'change_logentry'),(3,'Can delete log Entry',1,'delete_logentry'),(4,'Can view log Entry',1,'view_logentry'),(5,'Can add permission',2,'add_permission'),(6,'Can change permission',2,'change_permission'),(7,'Can delete permission',2,'delete_permission'),(8,'Can view permission',2,'view_permission'),(9,'Can add group',3,'add_group'),(10,'Can change group',3,'change_group'),(11,'Can delete group',3,'delete_group'),(12,'Can view group',3,'view_group'),(13,'Can add user',4,'add_user'),(14,'Can change user',4,'change_user'),(15,'Can delete user',4,'delete_user'),(16,'Can view user',4,'view_user'),(17,'Can add content type',5,'add_contenttype'),(18,'Can change content type',5,'change_contenttype'),(19,'Can delete content type',5,'delete_contenttype'),(20,'Can view content type',5,'view_contenttype'),(21,'Can add session',6,'add_session'),(22,'Can change session',6,'change_session'),(23,'Can delete session',6,'delete_session'),(24,'Can view session',6,'view_session');
/*!40000 ALTER TABLE `auth_permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_user`
--

DROP TABLE IF EXISTS `auth_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `auth_user` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `Password` varchar(128) NOT NULL,
  `Last_login` datetime(6) DEFAULT NULL,
  `Os_superuser` tinyint(1) NOT NULL,
  `Username` varchar(150) NOT NULL,
  `First_name` varchar(150) NOT NULL,
  `Last_name` varchar(150) NOT NULL,
  `Email` varchar(254) NOT NULL,
  `Is_staff` tinyint(1) NOT NULL,
  `Is_active` tinyint(1) NOT NULL,
  `Date_joined` datetime(6) NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `Username` (`Username`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_user`
--

LOCK TABLES `auth_user` WRITE;
/*!40000 ALTER TABLE `auth_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_user_groups`
--

DROP TABLE IF EXISTS `auth_user_groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `auth_user_groups` (
  `ID` bigint NOT NULL AUTO_INCREMENT,
  `User_id` int NOT NULL,
  `Group_id` int NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `auth_user_groups_user_id_group_id_94350c0c_uniq` (`User_id`,`Group_id`),
  KEY `auth_user_groups_user_id_6a12ed8b` (`User_id`),
  KEY `auth_user_groups_group_id_97559544` (`Group_id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_user_groups`
--

LOCK TABLES `auth_user_groups` WRITE;
/*!40000 ALTER TABLE `auth_user_groups` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth_user_groups` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_user_user_permissions`
--

DROP TABLE IF EXISTS `auth_user_user_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `auth_user_user_permissions` (
  `ID` bigint NOT NULL AUTO_INCREMENT,
  `User_id` int NOT NULL,
  `Permission_id` int NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `auth_user_user_permissions_user_id_permission_id_14a6b632_uniq` (`User_id`,`Permission_id`),
  KEY `auth_user_user_permissions_user_id_a95ead1b` (`User_id`),
  KEY `auth_user_user_permissions_permission_id_1fbb5f2c` (`Permission_id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_user_user_permissions`
--

LOCK TABLES `auth_user_user_permissions` WRITE;
/*!40000 ALTER TABLE `auth_user_user_permissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth_user_user_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `book_references`
--

DROP TABLE IF EXISTS `book_references`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `book_references` (
  `Book_referencesID` bigint NOT NULL AUTO_INCREMENT,
  `BooksID` bigint DEFAULT '0',
  `Pages` varchar(15) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Lesson_FilesID` bigint DEFAULT '0',
  `Volume` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`Book_referencesID`),
  KEY `WDIDX_book_references_BooksID` (`BooksID`),
  KEY `WDIDX_book_references_Lesson_FilesID` (`Lesson_FilesID`),
  KEY `WDIDX_book_references_WDIDX_book_references_WDIDX_Book_refe00002` (`BooksID`,`Pages`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `book_references`
--

LOCK TABLES `book_references` WRITE;
/*!40000 ALTER TABLE `book_references` DISABLE KEYS */;
/*!40000 ALTER TABLE `book_references` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `books`
--

DROP TABLE IF EXISTS `books`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `books` (
  `BooksID` bigint NOT NULL AUTO_INCREMENT,
  `SubjectID` bigint DEFAULT '0',
  `ISBN` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Publisher` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `BookTitle` varchar(300) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PublishedDate` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `JournalOrBook` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Edition` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `BookDescription` varchar(300) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  UNIQUE KEY `BooksID` (`BooksID`),
  KEY `WDIDX_books_SubjectID` (`SubjectID`),
  KEY `WDIDX_books_ISBN` (`ISBN`),
  KEY `WDIDX_books_BookTitle` (`BookTitle`),
  KEY `WDIDX_books_JournalOrBook` (`JournalOrBook`),
  KEY `WDIDX_books_WDIDX_books_SubjectIDBookTitle` (`SubjectID`,`BookTitle`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `books`
--

LOCK TABLES `books` WRITE;
/*!40000 ALTER TABLE `books` DISABLE KEYS */;
/*!40000 ALTER TABLE `books` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `candtcodes`
--

DROP TABLE IF EXISTS `candtcodes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `candtcodes` (
  `CandTCodesID` bigint NOT NULL,
  `DCode` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `OptionCode` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `OptionDesc` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`CandTCodesID`),
  KEY `WDIDX_candtcodes_DCode` (`DCode`),
  KEY `WDIDX_candtcodes_WDIDX_candtcodes_WDIDX_CandTCodes_DCodeOpt00003` (`DCode`,`OptionCode`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `candtcodes`
--

LOCK TABLES `candtcodes` WRITE;
/*!40000 ALTER TABLE `candtcodes` DISABLE KEYS */;
/*!40000 ALTER TABLE `candtcodes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cbds_data_items`
--

DROP TABLE IF EXISTS `cbds_data_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cbds_data_items` (
  `CBDS_Data_itemsID` bigint NOT NULL AUTO_INCREMENT,
  `CBDS_Level` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `CBDS_Module` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Identifier_1` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Identifier_2` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Data_Item_Name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Description` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Type_and_Format` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Code_set_Valid_Values` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Item_Level_Validation` varchar(1000) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `XML_Tag` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Class_CBDS` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ESCS_ISB_Compliance` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Status` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ESCS_Data_Model` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `SIF_Specification` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Collection_Notes` varchar(300) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Mandatory_Notes` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Output_Notes` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `History_Notes` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Multiplicity_Notes` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Source_s` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Item_CBDS_Level` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `RFC_Number` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Creation_Date` bigint DEFAULT NULL,
  `Implementation_Date` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Cut_Off_Date` bigint DEFAULT NULL,
  `Data_Item_review_date` bigint DEFAULT NULL,
  `Spring_School_Census_2016` tinyint DEFAULT '0',
  `Summer_School_Census_2016` int DEFAULT '0',
  `Early_Years_2016` int DEFAULT '0',
  `CIN_Cenus_2016_17` int DEFAULT '0',
  `Alt_Provision_2016` int DEFAULT '0',
  `EYFSP_2016` int DEFAULT '0',
  `Phonics_2016` int DEFAULT '0',
  `Keystage_1_2016` int DEFAULT '0',
  `School_Workforce_Census_2016` int DEFAULT '0',
  `Autumn_School_Census_2016` int DEFAULT '0',
  `Child_Social_Work_Worforce_2016` int DEFAULT '0',
  `Spring_School_Census_2017` int DEFAULT '0',
  `Summer_School_Census_2017` int DEFAULT '0',
  `CIN_Census_2017_18` int DEFAULT '0',
  `Alt_Provision_2017` int DEFAULT '0',
  `Early_Years_2017` int DEFAULT '0',
  `SLASC_2017` int DEFAULT '0',
  `EYFSP_2017` int DEFAULT '0',
  `Phonics_2017` int DEFAULT '0',
  `Keystage_1_2017` int DEFAULT '0',
  `School_Workforce_Census_2017` int DEFAULT '0',
  `Child_Social_Work_Worforce_2017` int DEFAULT '0',
  `Autumn_School_Census_2017` int DEFAULT '0',
  `Spring_School_Census_2018` int DEFAULT '0',
  `Summer_School_Census_2018` int DEFAULT '0',
  `Alt_Provision_2018` int DEFAULT '0',
  `ADT_for_admission_in_Sept_2016` int DEFAULT '0',
  `ALT_for_admission_in_Sept_2016` int DEFAULT '0',
  `ASL_for_admission_in_Sept_2016` int DEFAULT '0',
  `APT_for_admission_in_Sept_2016` int DEFAULT '0',
  `AOT_for_admission_in_Sept_2016` int DEFAULT '0',
  `CTF_16` int DEFAULT '0',
  `CTF_17` int DEFAULT '0',
  `Windev_tables_equiv` varchar(150) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`CBDS_Data_itemsID`),
  KEY `WDIDX_cbds_data_items_CBDS_Level` (`CBDS_Level`),
  KEY `WDIDX_cbds_data_items_CBDS_Module` (`CBDS_Module`),
  KEY `WDIDX_cbds_data_items_Identifier_1` (`Identifier_1`),
  KEY `WDIDX_cbds_data_items_Identifier_2` (`Identifier_2`),
  KEY `WDIDX_cbds_data_items_Data_Item_Name` (`Data_Item_Name`),
  KEY `WDIDX_cbds_data_items_XML_Tag` (`XML_Tag`),
  KEY `WDIDX_cbds_data_items_Class_CBDS` (`Class_CBDS`),
  KEY `WDIDX_cbds_data_items_ESCS_ISB_Compliance` (`ESCS_ISB_Compliance`),
  KEY `WDIDX_cbds_data_items_Status` (`Status`),
  KEY `WDIDX_cbds_data_items_SIF_Specification` (`SIF_Specification`),
  KEY `WDIDX_cbds_data_items_Mandatory_Notes` (`Mandatory_Notes`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cbds_data_items`
--

LOCK TABLES `cbds_data_items` WRITE;
/*!40000 ALTER TABLE `cbds_data_items` DISABLE KEYS */;
/*!40000 ALTER TABLE `cbds_data_items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `childreninneed`
--

DROP TABLE IF EXISTS `childreninneed`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `childreninneed` (
  `ChildrenInNeedID` bigint NOT NULL AUTO_INCREMENT,
  `Referral_Date` bigint NOT NULL DEFAULT '0',
  `Primary_Need_Code` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `CIN_Closure_Date` bigint NOT NULL DEFAULT '0',
  `Reason_for_Closure` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `IndividualsID` bigint DEFAULT '0',
  `Open_Case_Information` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Child_Protection_Plan_start_date` bigint NOT NULL DEFAULT '0',
  `Number_of_previous_Child_Protection_Plans` tinyint DEFAULT '0',
  `Target_Date_for_Initial_Child_Protection_Conference` bigint NOT NULL DEFAULT '0',
  `Date_of_Initial_Child_Protection_Conference` bigint NOT NULL DEFAULT '0',
  `ICPC_not_required` tinyint DEFAULT '0',
  `Referral_No_Further_Action_Flag` tinyint DEFAULT '0',
  `Child_Protection_Plan_End_Date` bigint NOT NULL DEFAULT '0',
  `Initial_Category_of_Abuse` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Plan_Review_Date` bigint NOT NULL DEFAULT '0',
  `Latest_Category_of_Abuse` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Section_47_Enquiry_Actual_Start_Date` bigint NOT NULL DEFAULT '0',
  `Seen_by_Social_Worker` tinyint DEFAULT '0',
  `Assessment_Actual_Start_Date` bigint NOT NULL DEFAULT '0',
  `Source_of_Referral` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Assessment_Authorisation_Date` bigint NOT NULL DEFAULT '0',
  `Assessment_Internal_Review_Date` bigint NOT NULL DEFAULT '0',
  `Factors_Identified_at_Assessment` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`ChildrenInNeedID`),
  KEY `WDIDX_childreninneed_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `childreninneed`
--

LOCK TABLES `childreninneed` WRITE;
/*!40000 ALTER TABLE `childreninneed` DISABLE KEYS */;
/*!40000 ALTER TABLE `childreninneed` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `collections`
--

DROP TABLE IF EXISTS `collections`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `collections` (
  `CollectionsID` bigint NOT NULL AUTO_INCREMENT,
  `Collection_Code` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Collection_Name` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`CollectionsID`),
  KEY `WDIDX_collections_Collection_Name` (`Collection_Name`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `collections`
--

LOCK TABLES `collections` WRITE;
/*!40000 ALTER TABLE `collections` DISABLE KEYS */;
/*!40000 ALTER TABLE `collections` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `contactrelationship`
--

DROP TABLE IF EXISTS `contactrelationship`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `contactrelationship` (
  `ContactRelationshipID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `ContactsID` bigint DEFAULT '0',
  `Contact_Relationship_to_pupil_child` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `WDIDX_contactrelationship_Archived` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Archived` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Priority` int DEFAULT NULL,
  PRIMARY KEY (`ContactRelationshipID`),
  KEY `WDIDX_contactrelationship_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_contactrelationship_ContactsID` (`ContactsID`),
  KEY `WDIDX_contactrelationship_Priority` (`Priority`),
  KEY `WDIDX_contactrelationship_WDIDX_contactrelationship_Individ00004` (`IndividualsID`,`ContactsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `contactrelationship`
--

LOCK TABLES `contactrelationship` WRITE;
/*!40000 ALTER TABLE `contactrelationship` DISABLE KEYS */;
/*!40000 ALTER TABLE `contactrelationship` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `contacts`
--

DROP TABLE IF EXISTS `contacts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `contacts` (
  `ContactsID` bigint NOT NULL AUTO_INCREMENT,
  `Title` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Surname` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Forename` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `MiddleNames` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Gender` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Responsibility` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Email_Address` varchar(260) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `ZipCode` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`ContactsID`),
  KEY `WDIDX_contacts_surname` (`Surname`),
  KEY `WDIDX_contacts_forename` (`Forename`),
  KEY `WDIDX_contacts_gender` (`Gender`),
  KEY `WDIDX_contacts_Email_Address` (`Email_Address`),
  KEY `WDIDX_contacts_WDIDX_contacts_forenamesurname` (`Forename`,`Surname`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `contacts`
--

LOCK TABLES `contacts` WRITE;
/*!40000 ALTER TABLE `contacts` DISABLE KEYS */;
/*!40000 ALTER TABLE `contacts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_plan_header`
--

DROP TABLE IF EXISTS `course_plan_header`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_plan_header` (
  `Course_plan_headerID` bigint NOT NULL AUTO_INCREMENT,
  `AllorSubOnly` int DEFAULT NULL,
  `Title` varchar(80) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `StartDate` bigint NOT NULL DEFAULT '0',
  `DeadlinesID` bigint DEFAULT '0',
  `lessonsAvail` int DEFAULT '0',
  `CourseID` bigint DEFAULT '0',
  `GroupAppliesTo` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `EndDate` bigint DEFAULT '0',
  UNIQUE KEY `Course_plan_headerID` (`Course_plan_headerID`),
  KEY `WDIDX_course_plan_header_StartDate` (`StartDate`),
  KEY `WDIDX_course_plan_header_CourseID` (`CourseID`),
  KEY `WDIDX_course_plan_header_GroupAppliesTo` (`GroupAppliesTo`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_plan_header`
--

LOCK TABLES `course_plan_header` WRITE;
/*!40000 ALTER TABLE `course_plan_header` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_plan_header` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_segment`
--

DROP TABLE IF EXISTS `course_segment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_segment` (
  `Course_segmentID` bigint NOT NULL AUTO_INCREMENT,
  `Course_plan_headerID` bigint DEFAULT '0',
  `TopicID` bigint DEFAULT '0',
  `Last_SessionsID` bigint DEFAULT '0',
  `LessonsUsed` int DEFAULT '0',
  UNIQUE KEY `Course_segmentID` (`Course_segmentID`),
  KEY `WDIDX_course_segment_course_plan_headerID` (`Course_plan_headerID`),
  KEY `WDIDX_course_segment_TopicsID` (`TopicID`),
  KEY `WDIDX_course_segment_Last_SessionsID` (`Last_SessionsID`),
  KEY `WDIDX_course_segment_WDIDX_Course_segment_Course_plan_heade00005` (`Course_plan_headerID`,`TopicID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_segment`
--

LOCK TABLES `course_segment` WRITE;
/*!40000 ALTER TABLE `course_segment` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_segment` ENABLE KEYS */;
UNLOCK TABLES;


--
-- Table structure for table `day_timings`
--

DROP TABLE IF EXISTS `day_timings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `day_timings` (
  `Period_Id` bigint NOT NULL AUTO_INCREMENT,
  `Dy` int DEFAULT '0',
  `Finish` time DEFAULT NULL,
  `AppliesToGroup` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '0',
  `PeriodName` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '0',
  `Strt` time DEFAULT NULL,
  `Academic_YearID` bigint DEFAULT '0',
  PRIMARY KEY (`Period_Id`),
  KEY `WDIDX_day_timings_dy` (`Dy`),
  KEY `WDIDX_day_timings_finish` (`Finish`),
  KEY `WDIDX_day_timings_appliesToGroup` (`AppliesToGroup`),
  KEY `WDIDX_day_timings_periodName` (`PeriodName`),
  KEY `WDIDX_day_timings_strt` (`Strt`),
  KEY `WDIDX_day_timings_Academic_YearID` (`Academic_YearID`),
  KEY `WDIDX_day_timings_WDIDX_day_timings_WDIDX_day_timings_dayap00006` (`Dy`,`AppliesToGroup`),
  KEY `WDIDX_day_timings_WDIDX_day_timings_WDIDX_day_timings_daype00007` (`Dy`,`PeriodName`,`AppliesToGroup`),
  KEY `WDIDX_day_timings_WDIDX_day_timings_WDIDX_day_timings_Optim00008` (`Dy`,`Finish`),
  KEY `WDIDX_day_timings_WDIDX_day_timings_WDIDX_day_timings_Optim00009` (`Strt`,`Dy`,`Finish`),
  KEY `WDIDX_day_timings_WDIDX_day_timings_WDIDX_day_timings_Optim00010` (`Strt`,`Period_Id`,`Dy`,`Finish`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `day_timings`
--

LOCK TABLES `day_timings` WRITE;
/*!40000 ALTER TABLE `day_timings` DISABLE KEYS */;
/*!40000 ALTER TABLE `day_timings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `deadcheck`
--

DROP TABLE IF EXISTS `deadcheck`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `deadcheck` (
  `DeadCheckID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `DeadlinesID` bigint DEFAULT '0',
  `Completed` bigint DEFAULT NULL,
  UNIQUE KEY `DeadCheckID` (`DeadCheckID`),
  UNIQUE KEY `DeadlinesID` (`DeadlinesID`),
  KEY `WDIDX_deadcheck_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_deadcheck_Completed` (`Completed`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `deadcheck`
--

LOCK TABLES `deadcheck` WRITE;
/*!40000 ALTER TABLE `deadcheck` DISABLE KEYS */;
/*!40000 ALTER TABLE `deadcheck` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `deadline_appliesto`
--

DROP TABLE IF EXISTS `deadline_appliesto`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `deadline_appliesto` (
  `deadline_appliesToID` bigint NOT NULL AUTO_INCREMENT,
  `Groupname` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Staff` tinyint DEFAULT '0',
  `Student` tinyint DEFAULT '0',
  `Parents` tinyint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  `DeadlinesID` bigint DEFAULT '0',
  PRIMARY KEY (`deadline_appliesToID`),
  KEY `WDIDX_deadline_appliesto_groupname` (`Groupname`),
  KEY `WDIDX_deadline_appliesto_Staff` (`Staff`),
  KEY `WDIDX_deadline_appliesto_Student` (`Student`),
  KEY `WDIDX_deadline_appliesto_Parents` (`Parents`),
  KEY `WDIDX_deadline_appliesto_DeadlinesID` (`DeadlinesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `deadline_appliesto`
--

LOCK TABLES `deadline_appliesto` WRITE;
/*!40000 ALTER TABLE `deadline_appliesto` DISABLE KEYS */;
/*!40000 ALTER TABLE `deadline_appliesto` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `deadlines`
--

DROP TABLE IF EXISTS `deadlines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `deadlines` (
  `DeadlinesID` bigint NOT NULL AUTO_INCREMENT,
  `Title` varchar(256) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `UserID` bigint DEFAULT '0',
  `DeadlineTime` time DEFAULT NULL,
  `DeadlineDate` bigint NOT NULL DEFAULT '0',
  `Completed` bigint NOT NULL DEFAULT '0',
  `Types_keyID` bigint DEFAULT '0',
  `MemoID` bigint DEFAULT '0',
  `Groupname` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`DeadlinesID`),
  KEY `WDIDX_deadlines_UserID` (`UserID`),
  KEY `WDIDX_deadlines_DeadlineDate` (`DeadlineDate`),
  KEY `WDIDX_deadlines_Completed` (`Completed`),
  KEY `WDIDX_deadlines_Types_keyID` (`Types_keyID`),
  KEY `WDIDX_deadlines_memoID` (`MemoID`),
  KEY `WDIDX_deadlines_groupname` (`Groupname`),
  KEY `WDIDX_deadlines_WDIDX_deadlines_WDIDX_deadlines_WDIDX_deadl00011` (`Title`,`DeadlineDate`),
  KEY `WDIDX_deadlines_WDIDX_deadlines_WDIDX_deadlines_WDIDX_Deadl00012` (`Title`,`DeadlineDate`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `deadlines`
--

LOCK TABLES `deadlines` WRITE;
/*!40000 ALTER TABLE `deadlines` DISABLE KEYS */;
/*!40000 ALTER TABLE `deadlines` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dismissedlast`
--

DROP TABLE IF EXISTS `dismissedlast`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dismissedlast` (
  `IndividualsID` bigint DEFAULT '0',
  `SessionsID` bigint DEFAULT '0',
  `LastDis` bigint DEFAULT '0',
  KEY `WDIDX_dismissedlast_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_dismissedlast_SessionsID` (`SessionsID`),
  KEY `WDIDX_dismissedlast_LastDis` (`LastDis`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dismissedlast`
--

LOCK TABLES `dismissedlast` WRITE;
/*!40000 ALTER TABLE `dismissedlast` DISABLE KEYS */;
/*!40000 ALTER TABLE `dismissedlast` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `django_admin_log`
--

DROP TABLE IF EXISTS `django_admin_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `django_admin_log` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `Action_time` datetime(6) NOT NULL,
  `Object_id` longtext,
  `Object_repr` varchar(200) NOT NULL,
  `Action_flag` smallint unsigned NOT NULL,
  `Change_message` longtext NOT NULL,
  `Content_type_id` int DEFAULT NULL,
  `User_id` int NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `django_admin_log_content_type_id_c4bce8eb` (`Content_type_id`),
  KEY `django_admin_log_user_id_c564eba6` (`User_id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `django_admin_log`
--

LOCK TABLES `django_admin_log` WRITE;
/*!40000 ALTER TABLE `django_admin_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `django_admin_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `django_content_type`
--

DROP TABLE IF EXISTS `django_content_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `django_content_type` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `App_label` varchar(100) NOT NULL,
  `Model` varchar(100) NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `django_content_type_app_label_model_76bd3d3b_uniq` (`App_label`,`Model`)
) ENGINE=MyISAM AUTO_INCREMENT=7 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `django_content_type`
--

LOCK TABLES `django_content_type` WRITE;
/*!40000 ALTER TABLE `django_content_type` DISABLE KEYS */;
INSERT INTO `django_content_type` VALUES (1,'admin','logentry'),(2,'auth','permission'),(3,'auth','group'),(4,'auth','user'),(5,'contenttypes','contenttype'),(6,'sessions','session');
/*!40000 ALTER TABLE `django_content_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `django_migrations`
--

DROP TABLE IF EXISTS `django_migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `django_migrations` (
  `ID` bigint NOT NULL AUTO_INCREMENT,
  `App` varchar(255) NOT NULL,
  `Name` varchar(255) NOT NULL,
  `Applied` datetime(6) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=MyISAM AUTO_INCREMENT=19 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `django_migrations`
--

LOCK TABLES `django_migrations` WRITE;
/*!40000 ALTER TABLE `django_migrations` DISABLE KEYS */;
INSERT INTO `django_migrations` VALUES (1,'contenttypes','0001_initial','2022-09-08 11:21:45.733617'),(2,'auth','0001_initial','2022-09-08 11:21:46.003561'),(3,'admin','0001_initial','2022-09-08 11:21:46.065457'),(4,'admin','0002_logentry_remove_auto_add','2022-09-08 11:21:46.086377'),(5,'admin','0003_logentry_add_action_flag_choices','2022-09-08 11:21:46.106316'),(6,'contenttypes','0002_remove_content_type_name','2022-09-08 11:21:46.154804'),(7,'auth','0002_alter_permission_name_max_length','2022-09-08 11:21:46.182491'),(8,'auth','0003_alter_user_email_max_length','2022-09-08 11:21:46.215485'),(9,'auth','0004_alter_user_username_opts','2022-09-08 11:21:46.237426'),(10,'auth','0005_alter_user_last_login_null','2022-09-08 11:21:46.264266'),(11,'auth','0006_require_contenttypes_0002','2022-09-08 11:21:46.278232'),(12,'auth','0007_alter_validators_add_error_messages','2022-09-08 11:21:46.307155'),(13,'auth','0008_alter_user_username_max_length','2022-09-08 11:21:46.337123'),(14,'auth','0009_alter_user_last_name_max_length','2022-09-08 11:21:46.365207'),(15,'auth','0010_alter_group_name_max_length','2022-09-08 11:21:46.396995'),(16,'auth','0011_update_proxy_permissions','2022-09-08 11:21:46.416943'),(17,'auth','0012_alter_user_first_name_max_length','2022-09-08 11:21:46.445847'),(18,'sessions','0001_initial','2022-09-08 11:21:46.476006');
/*!40000 ALTER TABLE `django_migrations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `django_session`
--

DROP TABLE IF EXISTS `django_session`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `django_session` (
  `Session_key` varchar(40) NOT NULL,
  `Session_data` longtext NOT NULL,
  `Expire_date` datetime(6) NOT NULL,
  PRIMARY KEY (`Session_key`),
  KEY `django_session_expire_date_a5c62663` (`Expire_date`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `django_session`
--

LOCK TABLES `django_session` WRITE;
/*!40000 ALTER TABLE `django_session` DISABLE KEYS */;
/*!40000 ALTER TABLE `django_session` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dwelling_occupants`
--

DROP TABLE IF EXISTS `dwelling_occupants`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dwelling_occupants` (
  `Dwelling_OccupantsID` bigint NOT NULL AUTO_INCREMENT,
  `DwellingsID` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  `ContactsID` bigint DEFAULT '0',
  PRIMARY KEY (`Dwelling_OccupantsID`),
  KEY `WDIDX_dwelling_occupants_DwellingsID` (`DwellingsID`),
  KEY `WDIDX_dwelling_occupants_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_dwelling_occupants_ContactsID` (`ContactsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dwelling_occupants`
--

LOCK TABLES `dwelling_occupants` WRITE;
/*!40000 ALTER TABLE `dwelling_occupants` DISABLE KEYS */;
/*!40000 ALTER TABLE `dwelling_occupants` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dwellings`
--

DROP TABLE IF EXISTS `dwellings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dwellings` (
  `DwellingsID` bigint NOT NULL AUTO_INCREMENT,
  `SubNameOrNumber` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `BuildNameNumber` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Street` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Locality` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Town` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AdminArea` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `PostTown` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `PostCode` varchar(8) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `PropertyEasting` decimal(24,6) DEFAULT '0.000000',
  `PropertyNorthing` decimal(24,6) DEFAULT '0.000000',
  `AddressLine1` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AddressLine2` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AddressLine3` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AddressLine4` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AddressLine5` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `PhoneContactID` int DEFAULT '0',
  PRIMARY KEY (`DwellingsID`),
  KEY `WDIDX_dwellings_PhoneContactID` (`PhoneContactID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dwellings`
--

LOCK TABLES `dwellings` WRITE;
/*!40000 ALTER TABLE `dwellings` DISABLE KEYS */;
/*!40000 ALTER TABLE `dwellings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `effortrecorder`
--

DROP TABLE IF EXISTS `effortrecorder`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `effortrecorder` (
  `EffortRecorderID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `Lesson_outlineID` bigint DEFAULT '0',
  `Effortscore` tinyint DEFAULT '0',
  UNIQUE KEY `EffortRecorderID` (`EffortRecorderID`),
  KEY `WDIDX_effortrecorder_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_effortrecorder_Lesson_outlineID` (`Lesson_outlineID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `effortrecorder`
--

LOCK TABLES `effortrecorder` WRITE;
/*!40000 ALTER TABLE `effortrecorder` DISABLE KEYS */;
/*!40000 ALTER TABLE `effortrecorder` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `emaillinkedarchive`
--
DROP TABLE IF EXISTS `microsoft_tokens`;

CREATE TABLE IF NOT EXISTS `microsoft_tokens` (
    `MicrosoftTokensID` bigint NOT NULL AUTO_INCREMENT,
    `IndividualsID` bigint NOT NULL,
    `access_token` LONGTEXT NOT NULL COMMENT 'Encrypted access token',
    `refresh_token` LONGTEXT NOT NULL COMMENT 'Encrypted refresh token',
    `expires_at` TIMESTAMP NOT NULL,
    `consent_given` BOOLEAN DEFAULT FALSE,
    `last_refreshed` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`MicrosoftTokensID`),
    UNIQUE KEY `IndividualsID` (`IndividualsID`),
    KEY `idx_expires_at` (`expires_at`),
    KEY `idx_consent_given` (`consent_given`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1 COMMENT='Stores encrypted Microsoft OAuth2 tokens for external students';
UNLOCK TABLES;


DROP TABLE IF EXISTS `emaillinkedarchive`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `emaillinkedarchive` (
  `EmailLinkedArchiveID` bigint NOT NULL AUTO_INCREMENT,
  `EmailID` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ConversationID` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `From` varchar(50) DEFAULT NULL,
  `Preview` varchar(200) DEFAULT NULL,
  `Received` bigint DEFAULT '0',
  `Subject` varchar(100) DEFAULT NULL,
  UNIQUE KEY `EmailLinkedArchiveID` (`EmailLinkedArchiveID`),
  KEY `WDIDX_emaillinkedarchive_EmailID` (`EmailID`),
  KEY `WDIDX_emaillinkedarchive_ConversationID` (`ConversationID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `emaillinkedarchive`
--

LOCK TABLES `emaillinkedarchive` WRITE;
/*!40000 ALTER TABLE `emaillinkedarchive` DISABLE KEYS */;
/*!40000 ALTER TABLE `emaillinkedarchive` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `entities`
--

DROP TABLE IF EXISTS `entities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `entities` (
  `EntitiesID` bigint NOT NULL AUTO_INCREMENT,
  `Ent_type` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Tier` int DEFAULT '0',
  PRIMARY KEY (`EntitiesID`),
  KEY `WDIDX_entities_Ent_type` (`Ent_type`),
  KEY `WDIDX_entities_tier` (`Tier`),
  KEY `WDIDX_entities_WDIDX_entities_WDIDX_Entities_OptimCompKey_t00013` (`Tier`,`EntitiesID`)
) ENGINE=MyISAM AUTO_INCREMENT=6 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `entities`
--

LOCK TABLES `entities` WRITE;
/*!40000 ALTER TABLE `entities` DISABLE KEYS */;
INSERT INTO `entities` VALUES (1,'student',1),(2,'teacher',1),(3,'class',2),(4,'year group',3),(5,'keystage',4);
/*!40000 ALTER TABLE `entities` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ethnicityDetails`
--

DROP TABLE IF EXISTS `ethnicitydetails`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ethnicitydetails` (
  `EthnicityID` bigint NOT NULL AUTO_INCREMENT,
  `ExtendedCode` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `ExtendedName` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `MainCode` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `SubCategory` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `CategoryNotes` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `MainCategory` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`EthnicityID`),
  KEY `WDIDX_ethnicitydetails_MainCode` (`MainCode`),
  KEY `WDIDX_ethnicitydetails_SubCategory` (`SubCategory`),
  KEY `WDIDX_ethnicitydetails_MainCategory` (`MainCategory`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ethnicitydetails`
--

LOCK TABLES `ethnicitydetails` WRITE;
/*!40000 ALTER TABLE `ethnicitydetails` DISABLE KEYS */;
/*!40000 ALTER TABLE `ethnicitydetails` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `externalfile`
--

DROP TABLE IF EXISTS `externalfile`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `externalfile` (
  `ExternalFileID` bigint NOT NULL AUTO_INCREMENT,
  `Filename` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  UNIQUE KEY `ExternalFileID` (`ExternalFileID`),
  KEY `WDIDX_externalfile_filename` (`Filename`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `externalfile`
--

LOCK TABLES `externalfile` WRITE;
/*!40000 ALTER TABLE `externalfile` DISABLE KEYS */;
/*!40000 ALTER TABLE `externalfile` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `file_types`
--

DROP TABLE IF EXISTS `file_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `file_types` (
  `File_TypesID` bigint NOT NULL AUTO_INCREMENT,
  `TypeName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Description` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`File_TypesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `file_types`
--

LOCK TABLES `file_types` WRITE;
/*!40000 ALTER TABLE `file_types` DISABLE KEYS */;
/*!40000 ALTER TABLE `file_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `filegroup`
--

DROP TABLE IF EXISTS `filegroup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `filegroup` (
  `FileGroupID` bigint NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`FileGroupID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `filegroup`
--

LOCK TABLES `filegroup` WRITE;
/*!40000 ALTER TABLE `filegroup` DISABLE KEYS */;
/*!40000 ALTER TABLE `filegroup` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `foc_reg_header`
--

DROP TABLE IF EXISTS `foc_reg_header`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `foc_reg_header` (
  `Foc_reg_headerID` bigint NOT NULL AUTO_INCREMENT,
  `Register_typesID` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  `Summary` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Conatct_name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Contact_email` varchar(260) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Foc_reg_headerID`),
  KEY `WDIDX_foc_reg_header_Register_typesID` (`Register_typesID`),
  KEY `WDIDX_foc_reg_header_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_foc_reg_header_WDIDX_foc_reg_header_WDIDX_Foc_reg_hea00014` (`Register_typesID`,`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `foc_reg_header`
--

LOCK TABLES `foc_reg_header` WRITE;
/*!40000 ALTER TABLE `foc_reg_header` DISABLE KEYS */;
/*!40000 ALTER TABLE `foc_reg_header` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `foc_reg_log`
--

DROP TABLE IF EXISTS `foc_reg_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `foc_reg_log` (
  `Foc_reg_logID` bigint NOT NULL AUTO_INCREMENT,
  `Foc_reg_headerID` bigint DEFAULT '0',
  `Date_of_entry` bigint NOT NULL DEFAULT '0',
  `Issue_or_support` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Entry` varchar(300) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Foc_reg_logID`),
  KEY `WDIDX_foc_reg_log_Foc_reg_headerID` (`Foc_reg_headerID`),
  KEY `WDIDX_foc_reg_log_Date_of_entry` (`Date_of_entry`),
  KEY `WDIDX_foc_reg_log_Issue_or_support` (`Issue_or_support`),
  KEY `WDIDX_foc_reg_log_WDIDX_foc_reg_log_WDIDX_Foc_reg_log_Date_00015` (`Date_of_entry`,`Foc_reg_headerID`,`Issue_or_support`),
  KEY `WDIDX_foc_reg_log_WDIDX_foc_reg_log_WDIDX_Foc_reg_log_Optim00016` (`Issue_or_support`,`Foc_reg_headerID`),
  KEY `WDIDX_foc_reg_log_WDIDX_foc_reg_log_WDIDX_Foc_reg_log_Optim00017` (`Issue_or_support`,`Entry`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `foc_reg_log`
--

LOCK TABLES `foc_reg_log` WRITE;
/*!40000 ALTER TABLE `foc_reg_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `foc_reg_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `focusregcolourcodes`
--

DROP TABLE IF EXISTS `focusregcolourcodes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `focusregcolourcodes` (
  `FocusRegColourCodesID` bigint NOT NULL AUTO_INCREMENT,
  `ColorCode` varchar(7) DEFAULT NULL,
  PRIMARY KEY (`FocusRegColourCodesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `focusregcolourcodes`
--

LOCK TABLES `focusregcolourcodes` WRITE;
/*!40000 ALTER TABLE `focusregcolourcodes` DISABLE KEYS */;
/*!40000 ALTER TABLE `focusregcolourcodes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `folderids`
--

DROP TABLE IF EXISTS `folderids`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `folderids` (
  `FolderIDsID` bigint NOT NULL AUTO_INCREMENT,
  `FolderName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `OneDriveFolderID` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`FolderIDsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `folderids`
--

LOCK TABLES `folderids` WRITE;
/*!40000 ALTER TABLE `folderids` DISABLE KEYS */;
/*!40000 ALTER TABLE `folderids` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gclientdata`
--

DROP TABLE IF EXISTS `gclientdata`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gclientdata` (
  `gClientDataID` bigint NOT NULL AUTO_INCREMENT,
  `access_token` varchar(1024) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `id_token` varchar(1024) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `expires_in` int DEFAULT '0',
  `token_type` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `refresh_token` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  UNIQUE KEY `gClientDataID` (`gClientDataID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gclientdata`
--

LOCK TABLES `gclientdata` WRITE;
/*!40000 ALTER TABLE `gclientdata` DISABLE KEYS */;
/*!40000 ALTER TABLE `gclientdata` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gradeboundaries`
--

DROP TABLE IF EXISTS `gradeboundaries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gradeboundaries` (
  `GradeBoundariesID` bigint NOT NULL AUTO_INCREMENT,
  `Lesson_FilesID` bigint DEFAULT '0',
  `Mark` int DEFAULT '0',
  `Grade` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`GradeBoundariesID`),
  KEY `WDIDX_gradeboundaries_Lesson_FilesID` (`Lesson_FilesID`),
  KEY `WDIDX_gradeboundaries_Mark` (`Mark`),
  KEY `WDIDX_gradeboundaries_grade` (`Grade`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gradeboundaries`
--

LOCK TABLES `gradeboundaries` WRITE;
/*!40000 ALTER TABLE `gradeboundaries` DISABLE KEYS */;
/*!40000 ALTER TABLE `gradeboundaries` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `group_hierarchy`
--

DROP TABLE IF EXISTS `group_hierarchy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `group_hierarchy` (
  `Parent_group` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '0',
  `Offspring_group` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '0',
  KEY `WDIDX_group_hierarchy_Parent_group` (`Parent_group`),
  KEY `WDIDX_group_hierarchy_Offspring_group` (`Offspring_group`),
  KEY `WDIDX_group_hierarchy_WDIDX_group_hierarchy_WDIDX_group_hie00018` (`Parent_group`,`Offspring_group`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `group_hierarchy`
--

LOCK TABLES `group_hierarchy` WRITE;
/*!40000 ALTER TABLE `group_hierarchy` DISABLE KEYS */;
/*!40000 ALTER TABLE `group_hierarchy` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `group_membership`
--

DROP TABLE IF EXISTS `group_membership`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `group_membership` (
  `GroupMembershipID` int NOT NULL AUTO_INCREMENT,
  `Groupname` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `IndividualsID` bigint NOT NULL DEFAULT '0',
  `Target` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `SubjectID` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`GroupMembershipID`),
  KEY `WDIDX_group_membership_groupname` (`Groupname`),
  KEY `WDIDX_group_membership_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_group_membership_SubjectID` (`SubjectID`),
  KEY `WDIDX_group_membership_WDIDX_group_membership_groupnameIndi00019` (`Groupname`,`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `group_membership`
--

LOCK TABLES `group_membership` WRITE;
/*!40000 ALTER TABLE `group_membership` DISABLE KEYS */;
/*!40000 ALTER TABLE `group_membership` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `group_totals`
--

DROP TABLE IF EXISTS `group_totals`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `group_totals` (
  `Group_totalsID` bigint NOT NULL AUTO_INCREMENT,
  `Groupname` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `TargetsID` bigint DEFAULT '0',
  `Totalpoints` tinyint DEFAULT '0',
  `Lesson_outlineID` bigint DEFAULT '0',
  `Threshold` int DEFAULT '0',
  `GroupOrEveryone` tinyint DEFAULT '0',
  UNIQUE KEY `Group_totalsID` (`Group_totalsID`),
  KEY `WDIDX_group_totals_groupname` (`Groupname`),
  KEY `WDIDX_group_totals_TargetsID` (`TargetsID`),
  KEY `WDIDX_group_totals_Lesson_outlineID` (`Lesson_outlineID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `group_totals`
--

LOCK TABLES `group_totals` WRITE;
/*!40000 ALTER TABLE `group_totals` DISABLE KEYS */;
/*!40000 ALTER TABLE `group_totals` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `helppoints`
--

DROP TABLE IF EXISTS `helppoints`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `helppoints` (
  `HelpPointsID` bigint NOT NULL AUTO_INCREMENT,
  `DateAndTime` timestamp NULL DEFAULT NULL,
  `IndividualsID` bigint DEFAULT '0',
  `RecordedBy` bigint DEFAULT '0',
  UNIQUE KEY `HelpPointsID` (`HelpPointsID`),
  KEY `WDIDX_helppoints_DateAndTime` (`DateAndTime`),
  KEY `WDIDX_helppoints_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_helppoints_RecordedBy` (`RecordedBy`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `helppoints`
--

LOCK TABLES `helppoints` WRITE;
/*!40000 ALTER TABLE `helppoints` DISABLE KEYS */;
/*!40000 ALTER TABLE `helppoints` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `image_assignments`
--

DROP TABLE IF EXISTS `image_assignments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `image_assignments` (
  `Image_assignmentsID` bigint NOT NULL AUTO_INCREMENT,
  `Major_Incident_imagesID` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  `Known` tinyint DEFAULT '0',
  UNIQUE KEY `Image_assignmentsID` (`Image_assignmentsID`),
  KEY `WDIDX_image_assignments_Major_Incident_imagesID` (`Major_Incident_imagesID`),
  KEY `WDIDX_image_assignments_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_image_assignments_WDIDX_Image_assignments_Major_Incid00020` (`Major_Incident_imagesID`,`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `image_assignments`
--

LOCK TABLES `image_assignments` WRITE;
/*!40000 ALTER TABLE `image_assignments` DISABLE KEYS */;
/*!40000 ALTER TABLE `image_assignments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `incident_unknowns`
--

DROP TABLE IF EXISTS `incident_unknowns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `incident_unknowns` (
  `Incident_unknownsID` bigint NOT NULL AUTO_INCREMENT,
  `Name` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Major_transgression_logID` bigint DEFAULT '0',
  UNIQUE KEY `Incident_unknownsID` (`Incident_unknownsID`),
  KEY `WDIDX_incident_unknowns_Name` (`Name`),
  KEY `WDIDX_incident_unknowns_Major_transgression_logID` (`Major_transgression_logID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `incident_unknowns`
--

LOCK TABLES `incident_unknowns` WRITE;
/*!40000 ALTER TABLE `incident_unknowns` DISABLE KEYS */;
/*!40000 ALTER TABLE `incident_unknowns` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `individuals`
--

DROP TABLE IF EXISTS `individuals`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `individuals` (
  `IndividualsID` bigint NOT NULL AUTO_INCREMENT,
  `Forename` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Surname` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `preferred` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Dob` bigint DEFAULT NULL,
  `Title` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Gender` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '0',
  `Photo` varchar(250) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT 'default-user-image.png',
  `Email_Address` varchar(260) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `EntitiesID` bigint DEFAULT '0',
  `Exam_number` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `LoginCode` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ID_numb` varchar(15) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Passw` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Reading_age` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Spelling_age` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `HasPhoto` tinyint DEFAULT '0',
  `MicrosoftID` varchar(100) DEFAULT NULL,
  `MiddleNames` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `FormerSurname` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `UPN_Unknown_Reason` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Learner_Support_Code` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Pupil_s_Former_UPN` varchar(13) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Pupil_Preferred_Surname` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Uniq_Learner_Number_ULN` varchar(11) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `LA_Child_ID` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Expected_Date_of_Birth` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Date_of_Death` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Unique_Candidate_Identifier_UCI` varchar(13) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Uniq_Pupil_Number_UPN` varchar(13) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `RestPassw` tinyint DEFAULT '0',
  `GgID` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `EthnicityID` bigint DEFAULT NULL,
  UNIQUE KEY `IndividualsID` (`IndividualsID`),
  KEY `WDIDX_individuals_forename` (`Forename`),
  KEY `WDIDX_individuals_surname` (`Surname`),
  KEY `WDIDX_individuals_preferred` (`preferred`),
  KEY `WDIDX_individuals_dob` (`Dob`),
  KEY `WDIDX_individuals_gender` (`Gender`),
  KEY `WDIDX_individuals_Photo` (`Photo`),
  KEY `WDIDX_individuals_Email_Address` (`Email_Address`),
  KEY `WDIDX_individuals_EntitiesID` (`EntitiesID`),
  KEY `WDIDX_individuals_exam_number` (`Exam_number`),
  KEY `WDIDX_individuals_loginCode` (`LoginCode`),
  KEY `WDIDX_individuals_ID_numb` (`ID_numb`),
  KEY `WDIDX_individuals_MicrosoftID` (`MicrosoftID`),
  KEY `WDIDX_individuals_FormerSurname` (`FormerSurname`),
  KEY `WDIDX_individuals_Pupil_s_Former_UPN` (`Pupil_s_Former_UPN`),
  KEY `WDIDX_individuals_Pupil_Preferred_Surname` (`Pupil_Preferred_Surname`),
  KEY `WDIDX_individuals_Uniq_Learner_Number_ULN` (`Uniq_Learner_Number_ULN`),
  KEY `WDIDX_individuals_LA_Child_ID` (`LA_Child_ID`),
  KEY `WDIDX_individuals_Expected_Date_of_Birth` (`Expected_Date_of_Birth`),
  KEY `WDIDX_individuals_Date_of_Death` (`Date_of_Death`),
  KEY `WDIDX_individuals_Unique_Candidate_Identifier_UCI` (`Unique_Candidate_Identifier_UCI`),
  KEY `WDIDX_individuals_Uniq_Pupil_Number_UPN` (`Uniq_Pupil_Number_UPN`),
  KEY `WDIDX_individuals_RestPassw` (`RestPassw`),
  KEY `WDIDX_individuals_GgID` (`GgID`),
  KEY `WDIDX_individuals_EthnicityID` (`EthnicityID`),
  KEY `WDIDX_individuals_WDIDX_individuals_WDIDX_Individuals_foren00021` (`Forename`,`Surname`,`Dob`),
  KEY `WDIDX_individuals_WDIDX_individuals_WDIDX_Individuals_foren00022` (`Forename`,`Surname`),
  KEY `WDIDX_individuals_WDIDX_individuals_WDIDX_Individuals_Optim00023` (`Surname`,`Forename`,`EntitiesID`,`preferred`),
  KEY `WDIDX_individuals_WDIDX_individuals_WDIDX_Individuals_Optim00024` (`IndividualsID`,`EntitiesID`),
  KEY `WDIDX_individuals_WDIDX_individuals_WDIDX_Individuals_Optim00025` (`Surname`,`EntitiesID`,`Forename`),
  KEY `WDIDX_individuals_WDIDX_individuals_WDIDX_Individuals_Optim00026` (`Surname`,`EntitiesID`,`LoginCode`),
  KEY `WDIDX_individuals_WDIDX_individuals_WDIDX_Individuals_Optim00027` (`IndividualsID`,`Surname`),
  KEY `WDIDX_individuals_WDIDX_individuals_WDIDX_Individuals_Optim00028` (`Surname`,`Forename`),
  KEY `WDIDX_individuals_WDIDX_individuals_WDIDX_Individuals_Optim00029` (`Surname`,`Forename`,`EntitiesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `individuals`
--

LOCK TABLES `individuals` WRITE;
/*!40000 ALTER TABLE `individuals` DISABLE KEYS */;
INSERT INTO `individuals` (`Forename`, `Surname`, `Email_Address`, `EntitiesID`, `Passw`) VALUES 
('John','Doe','test@example.com',1,'$2a$10$TPzWJbl4o7gRVHOXLtY0SOQkkXo9ZYYHTx//JvhDHD8TAWhgu98RO'),
('Anca','Mihaela Nichifor','CanadianEnglishTutoring@gmail.com',1,'$2a$10$TPzWJbl4o7gRVHOXLtY0SOQkkXo9ZYYHTx//JvhDHD8TAWhgu98RO'),
('Judy','Thompson','judythompson@hotmail.com',1,'$2a$10$TPzWJbl4o7gRVHOXLtY0SOQkkXo9ZYYHTx//JvhDHD8TAWhgu98RO');
/*!40000 ALTER TABLE `individuals` ENABLE KEYS */;
UNLOCK TABLES;
--
-- Table structure for table `la_educational_psychologists`
--

DROP TABLE IF EXISTS `la_educational_psychologists`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `la_educational_psychologists` (
  `Workforce_Educational_psychologistsID` bigint NOT NULL AUTO_INCREMENT,
  `Educational_Psychologists_Full_Time` int DEFAULT '0',
  `Educational_Psychologists_Part_Time` int DEFAULT '0',
  `Educational_Psychologists_FTE` decimal(24,6) DEFAULT '0.000000',
  PRIMARY KEY (`Workforce_Educational_psychologistsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `la_educational_psychologists`
--

LOCK TABLES `la_educational_psychologists` WRITE;
/*!40000 ALTER TABLE `la_educational_psychologists` DISABLE KEYS */;
/*!40000 ALTER TABLE `la_educational_psychologists` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `language_lookup`
--

DROP TABLE IF EXISTS `language_lookup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `language_lookup` (
  `Language_LookUpID` bigint NOT NULL AUTO_INCREMENT,
  `OriginalCode` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `OldLanguageName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `NewLanguageCode` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `NewLanguageName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `NewDialectCode` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `DialectCodeAvail` tinyint DEFAULT '0',
  `CodeChange` tinyint DEFAULT '0',
  PRIMARY KEY (`Language_LookUpID`),
  KEY `WDIDX_language_lookup_OriginalCode` (`OriginalCode`),
  KEY `WDIDX_language_lookup_OldLanguageName` (`OldLanguageName`),
  KEY `WDIDX_language_lookup_NewLanguageCode` (`NewLanguageCode`),
  KEY `WDIDX_language_lookup_NewLanguageName` (`NewLanguageName`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `language_lookup`
--

LOCK TABLES `language_lookup` WRITE;
/*!40000 ALTER TABLE `language_lookup` DISABLE KEYS */;
/*!40000 ALTER TABLE `language_lookup` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `languages`
--

DROP TABLE IF EXISTS `languages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `languages` (
  `LanguagesID` bigint NOT NULL AUTO_INCREMENT,
  `LanguageCode` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `LanguageName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `DialectCode` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`LanguagesID`),
  KEY `WDIDX_languages_LanguageCode` (`LanguageCode`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `languages`
--

LOCK TABLES `languages` WRITE;
/*!40000 ALTER TABLE `languages` DISABLE KEYS */;
/*!40000 ALTER TABLE `languages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `latest_file_io`
--

DROP TABLE IF EXISTS `latest_file_io`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `latest_file_io` (
  `thread` varchar(149) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `file` varchar(512) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `latency` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `operation` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `requested` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `latest_file_io`
--

LOCK TABLES `latest_file_io` WRITE;
/*!40000 ALTER TABLE `latest_file_io` DISABLE KEYS */;
/*!40000 ALTER TABLE `latest_file_io` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lea_clg`
--

DROP TABLE IF EXISTS `lea_clg`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lea_clg` (
  `LEA_CLGID` bigint NOT NULL AUTO_INCREMENT,
  `LEA_code` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `CLG_Code` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `CLG_name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`LEA_CLGID`),
  KEY `WDIDX_lea_clg_LEA_code` (`LEA_code`),
  KEY `WDIDX_lea_clg_CLG_Code` (`CLG_Code`),
  KEY `WDIDX_lea_clg_CLG_name` (`CLG_name`),
  KEY `WDIDX_lea_clg_WDIDX_lea_clg_WDIDX_LEA_CLG_LEA_codeCLG_Code` (`LEA_code`,`CLG_Code`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lea_clg`
--

LOCK TABLES `lea_clg` WRITE;
/*!40000 ALTER TABLE `lea_clg` DISABLE KEYS */;
/*!40000 ALTER TABLE `lea_clg` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lea_history`
--

DROP TABLE IF EXISTS `lea_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lea_history` (
  `LEA_HistoryID` bigint NOT NULL AUTO_INCREMENT,
  `New_LEA_Code` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `New_LEA_Name` varchar(150) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Old_LEA_Code` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Old_LEA_Name` varchar(150) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Year_changed` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`LEA_HistoryID`),
  KEY `WDIDX_lea_history_New_LEA_Code` (`New_LEA_Code`),
  KEY `WDIDX_lea_history_New_LEA_Name` (`New_LEA_Name`),
  KEY `WDIDX_lea_history_Old_LEA_Code` (`Old_LEA_Code`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lea_history`
--

LOCK TABLES `lea_history` WRITE;
/*!40000 ALTER TABLE `lea_history` DISABLE KEYS */;
/*!40000 ALTER TABLE `lea_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lea_names`
--

DROP TABLE IF EXISTS `lea_names`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lea_names` (
  `LEA_namesID` bigint NOT NULL AUTO_INCREMENT,
  `LEA_code` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `LEA_Name` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `LEA_Region_country` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `LEA_notes` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`LEA_namesID`),
  KEY `WDIDX_lea_names_LEA_Region_country` (`LEA_Region_country`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lea_names`
--

LOCK TABLES `lea_names` WRITE;
/*!40000 ALTER TABLE `lea_names` DISABLE KEYS */;
/*!40000 ALTER TABLE `lea_names` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lesson_file_favourites`
--

DROP TABLE IF EXISTS `lesson_file_favourites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lesson_file_favourites` (
  `Lesson_file_favouritesID` bigint NOT NULL AUTO_INCREMENT,
  `Lesson_FilesID` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  PRIMARY KEY (`Lesson_file_favouritesID`),
  KEY `WDIDX_lesson_file_favourites_Lesson_FilesID` (`Lesson_FilesID`),
  KEY `WDIDX_lesson_file_favourites_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_lesson_file_favourites_WDIDX_lesson_file_favourites_W00030` (`Lesson_FilesID`,`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lesson_file_favourites`
--

LOCK TABLES `lesson_file_favourites` WRITE;
/*!40000 ALTER TABLE `lesson_file_favourites` DISABLE KEYS */;
/*!40000 ALTER TABLE `lesson_file_favourites` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lesson_files`
--

DROP TABLE IF EXISTS `lesson_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lesson_files` (
  `Lesson_FilesID` bigint NOT NULL AUTO_INCREMENT,
  `Filename` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Description` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AudioFile` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `File_TypesID` bigint DEFAULT '0',
  `mimetypesID` bigint DEFAULT '0',
  `SkillsID` bigint DEFAULT '0',
  `timeReq` int DEFAULT '0',
  `LocalFoldersID` bigint DEFAULT '0',
  `Hidden` tinyint NOT NULL DEFAULT '0',
  `CloudID` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Lesson_file_favouritesID` bigint DEFAULT NULL,
  `MarksAvail` smallint DEFAULT '0',
  PRIMARY KEY (`Lesson_FilesID`),
  KEY `WDIDX_lesson_files_filename` (`Filename`),
  KEY `WDIDX_lesson_files_File_TypesID` (`File_TypesID`),
  KEY `WDIDX_lesson_files_mimetypesID` (`mimetypesID`),
  KEY `WDIDX_lesson_files_LocalFoldersID` (`LocalFoldersID`),
  KEY `WDIDX_lesson_files_hidden` (`Hidden`),
  KEY `WDIDX_lesson_files_CloudID` (`CloudID`),
  KEY `WDIDX_lesson_files_Lesson_file_favouritesID` (`Lesson_file_favouritesID`),
  KEY `WDIDX_lesson_files_WDIDX_lesson_files_WDIDX_Lesson_Files_Le00031` (`Lesson_FilesID`,`LocalFoldersID`),
  KEY `WDIDX_lesson_files_WDIDX_lesson_files_WDIDX_Lesson_Files_fi00032` (`Filename`,`LocalFoldersID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lesson_files`
--

LOCK TABLES `lesson_files` WRITE;
/*!40000 ALTER TABLE `lesson_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `lesson_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lesson_memos`
--

DROP TABLE IF EXISTS `lesson_memos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lesson_memos` (
  `Lesson_memosID` bigint NOT NULL AUTO_INCREMENT,
  `memo` varchar(250) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Lesson_FilesID` bigint DEFAULT '0',
  `Lesson_outlineID` bigint DEFAULT '0',
  PRIMARY KEY (`Lesson_memosID`),
  KEY `WDIDX_lesson_memos_Lesson_FilesID` (`Lesson_FilesID`),
  KEY `WDIDX_lesson_memos_Lesson_outlineID` (`Lesson_outlineID`),
  KEY `WDIDX_lesson_memos_WDIDX_lesson_memos_Lesson_FilesIDLesson_00033` (`Lesson_FilesID`,`Lesson_outlineID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lesson_memos`
--

LOCK TABLES `lesson_memos` WRITE;
/*!40000 ALTER TABLE `lesson_memos` DISABLE KEYS */;
/*!40000 ALTER TABLE `lesson_memos` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lesson_objectives`
--

DROP TABLE IF EXISTS `lesson_objectives`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lesson_objectives` (
  `lesson_objectivesID` bigint NOT NULL AUTO_INCREMENT,
  `LessonID` bigint DEFAULT '0',
  `objective` varchar(256) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`lesson_objectivesID`),
  KEY `WDIDX_lesson_objectives_LessonID` (`LessonID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lesson_objectives`
--

LOCK TABLES `lesson_objectives` WRITE;
/*!40000 ALTER TABLE `lesson_objectives` DISABLE KEYS */;
/*!40000 ALTER TABLE `lesson_objectives` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lesson_outline`
--

DROP TABLE IF EXISTS `lesson_outline`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lesson_outline` (
  `Lesson_outlineID` bigint NOT NULL AUTO_INCREMENT,
  `LessonID` bigint DEFAULT '0',
  `Timetable_entriesID` bigint DEFAULT '0',
  `SessionsID` bigint DEFAULT '0',
  `Topic_plan_headerID` bigint DEFAULT '0',
  `Lesson_PlansID` bigint DEFAULT NULL,
  PRIMARY KEY (`Lesson_outlineID`),
  KEY `WDIDX_lesson_outline_LessonID` (`LessonID`),
  KEY `WDIDX_lesson_outline_timetable_entriesID` (`Timetable_entriesID`),
  KEY `WDIDX_lesson_outline_SessionsID` (`SessionsID`),
  KEY `WDIDX_lesson_outline_topic_plan_headerID` (`Topic_plan_headerID`),
  KEY `WDIDX_lesson_outline_Lesson_PlansID` (`Lesson_PlansID`),
  KEY `WDIDX_lesson_outline_WDIDX_lesson_outline_WDIDX_Lesson_outl00034` (`LessonID`,`Timetable_entriesID`,`SessionsID`),
  KEY `WDIDX_lesson_outline_WDIDX_lesson_outline_WDIDX_Lesson_outl00035` (`Timetable_entriesID`,`SessionsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lesson_outline`
--

LOCK TABLES `lesson_outline` WRITE;
/*!40000 ALTER TABLE `lesson_outline` DISABLE KEYS */;
/*!40000 ALTER TABLE `lesson_outline` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lesson_plans`
--

DROP TABLE IF EXISTS `lesson_plans`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lesson_plans` (
  `Lesson_PlansID` bigint NOT NULL AUTO_INCREMENT,
  `Lesson_outlineID` bigint DEFAULT '0',
  `Lesson_FilesID` bigint DEFAULT '0',
  `LessonOrder` int DEFAULT '0',
  `PrintReq` int DEFAULT '0',
  `LessonStarterQuestionsID` bigint DEFAULT '0',
  `Lesson_memosID` bigint DEFAULT '0',
  PRIMARY KEY (`Lesson_PlansID`),
  KEY `WDIDX_lesson_plans_Lesson_outlineID` (`Lesson_outlineID`),
  KEY `WDIDX_lesson_plans_Lesson_FilesID` (`Lesson_FilesID`),
  KEY `WDIDX_lesson_plans_LessonOrder` (`LessonOrder`),
  KEY `WDIDX_lesson_plans_Lesson_memosID` (`Lesson_memosID`),
  KEY `WDIDX_lesson_plans_WDIDX_lesson_plans_WDIDX_Lesson_Plans_Le00036` (`Lesson_outlineID`,`Lesson_FilesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lesson_plans`
--

LOCK TABLES `lesson_plans` WRITE;
/*!40000 ALTER TABLE `lesson_plans` DISABLE KEYS */;
/*!40000 ALTER TABLE `lesson_plans` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lessonquestion`
--

DROP TABLE IF EXISTS `lessonquestion`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lessonquestion` (
  `LessonQuestionID` bigint NOT NULL AUTO_INCREMENT,
  `Question` varchar(400) DEFAULT NULL,
  `MarksAvailable` int DEFAULT '0',
  `Answer` varchar(400) DEFAULT NULL,
  `Lesson_FilesID` bigint DEFAULT '0',
  PRIMARY KEY (`LessonQuestionID`),
  KEY `WDIDX_LessonQuestion_Lesson_FilesID` (`Lesson_FilesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lessonquestion`
--

LOCK TABLES `lessonquestion` WRITE;
/*!40000 ALTER TABLE `lessonquestion` DISABLE KEYS */;
/*!40000 ALTER TABLE `lessonquestion` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lessons`
--

DROP TABLE IF EXISTS `lessons`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lessons` (
  `LessonID` bigint NOT NULL AUTO_INCREMENT,
  `lessonName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `TopicID` bigint DEFAULT '0',
  `lessonCode` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `CloudFolderID` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `Summary` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Hidden` int DEFAULT '0',
  PRIMARY KEY (`LessonID`),
  KEY `WDIDX_lessons_lessonName` (`lessonName`),
  KEY `WDIDX_lessons_TopicsID` (`TopicID`),
  KEY `WDIDX_lessons_CloudFolderID` (`CloudFolderID`),
  KEY `WDIDX_lessons_hidden` (`Hidden`),
  KEY `WDIDX_lessons_WDIDX_lessons_WDIDX_lessons_TopicsIDName` (`TopicID`,`lessonName`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lessons`
--

LOCK TABLES `lessons` WRITE;
/*!40000 ALTER TABLE `lessons` DISABLE KEYS */;
/*!40000 ALTER TABLE `lessons` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lessonstarterquestions`
--

DROP TABLE IF EXISTS `lessonstarterquestions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lessonstarterquestions` (
  `LessonStarterQuestionsID` bigint NOT NULL AUTO_INCREMENT,
  `LessonID` bigint DEFAULT '0',
  `StartQuestion` varchar(400) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Answer` varchar(400) DEFAULT NULL,
  `MarksAvail` smallint DEFAULT '0',
  `Skill_linkID` bigint DEFAULT '0',
  PRIMARY KEY (`LessonStarterQuestionsID`),
  KEY `WDIDX_lessonstarterquestions_LessonID` (`LessonID`),
  KEY `WDIDX_lessonstarterquestions_Skill_linkID` (`Skill_linkID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lessonstarterquestions`
--

LOCK TABLES `lessonstarterquestions` WRITE;
/*!40000 ALTER TABLE `lessonstarterquestions` DISABLE KEYS */;
/*!40000 ALTER TABLE `lessonstarterquestions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `linked_emails`
--

DROP TABLE IF EXISTS `linked_emails`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `linked_emails` (
  `id_linked_emails` int NOT NULL AUTO_INCREMENT,
  `graph_id` varchar(128) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`id_linked_emails`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `linked_emails`
--

LOCK TABLES `linked_emails` WRITE;
/*!40000 ALTER TABLE `linked_emails` DISABLE KEYS */;
/*!40000 ALTER TABLE `linked_emails` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `linked_files`
--

DROP TABLE IF EXISTS `linked_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `linked_files` (
  `Linked_filesID` bigint NOT NULL AUTO_INCREMENT,
  `FileA` bigint DEFAULT '0',
  `FileGroupID` bigint DEFAULT '0',
  `LessonID` bigint DEFAULT '0',
  PRIMARY KEY (`Linked_filesID`),
  KEY `WDIDX_linked_files_FileA` (`FileA`),
  KEY `WDIDX_linked_files_FileGroupID` (`FileGroupID`),
  KEY `WDIDX_linked_files_LessonID` (`LessonID`),
  KEY `WDIDX_linked_files_WDIDX_linked_files_WDIDX_Linked_files_Fi00037` (`FileA`,`FileGroupID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `linked_files`
--

LOCK TABLES `linked_files` WRITE;
/*!40000 ALTER TABLE `linked_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `linked_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `localfolders`
--

DROP TABLE IF EXISTS `localfolders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `localfolders` (
  `LocalFoldersID` bigint NOT NULL AUTO_INCREMENT,
  `FolderName` varchar(80) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `ParentID` bigint DEFAULT '0',
  `SubjectID` bigint DEFAULT '0',
  `Hidden` tinyint DEFAULT '0',
  `CloudID` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `TopicID` bigint DEFAULT NULL,
  `CourseID` bigint DEFAULT NULL,
  `Lesson_FilesID` bigint DEFAULT NULL,
  PRIMARY KEY (`LocalFoldersID`),
  KEY `WDIDX_localfolders_FolderName` (`FolderName`),
  KEY `WDIDX_localfolders_ParentID` (`ParentID`),
  KEY `WDIDX_localfolders_SubjectID` (`SubjectID`),
  KEY `WDIDX_localfolders_hidden` (`Hidden`),
  KEY `WDIDX_localfolders_CloudID` (`CloudID`),
  KEY `WDIDX_localfolders_TopicsID` (`TopicID`),
  KEY `WDIDX_localfolders_CourseID` (`CourseID`),
  KEY `WDIDX_localfolders_Lesson_FilesID` (`Lesson_FilesID`),
  KEY `WDIDX_localfolders_WDIDX_localfolders_WDIDX_LocalFolders_Fo00038` (`FolderName`,`ParentID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `localfolders`
--

LOCK TABLES `localfolders` WRITE;
/*!40000 ALTER TABLE `localfolders` DISABLE KEYS */;
/*!40000 ALTER TABLE `localfolders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `major_incident_involved`
--

DROP TABLE IF EXISTS `major_incident_involved`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `major_incident_involved` (
  `Major_Incident_InvolvedID` bigint NOT NULL AUTO_INCREMENT,
  `Major_transgression_logID` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  `Witness` tinyint DEFAULT '0',
  `Known` tinyint DEFAULT '0',
  UNIQUE KEY `Major_Incident_InvolvedID` (`Major_Incident_InvolvedID`),
  KEY `WDIDX_major_incident_involved_Major_transgression_logID` (`Major_transgression_logID`),
  KEY `WDIDX_major_incident_involved_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_major_incident_involved_Witness` (`Witness`),
  KEY `WDIDX_major_incident_involved_WDIDX_Major_Incident_Involved00039` (`Major_transgression_logID`,`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `major_incident_involved`
--

LOCK TABLES `major_incident_involved` WRITE;
/*!40000 ALTER TABLE `major_incident_involved` DISABLE KEYS */;
/*!40000 ALTER TABLE `major_incident_involved` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `major_incident_statement`
--

DROP TABLE IF EXISTS `major_incident_statement`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `major_incident_statement` (
  `Major_incident_statementID` bigint NOT NULL AUTO_INCREMENT,
  `Major_transgression_logID` bigint DEFAULT '0',
  `Statement` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `IndividualsID` bigint DEFAULT '0',
  `RecordedBy` bigint DEFAULT '0',
  `TimeRec` timestamp NULL DEFAULT NULL,
  `Known` tinyint DEFAULT '0',
  UNIQUE KEY `Major_incident_statementID` (`Major_incident_statementID`),
  KEY `WDIDX_major_incident_statement_Major_transgression_logID` (`Major_transgression_logID`),
  KEY `WDIDX_major_incident_statement_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_major_incident_statement_RecordedBy` (`RecordedBy`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `major_incident_statement`
--

LOCK TABLES `major_incident_statement` WRITE;
/*!40000 ALTER TABLE `major_incident_statement` DISABLE KEYS */;
/*!40000 ALTER TABLE `major_incident_statement` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `major_transgression_log`
--

DROP TABLE IF EXISTS `major_transgression_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `major_transgression_log` (
  `Major_transgression_logID` bigint NOT NULL AUTO_INCREMENT,
  `DateAndTime` timestamp NULL DEFAULT NULL,
  `Major_transgression_typesID` bigint DEFAULT '0',
  `Details` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `RecordedBy` bigint DEFAULT '0',
  `Resolved` tinyint DEFAULT '0',
  `RoomsID` bigint DEFAULT '0',
  UNIQUE KEY `Major_transgression_logID` (`Major_transgression_logID`),
  KEY `WDIDX_major_transgression_log_DateAndTime` (`DateAndTime`),
  KEY `WDIDX_major_transgression_log_Major_transgression_typesID` (`Major_transgression_typesID`),
  KEY `WDIDX_major_transgression_log_RecordedBy` (`RecordedBy`),
  KEY `WDIDX_major_transgression_log_Resolved` (`Resolved`),
  KEY `WDIDX_major_transgression_log_roomsID` (`RoomsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `major_transgression_log`
--

LOCK TABLES `major_transgression_log` WRITE;
/*!40000 ALTER TABLE `major_transgression_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `major_transgression_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `major_transgression_types`
--

DROP TABLE IF EXISTS `major_transgression_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `major_transgression_types` (
  `Major_transgression_typesID` bigint NOT NULL AUTO_INCREMENT,
  `Major_transgression_name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Sanction_typesID` bigint DEFAULT '0',
  UNIQUE KEY `Major_transgression_typesID` (`Major_transgression_typesID`),
  UNIQUE KEY `Major_transgression_name` (`Major_transgression_name`),
  KEY `WDIDX_major_transgression_types_Sanction_typesID` (`Sanction_typesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `major_transgression_types`
--

LOCK TABLES `major_transgression_types` WRITE;
/*!40000 ALTER TABLE `major_transgression_types` DISABLE KEYS */;
/*!40000 ALTER TABLE `major_transgression_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `markbook`
--

DROP TABLE IF EXISTS `markbook`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `markbook` (
  `MarkbookID` bigint NOT NULL AUTO_INCREMENT,
  `Mark` int DEFAULT '0',
  `Grade` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `WorkAssignmentID` bigint DEFAULT '0',
  PRIMARY KEY (`MarkbookID`),
  KEY `WDIDX_markbook_WorkAssignmentID` (`WorkAssignmentID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `markbook` WRITE;
/*!40000 ALTER TABLE `markbook` DISABLE KEYS */;
/*!40000 ALTER TABLE `markbook` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `markbookmarks`
--

DROP TABLE IF EXISTS `markbookmarks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `markbookmarks` (
  `MarkBookMarksID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `MarkSchemeAnswersID` bigint DEFAULT '0',
  `Mark` int DEFAULT '0',
  `MarkbookID` bigint DEFAULT NULL,
  PRIMARY KEY (`MarkBookMarksID`),
  KEY `WDIDX_markbookmarks_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_markbookmarks_MarkSchemeAnswersID` (`MarkSchemeAnswersID`),
  KEY `WDIDX_markbookmarks_Mark` (`Mark`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `markbookmarks`
--

LOCK TABLES `markbookmarks` WRITE;
/*!40000 ALTER TABLE `markbookmarks` DISABLE KEYS */;
/*!40000 ALTER TABLE `markbookmarks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `markschemeanswers`
--

DROP TABLE IF EXISTS `markschemeanswers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `markschemeanswers` (
  `MarkSchemeAnswersID` bigint NOT NULL AUTO_INCREMENT,
  `QuestionNo` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AnswerListed` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Skills_keyID` bigint DEFAULT '0',
  `SubTopic` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `MarkschemeHeaderID` bigint DEFAULT '0',
  `MarksAvail` smallint DEFAULT '0',
  UNIQUE KEY `MarkSchemeAnswersID` (`MarkSchemeAnswersID`),
  KEY `WDIDX_markschemeanswers_QuestionNo` (`QuestionNo`),
  KEY `WDIDX_markschemeanswers_Skills_keyID` (`Skills_keyID`),
  KEY `WDIDX_markschemeanswers_SubTopic` (`SubTopic`),
  KEY `WDIDX_markschemeanswers_MarkschemeHeaderID` (`MarkschemeHeaderID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `markschemeanswers`
--

LOCK TABLES `markschemeanswers` WRITE;
/*!40000 ALTER TABLE `markschemeanswers` DISABLE KEYS */;
/*!40000 ALTER TABLE `markschemeanswers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `markschemeheader`
--

DROP TABLE IF EXISTS `markschemeheader`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `markschemeheader` (
  `MarkschemeHeaderID` bigint NOT NULL AUTO_INCREMENT,
  `Lesson_FilesID` bigint DEFAULT '0',
  `MarksAvail` smallint DEFAULT '0',
  `SchemeFileID` bigint DEFAULT '0',
  PRIMARY KEY (`MarkschemeHeaderID`),
  KEY `WDIDX_markschemeheader_Lesson_FilesID` (`Lesson_FilesID`),
  KEY `WDIDX_markschemeheader_schemeFileID` (`SchemeFileID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `markschemeheader`
--

LOCK TABLES `markschemeheader` WRITE;
/*!40000 ALTER TABLE `markschemeheader` DISABLE KEYS */;
/*!40000 ALTER TABLE `markschemeheader` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `marktypes`
--

DROP TABLE IF EXISTS `marktypes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `marktypes` (
  `MarkTypesID` bigint NOT NULL AUTO_INCREMENT,
  `MarkTypeName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`MarkTypesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `marktypes`
--

LOCK TABLES `marktypes` WRITE;
/*!40000 ALTER TABLE `marktypes` DISABLE KEYS */;
/*!40000 ALTER TABLE `marktypes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `memo`
--

DROP TABLE IF EXISTS `memo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `memo` (
  `MemoID` bigint NOT NULL AUTO_INCREMENT,
  `Memo_content` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`MemoID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `memo`
--

LOCK TABLES `memo` WRITE;
/*!40000 ALTER TABLE `memo` DISABLE KEYS */;
/*!40000 ALTER TABLE `memo` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `memopairings`
--

DROP TABLE IF EXISTS `memopairings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `memopairings` (
  `MemoPairingsID` bigint NOT NULL AUTO_INCREMENT,
  `StandardizedMemoID` bigint DEFAULT '0',
  `StandardizeReplacementsID` bigint DEFAULT '0',
  PRIMARY KEY (`MemoPairingsID`),
  KEY `WDIDX_MemoPairings_StandardizedMemoID` (`StandardizedMemoID`),
  KEY `WDIDX_MemoPairings_StandardizeReplacementsID` (`StandardizeReplacementsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `memopairings`
--

LOCK TABLES `memopairings` WRITE;
/*!40000 ALTER TABLE `memopairings` DISABLE KEYS */;
/*!40000 ALTER TABLE `memopairings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `memory_by_host_by_current_bytes`
--

DROP TABLE IF EXISTS `memory_by_host_by_current_bytes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `memory_by_host_by_current_bytes` (
  `host` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `current_count_used` double DEFAULT NULL,
  `current_allocated` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `current_avg_alloc` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `current_max_alloc` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `total_allocated` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `memory_by_host_by_current_bytes`
--

LOCK TABLES `memory_by_host_by_current_bytes` WRITE;
/*!40000 ALTER TABLE `memory_by_host_by_current_bytes` DISABLE KEYS */;
/*!40000 ALTER TABLE `memory_by_host_by_current_bytes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `memory_global_total`
--

DROP TABLE IF EXISTS `memory_global_total`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `memory_global_total` (
  `total_allocated` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `memory_global_total`
--

LOCK TABLES `memory_global_total` WRITE;
/*!40000 ALTER TABLE `memory_global_total` DISABLE KEYS */;
/*!40000 ALTER TABLE `memory_global_total` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `metrics`
--

DROP TABLE IF EXISTS `metrics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `metrics` (
  `Variable_name` varchar(193) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Variable_value` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `Type` varchar(210) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Enabled` varchar(7) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `metrics`
--

LOCK TABLES `metrics` WRITE;
/*!40000 ALTER TABLE `metrics` DISABLE KEYS */;
/*!40000 ALTER TABLE `metrics` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mimetypes`
--

DROP TABLE IF EXISTS `mimetypes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `mimetypes` (
  `IDmimeTypes` bigint NOT NULL AUTO_INCREMENT,
  `TypeName` varchar(80) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`IDmimeTypes`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mimetypes`
--

LOCK TABLES `mimetypes` WRITE;
/*!40000 ALTER TABLE `mimetypes` DISABLE KEYS */;
/*!40000 ALTER TABLE `mimetypes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `minor_trangression_targets`
--

DROP TABLE IF EXISTS `minor_trangression_targets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `minor_trangression_targets` (
  `Minor_trangression_targetsID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `NumbOfSess` tinyint DEFAULT '0',
  `SanctID` bigint DEFAULT '0',
  `TLimit` int DEFAULT '0',
  UNIQUE KEY `Minor_trangression_targetsID` (`Minor_trangression_targetsID`),
  KEY `WDIDX_minor_trangression_targets_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `minor_trangression_targets`
--

LOCK TABLES `minor_trangression_targets` WRITE;
/*!40000 ALTER TABLE `minor_trangression_targets` DISABLE KEYS */;
/*!40000 ALTER TABLE `minor_trangression_targets` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `minor_transgression_log`
--

DROP TABLE IF EXISTS `minor_transgression_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `minor_transgression_log` (
  `Minor_transgression_logID` bigint NOT NULL AUTO_INCREMENT,
  `Minor_Transgression_typesID` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  `MinorPoints` tinyint DEFAULT '0',
  `Sanction_setID` bigint DEFAULT '0',
  `Lesson_outlineID` bigint DEFAULT '0',
  UNIQUE KEY `Minor_transgression_logID` (`Minor_transgression_logID`),
  KEY `WDIDX_minor_transgression_log_Minor_Transgression_typesID` (`Minor_Transgression_typesID`),
  KEY `WDIDX_minor_transgression_log_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_minor_transgression_log_Sanction_setID` (`Sanction_setID`),
  KEY `WDIDX_minor_transgression_log_Lesson_outlineID` (`Lesson_outlineID`),
  KEY `WDIDX_minor_transgression_log_WDIDX_Minor_transgression_log00040` (`IndividualsID`,`Lesson_outlineID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `minor_transgression_log`
--

LOCK TABLES `minor_transgression_log` WRITE;
/*!40000 ALTER TABLE `minor_transgression_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `minor_transgression_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `minor_transgression_types`
--

DROP TABLE IF EXISTS `minor_transgression_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `minor_transgression_types` (
  `Minor_Transgression_typesID` bigint NOT NULL AUTO_INCREMENT,
  `Minor_transgression_name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Transgression_points` int DEFAULT '0',
  `Default` varchar(1) DEFAULT NULL,
  UNIQUE KEY `Minor_Transgression_typesID` (`Minor_Transgression_typesID`),
  UNIQUE KEY `Minor_transgression_name` (`Minor_transgression_name`),
  KEY `WDIDX_minor_transgression_types_Transgression_points` (`Transgression_points`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `minor_transgression_types`
--

LOCK TABLES `minor_transgression_types` WRITE;
/*!40000 ALTER TABLE `minor_transgression_types` DISABLE KEYS */;
/*!40000 ALTER TABLE `minor_transgression_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mitigation`
--

DROP TABLE IF EXISTS `mitigation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `mitigation` (
  `MitigationID` bigint NOT NULL AUTO_INCREMENT,
  `Sanction_setID` bigint DEFAULT '0',
  `Text` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  UNIQUE KEY `MitigationID` (`MitigationID`),
  KEY `WDIDX_mitigation_Sanction_setID` (`Sanction_setID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mitigation`
--

LOCK TABLES `mitigation` WRITE;
/*!40000 ALTER TABLE `mitigation` DISABLE KEYS */;
/*!40000 ALTER TABLE `mitigation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `nation_states_and_countries`
--

DROP TABLE IF EXISTS `nation_states_and_countries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nation_states_and_countries` (
  `Nation_States_and_CountriesID` bigint NOT NULL AUTO_INCREMENT,
  `ISO_3166_1_Alpha_3` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Nation_Short_Name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Nation_Long_Name` varchar(150) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `ISO_3166_1_Alpha_2` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `ISO_3166_1_Code_Numeric` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Current` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Notes` varchar(300) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Nation_States_and_CountriesID`),
  KEY `WDIDX_nation_states_and_countries_ISO_3166_1_Alpha_2` (`ISO_3166_1_Alpha_2`),
  KEY `WDIDX_nation_states_and_countries_ISO_3166_1_Code_Numeric` (`ISO_3166_1_Code_Numeric`),
  KEY `WDIDX_nation_states_and_countries_Current` (`Current`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `nation_states_and_countries`
--

LOCK TABLES `nation_states_and_countries` WRITE;
/*!40000 ALTER TABLE `nation_states_and_countries` DISABLE KEYS */;
/*!40000 ALTER TABLE `nation_states_and_countries` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `newmembers`
--

DROP TABLE IF EXISTS `newmembers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `newmembers` (
  `NewMembersID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `SuppliedCode` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`NewMembersID`),
  KEY `WDIDX_newmembers_SuppliedCode` (`SuppliedCode`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `newmembers`
--

LOCK TABLES `newmembers` WRITE;
/*!40000 ALTER TABLE `newmembers` DISABLE KEYS */;
/*!40000 ALTER TABLE `newmembers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `nodes`
--

DROP TABLE IF EXISTS `nodes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nodes` (
  `NodesID` bigint NOT NULL AUTO_INCREMENT,
  `Node1` bigint DEFAULT '0',
  `ElementID1` bigint DEFAULT '0',
  `Node2` bigint DEFAULT '0',
  `ElementID2` bigint DEFAULT '0',
  PRIMARY KEY (`NodesID`),
  KEY `WDIDX_nodes_Node1` (`Node1`),
  KEY `WDIDX_nodes_ElementID1` (`ElementID1`),
  KEY `WDIDX_nodes_Node2` (`Node2`),
  KEY `WDIDX_nodes_ElementID2` (`ElementID2`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `nodes`
--

LOCK TABLES `nodes` WRITE;
/*!40000 ALTER TABLE `nodes` DISABLE KEYS */;
/*!40000 ALTER TABLE `nodes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `nodetable`
--

DROP TABLE IF EXISTS `nodetable`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nodetable` (
  `NodeTableID` bigint NOT NULL AUTO_INCREMENT,
  `NodeName` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`NodeTableID`),
  KEY `WDIDX_nodetable_NodeName` (`NodeName`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `nodetable`
--

LOCK TABLES `nodetable` WRITE;
/*!40000 ALTER TABLE `nodetable` DISABLE KEYS */;
/*!40000 ALTER TABLE `nodetable` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notices`
--

DROP TABLE IF EXISTS `notices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notices` (
  `NoticesID` bigint NOT NULL AUTO_INCREMENT,
  `MemoID` bigint DEFAULT '0',
  `StandardizedMemoID` bigint DEFAULT '0',
  `StudentID` bigint DEFAULT '0',
  `CarerID` bigint DEFAULT '0',
  `Acknowledged` bigint DEFAULT '0',
  PRIMARY KEY (`NoticesID`),
  KEY `WDIDX_Notices_memoID` (`MemoID`),
  KEY `WDIDX_Notices_StandardizedMemoID` (`StandardizedMemoID`),
  KEY `WDIDX_Notices_StudentID` (`StudentID`),
  KEY `WDIDX_Notices_CarerID` (`CarerID`),
  KEY `WDIDX_Notices_Acknowledged` (`Acknowledged`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notices`
--

LOCK TABLES `notices` WRITE;
/*!40000 ALTER TABLE `notices` DISABLE KEYS */;
/*!40000 ALTER TABLE `notices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `off_timetable_event`
--

DROP TABLE IF EXISTS `off_timetable_event`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `off_timetable_event` (
  `Off_timetable_eventID` bigint NOT NULL AUTO_INCREMENT,
  `StartDate` bigint NOT NULL DEFAULT '0',
  `EndTime` time DEFAULT NULL,
  `EndDate` bigint NOT NULL DEFAULT '0',
  `Name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AppliesTo` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Room_bookingsID` bigint NOT NULL DEFAULT '0',
  `StartTime` time DEFAULT NULL,
  PRIMARY KEY (`Off_timetable_eventID`),
  KEY `WDIDX_off_timetable_event_StartDate` (`StartDate`),
  KEY `WDIDX_off_timetable_event_EndTime` (`EndTime`),
  KEY `WDIDX_off_timetable_event_EndDate` (`EndDate`),
  KEY `WDIDX_off_timetable_event_Name` (`Name`),
  KEY `WDIDX_off_timetable_event_appliesTo` (`AppliesTo`),
  KEY `WDIDX_off_timetable_event_Room_bookingsID` (`Room_bookingsID`),
  KEY `WDIDX_off_timetable_event_StartTime` (`StartTime`),
  KEY `WDIDX_off_timetable_event_WDIDX_off_timetable_event_WDIDX_O00041` (`Name`,`StartDate`,`AppliesTo`),
  KEY `WDIDX_off_timetable_event_WDIDX_Off_timetable_event_NameSta00042` (`Name`,`StartDate`,`AppliesTo`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `off_timetable_event`
--

LOCK TABLES `off_timetable_event` WRITE;
/*!40000 ALTER TABLE `off_timetable_event` DISABLE KEYS */;
/*!40000 ALTER TABLE `off_timetable_event` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `off_timetable_sessions`
--

DROP TABLE IF EXISTS `off_timetable_sessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `off_timetable_sessions` (
  `Off_Timetable_sessionsID` bigint NOT NULL AUTO_INCREMENT,
  `Off_timetable_eventID` bigint DEFAULT '0',
  `SessionsID` bigint DEFAULT '0',
  PRIMARY KEY (`Off_Timetable_sessionsID`),
  KEY `WDIDX_off_timetable_sessions_Off_timetable_eventID` (`Off_timetable_eventID`),
  KEY `WDIDX_off_timetable_sessions_SessionsID` (`SessionsID`),
  KEY `WDIDX_off_timetable_sessions_WDIDX_off_timetable_sessions_W00043` (`SessionsID`,`Off_timetable_eventID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `off_timetable_sessions`
--

LOCK TABLES `off_timetable_sessions` WRITE;
/*!40000 ALTER TABLE `off_timetable_sessions` DISABLE KEYS */;
/*!40000 ALTER TABLE `off_timetable_sessions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `onedrive`
--

DROP TABLE IF EXISTS `onedrive`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `onedrive` (
  `Accesstoken` varchar(5000) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Refreshtoken` varchar(1100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `IndividualsFile` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `PhotosFile` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `GroupsFile` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `RegisterFile` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `TimetableFile` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `SetupFile` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `DateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Exercises` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `SchemeFile` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AppFolderID` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `PorB` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `IndividualsID` bigint DEFAULT NULL,
  UNIQUE KEY `IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `onedrive`
--

LOCK TABLES `onedrive` WRITE;
/*!40000 ALTER TABLE `onedrive` DISABLE KEYS */;
/*!40000 ALTER TABLE `onedrive` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `onetoonetuition`
--

DROP TABLE IF EXISTS `onetoonetuition`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `onetoonetuition` (
  `OneToOneTuitionID` bigint NOT NULL AUTO_INCREMENT,
  `Tuition_Attainment_Target_Focus` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `IndividualsID` bigint DEFAULT '0',
  `Tuition_Start_Date` bigint NOT NULL DEFAULT '0',
  `Tuition_Completion_Status` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Tuition_Completed_Hours` tinyint DEFAULT '0',
  PRIMARY KEY (`OneToOneTuitionID`),
  KEY `WDIDX_onetoonetuition_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `onetoonetuition`
--

LOCK TABLES `onetoonetuition` WRITE;
/*!40000 ALTER TABLE `onetoonetuition` DISABLE KEYS */;
/*!40000 ALTER TABLE `onetoonetuition` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `organisation_common_elements`
--

DROP TABLE IF EXISTS `organisation_common_elements`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `organisation_common_elements` (
  `Organisation_Common_ElementsID` bigint NOT NULL AUTO_INCREMENT,
  `Email_Address` varchar(260) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Sub_dwelling` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Dwelling` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Street` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Locality` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Town` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Administrative_Area` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Post_town` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `PostCode` varchar(8) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Address_Line_1` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Address_Line_2` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Address_Line_3` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Address_Line_4` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Address_Line_5` varchar(40) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Body_name` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Telephone` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `National_Curriculum_Year_Group` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Change_to_proprietor_indicator` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Organisation_Common_ElementsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `organisation_common_elements`
--

LOCK TABLES `organisation_common_elements` WRITE;
/*!40000 ALTER TABLE `organisation_common_elements` DISABLE KEYS */;
/*!40000 ALTER TABLE `organisation_common_elements` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pospoints`
--

DROP TABLE IF EXISTS `pospoints`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pospoints` (
  `PosPointsID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `Lesson_outlineID` bigint DEFAULT '0',
  `Points` tinyint DEFAULT '0',
  UNIQUE KEY `PosPointsID` (`PosPointsID`),
  KEY `WDIDX_pospoints_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_pospoints_Lesson_outlineID` (`Lesson_outlineID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pospoints`
--

LOCK TABLES `pospoints` WRITE;
/*!40000 ALTER TABLE `pospoints` DISABLE KEYS */;
/*!40000 ALTER TABLE `pospoints` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `post16learningaims`
--

DROP TABLE IF EXISTS `post16learningaims`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post16learningaims` (
  `Post16LearningAimsID` bigint NOT NULL AUTO_INCREMENT,
  `Learning_Aim_Start_Date` bigint DEFAULT NULL,
  `Learning_Aim_Planned_End_Date` bigint DEFAULT NULL,
  `Learning_Aim_Actual_End_Date` bigint DEFAULT NULL,
  `Learning_Aim_Status` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `IndividualsID` bigint DEFAULT '0',
  `Core_aim_indicator` tinyint DEFAULT '0',
  `Learning_Aim_Withdrawal_Reason` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Traineeship` tinyint DEFAULT '0',
  UNIQUE KEY `Post16LearningAimsID` (`Post16LearningAimsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post16learningaims`
--

LOCK TABLES `post16learningaims` WRITE;
/*!40000 ALTER TABLE `post16learningaims` DISABLE KEYS */;
/*!40000 ALTER TABLE `post16learningaims` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `previousschools`
--

DROP TABLE IF EXISTS `previousschools`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `previousschools` (
  `PreviousSchoolsID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `Previous_School_LA_number` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Previous_School_Name` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Previous_School_DfE_Establishment_Number` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pupil_StatusID` bigint DEFAULT '0',
  PRIMARY KEY (`PreviousSchoolsID`),
  KEY `WDIDX_previousschools_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `previousschools`
--

LOCK TABLES `previousschools` WRITE;
/*!40000 ALTER TABLE `previousschools` DISABLE KEYS */;
/*!40000 ALTER TABLE `previousschools` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ps_check_lost_instrumentation`
--

DROP TABLE IF EXISTS `ps_check_lost_instrumentation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ps_check_lost_instrumentation` (
  `variable_name` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `variable_value` varchar(1024) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ps_check_lost_instrumentation`
--

LOCK TABLES `ps_check_lost_instrumentation` WRITE;
/*!40000 ALTER TABLE `ps_check_lost_instrumentation` DISABLE KEYS */;
/*!40000 ALTER TABLE `ps_check_lost_instrumentation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pupil_admissions`
--

DROP TABLE IF EXISTS `pupil_admissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pupil_admissions` (
  `Pupil_AdmissionsID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `Admissions_entry_year_and_month` varchar(7) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Admissions_entry_year_group` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `ADT_File_status` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Text_file_header` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `Statemented` tinyint DEFAULT '0',
  `Verified_address` tinyint DEFAULT '0',
  `Preference_Rank` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '0',
  `Preference_reason` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '0',
  `Preference_reason_text` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `Aptitude_code` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Aptitude_text` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `Faith_text` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `Distance_from_school` int DEFAULT '0',
  `LA_pre_banding_information` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `LA_other_text` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `Local_reference` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Medical_or_social_information` tinyint DEFAULT '0',
  `Criterion` varchar(25) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Priority` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Admission_policy_criterion` varchar(25) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Notification_method` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Ranking_of_offers` varchar(25) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Date_Place_Required` bigint NOT NULL DEFAULT '0',
  `In_Year_Text` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `Date_Last_Attended` bigint NOT NULL DEFAULT '0',
  `Accept_Text` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `Offer_status` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Application_reference` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `In_Year` tinyint DEFAULT '0',
  `Council_Tax_Reference` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `LA_In_Year_Text` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `Faith` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Supplementary_Parental_Offer_Response` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Admissions_In_Care` tinyint DEFAULT '0',
  `Admissions_Care_Authority` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Alternative_Contact_Notes` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Pupil_AdmissionsID`),
  KEY `WDIDX_pupil_admissions_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pupil_admissions`
--

LOCK TABLES `pupil_admissions` WRITE;
/*!40000 ALTER TABLE `pupil_admissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `pupil_admissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pupil_assessment`
--

DROP TABLE IF EXISTS `pupil_assessment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pupil_assessment` (
  `Pupil_AssessmentID` bigint NOT NULL AUTO_INCREMENT,
  `Assessment_Year` int DEFAULT '0',
  `Assessment_Result_Indicator` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Assessment_Type` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Assessment_Component` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Component_Result` varchar(9) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Assessment_Locale` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `IndividualsID` bigint DEFAULT '0',
  `Assessment_Date` bigint NOT NULL DEFAULT '0',
  `Assessment_Subject` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Assessment_Identifier` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Component_Result_Type` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Pupil_AssessmentID`),
  KEY `WDIDX_pupil_assessment_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pupil_assessment`
--

LOCK TABLES `pupil_assessment` WRITE;
/*!40000 ALTER TABLE `pupil_assessment` DISABLE KEYS */;
/*!40000 ALTER TABLE `pupil_assessment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pupil_assessment_question_level_data`
--

DROP TABLE IF EXISTS `pupil_assessment_question_level_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pupil_assessment_question_level_data` (
  `Pupil_assessment_question_level_dataID` bigint NOT NULL AUTO_INCREMENT,
  `Pupil_AssessmentID` bigint DEFAULT '0',
  `Assessment_Focus` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Focus_Order` tinyint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  PRIMARY KEY (`Pupil_assessment_question_level_dataID`),
  KEY `WDIDX_pupil_assessment_question_level_data_Pupil_AssessmentID` (`Pupil_AssessmentID`),
  KEY `WDIDX_pupil_assessment_question_level_data_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pupil_assessment_question_level_data`
--

LOCK TABLES `pupil_assessment_question_level_data` WRITE;
/*!40000 ALTER TABLE `pupil_assessment_question_level_data` DISABLE KEYS */;
/*!40000 ALTER TABLE `pupil_assessment_question_level_data` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pupil_attendance`
--

DROP TABLE IF EXISTS `pupil_attendance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pupil_attendance` (
  `Pupil_AttendanceID` bigint NOT NULL AUTO_INCREMENT,
  `Attendance_Year` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Attendance_period_start_date` bigint NOT NULL DEFAULT '0',
  `Possible_Sessions` int DEFAULT '0',
  `Sessions_Attended` int DEFAULT '0',
  `Sessions_missed_due_to_Authorised_Absence` int DEFAULT '0',
  `Sessions_missed_due_to_Unauthorised_Absence` int DEFAULT '0',
  `Number_of_sessions_missed` int DEFAULT '0',
  `Funded_Hours` decimal(24,6) DEFAULT '0.000000',
  `Hours_at_Setting` decimal(24,6) DEFAULT '0.000000',
  `Unit_Contact_Time_Pupil` int DEFAULT '0',
  `Total_funded_spring_hours` decimal(24,6) DEFAULT '0.000000',
  `Attendance_Codes` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Attendance_marks_for_all_sessions` varchar(999) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `IndividualsID` bigint DEFAULT '0',
  PRIMARY KEY (`Pupil_AttendanceID`),
  KEY `WDIDX_pupil_attendance_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pupil_attendance`
--

LOCK TABLES `pupil_attendance` WRITE;
/*!40000 ALTER TABLE `pupil_attendance` DISABLE KEYS */;
/*!40000 ALTER TABLE `pupil_attendance` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pupil_exclusions`
--

DROP TABLE IF EXISTS `pupil_exclusions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pupil_exclusions` (
  `Pupil_ExclusionsID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `Exclusion_Category` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Exclusion_Reason` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Exclusion_Start_Date` bigint DEFAULT NULL,
  `Exclusion_Actual_Number_of_Sessions` int DEFAULT '0',
  UNIQUE KEY `Pupil_ExclusionsID` (`Pupil_ExclusionsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pupil_exclusions`
--

LOCK TABLES `pupil_exclusions` WRITE;
/*!40000 ALTER TABLE `pupil_exclusions` DISABLE KEYS */;
/*!40000 ALTER TABLE `pupil_exclusions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pupil_qualifications`
--

DROP TABLE IF EXISTS `pupil_qualifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pupil_qualifications` (
  `Pupil_QualificationsID` bigint NOT NULL AUTO_INCREMENT,
  `Qualification_Number` varchar(8) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Subject_Classification_Code` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `IndividualsID` bigint DEFAULT '0',
  PRIMARY KEY (`Pupil_QualificationsID`),
  KEY `WDIDX_pupil_qualifications_Qualification_Number` (`Qualification_Number`),
  KEY `WDIDX_pupil_qualifications_Subject_Classification_Code` (`Subject_Classification_Code`),
  KEY `WDIDX_pupil_qualifications_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pupil_qualifications`
--

LOCK TABLES `pupil_qualifications` WRITE;
/*!40000 ALTER TABLE `pupil_qualifications` DISABLE KEYS */;
/*!40000 ALTER TABLE `pupil_qualifications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pupil_removal_grounds`
--

DROP TABLE IF EXISTS `pupil_removal_grounds`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pupil_removal_grounds` (
  `Pupil_Removal_GroundsID` bigint NOT NULL AUTO_INCREMENT,
  `RemovalGroundsCode` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `RemovalReason` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `ReasonDescription` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `StatutoryRefernce` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Pupil_Removal_GroundsID`),
  KEY `WDIDX_pupil_removal_grounds_RemovalGroundsCode` (`RemovalGroundsCode`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pupil_removal_grounds`
--

LOCK TABLES `pupil_removal_grounds` WRITE;
/*!40000 ALTER TABLE `pupil_removal_grounds` DISABLE KEYS */;
/*!40000 ALTER TABLE `pupil_removal_grounds` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pupil_sen`
--

DROP TABLE IF EXISTS `pupil_sen`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pupil_sen` (
  `Pupil_SENID` bigint NOT NULL AUTO_INCREMENT,
  `Member_of_SEN_Unit_or_special_class_indicator` tinyint DEFAULT '0',
  `Member_of_resourced_provision_indicator` tinyint DEFAULT '0',
  `Pupil_date_placed_upon_stage` bigint NOT NULL DEFAULT '0',
  `Pupil_SEN_Type_ranking` int DEFAULT '0',
  `Pupil_Medical_Flag` tinyint DEFAULT '0',
  `Pupil_SEN_type_code` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `SEN_End_Date` bigint NOT NULL DEFAULT '0',
  `SEN_Need_Start_Date` bigint NOT NULL DEFAULT '0',
  `SEN_Need_End_Date` bigint NOT NULL DEFAULT '0',
  `SEN_Provision` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `IndividualsID` bigint DEFAULT '0',
  PRIMARY KEY (`Pupil_SENID`),
  KEY `WDIDX_pupil_sen_Member_of_SEN_Unit_or_special_class_indicator` (`Member_of_SEN_Unit_or_special_class_indicator`),
  KEY `WDIDX_pupil_sen_Member_of_resourced_provision_indicator` (`Member_of_resourced_provision_indicator`),
  KEY `WDIDX_pupil_sen_Pupil_date_placed_upon_stage` (`Pupil_date_placed_upon_stage`),
  KEY `WDIDX_pupil_sen_Pupil_SEN_Type_ranking` (`Pupil_SEN_Type_ranking`),
  KEY `WDIDX_pupil_sen_Pupil_Medical_Flag` (`Pupil_Medical_Flag`),
  KEY `WDIDX_pupil_sen_Pupil_SEN_type_code` (`Pupil_SEN_type_code`),
  KEY `WDIDX_pupil_sen_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pupil_sen`
--

LOCK TABLES `pupil_sen` WRITE;
/*!40000 ALTER TABLE `pupil_sen` DISABLE KEYS */;
/*!40000 ALTER TABLE `pupil_sen` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pupil_sen_provision`
--

DROP TABLE IF EXISTS `pupil_sen_provision`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pupil_sen_provision` (
  `Pupil_SEN_ProvisionID` bigint NOT NULL AUTO_INCREMENT,
  `ProvisionCode` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `ProvisionName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `ProvisionNotes` varchar(300) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Pupil_SEN_ProvisionID`),
  KEY `WDIDX_pupil_sen_provision_ProvisionCode` (`ProvisionCode`),
  KEY `WDIDX_pupil_sen_provision_ProvisionName` (`ProvisionName`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pupil_sen_provision`
--

LOCK TABLES `pupil_sen_provision` WRITE;
/*!40000 ALTER TABLE `pupil_sen_provision` DISABLE KEYS */;
/*!40000 ALTER TABLE `pupil_sen_provision` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pupil_status`
--

DROP TABLE IF EXISTS `pupil_status`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pupil_status` (
  `Pupil_StatusID` bigint NOT NULL AUTO_INCREMENT,
  `Pupil_child_Enrolment_status` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pupil_Date_of_entry` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pupil_Date_of_Leaving` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pupil_Part_time_Indicator` tinyint DEFAULT '0',
  `Pupil_Boarder_Indicator` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pupil_Class_Type_Indicator` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pupil_Entry_Date` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pupil_Leaving_Date` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Leaving_Reason` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pupil_s_Actual_National_Curriculum_Year_Group` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pupil_s_Actual_National_Curriculum_Year_Group_on_Leaving` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pupil_Removal_Grounds` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `IndividualsID` bigint DEFAULT '0',
  PRIMARY KEY (`Pupil_StatusID`),
  KEY `WDIDX_pupil_status_Pupil_child_Enrolment_status` (`Pupil_child_Enrolment_status`),
  KEY `WDIDX_pupil_status_Pupil_Date_of_entry` (`Pupil_Date_of_entry`),
  KEY `WDIDX_pupil_status_Pupil_Date_of_Leaving` (`Pupil_Date_of_Leaving`),
  KEY `WDIDX_pupil_status_Pupil_Part_time_Indicator` (`Pupil_Part_time_Indicator`),
  KEY `WDIDX_pupil_status_Pupil_Boarder_Indicator` (`Pupil_Boarder_Indicator`),
  KEY `WDIDX_pupil_status_Pupil_Class_Type_Indicator` (`Pupil_Class_Type_Indicator`),
  KEY `WDIDX_pupil_status_Pupil_Entry_Date` (`Pupil_Entry_Date`),
  KEY `WDIDX_pupil_status_Pupil_Leaving_Date` (`Pupil_Leaving_Date`),
  KEY `WDIDX_pupil_status_Leaving_Reason` (`Leaving_Reason`),
  KEY `WDIDX_pupil_status_Pupil_s_Actual_National_Curriculum_Year_Group` (`Pupil_s_Actual_National_Curriculum_Year_Group`),
  KEY `WDIDX_pupil_status_Pupil_s_Actual_National_Curriculum_Year_00044` (`Pupil_s_Actual_National_Curriculum_Year_Group_on_Leaving`),
  KEY `WDIDX_pupil_status_Pupil_Removal_Grounds` (`Pupil_Removal_Grounds`),
  KEY `WDIDX_pupil_status_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pupil_status`
--

LOCK TABLES `pupil_status` WRITE;
/*!40000 ALTER TABLE `pupil_status` DISABLE KEYS */;
/*!40000 ALTER TABLE `pupil_status` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pupilcharacteristics`
--

DROP TABLE IF EXISTS `pupilcharacteristics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pupilcharacteristics` (
  `PupilCharacteristicsID` bigint NOT NULL AUTO_INCREMENT,
  `Pupil_Free_School_Meal_Review_Date` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `In_Care_Indicator` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `In_Care_Caring_Authority_Code` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Language_Code` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Service_Children_in_Education_Indicator` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Type_of_Disability` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Language_Type` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Crown_Service` tinyint DEFAULT '0',
  `FSM_Eligibility_Start_Date` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `FSM_Eligibility_End_Date` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Youth_Support_Services_Agreement_Indicator` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Learner_Support_Code` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Top_up_funding_indicator` tinyint DEFAULT '0',
  `Full_time_employment_indicator` tinyint DEFAULT '0',
  `School_Lunch_Taken` tinyint DEFAULT '0',
  `Planned_Learning_Hours` int DEFAULT '0',
  `Planned_Employability_Enrichment_and_Pastorl_Hours` int DEFAULT '0',
  `Maths_GCSE_Highest_Prior_Attainment` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Maths_GCSE_Prior_Attainment_Year_Group` int DEFAULT '0',
  `English_GCSE_Highest_Prior_Attainment` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `English_GCSE_Prior_Attainment_Year_Group` int DEFAULT '0',
  `Early_years_pupil_premium_eligibility` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Early_years_pupil_premium_basis_for_funding` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Exempt_for_maths_GCSE_condition_of_funding` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Exempt_for_English_GCSE_condition_of_funding` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Child_Ethnicity` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Basis_for_funding` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Post_looked_after_arrangements` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pupil_Nationality` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pupil_Country_of_Birth` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Proficiency_in_English` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Proficiency_in_English_Date` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Source_of_Pupil_Ethnic_Code` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Source_for_Service_Child_Indicator` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Extended_childcare_hours` decimal(24,6) DEFAULT '0.000000',
  `Thirty_hour_code` varchar(11) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Disability_access_fund_indicator` tinyint DEFAULT '0',
  `URN_of_previous_school` varchar(6) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Primary_reason_for_placements` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Placement_Association` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Attendance_Pattern` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Sessions_per_week` int DEFAULT '0',
  `Partner_UKPRN` int DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  PRIMARY KEY (`PupilCharacteristicsID`),
  KEY `WDIDX_pupilcharacteristics_Pupil_Free_School_Meal_Review_Date` (`Pupil_Free_School_Meal_Review_Date`),
  KEY `WDIDX_pupilcharacteristics_In_Care_Indicator` (`In_Care_Indicator`),
  KEY `WDIDX_pupilcharacteristics_In_Care_Caring_Authority_Code` (`In_Care_Caring_Authority_Code`),
  KEY `WDIDX_pupilcharacteristics_Language_Code` (`Language_Code`),
  KEY `WDIDX_pupilcharacteristics_Service_Children_in_Education_In00045` (`Service_Children_in_Education_Indicator`),
  KEY `WDIDX_pupilcharacteristics_Type_of_Disability` (`Type_of_Disability`),
  KEY `WDIDX_pupilcharacteristics_Language_Type` (`Language_Type`),
  KEY `WDIDX_pupilcharacteristics_Crown_Service` (`Crown_Service`),
  KEY `WDIDX_pupilcharacteristics_FSM_Eligibility_Start_Date` (`FSM_Eligibility_Start_Date`),
  KEY `WDIDX_pupilcharacteristics_FSM_Eligibility_End_Date` (`FSM_Eligibility_End_Date`),
  KEY `WDIDX_pupilcharacteristics_Youth_Support_Services_Agreement00046` (`Youth_Support_Services_Agreement_Indicator`),
  KEY `WDIDX_pupilcharacteristics_Learner_Support_Code` (`Learner_Support_Code`),
  KEY `WDIDX_pupilcharacteristics_Top_up_funding_indicator` (`Top_up_funding_indicator`),
  KEY `WDIDX_pupilcharacteristics_Full_time_employment_indicator` (`Full_time_employment_indicator`),
  KEY `WDIDX_pupilcharacteristics_School_Lunch_Taken` (`School_Lunch_Taken`),
  KEY `WDIDX_pupilcharacteristics_Planned_Learning_Hours` (`Planned_Learning_Hours`),
  KEY `WDIDX_pupilcharacteristics_Planned_Employability_Enrichment00047` (`Planned_Employability_Enrichment_and_Pastorl_Hours`),
  KEY `WDIDX_pupilcharacteristics_Maths_GCSE_Highest_Prior_Attainment` (`Maths_GCSE_Highest_Prior_Attainment`),
  KEY `WDIDX_pupilcharacteristics_Maths_GCSE_Prior_Attainment_Year00048` (`Maths_GCSE_Prior_Attainment_Year_Group`),
  KEY `WDIDX_pupilcharacteristics_English_GCSE_Highest_Prior_Attainment` (`English_GCSE_Highest_Prior_Attainment`),
  KEY `WDIDX_pupilcharacteristics_English_GCSE_Prior_Attainment_Ye00049` (`English_GCSE_Prior_Attainment_Year_Group`),
  KEY `WDIDX_pupilcharacteristics_Early_years_pupil_premium_eligibility` (`Early_years_pupil_premium_eligibility`),
  KEY `WDIDX_pupilcharacteristics_Early_years_pupil_premium_basis_00050` (`Early_years_pupil_premium_basis_for_funding`),
  KEY `WDIDX_pupilcharacteristics_Exempt_for_maths_GCSE_condition_00051` (`Exempt_for_maths_GCSE_condition_of_funding`),
  KEY `WDIDX_pupilcharacteristics_Exempt_for_English_GCSE_conditio00052` (`Exempt_for_English_GCSE_condition_of_funding`),
  KEY `WDIDX_pupilcharacteristics_Child_Ethnicity` (`Child_Ethnicity`),
  KEY `WDIDX_pupilcharacteristics_Basis_for_funding` (`Basis_for_funding`),
  KEY `WDIDX_pupilcharacteristics_Post_looked_after_arrangements` (`Post_looked_after_arrangements`),
  KEY `WDIDX_pupilcharacteristics_Pupil_Nationality` (`Pupil_Nationality`),
  KEY `WDIDX_pupilcharacteristics_Pupil_Country_of_Birth` (`Pupil_Country_of_Birth`),
  KEY `WDIDX_pupilcharacteristics_Proficiency_in_English` (`Proficiency_in_English`),
  KEY `WDIDX_pupilcharacteristics_Proficiency_in_English_Date` (`Proficiency_in_English_Date`),
  KEY `WDIDX_pupilcharacteristics_Source_of_Pupil_Ethnic_Code` (`Source_of_Pupil_Ethnic_Code`),
  KEY `WDIDX_pupilcharacteristics_Source_for_Service_Child_Indicator` (`Source_for_Service_Child_Indicator`),
  KEY `WDIDX_pupilcharacteristics_Extended_childcare_hours` (`Extended_childcare_hours`),
  KEY `WDIDX_pupilcharacteristics_thirty_hour_code` (`Thirty_hour_code`),
  KEY `WDIDX_pupilcharacteristics_Disability_access_fund_indicator` (`Disability_access_fund_indicator`),
  KEY `WDIDX_pupilcharacteristics_URN_of_previous_school` (`URN_of_previous_school`),
  KEY `WDIDX_pupilcharacteristics_Primary_reason_for_placements` (`Primary_reason_for_placements`),
  KEY `WDIDX_pupilcharacteristics_placement_Association` (`Placement_Association`),
  KEY `WDIDX_pupilcharacteristics_Attendance_Pattern` (`Attendance_Pattern`),
  KEY `WDIDX_pupilcharacteristics_Sessions_per_week` (`Sessions_per_week`),
  KEY `WDIDX_pupilcharacteristics_Partner_UKPRN` (`Partner_UKPRN`),
  KEY `WDIDX_pupilcharacteristics_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pupilcharacteristics`
--

LOCK TABLES `pupilcharacteristics` WRITE;
/*!40000 ALTER TABLE `pupilcharacteristics` DISABLE KEYS */;
/*!40000 ALTER TABLE `pupilcharacteristics` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `reason_for_leaving`
--

DROP TABLE IF EXISTS `reason_for_leaving`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `reason_for_leaving` (
  `Reason_for_LeavingID` bigint NOT NULL AUTO_INCREMENT,
  `Reason_Code` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Reason_description` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Current` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Reason_for_LeavingID`),
  KEY `WDIDX_reason_for_leaving_Current` (`Current`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `reason_for_leaving`
--

LOCK TABLES `reason_for_leaving` WRITE;
/*!40000 ALTER TABLE `reason_for_leaving` DISABLE KEYS */;
/*!40000 ALTER TABLE `reason_for_leaving` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `register_types`
--

DROP TABLE IF EXISTS `register_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `register_types` (
  `Register_typesID` bigint NOT NULL DEFAULT '0',
  `RegName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `FocusRegColourCodesID` bigint DEFAULT '0',
  `Contact_name` varchar(50) DEFAULT NULL,
  `Contact_email` varchar(260) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  UNIQUE KEY `Register_typesID` (`Register_typesID`),
  KEY `WDIDX_register_types_RegName` (`RegName`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `register_types`
--

LOCK TABLES `register_types` WRITE;
/*!40000 ALTER TABLE `register_types` DISABLE KEYS */;
/*!40000 ALTER TABLE `register_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rewardslimits`
--

DROP TABLE IF EXISTS `rewardslimits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `rewardslimits` (
  `RewardsLimitsID` bigint NOT NULL AUTO_INCREMENT,
  `AttainGrp1Blk` tinyint DEFAULT '0',
  `AttainGrp2Blk` tinyint DEFAULT '0',
  `AttainGrpYr` tinyint DEFAULT '0',
  `AttainYr1Blk` tinyint DEFAULT '0',
  `AttainYrBlk` tinyint DEFAULT '0',
  `AttainYrYr` tinyint DEFAULT '0',
  `ProgGrp1Blk` tinyint DEFAULT '0',
  `ProgGrp2Blk` tinyint DEFAULT '0',
  `ProgGrpYr` tinyint DEFAULT '0',
  `ProgYr1Blk` tinyint DEFAULT '0',
  `ProgYr2Blk` tinyint DEFAULT '0',
  `ProgYrYr` tinyint DEFAULT '0',
  `AttainTopPercent` tinyint DEFAULT '0',
  `ProgTopPercent` tinyint DEFAULT '0',
  PRIMARY KEY (`RewardsLimitsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rewardslimits`
--

LOCK TABLES `rewardslimits` WRITE;
/*!40000 ALTER TABLE `rewardslimits` DISABLE KEYS */;
/*!40000 ALTER TABLE `rewardslimits` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `room_bookings`
--

DROP TABLE IF EXISTS `room_bookings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `room_bookings` (
  `Room_bookingsID` bigint NOT NULL AUTO_INCREMENT,
  `RoomsID` bigint DEFAULT '0',
  `Lesson_outlineID` bigint DEFAULT '0',
  `Room_name` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `AppointmentsID` bigint DEFAULT '0',
  UNIQUE KEY `Room_bookingsID` (`Room_bookingsID`),
  KEY `WDIDX_room_bookings_roomsID` (`RoomsID`),
  KEY `WDIDX_room_bookings_Room_name` (`Room_name`),
  KEY `WDIDX_room_bookings_AppointmentsID` (`AppointmentsID`),
  KEY `WDIDX_room_bookings_WDIDX_room_bookings_WDIDX_room_bookings00053` (`RoomsID`,`Lesson_outlineID`),
  KEY `WDIDX_room_bookings_WDIDX_room_bookings_WDIDX_Room_bookings00054` (`RoomsID`,`Lesson_outlineID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `room_bookings`
--

LOCK TABLES `room_bookings` WRITE;
/*!40000 ALTER TABLE `room_bookings` DISABLE KEYS */;
/*!40000 ALTER TABLE `room_bookings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `room_types`
--

DROP TABLE IF EXISTS `room_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `room_types` (
  `Room_typesID` bigint NOT NULL AUTO_INCREMENT,
  `Description` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Type_of_room` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Room_typesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `room_types`
--

LOCK TABLES `room_types` WRITE;
/*!40000 ALTER TABLE `room_types` DISABLE KEYS */;
/*!40000 ALTER TABLE `room_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rooms`
--

DROP TABLE IF EXISTS `rooms`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `rooms` (
  `RoomsID` bigint NOT NULL AUTO_INCREMENT,
  `Room_name` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Room_typesID` bigint DEFAULT '0',
  `Places` varchar(1000) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Seats_across` int DEFAULT '0',
  `Seat_down` int DEFAULT '0',
  `Room_description` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Capacity` int DEFAULT '0',
  PRIMARY KEY (`RoomsID`),
  KEY `WDIDX_rooms_Room_name` (`Room_name`),
  KEY `WDIDX_rooms_room_typesID` (`Room_typesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rooms`
--

LOCK TABLES `rooms` WRITE;
/*!40000 ALTER TABLE `rooms` DISABLE KEYS */;
/*!40000 ALTER TABLE `rooms` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rooms_freed`
--

DROP TABLE IF EXISTS `rooms_freed`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `rooms_freed` (
  `Rooms_FreedID` bigint NOT NULL AUTO_INCREMENT,
  `RoomsID` bigint DEFAULT '0',
  `SessionsID` bigint DEFAULT '0',
  UNIQUE KEY `Rooms_FreedID` (`Rooms_FreedID`),
  KEY `WDIDX_rooms_freed_roomsID` (`RoomsID`),
  KEY `WDIDX_rooms_freed_SessionsID` (`SessionsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rooms_freed`
--

LOCK TABLES `rooms_freed` WRITE;
/*!40000 ALTER TABLE `rooms_freed` DISABLE KEYS */;
/*!40000 ALTER TABLE `rooms_freed` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `safeguardissues`
--

DROP TABLE IF EXISTS `safeguardissues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `safeguardissues` (
  `safeguardIssuesID` bigint NOT NULL AUTO_INCREMENT,
  `initialDate` bigint DEFAULT '0',
  `issueDesc` varchar(200) DEFAULT NULL,
  `Current` varchar(1) DEFAULT NULL,
  `IndividualsID` bigint DEFAULT '0',
  `EnteredBy` bigint DEFAULT '0',
  PRIMARY KEY (`safeguardIssuesID`),
  KEY `WDIDX_safeguardIssues_Current` (`Current`),
  KEY `WDIDX_safeguardIssues_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_safeguardIssues_EnteredBy` (`EnteredBy`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `safeguardissues`
--

LOCK TABLES `safeguardissues` WRITE;
/*!40000 ALTER TABLE `safeguardissues` DISABLE KEYS */;
/*!40000 ALTER TABLE `safeguardissues` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `safeguardobs`
--

DROP TABLE IF EXISTS `safeguardobs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `safeguardobs` (
  `SafeguardObsID` bigint NOT NULL AUTO_INCREMENT,
  `DateTimeObs` bigint DEFAULT '0',
  `SafeguardObsType` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  `LeaderID` bigint DEFAULT '0',
  PRIMARY KEY (`SafeguardObsID`),
  KEY `WDIDX_SafeguardObs_DateTimeObs` (`DateTimeObs`),
  KEY `WDIDX_SafeguardObs_SafeguardObsType` (`SafeguardObsType`),
  KEY `WDIDX_SafeguardObs_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_SafeguardObs_leaderID` (`LeaderID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `safeguardobs`
--

LOCK TABLES `safeguardobs` WRITE;
/*!40000 ALTER TABLE `safeguardobs` DISABLE KEYS */;
/*!40000 ALTER TABLE `safeguardobs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `safeguardobstype`
--

DROP TABLE IF EXISTS `safeguardobstype`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `safeguardobstype` (
  `SafeguardObsTypeID` bigint NOT NULL AUTO_INCREMENT,
  `Description` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Level` tinyint DEFAULT '0',
  PRIMARY KEY (`SafeguardObsTypeID`),
  KEY `WDIDX_SafeguardObsType_Level` (`Level`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `safeguardobstype`
--

LOCK TABLES `safeguardobstype` WRITE;
/*!40000 ALTER TABLE `safeguardobstype` DISABLE KEYS */;
/*!40000 ALTER TABLE `safeguardobstype` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sanction_attendance`
--

DROP TABLE IF EXISTS `sanction_attendance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sanction_attendance` (
  `Sanction_attendanceID` bigint NOT NULL AUTO_INCREMENT,
  `Sanction_sessionsID` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  `Status` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  UNIQUE KEY `Sanction_attendanceID` (`Sanction_attendanceID`),
  KEY `WDIDX_sanction_attendance_Sanction_sessionsID` (`Sanction_sessionsID`),
  KEY `WDIDX_sanction_attendance_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_sanction_attendance_status` (`Status`),
  KEY `WDIDX_sanction_attendance_WDIDX_Sanction_attendance_Sanctio00055` (`Sanction_sessionsID`,`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sanction_attendance`
--

LOCK TABLES `sanction_attendance` WRITE;
/*!40000 ALTER TABLE `sanction_attendance` DISABLE KEYS */;
/*!40000 ALTER TABLE `sanction_attendance` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sanction_escalations`
--

DROP TABLE IF EXISTS `sanction_escalations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sanction_escalations` (
  `Sanction_escalationsID` bigint NOT NULL AUTO_INCREMENT,
  `Sanct1Level` int DEFAULT '0',
  `Sanction2ID` bigint DEFAULT '0',
  `HowMany` int DEFAULT '0',
  `Duration_weeks` int DEFAULT '0',
  `Sanction2_level` tinyint DEFAULT '0',
  UNIQUE KEY `Sanction_escalationsID` (`Sanction_escalationsID`),
  KEY `WDIDX_sanction_escalations_Sanct1Level` (`Sanct1Level`),
  KEY `WDIDX_sanction_escalations_Sanction2_level` (`Sanction2_level`),
  KEY `WDIDX_sanction_escalations_WDIDX_Sanction_escalations_Sanct00056` (`Sanct1Level`,`Sanction2ID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sanction_escalations`
--

LOCK TABLES `sanction_escalations` WRITE;
/*!40000 ALTER TABLE `sanction_escalations` DISABLE KEYS */;
/*!40000 ALTER TABLE `sanction_escalations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sanction_sessions`
--

DROP TABLE IF EXISTS `sanction_sessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sanction_sessions` (
  `Sanction_sessionsID` bigint NOT NULL AUTO_INCREMENT,
  `StartTime` timestamp NULL DEFAULT NULL,
  `FinishTime` timestamp NULL DEFAULT NULL,
  `RoomsID` bigint DEFAULT '0',
  `SessionName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Sanction_typesID` bigint DEFAULT '0',
  UNIQUE KEY `Sanction_sessionsID` (`Sanction_sessionsID`),
  KEY `WDIDX_sanction_sessions_StartTime` (`StartTime`),
  KEY `WDIDX_sanction_sessions_FinishTime` (`FinishTime`),
  KEY `WDIDX_sanction_sessions_roomsID` (`RoomsID`),
  KEY `WDIDX_sanction_sessions_SessionName` (`SessionName`),
  KEY `WDIDX_sanction_sessions_Sanction_typesID` (`Sanction_typesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sanction_sessions`
--

LOCK TABLES `sanction_sessions` WRITE;
/*!40000 ALTER TABLE `sanction_sessions` DISABLE KEYS */;
/*!40000 ALTER TABLE `sanction_sessions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sanction_set`
--

DROP TABLE IF EXISTS `sanction_set`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sanction_set` (
  `Sanction_setID` bigint NOT NULL AUTO_INCREMENT,
  `Major` tinyint DEFAULT '0',
  `TransgressionID` bigint DEFAULT '0',
  `Sanction_sessionsID` bigint DEFAULT '0',
  `DateAndTime` bigint DEFAULT NULL,
  `RecordedBy` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  `Sanction_typesID` bigint DEFAULT '0',
  UNIQUE KEY `Sanction_setID` (`Sanction_setID`),
  KEY `WDIDX_sanction_set_Sanction_sessionsID` (`Sanction_sessionsID`),
  KEY `WDIDX_sanction_set_DateAndTime` (`DateAndTime`),
  KEY `WDIDX_sanction_set_RecordedBy` (`RecordedBy`),
  KEY `WDIDX_sanction_set_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_sanction_set_Sanction_typesID` (`Sanction_typesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sanction_set`
--

LOCK TABLES `sanction_set` WRITE;
/*!40000 ALTER TABLE `sanction_set` DISABLE KEYS */;
/*!40000 ALTER TABLE `sanction_set` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sanction_staff`
--

DROP TABLE IF EXISTS `sanction_staff`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sanction_staff` (
  `Sanction_StaffID` bigint NOT NULL AUTO_INCREMENT,
  `Sanction_sessionsID` bigint DEFAULT '0',
  UNIQUE KEY `Sanction_StaffID` (`Sanction_StaffID`),
  KEY `WDIDX_sanction_staff_Sanction_sessionsID` (`Sanction_sessionsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sanction_staff`
--

LOCK TABLES `sanction_staff` WRITE;
/*!40000 ALTER TABLE `sanction_staff` DISABLE KEYS */;
/*!40000 ALTER TABLE `sanction_staff` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sanction_types`
--

DROP TABLE IF EXISTS `sanction_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sanction_types` (
  `Sanction_typesID` bigint NOT NULL AUTO_INCREMENT,
  `Sanction_name` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Sanction_level` int DEFAULT '0',
  UNIQUE KEY `Sanction_typesID` (`Sanction_typesID`),
  UNIQUE KEY `Sanction_name` (`Sanction_name`),
  KEY `WDIDX_sanction_types_Sanction_level` (`Sanction_level`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sanction_types`
--

LOCK TABLES `sanction_types` WRITE;
/*!40000 ALTER TABLE `sanction_types` DISABLE KEYS */;
/*!40000 ALTER TABLE `sanction_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `school_admissions`
--

DROP TABLE IF EXISTS `school_admissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `school_admissions` (
  `School_admissionsID` bigint NOT NULL AUTO_INCREMENT,
  `Admission_appeals_lodged` int DEFAULT '0',
  `Admission_appeals_withdrawn` int DEFAULT '0',
  `Appeals_heard_by_Ind_Admissions_Com` int DEFAULT '0',
  `Ind_Admissions_Com_decided_in_parent_s_favour` int DEFAULT '0',
  `Ind_Admissions_Com_rejected_appeal` int DEFAULT '0',
  `Admissions_entry_year_group` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`School_admissionsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `school_admissions`
--

LOCK TABLES `school_admissions` WRITE;
/*!40000 ALTER TABLE `school_admissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `school_admissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `school_attendance_codes`
--

DROP TABLE IF EXISTS `school_attendance_codes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `school_attendance_codes` (
  `School_attendance_codesID` bigint NOT NULL AUTO_INCREMENT,
  `AttendanceCode` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AttendanceCodeDescr` varchar(150) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `CensusIncusion` tinyint DEFAULT '0',
  PRIMARY KEY (`School_attendance_codesID`),
  KEY `WDIDX_school_attendance_codes_CensusIncusion` (`CensusIncusion`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `school_attendance_codes`
--

LOCK TABLES `school_attendance_codes` WRITE;
/*!40000 ALTER TABLE `school_attendance_codes` DISABLE KEYS */;
/*!40000 ALTER TABLE `school_attendance_codes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `school_childcare`
--

DROP TABLE IF EXISTS `school_childcare`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `school_childcare` (
  `School_childcareID` bigint NOT NULL AUTO_INCREMENT,
  `Childcare_provider` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Childcare_weeks_open` int DEFAULT '0',
  `Other_Schools` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Type_of_childcare` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Childcare_on_site` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Number_Of_Childcare_places` int DEFAULT '0',
  `Signposting_off_site_childcare_provision` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`School_childcareID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `school_childcare`
--

LOCK TABLES `school_childcare` WRITE;
/*!40000 ALTER TABLE `school_childcare` DISABLE KEYS */;
/*!40000 ALTER TABLE `school_childcare` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `school_class`
--

DROP TABLE IF EXISTS `school_class`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `school_class` (
  `School_ClassID` bigint NOT NULL AUTO_INCREMENT,
  `Class_Reference_Name` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Number_of_Teachers_in_the_Class` int DEFAULT '0',
  `Number_of_Adult_Non_Teachers_in_the_Class` int DEFAULT '0',
  `Class_Activity` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Number_of_Pupils_from_the_host_school_in_the_class` int DEFAULT '0',
  `Number_of_Pupils_from_other_schools_in_the_class` int DEFAULT '0',
  `Class_Type` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Class_Key_Stage` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Class_Year_Group` int DEFAULT '0',
  PRIMARY KEY (`School_ClassID`),
  KEY `WDIDX_school_class_Class_Reference_Name` (`Class_Reference_Name`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `school_class`
--

LOCK TABLES `school_class` WRITE;
/*!40000 ALTER TABLE `school_class` DISABLE KEYS */;
/*!40000 ALTER TABLE `school_class` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `school_detail`
--

DROP TABLE IF EXISTS `school_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `school_detail` (
  `SchoolName` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `RouterName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Password` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Timetable_length` int DEFAULT '0',
  `SMTP` varchar(80) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `LocationLat` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `LocationLong` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `CurrentWeek` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Teacher_name` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `LastLocCheck` bigint DEFAULT NULL,
  `Student_name` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `IMAP_port` int DEFAULT '0',
  `SMTP_port` int DEFAULT '0',
  `IMAP` varchar(80) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `TeacherEntID` bigint DEFAULT '0',
  `StudentEntID` bigint DEFAULT '0',
  `AccountType` varchar(6) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `GoogleAPIKey` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `OneDriveAPIKey` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `GoogleSecretKey` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `OneDriceSecretKey` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `DwellingsID` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  UNIQUE KEY `SchoolName` (`SchoolName`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `school_detail`
--

LOCK TABLES `school_detail` WRITE;
/*!40000 ALTER TABLE `school_detail` DISABLE KEYS */;
/*!40000 ALTER TABLE `school_detail` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `school_pupil_numbers`
--

DROP TABLE IF EXISTS `school_pupil_numbers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `school_pupil_numbers` (
  `School_pupil_NumbersID` bigint NOT NULL AUTO_INCREMENT,
  `Boy_19_over` int DEFAULT '0',
  `Boy_18` int DEFAULT '0',
  `Boy_17` int DEFAULT '0',
  `Boy_16` int DEFAULT '0',
  `Boy_15` int DEFAULT '0',
  `Boy_14` int DEFAULT '0',
  `Boy_13` int DEFAULT '0',
  `Boy_12` int DEFAULT '0',
  `Boy_11` int DEFAULT '0',
  `Boy_10` int DEFAULT '0',
  `Boy_9` int DEFAULT '0',
  `Boy_8` int DEFAULT '0',
  `Boy_7` int DEFAULT '0',
  `Boy_6` int DEFAULT '0',
  `Boy_5` int DEFAULT '0',
  `Boy_4_upper` int DEFAULT '0',
  `Boy_4_middle` int DEFAULT '0',
  `Boy_4_lower` int DEFAULT '0',
  `Boy_3_upper` int DEFAULT '0',
  `Boy_3_middle` int DEFAULT '0',
  `Boy_3_lower` int DEFAULT '0',
  `Boy_2` int DEFAULT '0',
  `Boy_1` int DEFAULT '0',
  `Boy_under_1` int DEFAULT '0',
  `Girl_19_over` int DEFAULT '0',
  `Boy_total` int DEFAULT '0',
  `Girl_18` int DEFAULT '0',
  `Girl_17` int DEFAULT '0',
  `Girl_16` int DEFAULT '0',
  `Girl_15` int DEFAULT '0',
  `Girl_14` int DEFAULT '0',
  `Girl_13` int DEFAULT '0',
  `Girl_12` int DEFAULT '0',
  `Girl_11` int DEFAULT '0',
  `Girl_10` int DEFAULT '0',
  `Girl_9` int DEFAULT '0',
  `Girl_8` int DEFAULT '0',
  `Girl_7` int DEFAULT '0',
  `Girl_6` int DEFAULT '0',
  `Girl_5` int DEFAULT '0',
  `Girl_4_upper` int DEFAULT '0',
  `Girl_4_middle` int DEFAULT '0',
  `Girl_4_lower` int DEFAULT '0',
  `Girl_3_upper` int DEFAULT '0',
  `Girl_3_middle` int DEFAULT '0',
  `Girl_3_lower` int DEFAULT '0',
  `Girl_2` int DEFAULT '0',
  `Girl_1` int DEFAULT '0',
  `Girl_under_1` int DEFAULT '0',
  `Girl_total` int DEFAULT '0',
  `Boy_PT_19_over` int DEFAULT '0',
  `Boy_PT_18` int DEFAULT '0',
  `Boy_Pt_17` int DEFAULT '0',
  `Boy_Pt_16` int DEFAULT '0',
  `Boy_Pt_15` int DEFAULT '0',
  `Boy_Pt_14` int DEFAULT '0',
  `Boy_Pt_13` int DEFAULT '0',
  `Boy_Pt_12` int DEFAULT '0',
  `Boy_Pt_11` int DEFAULT '0',
  `Boy_Pt_10` int DEFAULT '0',
  `Boy_Pt_9` int DEFAULT '0',
  `Boy_Pt_8` int DEFAULT '0',
  `Boy_Pt_7` int DEFAULT '0',
  `Boy_Pt_6` int DEFAULT '0',
  `Boy_Pt_5` int DEFAULT '0',
  `Boy_Pt_4_upper` int DEFAULT '0',
  `Boy_Pt_4_middle` int DEFAULT '0',
  `Boy_Pt_4_lower` int DEFAULT '0',
  `Boy_Pt_3_upper` int DEFAULT '0',
  `Boy_Pt_3_middle` int DEFAULT '0',
  `Boy_Pt_3_lower` int DEFAULT '0',
  `Boy_Pt_2` int DEFAULT '0',
  `Boy_Pt_1` int DEFAULT '0',
  `Boy_Pt_under_1` int DEFAULT '0',
  `Boy_Pt_total` int DEFAULT '0',
  `Girl_PT_19_over` int DEFAULT '0',
  `Girl_PT_18` int DEFAULT '0',
  `Girl_PT_17` int DEFAULT '0',
  `Girl_PT_16` int DEFAULT '0',
  `Girl_PT_15` int DEFAULT '0',
  `Girl_PT_14` int DEFAULT '0',
  `Girl_PT_13` int DEFAULT '0',
  `Girl_PT_12` int DEFAULT '0',
  `Girl_PT_11` int DEFAULT '0',
  `Girl_PT_10` int DEFAULT '0',
  `Girl_PT_9` int DEFAULT '0',
  `Girl_PT_8` int DEFAULT '0',
  `Girl_PT_7` int DEFAULT '0',
  `Girl_PT_6` int DEFAULT '0',
  `Girl_PT_5` int DEFAULT '0',
  `Girl_PT_4_upper` int DEFAULT '0',
  `Girl_PT_4_middle` int DEFAULT '0',
  `Girl_PT_4_lower` int DEFAULT '0',
  `Girl_PT_3_upper` int DEFAULT '0',
  `Girl_PT_3_middle` int DEFAULT '0',
  `Girl_PT_3_lower` int DEFAULT '0',
  `Girl_PT_2` int DEFAULT '0',
  `Girl_PT_1` int DEFAULT '0',
  `Girl_PT_Under_1` int DEFAULT '0',
  `Girl_PT_total` int DEFAULT '0',
  `Boarding_boys` int DEFAULT '0',
  `Boarding_girls` int DEFAULT '0',
  `Pupils_in_care` int DEFAULT '0',
  `Pupil_SEN_statement` int DEFAULT '0',
  `Pupil_SEN_no_statement` int DEFAULT '0',
  `Courses_boys_15_level_4_and_above` int DEFAULT '0',
  `Courses_boys_15_international_baccalaureate` int DEFAULT '0',
  `Course_boys15GCSE_AlevelPreUPrincGCSE_ASlevel_preUShort` int DEFAULT '0',
  `Courses_boys_15_other_level_3_equivalents` int DEFAULT '0',
  `Courses_boys_15_GCSE_IGCSE` int DEFAULT '0',
  `Courses_boys_15_other_level_2_courses` int DEFAULT '0',
  `Courses_boys_15_other_level_1_courses` int DEFAULT '0',
  `Courses_boys_15_other_courses` int DEFAULT '0',
  `Courses_boys_16_level_4_and_above` int DEFAULT '0',
  `Courses_boys_16_international_baccalaureate` int DEFAULT '0',
  `Course_boys16GCSE_AlevelPreUPrincGCSE_ASlevel_preUShort` int DEFAULT '0',
  `Courses_boys_16_other_level_3_equivalents` int DEFAULT '0',
  `Courses_boys_16_GCSE_IGCSE` int DEFAULT '0',
  `Courses_boys_16_other_level_2_courses` int DEFAULT '0',
  `Courses_boys_16_other_level_1_courses` int DEFAULT '0',
  `Courses_boys_16_other_courses` int DEFAULT '0',
  `Courses_boys_17_level_4_and_above` int DEFAULT '0',
  `Courses_boys_17_international_baccalaureate` int DEFAULT '0',
  `Course_boys17GCSE_AlevelPreUPrincGCSE_ASlevel_preUShort` int DEFAULT '0',
  `Courses_boys_17_other_level_3_equivalents` int DEFAULT '0',
  `Courses_boys_17_GCSE_IGCSE` int DEFAULT '0',
  `Courses_boys_17_other_level_2_courses` int DEFAULT '0',
  `Courses_boys_17_other_level_1_courses` int DEFAULT '0',
  `Courses_boys_17_other_courses` int DEFAULT '0',
  `Courses_boys_18_and_over_level_4_and_above` int DEFAULT '0',
  `Courses_boys_18_and_over_international_baccalaureate` int DEFAULT '0',
  `Course_boys18GCSE_AlevelPreUPrincGCSE_ASlevel_preUShort` int DEFAULT '0',
  `Courses_boys_18_and_over_other_level_3_equivalents` int DEFAULT '0',
  `Courses_boys_18_and_over_GCSE_IGCSE` int DEFAULT '0',
  `Courses_boys_18_and_over_other_level_2_courses` int DEFAULT '0',
  `Courses_boys_18_and_over_other_level_1_courses` int DEFAULT '0',
  `Courses_boys_18_and_over_other_courses` int DEFAULT '0',
  `Courses_girls_15_level_4_and_above` int DEFAULT '0',
  `Courses_girls_15_international_baccalaureate` int DEFAULT '0',
  `Course_girls15GCSE_AlevelPreUPrincGCSE_ASlevel_preUShort` int DEFAULT '0',
  `Courses_girls_15_other_level_3_equivalents` int DEFAULT '0',
  `Courses_girls_15_GCSE_IGCSE` int DEFAULT '0',
  `Courses_girls_15_other_level_2_courses` int DEFAULT '0',
  `Courses_girls_15_other_level_1_courses` int DEFAULT '0',
  `Courses_girls_15_other_courses` int DEFAULT '0',
  `Courses_girls_16_level_4_and_above` int DEFAULT '0',
  `Courses_girls_16_international_baccalaureate` int DEFAULT '0',
  `Course_girls16GCSE_AlevelPreUPrincGCSE_ASlevel_preUShort` int DEFAULT '0',
  `Courses_girls_16_other_level_3_equivalents` int DEFAULT '0',
  `Courses_girls_16_GCSE_IGCSE` int DEFAULT '0',
  `Courses_girls_16_other_level_2_courses` int DEFAULT '0',
  `Courses_girls_16_other_level_1_courses` int DEFAULT '0',
  `Courses_girls_16_other_courses` int DEFAULT '0',
  `Courses_girls_17_level_4_and_above` int DEFAULT '0',
  `Courses_girls_17_international_baccalaureate` int DEFAULT '0',
  `Course_girls17GCSE_AlevelPreUPrincGCSE_ASlevel_preUShort` int DEFAULT '0',
  `Courses_girls_17_other_level_3_equivalents` int DEFAULT '0',
  `Courses_girls_17_GCSE_IGCSE` int DEFAULT '0',
  `Courses_girls_17_other_level_2_courses` int DEFAULT '0',
  `Courses_girls_17_other_level_1_courses` int DEFAULT '0',
  `Courses_girls_17_other_courses` int DEFAULT '0',
  `Courses_girls_18_and_over_level_4_and_above` int DEFAULT '0',
  `Courses_girls_18_and_over_international_baccalaureate` int DEFAULT '0',
  `Course_girls18GCSE_AlevelPreUPrincGCSE_ASlevel_preUShort` int DEFAULT '0',
  `Courses_girls_18_and_over_other_level_3_equivalents` int DEFAULT '0',
  `Courses_girls_18_and_over_GCSE_IGCSE` int DEFAULT '0',
  `Courses_girls_18_and_over_other_level_2_courses` int DEFAULT '0',
  `Courses_girls_18_and_over_other_level_1_courses` int DEFAULT '0',
  `Courses_girls_18_and_over_other_courses` int DEFAULT '0',
  `Courses_boys_final_year_key_stage_4` int DEFAULT '0',
  `Courses_girls_final_year_key_stage_4` int DEFAULT '0',
  `Part_Time_pupils_not_at_school` int DEFAULT '0',
  `Private_Study_pupils` int DEFAULT '0',
  `Pupils_at_Another_School` int DEFAULT '0',
  `Pupils_on_Work_Experience` int DEFAULT '0',
  `Pupils_at_FE_Colleges` int DEFAULT '0',
  PRIMARY KEY (`School_pupil_NumbersID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `school_pupil_numbers`
--

LOCK TABLES `school_pupil_numbers` WRITE;
/*!40000 ALTER TABLE `school_pupil_numbers` DISABLE KEYS */;
/*!40000 ALTER TABLE `school_pupil_numbers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `school_pupil_statistics`
--

DROP TABLE IF EXISTS `school_pupil_statistics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `school_pupil_statistics` (
  `School_pupil_statisticsID` bigint NOT NULL AUTO_INCREMENT,
  `Free_School_Meals_Taken` int DEFAULT '0',
  `School_Term` int DEFAULT '0',
  `Number_of_3_year_olds` int DEFAULT '0',
  `Number_of_4_year_olds` int DEFAULT '0',
  `Number_of_2_year_olds` int DEFAULT '0',
  PRIMARY KEY (`School_pupil_statisticsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `school_pupil_statistics`
--

LOCK TABLES `school_pupil_statistics` WRITE;
/*!40000 ALTER TABLE `school_pupil_statistics` DISABLE KEYS */;
/*!40000 ALTER TABLE `school_pupil_statistics` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `school_reconciliation`
--

DROP TABLE IF EXISTS `school_reconciliation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `school_reconciliation` (
  `School_reconciliationID` bigint NOT NULL AUTO_INCREMENT,
  `Part_Time_pupils_not_at_school` int DEFAULT '0',
  PRIMARY KEY (`School_reconciliationID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `school_reconciliation`
--

LOCK TABLES `school_reconciliation` WRITE;
/*!40000 ALTER TABLE `school_reconciliation` DISABLE KEYS */;
/*!40000 ALTER TABLE `school_reconciliation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `school_slasc`
--

DROP TABLE IF EXISTS `school_slasc`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `school_slasc` (
  `School_SLASCID` bigint NOT NULL AUTO_INCREMENT,
  `School_accommodation_change_indicator` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `School_accommodation_change_details` varchar(4000) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Lowest_annual_fee_day_pupils` int DEFAULT '0',
  `Lowest_annual_fee_boarding_pupils` int DEFAULT '0',
  `Highest_annual_fee_day_pupils` int DEFAULT '0',
  `Highest_annual_fee_boarding_pupils` int DEFAULT '0',
  `Approved_places` int DEFAULT '0',
  `Boarding_pupils_reference_year_minus_2` int DEFAULT '0',
  `Boarding_pupils_reference_year_minus_1` int DEFAULT '0',
  `Boarding_pupils_forthcomming_year` int DEFAULT '0',
  `Accommodation_295_Days` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`School_SLASCID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `school_slasc`
--

LOCK TABLES `school_slasc` WRITE;
/*!40000 ALTER TABLE `school_slasc` DISABLE KEYS */;
/*!40000 ALTER TABLE `school_slasc` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `school_survey`
--

DROP TABLE IF EXISTS `school_survey`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `school_survey` (
  `School_surveyID` bigint NOT NULL AUTO_INCREMENT,
  `Contact_Telephone_Number` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Collection_Name` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Survey_Year` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Survey_Reference_Date` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Date_and_Time_when_the_message_file_is_produced` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Survey_Contact` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Headteacher_name` varchar(70) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Charity_name` varchar(70) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Charity_number` varchar(7) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Survey_Term` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Source_Level` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Software_Code` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Serial_No` int DEFAULT '0',
  `Release` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `X_Version` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `CBDS_Level` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Type_of_CTF_file` varchar(7) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Source_LEA` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Destination_LEA` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Data_descriptor_for_Partial_CTF` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Destination_School` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Source_School` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Supplier_ID` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Academic_Year_of_Transfer` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Type_of_Partial_CTF` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`School_surveyID`),
  KEY `WDIDX_school_survey_Collection_Name` (`Collection_Name`),
  KEY `WDIDX_school_survey_Survey_Year` (`Survey_Year`),
  KEY `WDIDX_school_survey_CBDS_Level` (`CBDS_Level`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `school_survey`
--

LOCK TABLES `school_survey` WRITE;
/*!40000 ALTER TABLE `school_survey` DISABLE KEYS */;
/*!40000 ALTER TABLE `school_survey` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `school_survey_staff_details`
--

DROP TABLE IF EXISTS `school_survey_staff_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `school_survey_staff_details` (
  `Staff_DetailsID` bigint NOT NULL AUTO_INCREMENT,
  `Teacher_full_time_men` int DEFAULT '0',
  `Teacher_full_time_women` int DEFAULT '0',
  `Teacher_part_time_men` int DEFAULT '0',
  `Teacher_part_time_men_hours` int DEFAULT '0',
  `Teacher_part_time_women` int DEFAULT '0',
  `Teacher_part_time_women_hours` int DEFAULT '0',
  `Newly_appointed_teachers_nil_return` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Newly_appointed_support_staff_nil_return` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Individual_proprietor_nil_return` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Proprietor_body_nil_return` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Newly_appointed_body_member_nil_return` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Post_held_indicator` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Leaving_teachers_nil_return` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Leaving_support_staff_nil_return` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Post_title` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Support_Head_Count` int DEFAULT '0',
  `Occasionals_QTS` int DEFAULT '0',
  `Occasionals_NOTQTS` int DEFAULT '0',
  `Occasionals_NOTKNWN` int DEFAULT '0',
  `Total_teaching_staff_at_establishment` int DEFAULT '0',
  `Category_of_Agency_TP_Support_staff` varchar(7) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Estab_Level_2` int DEFAULT '0',
  `Estab_Level_2_Math` int DEFAULT '0',
  `Estab_Level_2_English` int DEFAULT '0',
  `Estab_Level_3_Math` int DEFAULT '0',
  `Estab_Level_3_English` int DEFAULT '0',
  `Level_3_pre_September_2014` int DEFAULT '0',
  `Level_3_post_September_2014` int DEFAULT '0',
  `Total_staff_at_provider` int DEFAULT '0',
  `Number_of_Level_2_staff` int DEFAULT '0',
  `Number_of_Level_3_staff_non_managerial` int DEFAULT '0',
  `Number_of_Level_3_staff_managerial` int DEFAULT '0',
  `Numb_of_QTSstaff_who_work_with_children_under_5` int DEFAULT '0',
  `Numb_of_staff_with_EYTS_work_with_children_under_5` int DEFAULT '0',
  `Number_of_agency_workers_covering_vacancies_FTE` int DEFAULT '0',
  `Number_of_agency_workers_covering_vacancies_headcount` int DEFAULT '0',
  `Number_of_vacancies` int DEFAULT '0',
  `Reason_for_absence` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Staff_DetailsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `school_survey_staff_details`
--

LOCK TABLES `school_survey_staff_details` WRITE;
/*!40000 ALTER TABLE `school_survey_staff_details` DISABLE KEYS */;
/*!40000 ALTER TABLE `school_survey_staff_details` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `schools_connections`
--

DROP TABLE IF EXISTS `schools_connections`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `schools_connections` (
  `Schools_ConnectionsID` int NOT NULL AUTO_INCREMENT,
  `School_LA_Number` varchar(45) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `School_Name` varchar(45) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `School_DfE_Establishment_Number` varchar(45) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`Schools_ConnectionsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `schools_connections`
--

LOCK TABLES `schools_connections` WRITE;
/*!40000 ALTER TABLE `schools_connections` DISABLE KEYS */;
/*!40000 ALTER TABLE `schools_connections` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `seating_plans`
--

DROP TABLE IF EXISTS `seating_plans`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `seating_plans` (
  `Seating_plansID` bigint NOT NULL AUTO_INCREMENT,
  `Seating_arrangement` varchar(600) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '0',
  `RoomsID` bigint DEFAULT '0',
  `GroupID` bigint DEFAULT '0',
  PRIMARY KEY (`Seating_plansID`),
  KEY `WDIDX_seating_plans_roomsID` (`RoomsID`),
  KEY `WDIDX_seating_plans_groupsID` (`GroupID`),
  KEY `WDIDX_seating_plans_WDIDX_seating_plans_WDIDX_seating_plans00057` (`RoomsID`,`GroupID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `seating_plans`
--

LOCK TABLES `seating_plans` WRITE;
/*!40000 ALTER TABLE `seating_plans` DISABLE KEYS */;
/*!40000 ALTER TABLE `seating_plans` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `session_clash`
--

DROP TABLE IF EXISTS `session_clash`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `session_clash` (
  `Session_clashID` bigint NOT NULL AUTO_INCREMENT,
  `AttendanceID` bigint DEFAULT '0',
  `Off_timetable_eventID` bigint DEFAULT '0',
  PRIMARY KEY (`Session_clashID`),
  KEY `WDIDX_session_clash_AttendanceID` (`AttendanceID`),
  KEY `WDIDX_session_clash_Off_timetable_eventID` (`Off_timetable_eventID`),
  KEY `WDIDX_session_clash_WDIDX_session_clash_WDIDX_session_clash00058` (`AttendanceID`,`Off_timetable_eventID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `session_clash`
--

LOCK TABLES `session_clash` WRITE;
/*!40000 ALTER TABLE `session_clash` DISABLE KEYS */;
/*!40000 ALTER TABLE `session_clash` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `session_ssl_status`
--

DROP TABLE IF EXISTS `session_ssl_status`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `session_ssl_status` (
  `thread_id` bigint NOT NULL DEFAULT '0',
  `ssl_version` varchar(1024) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ssl_cipher` varchar(1024) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ssl_sessions_reused` varchar(1024) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `session_ssl_status`
--

LOCK TABLES `session_ssl_status` WRITE;
/*!40000 ALTER TABLE `session_ssl_status` DISABLE KEYS */;
/*!40000 ALTER TABLE `session_ssl_status` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sessions`
--

DROP TABLE IF EXISTS `sessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sessions` (
  `SessionsID` bigint NOT NULL AUTO_INCREMENT,
  `Period_Id` bigint DEFAULT '0',
  `DateAndStart` bigint DEFAULT NULL,
  `Teaching_BlocksID` bigint DEFAULT '0',
  `DateAndEnd` bigint DEFAULT NULL,
  PRIMARY KEY (`SessionsID`),
  KEY `WDIDX_sessions_Period_Id` (`Period_Id`),
  KEY `WDIDX_sessions_DateAndStart` (`DateAndStart`),
  KEY `WDIDX_sessions_Teaching_BlocksID` (`Teaching_BlocksID`),
  KEY `WDIDX_sessions_DateAndEnd` (`DateAndEnd`),
  KEY `WDIDX_sessions_WDIDX_sessions_WDIDX_sessions_WDIDX_sessions00059` (`Period_Id`,`DateAndStart`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sessions`
--

LOCK TABLES `sessions` WRITE;
/*!40000 ALTER TABLE `sessions` DISABLE KEYS */;
/*!40000 ALTER TABLE `sessions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `skill_attached`
--

DROP TABLE IF EXISTS `skill_attached`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `skill_attached` (
  `Skill_attachedID` bigint NOT NULL AUTO_INCREMENT,
  `Skills_keyID` bigint DEFAULT '0',
  `Lesson_FilesID` bigint DEFAULT '0',
  `CourseID` bigint DEFAULT '0',
  UNIQUE KEY `Skill_attachedID` (`Skill_attachedID`),
  KEY `WDIDX_skill_attached_Skills_keyID` (`Skills_keyID`),
  KEY `WDIDX_skill_attached_Lesson_FilesID` (`Lesson_FilesID`),
  KEY `WDIDX_skill_attached_CourseID` (`CourseID`),
  KEY `WDIDX_skill_attached_WDIDX_skill_attached_WDIDX_skill_attac00060` (`Skills_keyID`,`Lesson_FilesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `skill_attached`
--

LOCK TABLES `skill_attached` WRITE;
/*!40000 ALTER TABLE `skill_attached` DISABLE KEYS */;
INSERT INTO `skill_attached` (`Skills_keyID`, `Lesson_FilesID`, `CourseID`) VALUES
(1, 1, 3),
(2, 1, 3),
(3, 1, 3),
(4, 1, 3),
(5, 1, 3),
(6, 1, 3),
(7, 1, 3),
(8, 1, 3),
(9, 1, 3),
(10, 1, 3);
/*!40000 ALTER TABLE `skill_attached` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `skill_link`
--

DROP TABLE IF EXISTS `skill_link`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `skill_link` (
  `Skill_linkID` bigint NOT NULL AUTO_INCREMENT,
  `ParentSkillID` bigint DEFAULT '0',
  `OffspringSkillID` bigint DEFAULT '0',
  `CourseID` bigint DEFAULT '0',
  UNIQUE KEY `Skill_linkID` (`Skill_linkID`),
  KEY `WDIDX_skill_link_ParentSkillID` (`ParentSkillID`),
  KEY `WDIDX_skill_link_OffspringSkillID` (`OffspringSkillID`),
  KEY `WDIDX_skill_link_CourseID` (`CourseID`),
  KEY `WDIDX_skill_link_WDIDX_skill_link_WDIDX_skill_link_WDIDX_Sk00061` (`ParentSkillID`,`OffspringSkillID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `skill_link`
--

LOCK TABLES `skill_link` WRITE;
/*!40000 ALTER TABLE `skill_link` DISABLE KEYS */;
/*!40000 ALTER TABLE `skill_link` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `skills_key`
--

DROP TABLE IF EXISTS `skills_key`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `skills_key` (
  `Skills_keyID` bigint NOT NULL AUTO_INCREMENT,
  `Skill_name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Criteria1` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `Criteria2` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL  DEFAULT '',
  `Criteria3` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `Criteria4` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `Hidden` tinyint NOT NULL DEFAULT '0',
  `Criteria5` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `DevelopmentAge` double NOT NULL DEFAULT '0',
  `SkillDescription` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Skills_keyID`),
  KEY `WDIDX_skills_key_Skill_name` (`Skill_name`),
  KEY `WDIDX_skills_key_hidden` (`Hidden`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `skills_key`
--

LOCK TABLES `skills_key` WRITE;
/*!40000 ALTER TABLE `skills_key` DISABLE KEYS */;
/*!40000 ALTER TABLE `skills_key` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `skillsenable`
--

DROP TABLE IF EXISTS `skillsenable`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `skillsenable` (
  `SkillsEnableID` bigint NOT NULL AUTO_INCREMENT,
  `Activity` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Skills_keyID` bigint DEFAULT '0',
  PRIMARY KEY (`SkillsEnableID`),
  KEY `WDIDX_skillsenable_Skills_keyID` (`Skills_keyID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `skillsenable`
--

LOCK TABLES `skillsenable` WRITE;
/*!40000 ALTER TABLE `skillsenable` DISABLE KEYS */;
/*!40000 ALTER TABLE `skillsenable` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `skillsmapelements`
--

DROP TABLE IF EXISTS `skillsmapelements`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `skillsmapelements` (
  `SkillsMapElementID` bigint NOT NULL AUTO_INCREMENT,
  `CourseID` bigint DEFAULT '0',
  `Skills_keyID` bigint NOT NULL DEFAULT '0',
  `Skill_linkID` bigint NOT NULL DEFAULT '0',
  `OffspringSkillID` bigint DEFAULT '0',
  `PointAx` int DEFAULT '0',
  `PointAy` int DEFAULT '0',
  PRIMARY KEY (`SkillsMapElementID`),
  KEY `WDIDX_skillsmapelements_CourseID` (`CourseID`),
  KEY `WDIDX_skillsmapelements_Skills_keyID` (`Skills_keyID`),
  KEY `WDIDX_sskillsmapelements_Skill_linkID` (`Skill_linkID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `skillsmapelements`
--

LOCK TABLES `skillsmapelements` WRITE;
/*!40000 ALTER TABLE `skillsmapelements` DISABLE KEYS */;
/*!40000 ALTER TABLE `skillsmapelements` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `skillsmarkbook`
--

DROP TABLE IF EXISTS `skillsmarkbook`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;

CREATE TABLE `skillsmarkbook` (
  `SkillsMarkbookID` bigint NOT NULL AUTO_INCREMENT,
  `MarkbookID` bigint DEFAULT '0',
  `Skills_keyID` bigint DEFAULT '0',
  `Mark` int DEFAULT '0', -- 1-5 Proficiency Score
  PRIMARY KEY (`SkillsMarkbookID`),
  KEY `WDIDX_skillsmarkbook_MarkbookID` (`MarkbookID`),
  KEY `WDIDX_skillsmarkbook_Skills_keyID` (`Skills_keyID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `skillsmarkbook`
--

LOCK TABLES `skillsmarkbook` WRITE;
/*!40000 ALTER TABLE `skillsmarkbook` DISABLE KEYS */;
/*!40000 ALTER TABLE `skillsmarkbook` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `skillsreference`
--

DROP TABLE IF EXISTS `skillsreference`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `skillsreference` (
  `SkillsReferenceID` bigint NOT NULL AUTO_INCREMENT,
  `Skills_keyID` bigint DEFAULT '0',
  `File_TypesID` bigint DEFAULT '0',
  `Name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`SkillsReferenceID`),
  KEY `WDIDX_skillsreference_Skills_keyID` (`Skills_keyID`),
  KEY `WDIDX_skillsreference_File_TypesID` (`File_TypesID`),
  KEY `WDIDX_skillsreference_Name` (`Name`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `skillsreference`
--

LOCK TABLES `skillsreference` WRITE;
/*!40000 ALTER TABLE `skillsreference` DISABLE KEYS */;
/*!40000 ALTER TABLE `skillsreference` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `staff_details`
--

DROP TABLE IF EXISTS `staff_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `staff_details` (
  `Staff_DetailsID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `National_Insurance_Number` varchar(9) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Teacher_Number` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `QT_Status` tinyint DEFAULT '0',
  `HLTA_Status` tinyint DEFAULT '0',
  `QTS_Route` varchar(5) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Ethnic_Code` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Disability` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `HCPC_identifier` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Qualifying_institution` varchar(300) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Qualification_level` varchar(9) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Step_up_graduate` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Role_within_organisation` int DEFAULT '0',
  `Social_worker_origin` int DEFAULT '0',
  `Destination_of_leaver` int DEFAULT '0',
  `FTE_as_at_30_September_previous_census_year` decimal(24,6) DEFAULT '0.000000',
  `Number_of_cases_held_at_30_September` int DEFAULT '0',
  `Agency_worker_length_of_contract` int DEFAULT '0',
  `Reason_for_leaving` int DEFAULT '0',
  `Frontline_graduate` int DEFAULT '0',
  `Absent_on_30_September` int DEFAULT '0',
  `Agency_worker_indicator` int DEFAULT '0',
  PRIMARY KEY (`Staff_DetailsID`),
  KEY `WDIDX_staff_details_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_staff_details_National_Insurance_Number` (`National_Insurance_Number`),
  KEY `WDIDX_staff_details_Teacher_Number` (`Teacher_Number`),
  KEY `WDIDX_staff_details_QT_Status` (`QT_Status`),
  KEY `WDIDX_staff_details_HLTA_Status` (`HLTA_Status`),
  KEY `WDIDX_staff_details_QTS_Route` (`QTS_Route`),
  KEY `WDIDX_staff_details_Ethnic_Code` (`Ethnic_Code`),
  KEY `WDIDX_staff_details_Disability` (`Disability`),
  KEY `WDIDX_staff_details_HCPC_identifier` (`HCPC_identifier`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `staff_details`
--

LOCK TABLES `staff_details` WRITE;
/*!40000 ALTER TABLE `staff_details` DISABLE KEYS */;
/*!40000 ALTER TABLE `staff_details` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `staff_vacancies`
--

DROP TABLE IF EXISTS `staff_vacancies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `staff_vacancies` (
  `Staff_VacanciesID` bigint NOT NULL AUTO_INCREMENT,
  `Vacancy_Tenure` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Vacancy_Temporarily_Filled` tinyint DEFAULT '0',
  `Vacancy_Advertised` tinyint DEFAULT '0',
  `Vacancy_Subject` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Vacancy_Post` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Staff_VacanciesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `staff_vacancies`
--

LOCK TABLES `staff_vacancies` WRITE;
/*!40000 ALTER TABLE `staff_vacancies` DISABLE KEYS */;
/*!40000 ALTER TABLE `staff_vacancies` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `standard_structure`
--

DROP TABLE IF EXISTS `standard_structure`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `standard_structure` (
  `Tier2` varchar(15) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Tier3` varchar(15) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Tier4` varchar(15) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Tier5` varchar(15) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Non_timetable` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `standard_structure`
--

LOCK TABLES `standard_structure` WRITE;
/*!40000 ALTER TABLE `standard_structure` DISABLE KEYS */;
/*!40000 ALTER TABLE `standard_structure` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `standardizedmemo`
--

DROP TABLE IF EXISTS `standardizedmemo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `standardizedmemo` (
  `StandardizedMemoID` bigint NOT NULL AUTO_INCREMENT,
  `MemoName` varchar(50) DEFAULT NULL,
  `Type` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`StandardizedMemoID`),
  KEY `WDIDX_StandardizedMemo_MemoName` (`MemoName`),
  KEY `WDIDX_StandardizedMemo_Type` (`Type`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `standardizedmemo`
--

LOCK TABLES `standardizedmemo` WRITE;
/*!40000 ALTER TABLE `standardizedmemo` DISABLE KEYS */;
/*!40000 ALTER TABLE `standardizedmemo` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `standardizereplacements`
--

DROP TABLE IF EXISTS `standardizereplacements`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `standardizereplacements` (
  `StandardizeReplacementsID` bigint NOT NULL AUTO_INCREMENT,
  `FormReplacement` varchar(100) DEFAULT NULL,
  `FemaleForm` varchar(50) DEFAULT NULL,
  `MaleForm` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`StandardizeReplacementsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `standardizereplacements`
--

LOCK TABLES `standardizereplacements` WRITE;
/*!40000 ALTER TABLE `standardizereplacements` DISABLE KEYS */;
/*!40000 ALTER TABLE `standardizereplacements` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `std_homework`
--

DROP TABLE IF EXISTS `std_homework`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `std_homework` (
  `HomeworkName` varchar(250) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `std_homework`
--

LOCK TABLES `std_homework` WRITE;
/*!40000 ALTER TABLE `std_homework` DISABLE KEYS */;
/*!40000 ALTER TABLE `std_homework` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `subject`
--

DROP TABLE IF EXISTS `subject`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `subject` (
  `SubjectID` bigint NOT NULL AUTO_INCREMENT,
  `Subject_name` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Subject_code` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `LocalRootDir` bigint DEFAULT '0',
  `Subject_details` TEXT DEFAULT NULL,
  `Hidden` tinyint DEFAULT '0',
  `CloudFolderID` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `LocalFoldersID` bigint DEFAULT '0',
  PRIMARY KEY (`SubjectID`),
  KEY `WDIDX_subject_Subject_name` (`Subject_name`),
  KEY `WDIDX_subject_subjectcode` (`Subject_code`),
  KEY `WDIDX_subject_LocalRootDir` (`LocalRootDir`),
  KEY `WDIDX_subject_LocalFoldersID` (`LocalFoldersID`),
  KEY `WDIDX_subject_WDIDX_subject_WDIDX_Subject_OptimCompKey_subj00062` (`Subject_name`,`SubjectID`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `subject`
--





--
-- Table structure for table `courses`
--

DROP TABLE IF EXISTS `courses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `courses` (
  `CourseID` bigint NOT NULL AUTO_INCREMENT,
  `Course_Name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `Course_Code` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `SubjectID` bigint DEFAULT '0',
  `Details` TEXT CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `CloudFolderID` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '',
  `LocalFoldersID` bigint DEFAULT '0',
  `Node_size` int DEFAULT '0',
  `Hidden` int DEFAULT '0',
  PRIMARY KEY (`CourseID`),
  KEY `WDIDX_Courses_Course_Name` (`Course_Name`),
  KEY `WDIDX_Courses_Course_Code` (`Course_Code`),
  KEY `WDIDX_Courses_SubjectID` (`SubjectID`),
  KEY `WDIDX_Courses_LocalFoldersID` (`LocalFoldersID`),
  KEY `WDIDX_Courses_Hidden` (`Hidden`),
  KEY `WDIDX_Courses_CloudFolderID` (`CloudFolderID`),
  KEY `WDIDX_Courses_WDIDX_Courses_Name_SubjectID` (`Course_Name`,`SubjectID`)
) ENGINE=MyISAM AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;



DROP TABLE IF EXISTS `topics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
DROP TABLE IF EXISTS `topics`;
CREATE TABLE `topics` (
  `TopicID` bigint NOT NULL AUTO_INCREMENT,
  `Topic_name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `TopicLong` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Topic_code` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `ModuleID` bigint DEFAULT NULL,
  `Long_desc` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `CloudFolderID` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `Folder_id` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `Topic_block_tableID` bigint DEFAULT '0',
  `Hidden` int DEFAULT '0',
  `SpecificationRefID` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`TopicID`),
  KEY `WDIDX_topics_topicName` (`Topic_name`),
  KEY `WDIDX_topics_topicLong` (`TopicLong`),
  KEY `WDIDX_topics_Topic_code` (`Topic_code`),
  KEY `WDIDX_topics_ModuleID` (`ModuleID`),
  KEY `WDIDX_topics_CloudFolderID` (`CloudFolderID`),
  KEY `WDIDX_topics_folder_id` (`Folder_id`),
  KEY `WDIDX_topics_Topic_block_tableID` (`Topic_block_tableID`),
  KEY `WDIDX_topics_hidden` (`Hidden`),
  KEY `WDIDX_topics_WDIDX_topics_WDIDX_Topics_topicNameCourseID` (`Topic_name`,`ModuleID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;


DROP TABLE IF EXISTS `modules`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
-- Fix the modules table PRIMARY KEY
CREATE TABLE `modules` (
  `ModuleID` bigint NOT NULL AUTO_INCREMENT,
  `Module_name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `ModuleLong` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Module_code` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `CourseID` bigint NOT NULL DEFAULT '0',
  `Module_details` TEXT CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `CloudFolderID` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `Module_block_tableID` bigint NOT NULL DEFAULT '0',
  `Hidden` int DEFAULT '0',
  `SpecificationRefID` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`ModuleID`),
  KEY `WDIDX_modules_moduleName` (`Module_name`),
  KEY `WDIDX_modules_moduleLong` (`ModuleLong`),
  KEY `WDIDX_modules_module_code` (`Module_code`),
  KEY `WDIDX_modules_CourseID` (`CourseID`),
  KEY `WDIDX_modules_CloudFolderID` (`CloudFolderID`),
  KEY `WDIDX_modules_Module_block_tableID` (`Module_block_tableID`),
  KEY `WDIDX_modules_hidden` (`Hidden`),
  KEY `WDIDX_modules_WDIDX_modules_WDIDX_Modules_moduleNameCourseID` (`Module_name`,`CourseID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

/*!40101 SET character_set_client = @saved_cs_client */;



DROP TABLE IF EXISTS `sys_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_config` (
  `variable` varchar(128) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `value` varchar(128) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `set_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `set_by` varchar(128) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  UNIQUE KEY `variable` (`variable`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;


DROP TABLE IF EXISTS `targets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `targets` (
  `TargetsID` bigint NOT NULL AUTO_INCREMENT,
  `DateSet` bigint DEFAULT NULL,
  `SubjectID` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  `Target` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  UNIQUE KEY `TargetsID` (`TargetsID`),
  KEY `WDIDX_targets_SubjectID` (`SubjectID`),
  KEY `WDIDX_targets_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_targets_Target` (`Target`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;



DROP TABLE IF EXISTS `teachgroups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `teachgroups` (
  `GroupID` bigint NOT NULL AUTO_INCREMENT,
  `GroupEmail` varchar(100) DEFAULT NULL,
  `Groupname` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `EntitiesID` bigint DEFAULT '0',
  `LeaderID` bigint DEFAULT '0',
  `MicrosoftID` varchar(100) DEFAULT NULL,
  `ThreadID` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`GroupID`),
  KEY `WDIDX_teachgroups_groupname` (`Groupname`),
  KEY `WDIDX_teachgroups_EntitiesID` (`EntitiesID`),
  KEY `WDIDX_teachgroups_leaderID` (`LeaderID`),
  KEY `WDIDX_teachgroups_MicrosoftID` (`MicrosoftID`),
  KEY `WDIDX_teachgroups_groupname_EntitiesID` (`Groupname`,`EntitiesID`),
  KEY `WDIDX_teachgroups_groupname_leaderID` (`Groupname`,`LeaderID`)
) ENGINE=MyISAM AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `teachgroups`
--

LOCK TABLES `teachgroups` WRITE;
-- /*!40000 ALTER TABLE `teachgroups` DISABLE KEYS */;

-- /*!40000 ALTER TABLE `teachgroups` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `teaching_blocks`
--

DROP TABLE IF EXISTS `teaching_blocks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `teaching_blocks` (
  `Teaching_BlocksID` bigint NOT NULL AUTO_INCREMENT,
  `StartDate` bigint DEFAULT NULL,
  `FinishDate` bigint DEFAULT NULL,
  `Name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AppliesTo` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Academic_YearID` bigint DEFAULT '0',
  PRIMARY KEY (`Teaching_BlocksID`),
  KEY `WDIDX_teaching_blocks_StartDate` (`StartDate`),
  KEY `WDIDX_teaching_blocks_FinishDate` (`FinishDate`),
  KEY `WDIDX_teaching_blocks_Name` (`Name`),
  KEY `WDIDX_teaching_blocks_appliesTo` (`AppliesTo`),
  KEY `WDIDX_teaching_blocks_Academic_YearID` (`Academic_YearID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `teaching_blocks`
--

LOCK TABLES `teaching_blocks` WRITE;
/*!40000 ALTER TABLE `teaching_blocks` DISABLE KEYS */;
/*!40000 ALTER TABLE `teaching_blocks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `telephonecontact`
--

DROP TABLE IF EXISTS `telephonecontact`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `telephonecontact` (
  `TelephoneContactID` bigint NOT NULL AUTO_INCREMENT,
  `PhoneNumber` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `PhoneType` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `ContactsID` bigint DEFAULT '0',
  `IndividualsID` bigint DEFAULT '0',
  PRIMARY KEY (`TelephoneContactID`),
  KEY `WDIDX_telephonecontact_PhoneNumber` (`PhoneNumber`),
  KEY `WDIDX_telephonecontact_ContactsID` (`ContactsID`),
  KEY `WDIDX_telephonecontact_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `telephonecontact`
--

LOCK TABLES `telephonecontact` WRITE;
/*!40000 ALTER TABLE `telephonecontact` DISABLE KEYS */;
/*!40000 ALTER TABLE `telephonecontact` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `timetable_entries`
--

DROP TABLE IF EXISTS `timetable_entries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `timetable_entries` (
  `Timetable_entriesID` bigint NOT NULL AUTO_INCREMENT,
  `Period_Id` bigint DEFAULT '0',
  `GroupID` bigint DEFAULT '0',
  `RoomsID` bigint DEFAULT '0',
  `SubjectID` bigint DEFAULT '0',
  `Academic_YearID` bigint DEFAULT '0',
  PRIMARY KEY (`Timetable_entriesID`),
  KEY `WDIDX_timetable_entries_Period_Id` (`Period_Id`),
  KEY `WDIDX_timetable_entries_groupsID` (`GroupID`),
  KEY `WDIDX_timetable_entries_roomsID` (`RoomsID`),
  KEY `WDIDX_timetable_entries_SubjectID` (`SubjectID`),
  KEY `WDIDX_timetable_entries_Academic_YearID` (`Academic_YearID`),
  KEY `WDIDX_timetable_entries_WDIDX_timetable_entries_WDIDX_timet00063` (`Period_Id`,`GroupID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `timetable_entries`
--

LOCK TABLES `timetable_entries` WRITE;
/*!40000 ALTER TABLE `timetable_entries` DISABLE KEYS */;
/*!40000 ALTER TABLE `timetable_entries` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `topic_plan_header`
--
DROP TABLE IF EXISTS `enrollments`;

CREATE TABLE `enrollments` (
  `EnrollmentID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint NOT NULL DEFAULT '0',
  `CourseID` bigint NOT NULL DEFAULT '0',
  `start_date` TIMESTAMP NULL DEFAULT NULL,
  `end_date` TIMESTAMP NULL DEFAULT NULL,
  `status` VARCHAR(20) DEFAULT 'active',
  PRIMARY KEY (`EnrollmentID`),
  KEY `idx_individual` (`IndividualsID`),
  KEY `idx_course` (`CourseID`),
  KEY `idx_dates` (`start_date`, `end_date`),
  KEY `idx_status` (`status`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;
--
-- Dumping data for table `enrollments`
LOCK TABLES `enrollments` WRITE;
INSERT INTO `enrollments` (`IndividualsID`, `CourseID`, `start_date`, `end_date`, `status`) VALUES
(1, 3, '2024-01-15 09:00:00', NULL, 'active');
UNLOCK TABLES;


DROP TABLE IF EXISTS `topic_plan_header`;
CREATE TABLE `topic_plan_header` (
  `Topic_plan_headerID` bigint NOT NULL AUTO_INCREMENT,
  `Groupname` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `TopicID` bigint DEFAULT '0',
  `StartDate` bigint DEFAULT '0',
  `EndDate` bigint DEFAULT '0',
  PRIMARY KEY (`Topic_plan_headerID`),
  KEY `WDIDX_topic_plan_header_groupname` (`Groupname`),
  KEY `WDIDX_topic_plan_header_TopicID` (`TopicID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `topic_plan_header`
--

LOCK TABLES `topic_plan_header` WRITE;
/*!40000 ALTER TABLE `topic_plan_header` DISABLE KEYS */;
/*!40000 ALTER TABLE `topic_plan_header` ENABLE KEYS */;
UNLOCK TABLES;



--
-- Table structure for table `tt_clashes`
--

DROP TABLE IF EXISTS `tt_clashes`;
CREATE TABLE `tt_clashes` (
  `TT_clashesID` bigint NOT NULL AUTO_INCREMENT,
  `Timetable_entriesID` bigint DEFAULT '0',
  `field` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`TT_clashesID`),
  KEY `WDIDX_tt_clashes_timetable_entriesID` (`Timetable_entriesID`),
  KEY `WDIDX_tt_clashes_field` (`field`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tt_clashes`
--

LOCK TABLES `tt_clashes` WRITE;
/*!40000 ALTER TABLE `tt_clashes` DISABLE KEYS */;
/*!40000 ALTER TABLE `tt_clashes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `types_key`
--

DROP TABLE IF EXISTS `types_key`;
CREATE TABLE `types_key` (
  `Types_keyID` bigint NOT NULL AUTO_INCREMENT,
  `Name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Description` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`Types_keyID`),
  KEY `WDIDX_types_key_Name` (`Name`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `types_key`
--

LOCK TABLES `types_key` WRITE;
/*!40000 ALTER TABLE `types_key` DISABLE KEYS */;
/*!40000 ALTER TABLE `types_key` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `weblinks`
--

DROP TABLE IF EXISTS `weblinks`;
CREATE TABLE `weblinks` (
  `WeblinksID` bigint NOT NULL AUTO_INCREMENT,
  `SubjectID` bigint DEFAULT '0',
  `URL` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `LinkName` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Description` varchar(300) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`WeblinksID`),
  KEY `WDIDX_weblinks_SubjectID` (`SubjectID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `weblinks`
--

LOCK TABLES `weblinks` WRITE;
/*!40000 ALTER TABLE `weblinks` DISABLE KEYS */;
/*!40000 ALTER TABLE `weblinks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `weekcalcs`
--

DROP TABLE IF EXISTS `weekcalcs`;
CREATE TABLE `weekcalcs` (
  `WeekCalcsID` bigint NOT NULL AUTO_INCREMENT,
  `Teaching_BlocksID` bigint DEFAULT '0',
  `WeekNumber` int DEFAULT '0',
  `WeekType` varchar(1) DEFAULT NULL,
  PRIMARY KEY (`WeekCalcsID`),
  KEY `WDIDX_weekcalcs_Teaching_BlocksID` (`Teaching_BlocksID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `weekcalcs`
--

LOCK TABLES `weekcalcs` WRITE;
/*!40000 ALTER TABLE `weekcalcs` DISABLE KEYS */;
/*!40000 ALTER TABLE `weekcalcs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `wifinetwork`
--

DROP TABLE IF EXISTS `wifinetwork`;
CREATE TABLE `wifinetwork` (
  `WiFiNetworkID` bigint NOT NULL AUTO_INCREMENT,
  `WifiNameSSID` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Password` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`WiFiNetworkID`),
  KEY `WDIDX_wifinetwork_WifiNameSSID` (`WifiNameSSID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;


--
-- Dumping data for table `wifinetwork`
--

LOCK TABLES `wifinetwork` WRITE;
/*!40000 ALTER TABLE `wifinetwork` DISABLE KEYS */;
/*!40000 ALTER TABLE `wifinetwork` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `workassignment`
--




DROP TABLE IF EXISTS `workassignment`;
CREATE TABLE `workassignment` (
  `WorkAssignmentID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `Lesson_FilesID` bigint DEFAULT '0',
  `DateCollected` datetime DEFAULT NULL,
  `AssignedDate` datetime DEFAULT NULL,
  `DueDate` datetime DEFAULT NULL,
  `Status` varchar(1) DEFAULT 'N',
  PRIMARY KEY (`WorkAssignmentID`),
  KEY `WDIDX_workassignment_IndividualsID` (`IndividualsID`),
  KEY `WDIDX_workassignment_Lesson_FilesID` (`Lesson_FilesID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;


/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workassignment`
--

LOCK TABLES `workassignment` WRITE;
/*!40000 ALTER TABLE `workassignment` DISABLE KEYS */;
/*!40000 ALTER TABLE `workassignment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `workforce_absence`
--

DROP TABLE IF EXISTS `workforce_absence`;
CREATE TABLE `workforce_absence` (
  `Workforce_AbsenceID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `AbsenceDate` bigint DEFAULT '0',
  `AbsenceType` varchar(10) DEFAULT NULL,
  `Duration` decimal(4,2) DEFAULT '0.00',
  PRIMARY KEY (`Workforce_AbsenceID`),
  KEY `WDIDX_workforce_absence_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workforce_absence`
--

LOCK TABLES `workforce_absence` WRITE;
/*!40000 ALTER TABLE `workforce_absence` DISABLE KEYS */;
/*!40000 ALTER TABLE `workforce_absence` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `workforce_contracts`
--

DROP TABLE IF EXISTS `workforce_contracts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `workforce_contracts` (
  `Workforce_contractsID` bigint NOT NULL AUTO_INCREMENT,
  `Contract_Agreement_Type` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Start_Date` bigint DEFAULT NULL,
  `End_Date` bigint DEFAULT NULL,
  `Date_of_Arrival_in_School` bigint DEFAULT NULL,
  `Destination` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Hours_worked_per_week` int DEFAULT '0',
  `FTE_Hours_per_week` decimal(24,6) DEFAULT '0.000000',
  `Full_Time_Equivalence_in_Post` decimal(24,6) DEFAULT '0.000000',
  `Start_Date_in_Role` bigint DEFAULT NULL,
  `End_Date_in_Role` bigint DEFAULT NULL,
  `Safeguarded_Salary` tinyint DEFAULT '0',
  `Daily_Rate` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Additional_Payment_Amount` decimal(24,6) DEFAULT '0.000000',
  `Role_Identifier` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Origin` varchar(6) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Weeks_per_year` int DEFAULT '0',
  `Pay_Range` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pay_Framework` varchar(7) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Pay_Range_Minimum` decimal(24,6) DEFAULT '0.000000',
  `Pay_Range_Maximum` decimal(24,6) DEFAULT '0.000000',
  `Pay_Review_Date` bigint DEFAULT NULL,
  `Pay_Start_Date` bigint DEFAULT NULL,
  `Pay_End_Date` bigint DEFAULT NULL,
  `Category_of_Additional_Payment` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `PostName` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `LA_or_School_level` varchar(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Base_Pay` decimal(24,6) DEFAULT '0.000000',
  PRIMARY KEY (`Workforce_contractsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workforce_contracts`
--

LOCK TABLES `workforce_contracts` WRITE;
/*!40000 ALTER TABLE `workforce_contracts` DISABLE KEYS */;
/*!40000 ALTER TABLE `workforce_contracts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `workforce_curriculum`
--

DROP TABLE IF EXISTS `workforce_curriculum`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `workforce_curriculum` (
  `Workforce_curriculumID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `Hours` decimal(24,6) DEFAULT '0.000000',
  `Subject_Code` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Year_Group` varchar(2) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`Workforce_curriculumID`),
  KEY `WDIDX_workforce_curriculum_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workforce_curriculum`
--

LOCK TABLES `workforce_curriculum` WRITE;
/*!40000 ALTER TABLE `workforce_curriculum` DISABLE KEYS */;
/*!40000 ALTER TABLE `workforce_curriculum` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `workforce_qualifications`
--

DROP TABLE IF EXISTS `workforce_qualifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `workforce_qualifications` (
  `WorkforceID` bigint NOT NULL AUTO_INCREMENT,
  `IndividualsID` bigint DEFAULT '0',
  `Class_of_Degree` int DEFAULT '0',
  `Qualification_Code` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Subject_Code_1` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Subject_Code_2` varchar(4) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`WorkforceID`),
  KEY `WDIDX_workforce_qualifications_IndividualsID` (`IndividualsID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workforce_qualifications`
--

LOCK TABLES `workforce_qualifications` WRITE;
/*!40000 ALTER TABLE `workforce_qualifications` DISABLE KEYS */;
/*!40000 ALTER TABLE `workforce_qualifications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `workhierarchy`
--

DROP TABLE IF EXISTS `workhierarchy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `workhierarchy` (
  `WorkHierarchyID` bigint NOT NULL AUTO_INCREMENT,
  `GrID` varchar(8) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '0',
  `MemID` varchar(8) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `Tier` int DEFAULT '0',
  PRIMARY KEY (`WorkHierarchyID`),
  KEY `WDIDX_workhierarchy_GrID` (`GrID`),
  KEY `WDIDX_workhierarchy_MemID` (`MemID`),
  KEY `WDIDX_workhierarchy_tier` (`Tier`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workhierarchy`
--

LOCK TABLES `workhierarchy` WRITE;
/*!40000 ALTER TABLE `workhierarchy` DISABLE KEYS */;
/*!40000 ALTER TABLE `workhierarchy` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$host_summary`
--

DROP TABLE IF EXISTS `x$host_summary`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$host_summary` (
  `host` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `statements` double DEFAULT NULL,
  `statement_latency` double DEFAULT NULL,
  `statement_avg_latency` double DEFAULT NULL,
  `table_scans` double DEFAULT NULL,
  `file_ios` double DEFAULT NULL,
  `file_io_latency` double DEFAULT NULL,
  `current_connections` double DEFAULT NULL,
  `total_connections` double DEFAULT NULL,
  `unique_users` bigint NOT NULL DEFAULT '0',
  `current_memory` double DEFAULT NULL,
  `total_memory_allocated` double DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$host_summary`
--

LOCK TABLES `x$host_summary` WRITE;
/*!40000 ALTER TABLE `x$host_summary` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$host_summary` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$host_summary_by_file_io`
--

DROP TABLE IF EXISTS `x$host_summary_by_file_io`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$host_summary_by_file_io` (
  `host` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ios` double DEFAULT NULL,
  `io_latency` double DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$host_summary_by_file_io`
--

LOCK TABLES `x$host_summary_by_file_io` WRITE;
/*!40000 ALTER TABLE `x$host_summary_by_file_io` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$host_summary_by_file_io` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$host_summary_by_file_io_type`
--

DROP TABLE IF EXISTS `x$host_summary_by_file_io_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$host_summary_by_file_io_type` (
  `host` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `event_name` varchar(128) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `total` bigint NOT NULL DEFAULT '0',
  `total_latency` bigint NOT NULL DEFAULT '0',
  `max_latency` bigint NOT NULL DEFAULT '0'
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$host_summary_by_file_io_type`
--

LOCK TABLES `x$host_summary_by_file_io_type` WRITE;
/*!40000 ALTER TABLE `x$host_summary_by_file_io_type` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$host_summary_by_file_io_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$host_summary_by_stages`
--

DROP TABLE IF EXISTS `x$host_summary_by_stages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$host_summary_by_stages` (
  `host` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `event_name` varchar(128) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `total` bigint NOT NULL DEFAULT '0',
  `total_latency` bigint NOT NULL DEFAULT '0',
  `avg_latency` bigint NOT NULL DEFAULT '0'
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$host_summary_by_stages`
--

LOCK TABLES `x$host_summary_by_stages` WRITE;
/*!40000 ALTER TABLE `x$host_summary_by_stages` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$host_summary_by_stages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$host_summary_by_statement_latency`
--

DROP TABLE IF EXISTS `x$host_summary_by_statement_latency`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$host_summary_by_statement_latency` (
  `host` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `total` double DEFAULT NULL,
  `total_latency` double DEFAULT NULL,
  `max_latency` bigint DEFAULT NULL,
  `lock_latency` double DEFAULT NULL,
  `rows_sent` double DEFAULT NULL,
  `rows_examined` double DEFAULT NULL,
  `rows_affected` double DEFAULT NULL,
  `full_scans` double DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$host_summary_by_statement_latency`
--

LOCK TABLES `x$host_summary_by_statement_latency` WRITE;
/*!40000 ALTER TABLE `x$host_summary_by_statement_latency` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$host_summary_by_statement_latency` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$host_summary_by_statement_type`
--

DROP TABLE IF EXISTS `x$host_summary_by_statement_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$host_summary_by_statement_type` (
  `host` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `statement` varchar(128) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `total` bigint NOT NULL DEFAULT '0',
  `total_latency` bigint NOT NULL DEFAULT '0',
  `max_latency` bigint NOT NULL DEFAULT '0',
  `lock_latency` bigint NOT NULL DEFAULT '0',
  `rows_sent` bigint NOT NULL DEFAULT '0',
  `rows_examined` bigint NOT NULL DEFAULT '0',
  `rows_affected` bigint NOT NULL DEFAULT '0',
  `full_scans` bigint NOT NULL DEFAULT '0'
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$host_summary_by_statement_type`
--

LOCK TABLES `x$host_summary_by_statement_type` WRITE;
/*!40000 ALTER TABLE `x$host_summary_by_statement_type` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$host_summary_by_statement_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$innodb_buffer_stats_by_schema`
--

DROP TABLE IF EXISTS `x$innodb_buffer_stats_by_schema`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$innodb_buffer_stats_by_schema` (
  `object_schema` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `allocated` double DEFAULT NULL,
  `data` double DEFAULT NULL,
  `pages` bigint NOT NULL DEFAULT '0',
  `pages_hashed` bigint NOT NULL DEFAULT '0',
  `pages_old` bigint NOT NULL DEFAULT '0',
  `rows_cached` double NOT NULL DEFAULT '0'
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$innodb_buffer_stats_by_schema`
--

LOCK TABLES `x$innodb_buffer_stats_by_schema` WRITE;
/*!40000 ALTER TABLE `x$innodb_buffer_stats_by_schema` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$innodb_buffer_stats_by_schema` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$innodb_buffer_stats_by_table`
--

DROP TABLE IF EXISTS `x$innodb_buffer_stats_by_table`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$innodb_buffer_stats_by_table` (
  `object_schema` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `object_name` longtext CHARACTER SET ucs2 COLLATE ucs2_general_ci,
  `allocated` double DEFAULT NULL,
  `data` double DEFAULT NULL,
  `pages` bigint NOT NULL DEFAULT '0',
  `pages_hashed` bigint NOT NULL DEFAULT '0',
  `pages_old` bigint NOT NULL DEFAULT '0',
  `rows_cached` double NOT NULL DEFAULT '0'
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$innodb_buffer_stats_by_table`
--

LOCK TABLES `x$innodb_buffer_stats_by_table` WRITE;
/*!40000 ALTER TABLE `x$innodb_buffer_stats_by_table` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$innodb_buffer_stats_by_table` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$innodb_lock_waits`
--

DROP TABLE IF EXISTS `x$innodb_lock_waits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$innodb_lock_waits` (
  `wait_started` timestamp NULL DEFAULT NULL,
  `wait_age_secs` bigint DEFAULT NULL,
  `locked_table` varchar(1024) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `blocking_trx_age` time DEFAULT NULL,
  `locked_index` varchar(1024) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `locked_type` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `waiting_trx_id` varchar(18) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `waiting_trx_started` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `waiting_trx_rows_locked` bigint NOT NULL DEFAULT '0',
  `waiting_trx_rows_modified` bigint NOT NULL DEFAULT '0',
  `waiting_pid` bigint NOT NULL DEFAULT '0',
  `waiting_query` varchar(1024) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `waiting_lock_id` varchar(81) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `wait_age` time DEFAULT NULL,
  `waiting_lock_mode` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `blocking_trx_id` varchar(18) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `blocking_pid` bigint NOT NULL DEFAULT '0',
  `blocking_query` varchar(1024) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `blocking_lock_id` varchar(81) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `blocking_lock_mode` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `blocking_trx_started` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `waiting_trx_age` time DEFAULT NULL,
  `blocking_trx_rows_locked` bigint NOT NULL DEFAULT '0',
  `blocking_trx_rows_modified` bigint NOT NULL DEFAULT '0',
  `sql_kill_blocking_query` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `sql_kill_blocking_connection` varchar(26) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$innodb_lock_waits`
--

LOCK TABLES `x$innodb_lock_waits` WRITE;
/*!40000 ALTER TABLE `x$innodb_lock_waits` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$innodb_lock_waits` ENABLE KEYS */;
UNLOCK TABLES;

--
-- YouTube videos table - links to lesson_files for all context
DROP TABLE IF EXISTS `youtube_videos`;
CREATE TABLE `youtube_videos` (
  `YouTube_VideoID` bigint NOT NULL AUTO_INCREMENT,
  `VideoID` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT 'YouTube video ID from URL',
  `Duration` int DEFAULT 0 COMMENT 'Video duration in seconds',
  `Thumbnail_URL` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Channel_Name` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Published_Date` bigint DEFAULT NULL COMMENT 'Unix timestamp',
  `View_Count` bigint DEFAULT 0,
  `Lesson_FilesID` bigint NOT NULL DEFAULT 0 COMMENT 'Links to lesson_files for skill/course context',
  `Hidden` tinyint DEFAULT 0,
  `Created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `Updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`YouTube_VideoID`),
  UNIQUE KEY `VideoID` (`VideoID`),
  KEY `WDIDX_youtube_videos_VideoID` (`VideoID`),
  KEY `WDIDX_youtube_videos_Lesson_FilesID` (`Lesson_FilesID`),
  KEY `WDIDX_youtube_videos_Hidden` (`Hidden`),
  FOREIGN KEY (`Lesson_FilesID`) REFERENCES `lesson_files`(`Lesson_FilesID`) ON DELETE CASCADE
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='YouTube videos linked to lesson files';

-- YouTube videos table seed data (corrected structure)
LOCK TABLES `youtube_videos` WRITE;
INSERT INTO `youtube_videos` 
(`VideoID`, `Duration`, `Thumbnail_URL`, `Channel_Name`, `Published_Date`, `View_Count`, `Lesson_FilesID`, `Hidden`)
VALUES
  -- Strategic Foresight video (links to lesson_file 4)
  ('dQw4w9WgXcQ', 600, 'https://img.youtube.com/vi/dQw4w9WgXcQ/maxresdefault.jpg', 'Leadership Academy', UNIX_TIMESTAMP('2024-01-15'), 125000, 4, 0),
  
  -- Risk Intelligence video (links to lesson_file 5)
  ('jNQXAC9IVRw', 480, 'https://img.youtube.com/vi/jNQXAC9IVRw/maxresdefault.jpg', 'Business Skills Pro', UNIX_TIMESTAMP('2024-02-20'), 98000, 5, 0),
  
  -- Opportunity Recognition video (links to lesson_file 6)
  ('9bZkp7q19f0', 720, 'https://img.youtube.com/vi/9bZkp7q19f0/maxresdefault.jpg', 'Growth Mindset Channel', UNIX_TIMESTAMP('2024-03-10'), 156000, 6, 0),
  
  -- Intrapersonal Mastery video (links to lesson_file 7)
  ('kJQP7kiw5Fk', 540, 'https://img.youtube.com/vi/kJQP7kiw5Fk/maxresdefault.jpg', 'Self-Development Hub', UNIX_TIMESTAMP('2024-04-05'), 87000, 7, 0),
  
  -- Interpersonal Intelligence video (links to lesson_file 8)
  ('ScMzIvxBSi4', 660, 'https://img.youtube.com/vi/ScMzIvxBSi4/maxresdefault.jpg', 'Communication Mastery', UNIX_TIMESTAMP('2024-05-12'), 134000, 8, 0),
  
  -- Strategic Research video (links to lesson_file 9)
  ('oHg5SJYRHA0', 780, 'https://img.youtube.com/vi/oHg5SJYRHA0/maxresdefault.jpg', 'Research Methods Pro', UNIX_TIMESTAMP('2024-06-18'), 72000, 9, 0),
  
  -- Adaptive Planning video (links to lesson_file 10)
  ('fJ9rUzIMcZQ', 620, 'https://img.youtube.com/vi/fJ9rUzIMcZQ/maxresdefault.jpg', 'Planning Excellence', UNIX_TIMESTAMP('2024-07-22'), 91000, 10, 0),
  
  -- Integrative Leadership video (links to lesson_file 11)
  ('ZbZSe6N_BXs', 900, 'https://img.youtube.com/vi/ZbZSe6N_BXs/maxresdefault.jpg', 'Leadership Institute', UNIX_TIMESTAMP('2024-08-30'), 203000, 11, 0),
  
  -- Emotional Intelligence video (links to lesson_file 12)
  ('dPYE-7rkJCo', 720, 'https://img.youtube.com/vi/dPYE-7rkJCo/maxresdefault.jpg', 'EQ Development', UNIX_TIMESTAMP('2024-09-14'), 145000, 12, 0),
  
  -- Introduction to Strategic Thinking video (links to lesson_file 13)
  ('y8Kyi0WNg40', 900, 'https://img.youtube.com/vi/y8Kyi0WNg40/maxresdefault.jpg', 'Strategic Thinking Academy', UNIX_TIMESTAMP('2024-10-01'), 178000, 13, 0);
UNLOCK TABLES;


-- Table structure for table `x$io_global_by_file_by_bytes`
--

DROP TABLE IF EXISTS `x$io_global_by_file_by_bytes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$io_global_by_file_by_bytes` (
  `file` varchar(512) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `count_read` bigint NOT NULL DEFAULT '0',
  `total_read` bigint NOT NULL DEFAULT '0',
  `avg_read` decimal(23,4) NOT NULL DEFAULT '0.0000',
  `count_write` bigint NOT NULL DEFAULT '0',
  `total_written` bigint NOT NULL DEFAULT '0',
  `avg_write` decimal(23,4) NOT NULL DEFAULT '0.0000',
  `total` bigint NOT NULL DEFAULT '0',
  `write_pct` decimal(26,2) NOT NULL DEFAULT '0.00'
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$io_global_by_file_by_bytes`
--

LOCK TABLES `x$io_global_by_file_by_bytes` WRITE;
/*!40000 ALTER TABLE `x$io_global_by_file_by_bytes` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$io_global_by_file_by_bytes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$io_global_by_file_by_latency`
--

DROP TABLE IF EXISTS `x$io_global_by_file_by_latency`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$io_global_by_file_by_latency` (
  `file` varchar(512) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `total` bigint NOT NULL DEFAULT '0',
  `total_latency` bigint NOT NULL DEFAULT '0',
  `count_read` bigint NOT NULL DEFAULT '0',
  `read_latency` bigint NOT NULL DEFAULT '0',
  `count_write` bigint NOT NULL DEFAULT '0',
  `write_latency` bigint NOT NULL DEFAULT '0',
  `count_misc` bigint NOT NULL DEFAULT '0',
  `misc_latency` bigint NOT NULL DEFAULT '0'
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$io_global_by_file_by_latency`
--

LOCK TABLES `x$io_global_by_file_by_latency` WRITE;
/*!40000 ALTER TABLE `x$io_global_by_file_by_latency` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$io_global_by_file_by_latency` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$io_global_by_wait_by_bytes`
--

DROP TABLE IF EXISTS `x$io_global_by_wait_by_bytes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$io_global_by_wait_by_bytes` (
  `event_name` varchar(128) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `total` bigint NOT NULL DEFAULT '0',
  `total_latency` bigint NOT NULL DEFAULT '0',
  `min_latency` bigint NOT NULL DEFAULT '0',
  `avg_latency` bigint NOT NULL DEFAULT '0',
  `max_latency` bigint NOT NULL DEFAULT '0',
  `count_read` bigint NOT NULL DEFAULT '0',
  `total_read` bigint NOT NULL DEFAULT '0',
  `avg_read` decimal(23,4) NOT NULL DEFAULT '0.0000',
  `count_write` bigint NOT NULL DEFAULT '0',
  `total_written` bigint NOT NULL DEFAULT '0',
  `avg_written` decimal(23,4) NOT NULL DEFAULT '0.0000',
  `total_requested` bigint NOT NULL DEFAULT '0'
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$io_global_by_wait_by_bytes`
--

LOCK TABLES `x$io_global_by_wait_by_bytes` WRITE;
/*!40000 ALTER TABLE `x$io_global_by_wait_by_bytes` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$io_global_by_wait_by_bytes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$io_global_by_wait_by_latency`
--

DROP TABLE IF EXISTS `x$io_global_by_wait_by_latency`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$io_global_by_wait_by_latency` (
  `event_name` varchar(128) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `total` bigint NOT NULL DEFAULT '0',
  `total_latency` bigint NOT NULL DEFAULT '0',
  `avg_latency` bigint NOT NULL DEFAULT '0',
  `max_latency` bigint NOT NULL DEFAULT '0',
  `read_latency` bigint NOT NULL DEFAULT '0',
  `write_latency` bigint NOT NULL DEFAULT '0',
  `misc_latency` bigint NOT NULL DEFAULT '0',
  `count_read` bigint NOT NULL DEFAULT '0',
  `total_read` bigint NOT NULL DEFAULT '0',
  `avg_read` decimal(23,4) NOT NULL DEFAULT '0.0000',
  `count_write` bigint NOT NULL DEFAULT '0',
  `total_written` bigint NOT NULL DEFAULT '0',
  `avg_written` decimal(23,4) NOT NULL DEFAULT '0.0000'
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$io_global_by_wait_by_latency`
--

LOCK TABLES `x$io_global_by_wait_by_latency` WRITE;
/*!40000 ALTER TABLE `x$io_global_by_wait_by_latency` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$io_global_by_wait_by_latency` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$latest_file_io`
--

DROP TABLE IF EXISTS `x$latest_file_io`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$latest_file_io` (
  `thread` varchar(149) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `file` varchar(512) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `latency` bigint DEFAULT NULL,
  `operation` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `requested` bigint DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$latest_file_io`
--

LOCK TABLES `x$latest_file_io` WRITE;
/*!40000 ALTER TABLE `x$latest_file_io` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$latest_file_io` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `x$memory_by_host_by_current_bytes`
--

DROP TABLE IF EXISTS `x$memory_by_host_by_current_bytes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x$memory_by_host_by_current_bytes` (
  `host` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `current_count_used` double DEFAULT NULL,
  `current_allocated` double DEFAULT NULL,
  `current_avg_alloc` double NOT NULL DEFAULT '0',
  `current_max_alloc` bigint DEFAULT NULL,
  `total_allocated` double DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `x$memory_by_host_by_current_bytes`
--

LOCK TABLES `x$memory_by_host_by_current_bytes` WRITE;
/*!40000 ALTER TABLE `x$memory_by_host_by_current_bytes` DISABLE KEYS */;
/*!40000 ALTER TABLE `x$memory_by_host_by_current_bytes` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `entities` WRITE, `teachgroups` WRITE;
-- reset entities table auto increment
ALTER TABLE entities AUTO_INCREMENT = 1;

INSERT INTO entities (Ent_type, Tier) VALUES 
('student',1),
('teacher',1),
('class',2),
('year group',3),
('key stage',4),
('school',5),
('subject year',3);


INSERT INTO teachgroups (Groupname, entitiesID) VALUES ('ALL', (SELECT entities.entitiesID from entities WHERE entities.Ent_type='school' LIMIT 1)); 

UNLOCK TABLES;

LOCK TABLES `skills_key` WRITE;
-- reset skills_key table auto increment
ALTER TABLE skills_key AUTO_INCREMENT = 1;

UNLOCK TABLES;



LOCK TABLES `mimetypes` WRITE, `file_types` WRITE, `marktypes` WRITE;
-- reset mimetypes table auto increment

-- Insert standard MIME types for lesson files and resources
INSERT INTO `mimetypes` (`TypeName`) VALUES
('application/pdf'),
('application/msword'),
('application/vnd.openxmlformats-officedocument.wordprocessingml.document'),
('application/vnd.ms-excel'),
('application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'),
('application/vnd.ms-powerpoint'),
('application/vnd.openxmlformats-officedocument.presentationml.presentation'),
('text/plain'),
('text/csv'),
('text/html'),
('image/jpeg'),
('image/png'),
('image/gif'),
('image/webp'),
('video/mp4'),
('video/webm'),
('video/quicktime'),
('audio/mpeg'),
('audio/wav'),
('audio/ogg'),
('application/zip'),
('application/x-rar-compressed'),
('application/json'),
('application/xml');

-- Insert standard file types for categorization
INSERT INTO `file_types` (`TypeName`, `Description`) VALUES
('Document', 'Word documents, PDFs, and text files'),
('Spreadsheet', 'Excel spreadsheets and CSV files'),
('Presentation', 'PowerPoint presentations'),
('Image', 'Images in various formats (JPG, PNG, GIF)'),
('Video', 'Video files for multimedia lessons'),
('Audio', 'Audio files for listening exercises'),
('Archive', 'Compressed files (ZIP, RAR)'),
('Web', 'HTML and web-based resources'),
('Data', 'JSON and XML data files');

-- Insert standard mark types for assessment
INSERT INTO `marktypes` (`MarkTypeName`) VALUES
('Formative Assessment'),
('Summative Assessment'),
('Self Assessment'),
('Peer Assessment'),
('Teacher Feedback'),
('Skill Proficiency');

UNLOCK TABLES;

LOCK TABLES `skills_key` WRITE;
-- reset skills_key table auto increment

-- Insert skills from leadership.json into skills_key table with criteria
INSERT INTO skills_key (Skill_name, SkillDescription, DevelopmentAge, Hidden, Criteria1, Criteria2, Criteria3, Criteria4, Criteria5) VALUES
('Strategic Foresight', 'The ability to anticipate future trends, challenges, and opportunities by analyzing patterns and projecting current trajectories.', 16, 0, 
 'Identifies obvious immediate consequences of current decisions',
 'Projects short-term outcomes based on recent patterns and data',
 'Anticipates medium-term scenarios and prepares contingency plans',
 'Forecasts long-term trends and identifies emerging opportunities',
 'Shapes future outcomes through proactive strategy and innovation'),

('Risk Intelligence', 'The ability to identify, assess, and mitigate potential threats while balancing risk and reward in decision-making.', 18, 0,
 'Recognizes obvious risks in familiar situations',
 'Assesses probability and impact of identified risks',
 'Develops mitigation strategies for complex risk scenarios',
 'Anticipates emergent risks and creates proactive safeguards',
 'Transforms risks into strategic advantages through innovative approaches'),

('Opportunity Recognition', 'The ability to identify and capitalize on emerging possibilities, resources, and advantageous situations.', 16, 0,
 'Notices obvious opportunities when pointed out by others',
 'Identifies clear opportunities in familiar domains',
 'Recognizes subtle opportunities across different contexts',
 'Creates opportunities through strategic positioning and networking',
 'Generates novel opportunities by connecting disparate domains and trends'),

('Intrapersonal Mastery', 'The ability to understand and manage one''s own thoughts, emotions, motivations, and behavioral patterns.', 14, 0,
 'Recognizes basic emotions and their immediate triggers',
 'Understands personal strengths, weaknesses, and patterns',
 'Manages emotional responses and aligns actions with values',
 'Uses self-awareness to optimize performance and decision-making',
 'Transforms personal limitations into strengths through metacognitive practice'),

('Interpersonal Intelligence', 'The ability to understand, communicate with, and influence others effectively across diverse contexts and personalities.', 15, 0,
 'Recognizes basic emotions and perspectives in others',
 'Adapts communication style to different personalities',
 'Builds rapport and manages group dynamics effectively',
 'Resolves complex conflicts and facilitates collaboration',
 'Inspires and transforms relationships through deep understanding and trust'),

('Strategic Research', 'The ability to gather, analyze, and synthesize information to inform decision-making and identify patterns others miss.', 16, 0,
 'Gathers basic information from obvious sources',
 'Uses multiple sources and identifies relevant data',
 'Synthesizes information from diverse sources to identify patterns',
 'Conducts sophisticated analysis to reveal Hidden insights',
 'Generates novel knowledge through innovative research methodologies'),

('Adaptive Planning', 'The ability to create flexible plans that can evolve with changing circumstances while maintaining strategic direction.', 17, 0,
 'Creates simple linear plans for straightforward tasks',
 'Develops plans with basic contingency options',
 'Designs flexible plans that adapt to multiple scenarios',
 'Creates dynamic planning systems that evolve with new information',
 'Designs self-adjusting planning frameworks that anticipate and respond to change'),

('Integrative Leadership', 'The master skill of combining strategic foresight, interpersonal intelligence, and opportunity recognition to guide others effectively.', 20, 0,
 'Manages simple tasks with direct oversight of others',
 'Coordinates team efforts and communicates clear expectations',
 'Inspires team performance and navigates complex challenges',
 'Transforms organizational culture and achieves strategic objectives',
 'Creates legacy-level impact through visionary leadership and systemic change'),
 ('Emotional Intelligence', 'The ability to recognize, understand, and manage one''s own emotions as well as the emotions of others to foster effective relationships and decision-making.', 15, 0,
 'Identifies basic emotions in self and others',
 'Understands the causes and effects of emotions',
 'Manages emotional responses in various situations',
 'Utilizes emotional insights to enhance relationships and decision-making',
 'Leverages emotional intelligence to lead and inspire others');

UNLOCK TABLES;
-- Create skill relationships (parent-child dependencies)
-- Strategic Foresight (1) and Risk Intelligence (2) are parents of Adaptive Planning (7)
-- Strategic Foresight (1) and Risk Intelligence (2) are parents of Adaptive Planning (7)
-- Strategic Foresight (1), Risk Intelligence (2), Opportunity Recognition (3), 
-- Intrapersonal Mastery (4), Interpersonal Intelligence (5), Strategic Research (6) 
-- are all parents of Integrative Leadership (8)
-- Emotional Intelligence (9) is also a parent of Integrative Leadership (8)
-- Adaptive Planning (7) is also a parent of Integrative Leadership (8)

LOCK TABLES `skill_link` WRITE;
-- reset skill_link table auto increment
ALTER TABLE skill_link AUTO_INCREMENT = 1;

  INSERT INTO skill_link (ParentSkillID, OffspringSkillID) VALUES
  (3, 7),
  (6, 7),
  (4,9),
  (5,9),
  (9,8),
  (7,8),
  (1, 7),
  (2, 7),
  (	52,	62),
  (	70,	38),
  (	44,	55),
  (	46,	50),
  (	64, 28),
  (	60,	61),
  (	12,	57),
  (	66,	55),
  (	52,	38),
  (	28,	52),
  (	50,	51),
  (	57,	55),
  (	56,	53),
  (	59,	45),
  (	65,	45),
  (	45,	50),
  (	53,	66),
  (	61,	49),
  (	58,	52),
  (	63,	52),
  (	65,	46),
  (	56,	28),
  (	51,	54),
  (	60,	52),
  (	65,	44),
  (	65,	47),
  (	47,	66);
  UNLOCK TABLES;


LOCK TABLES `subject` WRITE;
INSERT INTO `subject` (`SubjectID`, `Subject_name`, `Subject_code`) VALUES
(1, 'GCSE Physics', '8463'),
(2, 'GCSE Biology', '8461'),
(3, 'GCSE Chemistry', '8462'),
(4, 'Leadership', 'LEAD'),
(5, 'STC Courses', 'STC');
UNLOCK TABLES;

LOCK TABLES `courses` WRITE;
-- reset the auto increment to 1
ALTER TABLE `courses` AUTO_INCREMENT = 1;
INSERT INTO `courses` (`Course_Name`, `Course_Code`, `SubjectID`) VALUES 
/*1*/('Corporate Leadership', 'CORP', 4),
/*2*/('Community Development', 'COMM', 4),
/*3*/('School Based Leadership', 'SBL', 4);
INSERT INTO `courses` (`Course_Name`, `Course_Code`, `SubjectID`, `Details`) VALUES
/*4*/('Physics Course', 'PHY01', 1, 'Detailed course content for GCSE Physics including topics such as energy transfers, forces, and space physics.'),
/*5*/('Chemistry Course', 'CHEM01', 3, 'Detailed course content for GCSE Chemistry including topics such as atomic structure, chemical changes, and organic chemistry.'),
/*6*/('Biology Course', 'BIO01', 2, 'Detailed course content for GCSE Biology including topics such as cell biology, bioenergetics, and ecology.'),
/*7*/('ETCAL Level 1 Award/ Certificate/ Diploma in Personal Social Skills', 'ETCAL01', 5, 'A course designed to develop personal and social skills for learners at Level 1.');
UNLOCK TABLES;

LOCK TABLES `modules` WRITE;
-- reset the auto increment to 1
ALTER TABLE `modules` AUTO_INCREMENT = 1;

INSERT INTO `modules` (`Module_name`, `Module_code`, `CourseID`, `Module_details`) VALUES
/*1*/('Energy', '4.1', 4, 'Study of energy transfers, conservation, and resources.'),
/*2*/('Electricity', '4.2', 4, 'Understanding electrical circuits, current, and potential difference.'),
/*3*/('Particle model of matter', '4.3', 4, 'Exploration of states of matter and particle behavior.'),
/*4*/('Atomic structure', '4.4', 4, 'Study of atomic models, isotopes, and nuclear reactions. '),
/*5*/('Forces', '4.5', 4, 'Analysis of forces, motion, and their interactions. '),
/*6*/('Waves', '4.6', 4, 'Study of wave properties, behavior, and applications.'),
/*7*/('Magnetism and electromagnetism', '4.7', 4, 'Understanding magnetic fields and electromagnetic effects.'),
/*8*/('Space physics', '4.8', 4, 'Study of the solar system, stars, and the universe.');
INSERT INTO `modules` (`Module_name`, `Module_code`, `CourseID`, `Module_details`) VALUES
/*9*/('Cell biology', '4.1', 6, 'Study of cells, their structure, division, and functions. '),
/*10*/('Organisation', '4.2', 6, 'Exploration of the human body systems and their organisation.'),
/*11*/('Infection and response', '4.3', 6, 'Understanding pathogens, diseases, and the immune system.'),
/*12*/('Bioenergetics', '4.4', 6, 'Study of photosynthesis and respiration processes. '),
/*13*/('Homeostasis and response', '4.5', 6, 'Regulation of internal conditions and responses to changes.'),
/*14*/('Inheritance, variation and evolution', '4.6', 6, 'Genetics, evolution, and variation in organisms.'),
/*15*/('Ecology', '4.7', 6, 'Study of ecosystems, biodiversity, and environmental interactions.');
INSERT INTO `modules` (`Module_name`, `Module_code`, `CourseID`, `Module_details`) VALUES
/*16*/('Atomic structure and the periodic table', '4.1', 5, 'Study of atoms, elements, compounds, and the periodic table.'),
/*17*/('Bonding, structure, and the properties of matter', '4.2', 5, 'Exploration of chemical bonding and material properties.'),
/*18*/('Quantitative chemistry', '4.3', 5, 'Quantitative analysis of chemical reactions and equations.'),
/*19*/('Chemical changes', '4.4', 5, 'Understanding chemical reactions and changes in matter.'),
/*20*/('Energy changes', '4.5', 5, 'Study of energy transfer in chemical reactions. '),
/*21*/('The rate and extent of chemical change', '4.6', 5, 'Factors affecting reaction rates and equilibrium.'),
/*22*/('Organic chemistry', '4.7', 5, 'Study of carbon-based compounds and their reactions.'),
/*23*/('Chemical analysis', '4.8', 5, 'Qualitative and instrumental methods for chemical analysis. '),
/*24*/('Chemistry of the atmosphere', '4.9', 5, 'Study of atmospheric chemistry and environmental impact.'),
/*25*/('Using resources', '4.10', 5, 'Sustainable use of resources and their applications.');

UNLOCK TABLES;

LOCK TABLES `topics` WRITE;
-- reset the auto increment to 1
ALTER TABLE `topics` AUTO_INCREMENT = 1;

INSERT INTO `topics` (`Topic_name`, `Topic_code`, `ModuleID`) VALUES
('Energy changes in a system, and the ways energy is stored before and after such changes', '4.1.1', 1),
('National and global energy resources', '4.1.3', 1),
('Current, potential difference and resistance', '4.2.1', 2),
('Series and parallel circuits', '4.2.2', 2), 
('Domestic uses and safety', '4.2.3', 2), 
('Energy transfers', '4.2.4', 2),
('Static electricity', '4.2.5', 2), 
('Changes of state and the particle model', '4.3.1', 3),
('Internal energy and energy transfers', '4.3.2', 3), 
('Particle model and pressure', '4.3.3', 3), 
('Atoms and isotopes', '4.4.1', 4), 
('Atoms and nuclear radiation', '4.4.2', 4),
('Hazards and uses of radioactive emissions and of background
radiation', '4.4.3', 4), 
('Nuclear fission and fusion', '4.4.4', 4), 
('Forces and their interactions', '4.5.1', 5), 
('Work done and energy transfer', '4.5.2', 5), 
('Forces and elasticity', '4.5.3', 5),
('Moments, levers and gears', '4.5.4', 5),
('Pressure and pressure differences in fluids', '4.5.5', 5), 
('Forces and motion', '4.5.6', 5),
('Momentum', '4.5.7', 5),
('Waves in air, fluids and solids', '4.6.1', 6), 
('Electromagnetic waves', '4.6.2', 6), 
('Black body radiation', '4.6.3', 6),
('Permanent and induced magnetism, magnetic forces and fields', '4.7.1', 7), 
('The motor effect', '4.7.2', 7),
('Induced potential, transformers and the National Grid', '4.7.3', 7),
('Solar system; stability of orbital motions; satellites', '4.8.1', 8), 
('Red-shift', '4.8.2', 8);


INSERT INTO `topics` ( `ModuleID`,`Topic_code`, `Topic_name`) VALUES
(9,'4.1.1','Cell structure'),
(9,'4.1.2','Cell division'),
(9,'4.1.3','Transport in cells'),
(10,'4.2.1','Principles of organisation'),
(10,'4.2.2','Animal tissues, organs and organ systems'),
(10,'4.2.3','Plant tissues, organs and systems'),
(11,'4.3.1','Communicable diseases'),
(11,'4.3.2','Monoclonal antibodies'),
(11,'4.3.3','Plant diseases'),
(12,'4.4.1','Photosynthesis'),
(12,'4.4.2','Respiration'),
(12,'4.4.3','Homeostasis'),
(13,'4.5.1','Inheritance'),
(13,'4.5.2','The human nervous system'),
(13,'4.5.3','Hormonal coordination in humans'),
(13,'4.5.4','Plant hormones'),
(14,'4.6.1','Reproduction'),
(14,'4.6.2','Variation and evolution'),
(14,'4.6.3','The development of understanding of genetics and evolution'),
(14,'4.6.4','Classification of living organisms'),
(15,'4.7.1','Adaptations, interdependence and competition'),
(15,'4.7.2','Organisation of an ecosystem'),
(15,'4.7.3','Biodiversity and the effect of human interaction on ecosystems'),
(15,'4.7.4','Trophic levels in an ecosystem'),
(15,'4.7.5','Food production');


INSERT INTO `topics` ( `ModuleID`,`Topic_code`, `Topic_name`) VALUES
(16, '4.1.1', 'A simple model of the atom, symbols, relative atomic mass, electronic charge and isotopes'),
(16, '4.1.2', 'The periodic table'),
(16, '4.1.3', 'Properties of transition metals'),
(17, '4.2.1', 'Chemical bonds, ionic, covalent and metallic'),
(17, '4.2.2', 'How bonding and structure are related to the properties of
substances'),
(17, '4.2.3', 'Structure and bonding of carbon'),
(17, '4.2.4', 'Bulk and surface properties of matter including nanoparticles'),
(18, '4.3.1', 'Chemical measurements, conservation of mass and the
quantitative interpretation of chemical equations'),
(18, '4.3.2', 'Use of amount of substance in relation to masses of pure
substances'),
(18, '4.3.3', 'Yield and atom economy of chemical reactions'),
(18, '4.3.4', 'Using concentrations of solutions in mol/dm3'),
(18, '4.3.5', 'Use of amount of substance in relation to volumes of gases'),
(19, '4.4.1', 'Reactivity of metals'),
(19, '4.4.2', 'Reactions of acids'),
(19, '4.4.3', 'Electrolysis'),
(20, '4.5.1', 'Exothermic and endothermic reactions'),
(20, '4.5.2', 'Chemical cells and fuel cells'),
(21, '4.6.1', 'Rate of reaction'),
(21, '4.6.2', 'Reversible reactions and dynamic equilibrium'),
(22, '4.7.1', 'Carbon compounds as fuels and feedstock'),
(22, '4.7.2', 'Reactions of alkenes and alcohols'),
(22, '4.7.3', 'Synthetic and naturally occurring polymers'),
(23, '4.8.1', 'Purity, formulations and chromatography'),
(23, '4.8.2', 'Identification of common gases'),
(23, '4.8.3', 'Identification of ions by chemical and spectroscopic means'),
(24, '4.9.1', 'The composition and evolution of the Earth''s atmosphere'),
(24, '4.9.2', 'Carbon dioxide and methane as greenhouse gases'),
(24, '4.9.3', 'Common atmospheric pollutants and their sources'),
(25, '4.10.1', 'Using the Earth''s resources and obtaining potable water'),
(25, '4.10.2', 'Life cycle assessment and recycling'),
(25, '4.10.3', 'Using materials'),
(25, '4.10.4', 'The Haber process and the use of NPK fertilisers');

UNLOCK TABLES;

LOCK TABLES `lessons` WRITE;
-- reset the auto increment to 1
ALTER TABLE `lessons` AUTO_INCREMENT = 1;

INSERT INTO `lessons` (`TopicID`, `lessonName`, `lessonCode`) VALUES 
(1,'Energy stores and systems', '4.1.1.1'), 
(1, 'Changes in energy', '4.1.1.2'), 
(1, 'Energy changes in systems', '4.1.1.3'), 
(1, 'Power', '4.1.1.4'),
(2, 'Energy transfers in a system', '4.1.2.1'), 
(2, 'Efficiency', '4.1.2.2'), 
(3, 'National and global energy resources', '4.1.3.1'), 
(4, 'Standard circuit diagram symbols', '4.2.1.1'), 
(4, 'Electrical charge and current', '4.2.1.2'), 
(4, 'Current, resistance and potential difference', '4.2.1.3'), 
(4, 'Resistors', '4.2.1.4'), 
(5, 'Series and parallel circuits', '4.2.2'),
(6, 'Direct and alternating potential difference', '4.2.3.1'),
(6, 'Mains electricity', '4.2.3.2'), 
(7, 'Power', '4.2.4.1'), 
(7, 'Energy transfers in everyday appliances', '4.2.4.2'), 
(7, 'The National Grid', '4.2.4.3'),
(8, 'Static charge', '4.2.5.1'),
(8, 'Electric fields', '4.2.5.2'),
(9, 'Density of materials', '4.3.1.1'), 
(9, 'Changes of state', '4.3.1.2'), 
(10, 'Internal energy', '4.3.2.1'), 
(10, 'Temperature changes in a system and specific heat capacity', '4.3.2.2'), 
(10, 'Changes of state and specific latent heat', '4.3.2.3'), 
(11, 'Particle motion in gases', '4.3.3.1'), 
(11, 'Pressure in gases', '4.3.3.2'),
(11, 'Increasing the pressure of a gas', '4.3.3.3'), 
(12, 'The structure of an atom', '4.4.1.1'), 
(12, 'Mass number, atomic number and isotopes', '4.4.1.2'), 
(12, 'The development of the model of the atom', '4.4.1.3'), 
(13, 'Radioactive decay and nuclear radiation', '4.4.2.1'), 
(13, 'Nuclear equations', '4.4.2.2'), 
(13, 'Half-lives and the random nature of radioactive decay', '4.4.2.3'), 
(13, 'Radioactive contamination', '4.4.2.4'), 
(14, 'Background radiation', '4.4.3.1'), 
(14, 'Different half-lives of radioactive isotopes', '4.4.3.2'), 
(14, 'Uses of nuclear radiation', '4.4.3.3'), 
(15, 'Nuclear fission', '4.4.4.1'), 
(15, 'Nuclear fusion', '4.4.4.2'),
(16, 'Scalar and vector quantities', '4.5.1.1'), 
(16, 'Contact and non-contact forces', '4.5.1.2'), 
(16, 'Gravity', '4.5.1.3'), 
(16, 'Resultant forces', '4.5.1.4'), 
(17, 'Work done and energy transfer', '4.5.2'), 
(18, 'Forces and elasticity', '4.5.3'), 
(19, 'Moments, levers and gears', '4.5.4'), 
(20, 'Pressure in a fluid', '4.5.5.1.1'), 
(20, 'Pressure in a fluid 2', '4.5.5.1.2'), 
(20, 'Atmospheric pressure', '4.5.5.2'), 
(21, 'Distance and displacement', '4.5.6.1.1'), 
(21, 'Speed', '4.5.6.1.2'),
(21, 'Velocity', '4.5.6.1.3'), 
(21, 'The distance–time relationship', '4.5.6.1.4'), 
(21, 'Acceleration', '4.5.6.1.5'), 
(21, 'Newton''s First Law', '4.5.6.2.1'), 
(21, 'Newton''s Second Law', '4.5.6.2.2'), 
(21, 'Newton''s Third Law', '4.5.6.2.3'), 
(21, 'Stopping distances', '4.5.6.3.1'),
(21, 'Reaction time', '4.5.6.3.2'), 
(21, 'Factors affecting braking distance 1', '4.5.6.3.3'), 
(21, 'Factors affecting braking distance 2', '4.5.6.3.4'), 
(22, 'Momentum is a property of moving objects', '4.5.7.1'), 
(22, 'Conservation of momentum', '4.5.7.2'), 
(22, 'Changes in momentum', '4.5.7.3'), 
(23, 'Transverse and longitudinal waves', '4.6.1.1'), 
(23, 'Properties of waves', '4.6.1.2'), 
(23, 'Reflection of waves', '4.6.1.3'), 
(23, 'Sound waves', '4.6.1.4'), 
(23, 'Waves for detection and exploration', '4.6.1.5'), 
(24, 'Types of electromagnetic waves', '4.6.2.1'), 
(24, 'Properties of electromagnetic waves 1', '4.6.2.2'), 
(24, 'Properties of electromagnetic waves 2', '4.6.2.3'), 
(24, 'Uses and applications of electromagnetic waves', '4.6.2.4'),
(25, 'Lenses', '4.6.2.5'), 
(25, 'Visible light', '4.6.2.6'),
(26, 'Emission and absorption of infrared radiation', '4.6.3.1'), 
(26, 'Perfect black bodies and radiation', '4.6.3.2'), 
(27, 'Poles of a magnet', '4.7.1.1'), 
(27, 'Magnetic fields', '4.7.1.2'), 
(28, 'Electromagnetism', '4.7.2.1'), 
(28, 'Fleming''s left-hand rule', '4.7.2.2'), 
(28, 'Electric motors', '4.7.2.3'), 
(28, 'Loudspeakers', '4.7.2.4'), 
(29, 'Induced potential', '4.7.3.1'), 
(29, 'Uses of the generator effect', '4.7.3.2'), 
(29, 'Microphones', '4.7.3.3'),
(29, 'Transformers', '4.7.3.4'),
(30, 'Our solar system', '4.8.1.1'), 
(30, 'The life cycle of a star', '4.8.1.2'), 
(30, 'Orbital motion, natural and artificial satellites', '4.8.1.3'), 
(30, 'Red-shift', '4.8.2'); 


INSERT INTO `lessons` ( `TopicID`, `lessonCode`, `lessonName`) VALUES
(32,'4.1.1.1','Eukaryotes and prokaryotes'),
(32,'4.1.1.2','Animal and plant cells'),
(32,'4.1.1.3','Cell specialisation'),
(32,'4.1.1.4','Cell differentiation'),
(32,'4.1.1.5','Microscopy'),
(32,'4.1.1.6','Culturing microorganisms'),
(33,'4.1.2.2','Mitosis and the cell cycle'),
(33,'4.1.2.3','Stem cells'),
(34,'4.1.3.1','Diffusion'),
(34,'4.1.3.2','Osmosis'),
(34,'4.1.3.3','Active transport'),
(35,'4.2.2.1','The human digestive system'),
(35,'4.2.2.2','The heart and blood vessels'),
(35,'4.2.2.3','Blood'),
(35,'4.2.2.4','Coronary heart disease: a non-communicable disease'),
(35,'4.2.2.5','Health issues'),
(35,'4.2.2.6','The effect of lifestyle on some non-communicable diseases'),
(35,'4.2.2.7','Cancer'),
(36,'4.2.3.1','Plant tissues'),
(36,'4.2.3.2','Plant organ system'),
(37,'4.3.1.1','Communicable (infectious) diseases'),
(37,'4.3.1.2','Viral diseases'),
(37,'4.3.1.3','Bacterial diseases'),
(37,'4.3.1.4','Fungal diseases'),
(37,'4.3.1.5','Protist diseases'),
(37,'4.3.1.6','Human defence systems'),
(37,'4.3.1.7','Vaccination'),
(37,'4.3.1.8','Antibiotics and painkillers'),
(37,'4.3.1.9','Discovery and development of drugs'),
(38,'4.3.2.1','Producing monoclonal antibodies'),
(38,'4.3.2.2','Uses of monoclonal antibodies'),
(39,'4.3.3.1','Detection and identification of plant diseases'),
(39,'4.3.3.2','Plant defence responses'),
(40,'4.4.1.1','Photosynthetic reaction'),
(40,'4.4.1.2','Rate of photosynthesis'),
(40,'4.4.1.3','Uses of glucose from photosynthesis'),
(41,'4.4.2.1','Aerobic and anaerobic respiration'),
(41,'4.4.2.2','Response to exercise'),
(41,'4.4.2.3','Metabolism'),
(42,'4.5.2.1','Structure and function'),
(42,'4.5.2.2','The brain'),
(42,'4.5.2.3','The eye'),
(42,'4.5.2.4','Control of body temperature'),
(43,'4.5.3.1','Human endocrine system'),
(43,'4.5.3.2','Control of blood glucose concentration'),
(43,'4.5.3.3','Maintaining water and nitrogen balance in the body'),
(43,'4.5.3.4','Hormones in human reproduction'),
(43,'4.5.3.5','Contraception'),
(43,'4.5.3.6','The use of hormones to treat infertility'),
(43,'4.5.3.7','Negative feedback'),
(44,'4.5.4.1','Control and coordination'),
(44,'4.5.4.2','Use of plant hormones'),
(45,'4.6.1.1','Sexual and asexual reproduction'),
(45,'4.6.1.2','Meiosis'),
(45,'4.6.1.3','Advantages and disadvantages of sexual and asexual reproduction'),
(45,'4.6.1.4','DNA and the genome'),
(45,'4.6.1.5','DNA structure'),
(45,'4.6.1.6','Genetic inheritance'),
(45,'4.6.1.7','Inherited disorders'),
(45,'4.6.1.8','Sex determination'),
(46,'4.6.2.1','Variation'),
(46,'4.6.2.2','Evolution'),
(46,'4.6.2.3','Selective breeding'),
(46,'4.6.2.4','Genetic engineering'),
(46,'4.6.2.5','Cloning'),
(47,'4.6.3.1','Theory of evolution'),
(47,'4.6.3.2','Speciation'),
(47,'4.6.3.3','The understanding of genetics'),
(47,'4.6.3.4','Evidence for evolution'),
(47,'4.6.3.5','Fossils'),
(47,'4.6.3.6','Extinction'),
(47,'4.6.3.7','Resistant bacteria'),
(48,'4.7.1.1','Communities'),
(48,'4.7.1.2','Abiotic factors'),
(48,'4.7.1.3','Biotic factors'),
(48,'4.7.1.4','Adaptations'),
(49,'4.7.2.1','Levels of organisation'),
(49,'4.7.2.2','How materials are cycled'),
(49,'4.7.2.3','Decomposition'),
(49,'4.7.2.4','Impact of environmental change'),
(50,'4.7.3.1','Biodiversity'),
(50,'4.7.3.2','Waste management'),
(50,'4.7.3.3','Land use'),
(50,'4.7.3.4','Deforestation'),
(50,'4.7.3.5','Global warming'),
(50,'4.7.3.6','Maintaining biodiversity'),
(51,'4.7.4.1','Trophic levels'),
(51,'4.7.4.2','Pyramids of biomass'),
(51,'4.7.4.3','Transfer of biomass'),
(51,'4.7.5.1','Factors affecting food security'),
(51,'4.7.5.2','Farming techniques'),
(51,'4.7.5.3','Sustainable fisheries'),
(51,'4.7.5.4','Role of biotechnology');


INSERT INTO `lessons` (`lessonName`, `TopicID`, `lessonCode`) VALUES
('Atoms, elements and compounds', 52, '4.1.1.1'),
('Mixtures', 52, '4.1.1.2'),
('The development of the model of the atom', 52, '4.1.1.3'),
('Relative electrical charges of subatomic particles', 52, '4.1.1.4'),
('Size and mass of atoms', 52, '4.1.1.5'),
('Relative atomic mass', 52, '4.1.1.6'),
('Electronic structure', 52, '4.1.1.7'),
('The periodic table', 53, '4.1.2.1'),
('Development of the periodic table', 53, '4.1.2.2'),
('Metals and non-metals', 53, '4.1.2.3'),
('Group 0', 53, '4.1.2.4'),
('Group 1', 53, '4.1.2.5'),
('Group 7', 53, '4.1.2.6'),
('Comparison with Group 1 elements', 54, '4.1.3.1'),
('Typical properties of transition metals', 54, '4.1.3.2'),
('Chemical bonds', 55, '4.2.1.1'),
('Ionic bonding', 55, '4.2.1.2'),
('Ionic compounds', 55, '4.2.1.3'),
('Covalent bonding', 55, '4.2.1.4'),
('Metallic bonding', 55, '4.2.1.5'),
('The three states of matter', 56, '4.2.2.1'),
('State symbols', 56, '4.2.2.2'),
('Properties of ionic compounds', 56, '4.2.2.3'),
('Properties of small molecules', 56, '4.2.2.4'),
('Polymers', 56, '4.2.2.5'),
('Giant covalent structures', 56, '4.2.2.6'),
('Properties of metals and alloys', 56, '4.2.2.7'),
('Metals as conductors', 56, '4.2.2.8'),
('Diamond', 57, '4.2.3.1'),
('Graphite', 57, '4.2.3.2'),
('Graphene and fullerenes', 57, '4.2.3.3.'),
('Sizes of particles and their properties', 58, '4.2.4.1'),
('Uses of nanoparticles', 58, '4.2.4.2'),
('Conservation of mass and balanced chemical equations', 59, '4.3.1.1'),
('Relative formula mass', 59, '4.3.1.2'),
('Mass changes when a reactant or product is a gas', 59, '4.3.1.3'),
('Chemical measurements', 59, '4.3.1.4'),
('Moles', 60, '4.3.2.1'),
('Amount of substance in equations', 60, '4.3.2.2'),
('Using moles to balance equations', 60, '4.3.2.3'),
('Limiting reactants', 60, '4.3.2.4'),
('Concentrations of solutions', 60, '4.3.2.5'),
('Percentage yield', 61, '4.3.3.1'),
('Atom economy', 61, '4.3.3.2'),
('Using concentrations of solutions in mol/dm3', 62, '4.3.4'),
('Use of amount of substance in relation to volumes of gases', 63, '4.3.5'),
('Metal oxides', 64, '4.4.1.1'),
('The reactivity series', 64, '4.4.1.2'),
('Extractiion of metals and reductions', 64, '4.4.1.3'),
('Oxidation and reduction in terms of electrons', 64, '4.4.1.4'),
('Reactions of acids with metals', 65, '4.4.2.1'),
('Neutralisation of acids and salt production', 65, '4.4.2.2'),
('Soluble salts', 65, '4.4.2.3'),
('The pH scale and neutralisation', 65, '4.4.2.4'),
('Titrations', 65, '4.4.2.5'),
('Strong and weak acids', 65, '4.4.2.6'),
('The process of electrolysis', 66, '4.4.3.1'),
('Electrolysis of molten ionic compounds', 66, '4.4.3.2'),
('Using electrolysis to extract metals', 66, '4.4.3.3'),
('Electrolysis of aqueous solutions', 66, '4.4.3.4'),
('Representation of reactions at electrodes as half equations', 66, '4.4.3.5'),
('Energy transfer during exothermic and endothermic reactions', 67, '4.5.1.1'),
('Reaction profiles', 67, '4.5.1.2'),
('The energy change of reactions', 67, '4.5.1.3'),
('Cells and batteries', 68, '4.5.2.1'),
('Fuel cells', 68, '4.5.2.2'),
('Calculating rates of reactions', 69, '4.6.1.1'),
('Factors affecting the rate of reaction', 69, '4.6.1.2'),
('Collision theory and activation energy', 69, '4.6.1.3'),
('Catalysts', 69, '4.6.1.4'),
('Reversible reactions', 70, '4.6.2.1'),
('Energy changes and reversible reactions', 70, '4.6.2.2'),
('Equilibrium', 70, '4.6.2.3'),
('The effect of changing conditions on equilibrium', 70, '4.6.2.4'),
('The effect of changing concentration', 70, '4.6.2.5'),
('The effect of temperature changes on equilibrium', 70, '4.6.2.6'),
('The effect of pressure changes on equilibrium', 70, '4.6.2.7'),
('Crude oil, hydrocarbons and alkanes', 71, '4.7.1.1'),
('Fractional distillation and petrochemicals', 71, '4.7.1.2'),
('Properties of hydrocarbons', 71, '4.7.1.3'),
('Cracking and alkenes', 71, '4.7.1.4'),
('Cracking and alkenes', 72, '4.7.2.1'),
('Reactions of alkenes', 72, '4.7.2.2'),
('Alcohols', 72, '4.7.2.3'),
('Carboxylic acids', 72, '4.7.2.4'),
('Addition polymerisation', 73, '4.7.3.1'),
('Condensation polymerisation', 73, '4.7.3.2'),
('Amino acids', 73, '4.7.3.3'),
('Pure substances', 74, '4.8.1.1'),
('Formulations', 74, '4.8.1.2'),
('Chromatography', 74, '4.8.1.3'),
('Test for hydrogen', 75, '4.8.2.1'),
('Test for oxygen', 75, '4.8.2.2'),
('Test for carbon dioxide', 75, '4.8.2.3'),
('Test for chlorine', 75, '4.8.2.4'),
('Flame tests', 76, '4.8.3.1'),
('Metal hydroxides', 76, '4.8.3.2'),
('Carbonates', 76, '4.8.3.3'),
('Halides', 76, '4.8.3.4'),
('Sulfates', 76, '4.8.3.5'),
('Instrumental methods', 76, '4.8.3.6'),
('Flame emission spectroscopy', 76, '4.8.3.7'),
('The proportions of different gases in the atmosphere', 77, '4.9.1.1'),
('The Earth''s early atmosphere', 77, '4.9.1.2'),
('How oxygen increased', 77, '4.9.1.3'),
('How carbon dioxide decreased', 77, '4.9.1.4'),
('Greenhouse gases', 78, '4.9.2.1'),
('Human activities which contribute to an increase in greenhouse gases in the atmosphere', 78, '4.9.2.2'),
('Global climate change', 78, '4.9.2.3'),
('The carbon footprint and its reduction', 78, '4.9.2.4'),
('Atmospheric pollutants from fuels', 79, '4.9.3.1'),
('Properties and effects of atmospheric pollutants', 79, '4.9.3.2'),
('Using the Earth''s resources and sustainable development', 80, '4.10.1.1'),
('Potable water', 80, '4.10.1.2'),
('Waste water treatment', 80, '4.10.1.3'),
('Alternative methods of extracting metals', 80, '4.10.1.4'),
('Life cycle assessment', 81, '4.10.2.1'),
('Ways of reducing the use of resources', 82, '4.10.2.2'),
('Corrosion and its prevention', 83, '4.10.3.1'),
('Alloys as useful materials', 83, '4.10.3.2'),
('Ceramics, polymers and composites', 83, '4.10.3.3'),
('The Haber process', 32, '4.10.4.1'),
('Production and uses of NPK fertilisers', 32, '4.10.4.2');


UNLOCK TABLES;

LOCK TABLES `skills_key` WRITE;
-- reset the auto increment to 1
ALTER TABLE `skills_key` AUTO_INCREMENT = 1;
INSERT INTO `skills_key` (`Skills_keyID`, `Skill_name`, `Criteria1`, `Criteria2`, `Criteria3`, `Criteria4`, `hidden`, `Criteria5`, `DevelopmentAge`, `SkillDescription`) VALUES
(11,	'Working Memory',	'Can remember 1 item',	'Can remember 2 items',	'Can remember 3 items',	'Can remember 4 items',	0,	'Can remember 5 items',	3,	'Working memory is the ability to remember information for a short time.  It is used when we are thinking about something in the present moment.  It is used when we are thinking about something in the present moment.  It is used when we are thinking about something in the present moment.'),
(12,	'Attention',	'Can attend for 1 minute',	'Can attend for 2 minutes',	'Can attend for 3 minutes',	'Can attend for 4 minutes',	0,	'Can attend for 5 minutes',	3,	'Attention is the ability to focus on something for a period of time.  It is used when we are thinking about something in the present moment.  It is used when we are thinking about something in the present moment.  It is used when we are thinking about something in the present moment.'),
(13,	'Audio Processing',	'Can hear speach in a quiet room',	'Can hear speach in a room with some background noise',	'Can hear speach in a room with a lot of background noise',	'Can hear speach in a room with a lot of background noise and other people talking',	0,	'Can hear speach in a room with a lot of background noise and other people talking and music playing',	3,	'Audio processing is the ability to hear and understand speach.  It is used when we are thinking about something in the present moment.  It is used when we are thinking about something in the present moment.  It is used when we are thinking about something in the present moment.'),
(14,	'Critical Thinking',	'Can understand the context of a situation',	'Can understand the context of a situation and the chain of events that lead to the situation',	'Can understand the context of a situation and the chain of events that lead to the situation and the possible outcomes of the situation',	'Can understand the context of a situation and the chain of events that lead to the situation and the possible outcomes of the situation and the consequences of the situation',	0,	'Can understand the context of a situation and the chain of events that lead to the situation and the possible outcomes of the situation and the consequences of the situation and the consequences of the situation',	3,	'Critical thinking is the ability to understand the context of a situation and the chain of events that lead to the situation and the possible outcomes of the situation and the consequences of the situation.'),
(15,	'Problem Solving',	'Can solve a problem with one step',	'Can solve a problem with two steps',	'Can solve a problem with three steps',	'Can solve a problem with four steps',	0,	'Can solve a problem with five steps',	3,	'Problem solving is the ability to solve a problem with multiple steps.'),
(16,	'Active Listening',	'Can listen to someone and understand what they are saying',	'Can listen to someone and understand what they are saying and ask questions to clarify what they are saying',	'Can listen to someone and understand what they are saying and ask questions to clarify what they are saying and to contextualise what they are saying',	'Can listen to someone and understand what they are saying and ask questions to clarify what they are saying and to contextualise what they are saying and to make predictions about what they are saying',	0,	'Can listen to someone and understand what they are saying and ask questions to clarify what they are saying and to contextualise what they are saying and to make predictions about what they are saying and to make predictions about what they are saying',	3,	'Active listening is the ability to listen to someone and understand what they are saying and to ask questions to clarify what they are saying and to contextualise what they are saying and to make predictions about what they are saying.'),
(17,	'Leadership',	'The child is able to make decisions and take responsibility for the consequences of those decisions.',	'The child is able to organise and motivate others to achieve a common goal.',	'The child is able to make decisions and take responsibility for the consequences of those decisions.',	'The child is able to organise and motivate others to achieve a common goal.',	0,	'The child is able to make decisions and take responsibility for the consequences of those decisions.',	3,	'This skill involves the ability to make decisions and to take responsibility for the consequences of those decisions.  It also involves the ability to organise and motivate others to achieve a common goal.'),
(18,	'Vocabulary',	'1.  Uses a limited range of words',	'2.  Uses a range of words',	'3.  Uses a wide range of words',	'4.  Uses a wide range of words and specialist vocabulary',	0,	'5.  Uses a wide range of words and specialist vocabulary',	3,	'This skill involves the ability to use a wide range of words and specialist vocabulary.  It also involves the ability to use words in the correct context.'),
(19,	'Non-verbal communication',	'1.  Uses a limited range of non-verbal communication',	'2.  Uses a range of non-verbal communication',	'3.  Uses a wide range of non-verbal communication',	'4.  Uses a wide range of non-verbal communication',	0,	'5.  Uses a wide range of non-verbal communication',	3,	'This skill involves the ability to use a wide range of non-verbal communication.  It also involves the ability to use non-verbal communication in the correct context.'),
(20,	'Technology',	'1.  Uses a limited range of technology',	'2.  Uses a range of technology',	'3.  Uses a wide range of technology',	'4.  Uses a wide range of technology',	0,	'5.  Uses a wide range of technology',	3,	'This skill involves the ability to use a wide range of technology.  It also involves the ability to use technology in the correct context.'),
(21,	'Comprehension',	'1.  Has a limited understanding of what is being said',	'2.  Has a good understanding of what is being said',	'3.  Has a good understanding of what is being said and can answer questions about it',	'4.  Has a good understanding of what is being said and can answer questions about it',	0,	'5.  Has a good understanding of what is being said and can answer questions about it',	3,	'This skill involves the ability to understand what is being said.  It also involves the ability to answer questions about what is being said.'),
(22,	'Speaking',	'1.  Speaks in short sentences',	'2.  Speaks in sentences',	'3.  Speaks in sentences and can answer questions about what has been said',	'4.  Speaks in sentences and can answer questions about what has been said',	0,	'5.  Speaks in sentences and can answer questions about what has been said',	3,	'This skill involves the ability to speak in sentences.  It also involves the ability to answer questions about what has been said.'),
(23,	'Reading',	'1.  25% behind reading age',	'2.  10% behind reading age',	'3.  Reading age',	'4.  10% ahead of reading age',	0,	'5.  25% ahead of reading age',	3,	'This skill involves the ability to read.'),
(24,	'Phonics',	'1.  Recognises some sounds',	'2.  Recognises most sounds',	'3.  Recognises all sounds',	'4.  Recognises all sounds and can read words where the sounds are not as expected',	0,	'5.  Recognises all sounds and can read words where the sounds are not as expected',	3,	'This skill involves the ability to recognise sounds and to blend those sounds together to form words.  It also involves the ability to read words where the sounds are not as expected.'),
(25,	'Grammar',	'1.  Uses a limited range of grammar',	'2.  Uses a range of grammar',	'3.  Uses a wide range of grammar',	'4.  Uses a wide range of grammar',	0,	'5.  Uses a wide range of grammar',	3,	'This skill involves the ability to use a wide range of grammar.  It also involves the ability to use grammar in the correct context.'),
(27,	'Spelling',	'1.  Spelling age',	'2.  10% behind spelling age',	'3.  Spelling age',	'4.  10% ahead of spelling age',	0,	'5.  25% ahead of spelling age',	3,	'This skill involves the ability to spell.'),
(28,	'Handwriting',	'1.  Writes legibly',	'2.  Writes legibly and fluently',	'3.  Writes legibly and fluently',	'4.  Writes legibly and fluently',	0,	'5.  Writes legibly and fluently',	3,	'This skill involves the ability to write legibly and fluently.  It also involves the ability to write legibly and fluently in the correct context.'),
(29,	'Persuading',	'The student has demonstrated only a basic grasp of some of the parent skills',	'The student has demonstrated a basic grasp of some of the parent skills',	'The student has demonstrated a good grasp of some of the parent skills',	'The student has demonstrated a very good grasp of most of the parent skills',	0,	'The student has demonstrated an excellent grasp of all of the parent skills',	0,	'Persuading is the ability to convince others to accept your point of view.  It involves the ability to research a topic and to present a case for or against the topic.  It also involves the ability to listen to the opposing case and understand the reasons and emotional to respond to it.  It also involves the ability to think on your feet and to respond to questions from the audience.  The key elements of persuading are: Public Speaking, critical thinking, interpersonal skills, listening, thinkin'),
(30,	'Questioning',	'The student has demonstrated only a basic grasp of some of the parent skills',	'The student has demonstrated a basic grasp of some of the parent skills',	'The student has demonstrated a good grasp of some of the parent skills',	'The student has demonstrated a very good grasp of most of the parent skills',	0,	'The student has demonstrated an excellent grasp of all of the parent skills',	0,	'The skill of questioning requires the ability to understand the area to a depth to lead a discussion on the topic.  It also requires the ability to ask questions that will lead the discussion in a particular direction.  It also requires the ability to listen to the answers and to ask further questions based on the answers.  It also requires the ability to think on your feet and to respond to questions from the audience.  The key elements of questioning are: Public Speaking, critical thinking, in'),
(31,	'Debating',	'The student has demonstrated only a basic grasp of some of the parent skills',	'The student has demonstrated a basic grasp of some of the parent skills',	'The student has demonstrated a good grasp of some of the parent skills',	'The student has demonstrated a very good grasp of most of the parent skills',	0,	'The student has demonstrated an excellent grasp of all of the parent skills',	0,	'Debating is the ability to argue a point of view in a formal setting.  It involves the ability to research a topic and to present a case for or against the topic.  It also involves the ability to listen to the opposing case and to respond to it.  It also involves the ability to think on your feet and to respond to questions from the audience.  The key elements of debating are: Public Speaking, critical thinking, persuading'),
(32,	'Poetry',	'The student has demonstrated only a basic grasp of some of the parent skills',	'The student has demonstrated a basic grasp of some of the parent skills',	'The student has demonstrated a good grasp of some of the parent skills',	'The student has demonstrated a very good grasp of most of the parent skills',	0,	'The student has demonstrated an excellent grasp of some of the parent skills',	0,	'The skill of poetry is the abilty to write and analyse poetry.  It involves the ability to understand the meaning of a poem and to convey that meaning to an audience.  It also involves the ability to write poetry that conveys a meaning to an audience.  It also involves the ability to understand the technical aspects of poetry such as rhyme, rhythm, alliteration, assonance, onomatopoeia, metaphor, simile, personification, hyperbole, symbolism, imagery, irony, oxymoron, paradox, pun, allusion, eup'),
(34,	'Drama',	'The student has demonstrated only a basic grasp of some of the parent skills',	'The student has demonstrated a basic grasp of some of the parent skills',	'The student has demonstrated a good grasp of some of the parent skills',	'The student has demonstrated a very good grasp of most of the parent skills',	0,	'The student has demonstrated an excellent grasp of all of the parent skills',	0,	'This is a composite skill involving the ability to emote and convey emotion to an audience in a number of contexts.  The contexts involve range of roles, range of emotions as well as the tecnical skills involved in film, stage or role playing.'),
(35,	'Storytelling',	'The student has demonstrated only a basic grasp of some of the parent skills',	'The student has demonstrated a basic grasp of some of the parent skills',	'The student has demonstrated a good grasp of some of the parent skills',	'The student has demonstrated a very good grasp of most of the parent skills',	0,	'The student has demonstrated an excellent grasp of all of the parent skills',	0,	'This is a composite skill involving the ability to tell a story in a number of contexts.  The contexts involve range of roles, range of emotions as well as the tecnical skills involved in film, stage or role playing.'),
(67,	'Audience Pitching',	'',	'',	'',	'',	0,	'',	0,	''),
(68,	'Summarizing',	'',	'',	'',	'',	0,	'',	0,	''),
(69,	'Contextualization',	'',	'',	'',	'',	0,	'',	0,	''),
(70,	'Deconstruction.',	'',	'',	'',	'',	0,	'',	0,	''),
(38,	'Explaining',	'The student has demonstrated only a basic grasp of some of the parent skills',	'The student has demonstrated a basic grasp of some of the parent skills',	'The student has demonstrated a good grasp of some of the parent skills',	'The student has demonstrated a very good grasp of most of the parent skills',	0,	'The student has demonstrated an excellent grasp of all of the parent skills',	0,	'The skill of explaining requires the ability to understand the area to a depth on the topic.  It also requires the ability to break down a topic so that it can be understandable to the audience.  It may require the use of analogies, or story telling.  It also requires the ability to listen to the answers and to ask further questions based on the answers.  It also requires the ability to think on your feet and to respond to questions or guage the responses from the audience.'),
(41,	'Self-regulation',	'1. The student will only be able to describe their emotions after they have acted on them',	'2. The student will be able to describe their emotions before they act on them, but have limited control of their actins on those emotions',	'3. The student will be able to describe their emotions before they act on them and will, in the main part, be able to control their actions',	'4. The student will have some control over their emotional state using such techniques as mindfulness',	0,	'5. The student will have a high level of control over their emotional state using such techniques as mindfulness and self-motivation',	0,	'The skill of self-regulation is the ability to control your emotions.  This involves being able to think before you act and the ability to change your emotional state to suit the situation, be it to calm yourself down, to get yourself motivated or to get yourself into a state of mind to perform a task.'),
(42,	'Empathy',	'1. The student will not be able to understand the emotions of others',	'2. The student will be able to understand the emotions of others, but will not be able to understand how their actions or circumstances will affect the emotions of others',	'3. The student will be able to understand the emotions of others and will be able to understand how their actions or circumstances will affect the emotions of others',	'4. The student will be able to understand the emotions of others and will be able to understand how their actions or circumstances will affect the emotions of others and how the emotions of others will affect their behaviour',	0,	'5. The student will be able to understand the emotions of others and will be able to understand how their actions or circumstances will affect the emotions of others and how the emotions of others will affect their behaviour and will be able to use this understanding to influence the emotions of others',	0,	'The skill of empathy is the ability to understand the emotions of others.  It involves the ability to understand the emotions of others and to understand how your actions or circumstances will affect the emotions of others.  It also involves the ability to understand how the emotions of others will affect their behaviour.'),
(43,	'Self-Motivation',	'1. The student will not be able to set goals or to work towards achieving those goals',	'2. The student will be able to set goals and to work towards achieving those goals, but will not be able to understand the steps that are required to achieve the goal',	'3. The student will be able to set goals and to work towards achieving those goals and will be able to understand the steps that are required to achieve the goal',	'4. The student will be able to set goals and to work towards achieving those goals and will be able to understand the steps that are required to achieve the goal and will be able to understand the consequences of not achieving the goal',	NULL,	'5. The student will be able to set goals and to work towards achieving those goals and will be able to understand the steps that are required to achieve the goal and will be able to understand the consequences of not achieving the goal and will be able to motivate themselves to achieve a goal that is difficult to achieve, or steps that the student does not necessarily enjoy',	0,	'Self-motivation is the ability to motivate yourself to achieve a goal.  It involves the ability to set goals and to work towards achieving those goals.  It also involves the ability to understand the steps that are required to achieve the goal and to understand the consequences of not achieving the goal.  Strong self-motivation involves the ability to motivate yourself to achieve a goal that is difficult to achieve, or steps that the student does not necessarily enjoy.'),
(44,	'Graph - Drawing',	'',	'',	'',	'',	0,	'',	0,	''),
(45,	'Graph - Reading',	'',	'',	'',	'',	0,	'',	0,	''),
(46,	'Data Tables - Reading',	'',	'',	'',	'',	0,	'',	0,	''),
(47,	'Data Tables - Producing',	'',	'',	'',	'',	0,	'',	0,	''),
(48,	'Application of Ideas',	'',	'',	'',	'',	0,	'',	0,	''),
(49,	'Categorization',	'',	'',	'',	'',	0,	'',	0,	''),
(50,	'Data Handling',	'',	'',	'',	'',	0,	'',	0,	''),
(51,	'Data Contextualization',	'',	'',	'',	'',	0,	'',	0,	''),
(52,	'Description',	'',	'',	'',	'',	0,	'',	0,	''),
(53,	'Equipment Handling',	'',	'',	'',	'',	0,	'',	0,	''),
(54,	'Experiment Design',	'',	'',	'',	'',	0,	'',	0,	''),
(55,	'Experiment Write Up',	'',	'',	'',	'',	0,	'',	0,	''),
(56,	'Fine Motor Skills',	'',	'',	'',	'',	0,	'',	0,	''),
(57,	'Following Protocols',	'',	'',	'',	'',	0,	'',	0,	''),
(58,	'Recall',	'',	'',	'',	'',	0,	'',	0,	''),
(59,	'Internet Searching',	'',	'',	'',	'',	0,	'',	0,	''),
(60,	'Observation',	'',	'',	'',	'',	0,	'',	0,	''),
(61,	'Organizing Information',	'',	'',	'',	'',	0,	'',	0,	''),
(62,	'Public Speaking',	'',	'',	'',	'',	0,	'',	0,	''),
(63,	'Punctuation',	'1.  Uses a limited range of punctuation',	'2.  Uses a range of punctuation',	'3.  Uses a wide range of punctuation',	'4.  Uses a wide range of punctuation',	0,	'5.  Uses a wide range of punctuation',	0,	'This skill involves the ability to use a wide range of punctuation.  It also involves the ability to use punctuation in the correct context.'),
(64,	'Social Confidence',	'',	'',	'',	'',	0,	'',	0,	''),
(65,	'Understanding Variables',	'',	'',	'',	'',	0,	'',	0,	''),
(66,	'Taking Measurements',	'',	'',	'',	'',	0,	'',	0,	'');

UNLOCK TABLES;

DROP TABLE IF EXISTS `mind_map_links`;
CREATE TABLE `mind_map_links` (
  `Mind_map_linkID` bigint NOT NULL AUTO_INCREMENT,
  `NodeA` bigint NOT NULL,
  `NodeB` bigint NOT NULL,
  `CourseID` bigint NOT NULL,
  PRIMARY KEY (`Mind_map_linkID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;


DROP TABLE IF EXISTS `mind_map_nodes`;
CREATE TABLE `mind_map_nodes` (
  `MindmapID` bigint NOT NULL AUTO_INCREMENT,
  `NodeID` bigint NOT NULL,
  `NodeX` int NOT NULL,
  `NodeY` int NOT NULL,
  `Node_name` varchar(150) NOT NULL,
  `Node_type` VARCHAR(20) NOT NULL DEFAULT 'unknown',
  `CourseID` bigint NOT NULL,
  PRIMARY KEY (`MindmapID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

UNLOCK TABLES;

LOCK TABLES `workassignment` WRITE, `lesson_files` WRITE, `skill_attached` WRITE, `markbook` WRITE, `skillsmarkbook` WRITE;

-- Seed test data for Skills Markbook (Competency Proficiency Scores)
-- This allows testing the Skills Tree color coding based on average proficiency

-- First, ensure we have a test work assignment for the student (IndividualsID = 1)
INSERT INTO `workassignment` (`IndividualsID`, `DateCollected`, `status`, `Lesson_FilesID`)
VALUES 
  (1, '2025-11-30 00:00:00','Y', 1),  -- WorkAssignmentID = 1
  (1, '2025-12-07 00:00:00','Y', 2),  -- WorkAssignmentID = 2
  (1, '2025-12-15 00:00:00','Y', 3),  -- WorkAssignmentID = 3
  (1, '2025-12-22 00:00:00','Y', 4);  -- WorkAssignmentID = 4


-- Now add skill proficiency scores (1-5) in skillsmarkbook
-- These represent the teacher's assessment of how well the student demonstrated each skill

-- Optional: Create additional markbook entries for longitudinal tracking
-- This shows how averaging works across multiple assignments
INSERT INTO `markbook` (`MarkbookID`, `Mark`, `WorkAssignmentID`, `Grade`)
VALUES 
(1, 85, 1, 'A'),
  (2, 78, 2, 'B'),
  (3, 90, 3, 'A'),
  (4, 88, 4, 'A');

INSERT INTO `skillsmarkbook` (`MarkbookID`, `Skills_keyID`, `Mark`) VALUES
  -- Same student, second assessment (creates average for color coding)
  (1, 1, 5),
  (1, 2, 4),
  (1, 3, 5),
  (1, 4, 3),
  (1, 5, 2),
  (1, 6, 1),
  (1, 7, 2),
  (1, 8, 4),
  (1, 9, 3),
  (2, 1, 4),  -- Leadership: 5 + 4 = avg 4.5 (Green)
  (2, 2, 5),  -- Communication: 4 + 5 = avg 4.5 (Green)
  (2, 3, 4),  -- Teamwork: 5 + 4 = avg 4.5 (Green)
  (2, 4, 2),  -- Problem-Solving: 3 + 2 = avg 2.5 (Orange)
  (2, 5, 2),  -- Creativity: 2 + 2 = avg 2 (Orange)
  (2, 6, 1),  -- Time Management: 1 + 1 = avg 1 (Red)
  (2, 7, 3),  -- Research: 2 + 3 = avg 2.5 (Orange)
  (3, 8, 5),  -- Collaboration: 4 + 5 = avg 4.5 (Green)
  (3, 9, 3),  -- Critical Thinking: 3 + 3 = avg 3 (Orange)
  (3, 10, 5), -- Initiative: 5 + 5 = avg 5 (Green)
  (4, 11, 4), -- Adaptability: 4 + 4 = avg 4 (Green)
  (4, 12, 2); -- Emotional Intelligence: 2 + 2 = avg 2 (Orange)


LOCK TABLES `lesson_files` WRITE, `skill_attached` WRITE;

INSERT INTO `lesson_files` 
(`Lesson_FilesID`, `Filename`, `Description`, `AudioFile`, `File_TypesID`, `mimetypesID`, `timeReq`, `LocalFoldersID`, `Hidden`, `CloudID`, `Lesson_file_favouritesID`, `MarksAvail`)
VALUES
  (4,  'strategic_foresight_guide.pdf',     'Strategic Foresight guide',             '', 1, 1, 30, 0, 0, 'cloud-004', NULL, 10),
  (5,  'risk_intelligence_guide.pdf',       'Risk Intelligence guide',               '', 1, 1, 25, 0, 0, 'cloud-005', NULL, 9),
  (6,  'opportunity_recognition_guide.pdf', 'Opportunity Recognition guide',         '', 1, 1, 20, 0, 0, 'cloud-006', NULL, 8),
  (7,  'intrapersonal_mastery_guide.pdf',   'Intrapersonal Mastery guide',           '', 1, 1, 25, 0, 0, 'cloud-007', NULL, 8),
  (8,  'interpersonal_intelligence.pdf',    'Interpersonal Intelligence resources',  '', 1, 1, 25, 0, 0, 'cloud-008', NULL, 9),
  (9,  'strategic_research_guide.pdf',      'Strategic Research materials',          '', 1, 1, 30, 0, 0, 'cloud-009', NULL, 10),
  (10, 'adaptive_planning_guide.pdf',       'Adaptive Planning guide',               '', 1, 1, 22, 0, 0, 'cloud-010', NULL, 9),
  (11, 'integrative_leadership_guide.pdf',  'Integrative Leadership resources',      '', 1, 1, 35, 0, 0, 'cloud-011', NULL, 12),
  (12, 'emotional_intelligence_guide.pdf',  'Emotional Intelligence guide',          '', 1, 1, 26, 0, 0, 'cloud-012', NULL, 10),
  (13, 'Introduction to Strategic Thinking', 'Learn the fundamentals of strategic thinking and planning.', '', 1, 1, 15, 0, 0, 'cloud-013', NULL, 5);

-- First, insert lesson_files entries for the YouTube videos
LOCK TABLES `lesson_files` WRITE;

INSERT INTO `lesson_files` 
(`Lesson_FilesID`, `Filename`, `Description`, `AudioFile`, `File_TypesID`, `mimetypesID`, `timeReq`, `LocalFoldersID`, `Hidden`, `CloudID`, `MarksAvail`)
VALUES
-- Strategic Foresight video lesson file (links to skill 1)
(14, 'Strategic Foresight - Video Resource', 'Learn to anticipate future trends and challenges through strategic foresight analysis', '', 5, 15, 10, 0, 0, 'yt-dQw4w9WgXcQ', 5),

-- Risk Intelligence video lesson file (links to skill 2)
(15, 'Risk Intelligence - Video Resource', 'Master risk assessment and mitigation strategies in complex scenarios', '', 5, 15, 8, 0, 0, 'yt-jNQXAC9IVRw', 5),

-- Opportunity Recognition video lesson file (links to skill 3)
(16, 'Opportunity Recognition - Video Resource', 'Develop skills to identify and capitalize on emerging opportunities', '', 5, 15, 12, 0, 0, 'yt-9bZkp7q19f0', 5),

-- Intrapersonal Mastery video lesson file (links to skill 4)
(17, 'Intrapersonal Mastery - Video Resource', 'Build self-awareness and emotional intelligence for personal growth', '', 5, 15, 9, 0, 0, 'yt-kJQP7kiw5Fk', 5),

-- Interpersonal Intelligence video lesson file (links to skill 5)
(18, 'Interpersonal Intelligence - Video Resource', 'Enhance communication and relationship-building skills', '', 5, 15, 11, 0, 0, 'yt-ScMzIvxBSi4', 5),

-- Strategic Research video lesson file (links to skill 6)
(19, 'Strategic Research - Video Resource', 'Learn advanced research methodologies and analysis techniques', '', 5, 15, 13, 0, 0, 'yt-oHg5SJYRHA0', 5),

-- Adaptive Planning video lesson file (links to skill 7)
(20, 'Adaptive Planning - Video Resource', 'Create flexible plans that evolve with changing circumstances', '', 5, 15, 10, 0, 0, 'yt-fJ9rUzIMcZQ', 5),

-- Integrative Leadership video lesson file (links to skill 8)
(21, 'Integrative Leadership - Video Resource', 'Master comprehensive leadership skills and organizational transformation', '', 5, 15, 15, 0, 0, 'yt-ZbZSe6N_BXs', 5),

-- Emotional Intelligence video lesson file (links to skill 9)
(22, 'Emotional Intelligence - Video Resource', 'Develop emotional awareness and relationship management skills', '', 5, 15, 12, 0, 0, 'yt-dPYE-7rkJCo', 5),

-- Introduction to Strategic Thinking video lesson file (general)
(23, 'Introduction to Strategic Thinking - Video Resource', 'Foundational concepts in strategic thinking and decision-making', '', 5, 15, 15, 0, 0, 'yt-y8Kyi0WNg40', 5);

UNLOCK TABLES;

-- Now attach these lesson files to their corresponding skills
LOCK TABLES `skill_attached` WRITE;


-- Attach the files to the nine skills (CourseID = 3)
INSERT INTO `skill_attached` (`Skills_keyID`, `Lesson_FilesID`, `CourseID`) VALUES
  (1,  4,  3),
  (2,  5,  3),
  (3,  6,  3),
  (4,  7,  3),
  (5,  8,  3),
  (6,  9,  3),
  (7,  10, 3),
  (8,  11, 3),
  (9,  12, 3);

UNLOCK TABLES;


/* Seed Lesson Files */
INSERT INTO `lesson_files` (`Lesson_FilesID`, `Filename`, `Description`, `AudioFile`, `File_TypesID`, `mimetypesID`, `timeReq`, `LocalFoldersID`, `Hidden`, `CloudID`, `Lesson_file_favouritesID`, `MarksAvail`)
VALUES
  (1, 'intro_to_strategy.pdf', 'Intro to Strategic Foresight', '', 1, 1, 30, 0, 0, 'cloud-001', NULL, 10),
  (2, 'risk_intelligence.pdf', 'Risk Intelligence overview', '', 1, 1, 25, 0, 0, 'cloud-002', NULL, 8),
  (3, 'opportunity_recognition.pdf', 'Opportunity Recognition', '', 1, 1, 20, 0, 0, 'cloud-003', NULL, 9);

/* Attach the files to skills (so skill_attached & map usage reflect the file ids above) */
INSERT INTO `skill_attached` (`Skills_keyID`, `Lesson_FilesID`, `CourseID`) VALUES
  (1, 1, 3),
  (2, 2, 3),
  (3, 3, 3),
  (4, 0, 7),
  (5, 0, 7),
  (6, 0, 7),
  (7, 0, 7),
  (8, 0, 7),
  (9, 0, 7),
  (111,0,7),
  (104,0,7),
  (105,0,7),
  (113,0,7),
  (113,0,7);


UNLOCK TABLES;
-- ...existing code...

LOCK TABLES `skill_link` WRITE;

INSERT INTO `skill_link` (`ParentSkillID`, `OffspringSkillID`, `CourseID`) VALUES
-- 1. Connecting Foundational Skills to Intermediate Skills
(111, 7, 7),  -- Goal Setting (111) -> Adaptive Planning (7)
(104, 5, 7),  -- Managing Conflict (104) -> Interpersonal Intelligence (5)
(105, 4, 7),  -- Self-Advocacy (105) -> Intrapersonal Mastery (4)

-- 2. Connecting Intermediate Skills to Capstone Skills (The course goals: 8 & 9)
(7, 8, 7),    -- Adaptive Planning (7) -> Integrative Leadership (8)
(4, 9, 7),    -- Intrapersonal Mastery (4) -> Emotional Intelligence (9)
(5, 9, 7),    -- Interpersonal Intelligence (5) -> Emotional Intelligence (9)
-- EMotional intelligence to zen warrior
(113, 9, 7),  -- Zen Warrior Mindset (113) -> Emotional Intelligence (9)
(2,7,7),
(3,7,7),
(6,8,7);
UNLOCK TABLES;

LOCK TABLES `skill_attached` WRITE;

-- Attaching key Capstone Skills (The ultimate goals) to ETCAL (CourseID 7)
-- We attach them to a high-level file (e.g., the Course Overview or Final Project file, ID 1001)

INSERT INTO `skill_attached` (`Skills_keyID`, `Lesson_FilesID`, `CourseID`) VALUES
-- Strategic Foresight (1): High-level goal skill
(1, 1001, 7),
-- Intrapersonal Mastery (4): High-level goal skill
(4, 1001, 7),
-- Interpersonal Intelligence (5): High-level goal skill
(5, 1001, 7),
-- Integrative Leadership (8): High-level goal skill
(8, 1001, 7),
-- Emotional Intelligence (9): High-level goal skill
(9, 1001, 7),


-- Attach a second link for redundancy/completeness, using another high-level file ID
(1, 1024, 7),
(4, 1024, 7),
(5, 1024, 7),
(8, 1024, 7),
(9, 1024, 7);

UNLOCK TABLES;

LOCK TABLES `skills_key` WRITE;

INSERT INTO `skills_key` (`Skills_keyID`, `Skill_name`, `Criteria1`, `Criteria2`, `Criteria3`, `Criteria4`, `Criteria5`, `hidden`, `DevelopmentAge`, `SkillDescription`) VALUES
(101, 'Personal Finance Awareness', 'Can identify different types of income.', 'Can create a basic personal budget.', 'Can distinguish needs from wants.', '', '', 0, 'Level 1', 'Understanding and managing basic personal finances.'),
(102, 'Basic Digital Safety', 'Can create a strong password.', 'Can identify a phishing attempt.', 'Understands safe social media sharing.', '', '', 0, 'Level 1', 'Knowledge of fundamental online security and safety practices.'),
(103, 'Healthy Lifestyle Choices', 'Can list components of a balanced diet.', 'Can track sleep and hydration.', 'Can identify stress symptoms.', '', '', 0, 'Level 1', 'Awareness of choices that support physical and mental wellbeing.'),
(104, 'Managing Personal Conflict', 'Can define mediation.', 'Can identify non-aggressive communication.', 'Can seek help from an appropriate source.', '', '', 0, 'Level 1', 'Skills for resolving disagreements constructively.'),
(105, 'Self-Advocacy', 'Can articulate a personal need clearly.', 'Can state a preference politely.', 'Can ask clarifying questions.', '', '', 0, 'Level 1', 'Ability to communicate one''s needs and rights effectively.'),
(106, 'Recognising Diversity', 'Can define diversity (e.g., race, family).', 'Can state two forms of respectful behavior.', 'Can identify cultural differences.', '', '', 0, 'Level 1', 'Understanding and respect for varied social and cultural backgrounds.'),
(107, 'Active Citizenship', 'Can name a local civic body.', 'Can list three community responsibilities.', 'Understands basic voting rights.', '', '', 0, 'Level 1', 'Awareness of rights, responsibilities, and civic participation.'),
(108, 'Media Influence Awareness', 'Can distinguish news from opinion.', 'Can identify an advertisement.', 'Understands emotional impact of social media.', '', '', 0, 'Level 1', 'Ability to critically analyze information sources.'),
(109, 'Establishing Boundaries', 'Can define personal space.', 'Can appropriately say "no".', 'Can identify inappropriate behavior.', '', '', 0, 'Level 1', 'Setting and communicating limits in relationships.'),
(110, 'Coping Strategies', 'Can name three stress-relief methods.', 'Can identify a trusted person to talk to.', 'Can practice a relaxation technique.', '', '', 0, 'Level 1', 'Techniques for dealing with stress, failure, and disappointment.'),
(111, 'Goal Setting (Personal)', 'Can set a short-term goal.', 'Can break a goal into two steps.', 'Can track progress toward a goal.', '', '', 0, 'Level 1', 'Ability to define and pursue personal targets.'),
(112, 'Identifying Support Needs', 'Can name different types of support (e.g., medical, academic).', 'Can ask for help.', 'Knows two sources of professional help.', '', '', 0, 'Level 1', 'Recognizing and obtaining necessary help.'),
(113, 'Zen Warrior Mindset', 'Can describe mindfulness.', 'Can practice deep breathing.', 'Can identify positive self-talk.', '', '', 0, 'Level 1', 'Cultivating resilience and a positive mindset through mindfulness.'),
(114, 'Time Management Basics', 'Can list daily tasks.', 'Can prioritize tasks.', 'Can use a simple planner.', '', '', 0, 'Level 1', 'Fundamental skills for organizing and managing time effectively.'),
(115, 'Financial Awareness', 'Can identify basic financial terms (e.g., saving, spending).', 'Can create a simple budget.', 'Understands the concept of interest.', '', '', 0, 'Level 1', 'Basic understanding of personal finance and money management.');

UNLOCK TABLES;