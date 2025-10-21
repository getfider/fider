-- 1. Add language column
ALTER TABLE posts ADD COLUMN language regconfig;

-- 2. Create language mapping function
CREATE OR REPLACE FUNCTION map_language_to_tsvector(lang TEXT)
RETURNS regconfig AS $$
BEGIN
	RETURN CASE lang
        WHEN 'ar' THEN 'arabic'::regconfig
        WHEN 'cs' THEN 'simple'::regconfig   -- Czech not supported, fallback to simple
        WHEN 'de' THEN 'german'::regconfig
        WHEN 'el' THEN 'simple'::regconfig   -- Greek not supported, fallback to simple
        WHEN 'en' THEN 'english'::regconfig
        WHEN 'es-ES' THEN 'spanish'::regconfig
        WHEN 'fa' THEN 'simple'::regconfig    -- Farsi not supported, fallback to simple
        WHEN 'fr' THEN 'french'::regconfig
        WHEN 'it' THEN 'italian'::regconfig
        WHEN 'ja' THEN 'simple'::regconfig    -- Japanese not supported, fallback to simple
        WHEN 'ko' THEN 'simple'::regconfig    -- Korean not supported, fallback to simple
        WHEN 'nl' THEN 'dutch'::regconfig
        WHEN 'pl' THEN 'simple'::regconfig    -- Polish not supported, fallback to simple
        WHEN 'pt-BR' THEN 'portuguese'::regconfig
        WHEN 'ru' THEN 'russian'::regconfig
        WHEN 'sk' THEN 'simple'::regconfig    -- Slovak not supported, fallback to simple
        WHEN 'si-LK' THEN 'simple'::regconfig   -- Sinhala not supported, fallback to simple
        WHEN 'sv-SE' THEN 'swedish'::regconfig
        WHEN 'tr' THEN 'turkish'::regconfig
        WHEN 'zh-CN' THEN 'simple'::regconfig -- Chinese not supported, fallback to simple
		ELSE 'simple'::regconfig
	END;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- 3. Update language column using tenant's language (assuming you have a way to join or fetch it)
UPDATE posts
SET language = map_language_to_tsvector(
    (SELECT language::TEXT FROM tenants WHERE tenants.id = posts.tenant_id)
);

-- 4. Add search as a generated column
ALTER TABLE posts ADD search tsvector GENERATED ALWAYS AS (
    CASE WHEN language <> 'simple'::regconfig THEN
        setweight(to_tsvector(language, title::TEXT), 'A') ||
        setweight(to_tsvector(language, description), 'B') ||
        setweight(to_tsvector('simple'::regconfig, title::TEXT), 'C') ||
        setweight(to_tsvector('simple'::regconfig, description), 'D')::tsvector
    ELSE
        setweight(to_tsvector('simple'::regconfig, title::TEXT), 'A') ||
        setweight(to_tsvector('simple'::regconfig, description), 'B')::tsvector
    END
) STORED;

-- 5. Create GIN index
CREATE INDEX idx_posts_search_gin ON posts USING GIN (search);
