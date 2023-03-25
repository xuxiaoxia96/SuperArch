/*
 Navicat Premium Data Transfer

 Source Server         : mypostgres
 Source Server Type    : PostgreSQL
 Source Server Version : 140004
 Source Host           : localhost:5432
 Source Catalog        : admin
 Source Schema         : superarch

 Target Server Type    : PostgreSQL
 Target Server Version : 140004
 File Encoding         : 65001

 Date: 12/11/2022 09:39:51
*/


-- ----------------------------
-- Table structure for TaskType
-- ----------------------------
DROP TABLE IF EXISTS "superarch"."TaskType";
CREATE TABLE "superarch"."TaskType" (
  "id" serial NOT NULL,
  "module" varchar(64) COLLATE "pg_catalog"."default",
  "settings" jsonb
)
;
ALTER TABLE "superarch"."TaskType" OWNER TO "admin";

-- ----------------------------
-- Primary Key structure for table TaskType
-- ----------------------------
ALTER TABLE "superarch"."TaskType" ADD CONSTRAINT "TaskType_pkey" PRIMARY KEY ("id");
