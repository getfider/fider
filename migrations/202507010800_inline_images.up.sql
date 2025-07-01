-- Update comments with attachment images
UPDATE comments 
SET content = CASE
    WHEN content IS NULL OR content = '' THEN '![](fider_image:' || a.attachment_bkey || ')'
    ELSE content || E'\n\n![](fider_image:' || a.attachment_bkey || ')'
END
FROM attachments a
WHERE comments.id = a.comment_id 
AND a.comment_id IS NOT NULL;

-- Update posts with attachment images  
UPDATE posts
SET description = CASE
    WHEN description IS NULL OR description = '' THEN '![](fider_image:' || a.attachment_bkey || ')'
    ELSE description || E'\n\n![](fider_image:' || a.attachment_bkey || ')'
END
FROM attachments a
WHERE posts.id = a.post_id
AND a.comment_id IS NULL;