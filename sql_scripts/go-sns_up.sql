CREATE TABLE IF NOT EXISTS "events" (
ID  SERIAL PRIMARY KEY 
,  "inserted_at" TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
, "ip_addr" VARCHAR NOT NULL
, "mac_addr" VARCHAR NOT NULL
, "subject" VARCHAR NOT NULL
, "message" VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS "messages"(
    ID SERIAL PRIMARY KEY
    , "code" VARCHAR NOT NULL
    , "subject" VARCHAR NOT NULL
    , "body" VARCHAR NOT NULL
);