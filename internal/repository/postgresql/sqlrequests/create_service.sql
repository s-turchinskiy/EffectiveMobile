INSERT INTO effectivemobile.services (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING;