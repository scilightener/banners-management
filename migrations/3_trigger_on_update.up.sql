DROP TRIGGER IF EXISTS unique_feature_tag_trigger ON banner;

CREATE CONSTRAINT TRIGGER unique_feature_tag_trigger
    AFTER INSERT OR UPDATE ON banner
    INITIALLY DEFERRED
    FOR EACH ROW
EXECUTE FUNCTION check_feature_tag_unique();