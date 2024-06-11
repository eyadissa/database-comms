DROP DATABASE IF EXISTS MSDS;
CREATE DATABASE MSDS;

\c msds;

DROP TABLE IF EXISTS classData;

CREATE TABLE classData (
                           ID SERIAL,
                           CID VARCHAR(3) PRIMARY KEY,
                           Name VARCHAR(100),
                           Prereq VARCHAR(100)
);
