-- PostgreSQL Direct COPY Import Script for Photo Index data
-- This script updates photo_key in ad_df_listing table based on data from Photo_Index_ListingAnalytics_TAB.txt
--
-- Test locally with included sample file:
--   psql -h postgres.abodemine.local -U abodemine -d datapipe -f 20250514000001-update-photokeys.sql 
--
-- For production with full dataset:
--   1. Replace Photo_Index_ListingAnalytics_TAB.txt in the same directory as this script
--   2. Run: psql -d your_database -f update_photo_keys.sql

BEGIN;

-- Create a temporary table that matches the TXT file structure
CREATE TEMPORARY TABLE temp_photo_index (
    mls_listing_id BIGINT,
    photo_key TEXT
);

-- Import data from TXT file using COPY
\COPY temp_photo_index FROM 'Photo_Index_ListingAnalytics_TAB.txt' WITH (FORMAT text, DELIMITER E'\t', HEADER true);

-- Report the number of rows loaded
SELECT COUNT(*) AS "Rows loaded into temporary table" FROM temp_photo_index;

-- Preview the data to ensure correct import
SELECT * FROM temp_photo_index LIMIT 5;

-- Create index on temporary table for faster join with large datasets
CREATE INDEX idx_temp_photo_index_mls_listing_id ON temp_photo_index (mls_listing_id);

-- Update the ad_df_listing table with photo_key values from the imported data
UPDATE ad_df_listing
SET photo_key = temp_photo_index.photo_key
FROM temp_photo_index
WHERE ad_df_listing.mls_listing_id = temp_photo_index.mls_listing_id;

-- Report the number of rows updated
SELECT COUNT(*) AS "Rows updated in ad_df_listing" 
FROM ad_df_listing
JOIN temp_photo_index ON ad_df_listing.mls_listing_id = temp_photo_index.mls_listing_id;

-- Drop the temp table
DROP TABLE temp_photo_index;

COMMIT;