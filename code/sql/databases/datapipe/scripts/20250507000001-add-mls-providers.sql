-- PostgreSQL Direct COPY Import Script for mls_providers table
-- This script creates the mls_providers table and directly imports data from mls_attom_reso_matching.csv
-- Source file: https://rentnectar.sharepoint.com/:x:/s/RentNectarManagement/ETKx3BTYYQpAvo4cOzlnRHgBt491V61IleHkXkcai0lzDg?e=91zggb

BEGIN;

CREATE TABLE IF NOT EXISTS mls_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ouid TEXT NOT NULL,
    mls_source TEXT NOT NULL
);

-- Create a temporary table that matches the CSV structure
CREATE TEMPORARY TABLE temp_reso_import (
    "MLS Source" TEXT,
    "OUID" TEXT
);

-- Import data from CSV file using COPY
\COPY temp_reso_import FROM './mls_attom_reso_matching.csv' WITH (FORMAT csv, HEADER true);

-- Check if the CSV was loaded properly
SELECT COUNT(*) AS "Rows loaded into temporary table" FROM temp_reso_import;

-- Truncate the mls_providers table to remove old data
TRUNCATE TABLE mls_providers;

-- Insert the data into the mls_providers table
INSERT INTO mls_providers (ouid, mls_source)
SELECT 
    "OUID" AS ouid,
    "MLS Source" AS mls_source
FROM temp_reso_import
WHERE "OUID" IS NOT NULL AND "MLS Source" IS NOT NULL;

-- Add indexes 
CREATE INDEX IF NOT EXISTS idx_mls_providers_ouid ON mls_providers (ouid);
CREATE INDEX IF NOT EXISTS idx_mls_providers_mls_source ON mls_providers (mls_source);

-- Drop the temp table
DROP TABLE temp_reso_import;

-- Report the number of rows imported
SELECT COUNT(*) AS "Rows Imported" FROM mls_providers;

COMMIT;