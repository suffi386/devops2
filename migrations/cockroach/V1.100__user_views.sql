ALTER TABLE adminapi.iam_members ADD COLUMN machine_description STRING;
ALTER TABLE auth.user_grants ADD COLUMN machine_description STRING;
ALTER TABLE authz.user_grants ADD COLUMN machine_description STRING;
ALTER TABLE management.user_grants ADD COLUMN machine_description STRING;
ALTER TABLE management.project_grant_members ADD COLUMN machine_description STRING;
ALTER TABLE management.org_members ADD COLUMN machine_description STRING;
ALTER TABLE management.project_members ADD COLUMN machine_description STRING;

--TODO: (adminapi)iam_members
--TODO: (auth,authz,management)user_grants --authz doesn't use this fields
--TODO: (auth)user_sessions -> maybe no changes needed --are user sessions also for machines?
--TODO: (management)project_grant_members 
--TODO: (management)org_members
--TODO: (management)project_members