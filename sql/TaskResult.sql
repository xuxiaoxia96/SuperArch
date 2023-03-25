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

 Date: 12/11/2022 09:39:43
*/


-- ----------------------------
-- Table structure for TaskResult
-- ----------------------------
DROP TABLE IF EXISTS "superarch"."TaskResult";
CREATE TABLE "superarch"."TaskResult" (
  "id" serial NOT NULL,
  "request_id" varchar(64) COLLATE "pg_catalog"."default",
  "module" varchar(64) COLLATE "pg_catalog"."default",
  "result" jsonb,
  "insert_time" date DEFAULT CURRENT_TIMESTAMP
)
;
ALTER TABLE "superarch"."TaskResult" OWNER TO "admin";

-- ----------------------------
-- Primary Key structure for table TaskResult
-- ----------------------------
ALTER TABLE "superarch"."TaskResult" ADD CONSTRAINT "Task_pkey" PRIMARY KEY ("id");
