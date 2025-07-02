-- Add is_approved column to posts table
ALTER TABLE posts ADD is_approved BOOLEAN NULL;
UPDATE posts SET is_approved = true;
ALTER TABLE posts ALTER COLUMN is_approved SET NOT NULL;

-- Add is_approved column to comments table  
ALTER TABLE comments ADD is_approved BOOLEAN NULL;
UPDATE comments SET is_approved = true;
ALTER TABLE comments ALTER COLUMN is_approved SET NOT NULL;