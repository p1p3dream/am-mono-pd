# Create

```sql
create user saas_production_admin with password '';
create database saas_production owner saas_production_admin;
grant all privileges on database saas_production to saas_production_admin;

create user saas_testing_admin with password '';
create database saas_testing owner saas_testing_admin;
grant all privileges on database saas_testing to saas_testing_admin;
```
