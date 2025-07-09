-- +migrate Up

update current_listings cl set 
 mls_number = substring(adl.mls_number from 4) 
from ad_df_listing adl
where cl.am_listing_id = adl.am_id 
and cl.ouid = 'M00000146' 
and length(cl.mls_number) > 3;

-- +migrate Down
