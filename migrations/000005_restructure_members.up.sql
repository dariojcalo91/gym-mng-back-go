DO $$ BEGIN
    CREATE TYPE member_type AS ENUM ('monthly', 'occasional');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

ALTER TABLE members
    ADD COLUMN gym_id UUID REFERENCES gyms(id),
    ADD COLUMN phone VARCHAR(30),
    ADD COLUMN type member_type NOT NULL DEFAULT 'occasional',
    ADD COLUMN membership_start TIMESTAMP WITH TIME ZONE,
    ADD COLUMN membership_end TIMESTAMP WITH TIME ZONE;

DROP INDEX IF EXISTS idx_members_email;
ALTER TABLE members DROP COLUMN IF EXISTS email;
ALTER TABLE members DROP COLUMN IF EXISTS plan;
ALTER TABLE members DROP COLUMN IF EXISTS status;

DROP TYPE IF EXISTS member_status;
