DO $$ BEGIN
    CREATE TYPE member_status AS ENUM ('active', 'inactive', 'pending');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

ALTER TABLE members
    ADD COLUMN email TEXT,
    ADD COLUMN plan TEXT,
    ADD COLUMN status member_status NOT NULL DEFAULT 'active';

CREATE INDEX IF NOT EXISTS idx_members_email ON members(email);

ALTER TABLE members
    DROP COLUMN IF EXISTS gym_id,
    DROP COLUMN IF EXISTS phone,
    DROP COLUMN IF EXISTS type,
    DROP COLUMN IF EXISTS membership_start,
    DROP COLUMN IF EXISTS membership_end;

DROP TYPE IF EXISTS member_type;
