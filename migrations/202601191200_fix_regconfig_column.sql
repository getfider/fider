-- Fix: Change language column from regconfig to TEXT to allow pg_upgrade
-- The regconfig type stores OIDs that change between PostgreSQL versions,
-- breaking pg_upgrade. Store as TEXT and convert at query time instead.

-- 1. Drop the generated column that depends on language
ALTER TABLE posts DROP COLUMN search;

-- 2. Change language column from regconfig to TEXT
-- The regconfig values like 'english'::regconfig convert to their name when cast to TEXT
ALTER TABLE posts ALTER COLUMN language TYPE TEXT USING language::TEXT;

-- 3. Drop and recreate the mapping function to accept config names directly
-- Must be IMMUTABLE for use in generated columns
DROP FUNCTION IF EXISTS map_language_to_tsvector(TEXT);
CREATE FUNCTION map_language_to_tsvector(config_name TEXT)
RETURNS regconfig AS $$
BEGIN
    RETURN CASE config_name
        WHEN 'arabic' THEN 'arabic'::regconfig
        WHEN 'dutch' THEN 'dutch'::regconfig
        WHEN 'english' THEN 'english'::regconfig
        WHEN 'french' THEN 'french'::regconfig
        WHEN 'german' THEN 'german'::regconfig
        WHEN 'italian' THEN 'italian'::regconfig
        WHEN 'portuguese' THEN 'portuguese'::regconfig
        WHEN 'russian' THEN 'russian'::regconfig
        WHEN 'spanish' THEN 'spanish'::regconfig
        WHEN 'swedish' THEN 'swedish'::regconfig
        WHEN 'turkish' THEN 'turkish'::regconfig
        -- Unsupported languages fall back to 'simple':
        -- Chinese, Czech, Greek, Japanese, Korean, Persian, Polish, Sinhala, Slovak
        ELSE 'simple'::regconfig
    END;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- 4. Recreate the search column using the IMMUTABLE function
ALTER TABLE posts ADD search tsvector GENERATED ALWAYS AS (
    CASE WHEN language <> 'simple' THEN
        setweight(to_tsvector(map_language_to_tsvector(language), title::TEXT), 'A') ||
        setweight(to_tsvector(map_language_to_tsvector(language), description), 'B') ||
        setweight(to_tsvector('simple'::regconfig, title::TEXT), 'C') ||
        setweight(to_tsvector('simple'::regconfig, description), 'D')::tsvector
    ELSE
        setweight(to_tsvector('simple'::regconfig, title::TEXT), 'A') ||
        setweight(to_tsvector('simple'::regconfig, description), 'B')::tsvector
    END
) STORED;

-- 5. Recreate the GIN index
CREATE INDEX idx_posts_search_gin ON posts USING GIN (search);
