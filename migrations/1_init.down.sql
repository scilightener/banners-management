DROP TRIGGER IF EXISTS unique_feature_tag_trigger ON banner;

DROP INDEX IF EXISTS banner_tag_tag_id;

DROP TABLE IF EXISTS banner_tag;

DROP INDEX IF EXISTS banner_feature_id;

DROP TABLE IF EXISTS banner;
DROP TABLE IF EXISTS feature;
DROP TABLE IF EXISTS tag;
