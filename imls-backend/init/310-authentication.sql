-- migrate:up

-- We put things inside the basic_auth schema to hide
-- them from public view. Certain public procs/views will
-- refer to helpers and tables inside.


-- We don't need to know the JWT secret on the client side. The server needs to know
-- it to generate tokens and check tokens that come back to us. The client, however,
-- never generates a token. Therefore, this can be randomly generated every time we start the server.
-- If we catch someone between requesting a token and making a POST, they'll be screwed. However,
-- if they back off and try again (including grabbing the token), they would then get a new token.
-- So, ourbackoff should be 3-5m.
-- ALTER DATABASE imls SET "app.jwt_secret" TO select substr(md5(random()::text), 0, 25);


create or replace function
basic_auth.check_role_exists() returns trigger as $$
begin
  if not exists (select 1 from pg_roles as r where r.rolname = new.role) then
    raise foreign_key_violation using message =
      'unknown database role: ' || new.role;
    return null;
  end if;
  return new;
end
$$ language plpgsql;

drop trigger if exists ensure_user_role_exists on basic_auth.users;
create constraint trigger ensure_user_role_exists
  after insert or update on basic_auth.users
  for each row
  execute procedure basic_auth.check_role_exists();

create extension if not exists pgcrypto;

create or replace function
basic_auth.encrypt_pass() returns trigger as $$
begin
  if tg_op = 'INSERT' or new.api_key <> old.api_key then
    new.api_key = crypt(new.api_key, gen_salt('bf'));
  end if;
  return new;
end
$$ language plpgsql;

drop trigger if exists encrypt_pass on basic_auth.users;
create trigger encrypt_pass
  before insert or update on basic_auth.users
  for each row
  execute procedure basic_auth.encrypt_pass();

CREATE OR REPLACE FUNCTION basic_auth.user_role(fscs_id text, api_key text)
 RETURNS name
 LANGUAGE plpgsql
AS $function$
begin
  	return (
  		select role from basic_auth.users
   			where users.fscs_id = user_role.fscs_id
     			and users.api_key = crypt(user_role.api_key, users.api_key)
			);
end;
$function$
;

-- add type
CREATE TYPE basic_auth.jwt_token AS (
  token text
);

-- login should be on your exposed schema
create or replace function
api.login(fscs_id text, api_key text) returns basic_auth.jwt_token as $$
declare
  _role name;
  result basic_auth.jwt_token;
begin
  -- check email and password
  select basic_auth.user_role(login.fscs_id, login.api_key) into _role;
  if _role is null then
    raise invalid_password using message = 'invalid user or password';
  end if;

  select sign(
      row_to_json(r), current_setting('app.jwt_secret')
    ) as token
    from (
      select _role as role, login.fscs_id as fscs_id,
         extract(epoch from now())::integer + 60*60 as exp
    ) r
    into result;
  return result;
end;
$$ language plpgsql security definer;

-- the names "anon" and "authenticator" are configurable and not
-- sacred, we simply choose them for clarity

