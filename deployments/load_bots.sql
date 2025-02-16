\c rglua;

INSERT INTO bots (id, name, script, userId, userName)
  VALUES (:v1, :v2, load_bot(CAST(:v3 AS TEXT)), :v4, 'Examples');
