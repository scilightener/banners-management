INSERT INTO feature (id)
SELECT generate_series(1, 1000);

INSERT INTO tag (id)
SELECT generate_series(1, 1000);