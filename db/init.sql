CREATE TABLE url (
  mapping_id SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  old_url    VARCHAR(51) NOT NULL UNIQUE,
  shortkey   VARCHAR(51) NOT NULL UNIQUE
); 
