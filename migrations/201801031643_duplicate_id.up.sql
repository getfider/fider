ALTER TABLE ideas ADD duplicate_id INT NULL REFERENCES ideas(id);
