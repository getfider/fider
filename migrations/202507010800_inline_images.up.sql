-- Update comments with all attachment images
UPDATE comments 
SET content = CASE
    WHEN content IS NULL OR content = '' THEN image_markdown
    ELSE content || E'\n\n' || image_markdown
END
FROM (
    SELECT 
        comment_id,
        string_agg('![](fider-image:' || attachment_bkey || ')', E'\n') as image_markdown
    FROM attachments 
    WHERE comment_id IS NOT NULL
    GROUP BY comment_id
) a
WHERE comments.id = a.comment_id;

-- Update posts with all attachment images  
UPDATE posts
SET description = CASE
    WHEN description IS NULL OR description = '' THEN image_markdown
    ELSE description || E'\n\n' || image_markdown
END
FROM (
    SELECT 
        post_id,
        string_agg('![](fider-image:' || attachment_bkey || ')', E'\n') as image_markdown
    FROM attachments 
    WHERE comment_id IS NULL
    GROUP BY post_id
) a
WHERE posts.id = a.post_id;