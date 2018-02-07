INSERT INTO idea_subscribers (user_id, idea_id, created_on, updated_on, status) 
SELECT user_id, idea_id, created_on, created_on, 1 FROM idea_supporters  ON CONFLICT DO NOTHING;

INSERT INTO idea_subscribers (user_id, idea_id, created_on, updated_on, status) 
SELECT user_id, id, created_on, created_on, 1 FROM ideas ON CONFLICT DO NOTHING;

INSERT INTO idea_subscribers (user_id, idea_id, created_on, updated_on, status) 
SELECT user_id, idea_id, created_on, created_on, 1 FROM comments ON CONFLICT DO NOTHING;