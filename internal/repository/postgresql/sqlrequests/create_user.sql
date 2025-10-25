INSERT INTO effectivemobile.users (uuid)
VALUES ($1)
ON CONFLICT (uuid) DO NOTHING;