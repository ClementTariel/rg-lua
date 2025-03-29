CREATE USER referee_user WITH PASSWORD 'referee_temporary_password';
CREATE USER matchmaker_user WITH PASSWORD 'matchmaker_temporary_password';
CREATE USER rglua_user WITH PASSWORD 'rglua_temporary_password';
CREATE DATABASE rglua;

\c rglua;

CREATE OR REPLACE FUNCTION load_bot(path TEXT) RETURNS TEXT AS $$
  SELECT CAST(pg_read_file($1) AS TEXT);
$$ LANGUAGE sql SECURITY DEFINER;
ALTER FUNCTION load_bot(TEXT) OWNER TO postgres;

CREATE TABLE users (id UUID PRIMARY KEY, name VARCHAR (20));
CREATE TABLE bots (id UUID PRIMARY KEY, name VARCHAR (50), script TEXT, userId UUID REFERENCES users (id) , userName VARCHAR (20));
-- TODO: WIP create table matchs

GRANT SELECT ON TABLE bots TO referee_user;
GRANT SELECT ON TABLE bots TO matchmaker_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE bots TO rglua_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE users TO rglua_user;
INSERT INTO users(id, name) VALUES('f58ee1e8-b9b8-4a63-b6a1-4188f9e1e4f8', 'Examples');
INSERT INTO bots(id, name, script, userId, userName)  VALUES('3cd678c8-7864-477d-b93a-8347b7a61392', 'paper', load_bot(CAST('/bots/paper.lua' AS TEXT)), 'f58ee1e8-b9b8-4a63-b6a1-4188f9e1e4f8', 'Examples');
INSERT INTO bots(id, name, script, userId, userName)  VALUES('d3954130-d94a-4bd8-be57-e6f424d7d43b', 'random', load_bot(CAST('/bots/random.lua' AS TEXT)), 'f58ee1e8-b9b8-4a63-b6a1-4188f9e1e4f8', 'Examples');
INSERT INTO bots(id, name, script, userId, userName)  VALUES('5a6b0f62-af13-4f13-a9e7-8e5f63494c09', 'rock', load_bot(CAST('/bots/rock.lua' AS TEXT)), 'f58ee1e8-b9b8-4a63-b6a1-4188f9e1e4f8', 'Examples');
INSERT INTO bots(id, name, script, userId, userName)  VALUES('0b40a114-1556-47f8-be57-cb1000d0c31e', 'scissors', load_bot(CAST('/bots/scissors.lua' AS TEXT)), 'f58ee1e8-b9b8-4a63-b6a1-4188f9e1e4f8', 'Examples');
