-- Create the Datadog User
CREATE USER datadog@'%' IDENTIFIED by 'datadog123';
ALTER USER datadog@'%' WITH MAX_USER_CONNECTIONS 5;
GRANT REPLICATION CLIENT ON *.* TO datadog@'%';
GRANT PROCESS ON *.* TO datadog@'%';
GRANT SELECT ON performance_schema.* TO datadog@'%';

-- Create Schema
CREATE SCHEMA IF NOT EXISTS datadog;
GRANT EXECUTE ON datadog.* to datadog@'%';
GRANT CREATE TEMPORARY TABLES ON datadog.* TO datadog@'%';

-- Create explain_statement
DELIMITER $$
CREATE PROCEDURE datadog.explain_statement(IN query TEXT)
    SQL SECURITY DEFINER
BEGIN
    SET @explain := CONCAT('EXPLAIN FORMAT=json ', query);
    PREPARE stmt FROM @explain;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END $$
DELIMITER ;


-- Explain Plan Procedure for Custom Schema
-- DELIMITER $$
-- CREATE PROCEDURE <YOUR_SCHEMA>.explain_statement(IN query TEXT)
--     SQL SECURITY DEFINER
-- BEGIN
--     SET @explain := CONCAT('EXPLAIN FORMAT=json ', query);
--     PREPARE stmt FROM @explain;
--     EXECUTE stmt;
--     DEALLOCATE PREPARE stmt;
-- END $$
-- DELIMITER ;
-- GRANT EXECUTE ON PROCEDURE <YOUR_SCHEMA>.explain_statement TO datadog@'%';


-- Runtime Setup Consumers
DELIMITER $$
CREATE PROCEDURE datadog.enable_events_statements_consumers()
    SQL SECURITY DEFINER
BEGIN
    UPDATE performance_schema.setup_consumers SET enabled='YES' WHERE name LIKE 'events_statements_%';
    UPDATE performance_schema.setup_consumers SET enabled='YES' WHERE name = 'events_waits_current';
END $$
DELIMITER ;
GRANT EXECUTE ON PROCEDURE datadog.enable_events_statements_consumers TO datadog@'%';

