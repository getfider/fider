-- Add a tsvector column derived from posts.response so that collaborators/administrators
-- can search across status responses (e.g. "wpada w v2.3.1").
-- Mirrors the structure of posts.search (migration 202601191200) but operates on response.

ALTER TABLE posts ADD search_response tsvector GENERATED ALWAYS AS (
    CASE
        WHEN response IS NULL THEN NULL
        WHEN language <> 'simple' THEN
            setweight(to_tsvector(map_language_to_tsvector(language), response), 'A') ||
            setweight(to_tsvector('simple'::regconfig, response), 'B')
        ELSE
            setweight(to_tsvector('simple'::regconfig, response), 'A')
    END
) STORED;

CREATE INDEX idx_posts_search_response_gin ON posts USING GIN (search_response);
