CREATE TABLE url (
  mapping_id SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  last_used  TIMESTAMP,
  old_url    VARCHAR(50) NOT NULL,
  new_url    VARCHAR(50) NOT NULL
); 
