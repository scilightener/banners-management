CREATE TABLE feature (
    id INT PRIMARY KEY
);

CREATE TABLE tag (
    id INT PRIMARY KEY
);

CREATE TABLE banner (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title TEXT NOT NULL,
    text TEXT,
    url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    feature_id INT REFERENCES feature(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX banner_feature_id ON banner(feature_id);

CREATE TABLE banner_tag (
    banner_id INT REFERENCES banner(id) ON DELETE CASCADE,
    tag_id INT REFERENCES tag(id) ON DELETE CASCADE,
    PRIMARY KEY (banner_id, tag_id)
);

CREATE INDEX banner_tag_tag_id ON banner_tag(tag_id);

CREATE OR REPLACE FUNCTION check_feature_tag_unique()
    RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM banner b JOIN banner_tag bt ON b.id = bt.banner_id
        WHERE b.feature_id = NEW.feature_id
        GROUP BY bt.tag_id
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION 'Duplicate banner tags found for the given feature_id';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE CONSTRAINT TRIGGER unique_feature_tag_trigger
    AFTER INSERT ON banner
    INITIALLY DEFERRED
    FOR EACH ROW
EXECUTE FUNCTION check_feature_tag_unique();

