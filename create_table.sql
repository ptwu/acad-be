CREATE TABLE users (
  userid TEXT PRIMARY KEY,
  streak INT,
  highestStreak INT,
  totalLearned INT,
  reviewPoints INT,
  lastLearned BIGINT,
  usesTraditional BOOLEAN,
  age INT,
  location TEXT
);