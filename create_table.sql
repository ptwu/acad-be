CREATE TABLE users (
  userid TEXT PRIMARY KEY,
  streak INT,
  higheststreak INT,
  totallearned INT,
  reviewpoints INT,
  lastlearned BIGINT,
  usestraditional BOOLEAN
);