-- +goose Up
-- Pro Outfitters' "Plan Your Trip" inquiry form captures more than the
-- generic contact form: what the guest is interested in, when they'd like
-- to come, and how many are in the party. Stored as TEXT NOT NULL DEFAULT
-- '' so existing read paths scan cleanly into plain strings.

ALTER TABLE inquiries ADD COLUMN interest        TEXT NOT NULL DEFAULT '';
ALTER TABLE inquiries ADD COLUMN preferred_dates TEXT NOT NULL DEFAULT '';
ALTER TABLE inquiries ADD COLUMN group_size      TEXT NOT NULL DEFAULT '';

-- +goose Down

ALTER TABLE inquiries DROP COLUMN group_size;
ALTER TABLE inquiries DROP COLUMN preferred_dates;
ALTER TABLE inquiries DROP COLUMN interest;
