ALTER TABLE users
    ALTER COLUMN password_hash TYPE CHAR(60);
ALTER TABLE users
    ALTER COLUMN public_key TYPE TEXT;