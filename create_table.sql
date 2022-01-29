CREATE TABLE users (
  userid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  streak INT,
  higheststreak INT,
  totallearned INT,
  reviewpoints INT,
  lastlearned BIGINT,
  usestraditional BOOLEAN
);