# Create

```sql
create user api_production_admin with password '';
create database api_production owner api_production_admin;
grant all privileges on database api_production to api_production_admin;

create user api_testing_admin with password '';
create database api_testing owner api_testing_admin;
grant all privileges on database api_testing to api_testing_admin;
```
