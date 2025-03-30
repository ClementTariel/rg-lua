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
INSERT INTO users(id, name) VALUES('4679ec66-8a17-409d-a0dd-7074e6dae3d0', 'Examples');
INSERT INTO bots(id, name, script, userId, userName)  VALUES('e58d370f-4155-4b55-a06e-2d86022156b7', 'paper', load_bot(CAST('/bots/paper.lua' AS TEXT)), '4679ec66-8a17-409d-a0dd-7074e6dae3d0', 'Examples');
INSERT INTO bots(id, name, script, userId, userName)  VALUES('8f66a094-497b-4e11-adf7-1dcd15b32c63', 'random', load_bot(CAST('/bots/random.lua' AS TEXT)), '4679ec66-8a17-409d-a0dd-7074e6dae3d0', 'Examples');
INSERT INTO bots(id, name, script, userId, userName)  VALUES('2848cbbc-3789-4fea-a0e2-1951fed76d44', 'rock', load_bot(CAST('/bots/rock.lua' AS TEXT)), '4679ec66-8a17-409d-a0dd-7074e6dae3d0', 'Examples');
INSERT INTO bots(id, name, script, userId, userName)  VALUES('a697d721-4242-49aa-aff9-887805769d6d', 'scissors', load_bot(CAST('/bots/scissors.lua' AS TEXT)), '4679ec66-8a17-409d-a0dd-7074e6dae3d0', 'Examples');
