-- PostgreSQL initialization (runs once on first container start)

CREATE TABLE IF NOT EXISTS account (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS events (
  id SERIAL PRIMARY KEY,
  event TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_events (
  account_id INT NOT NULL,
  event_id INT NOT NULL,
  PRIMARY KEY (account_id, event_id),
  CONSTRAINT fk_user_events_account FOREIGN KEY (account_id) REFERENCES account(id),
  CONSTRAINT fk_user_events_events FOREIGN KEY (event_id) REFERENCES events(id)
);

INSERT INTO account (email) VALUES
  ('foo@example.com'),
  ('bar@example.com'),
  ('baz@example.com')
ON CONFLICT (email) DO NOTHING;

INSERT INTO events (event) VALUES
  ('User signed up'),
  ('User logged in'),
  ('User updated profile')
ON CONFLICT (event) DO NOTHING;

-- Use upsert-safe mapping with CTEs
WITH pairs AS (
  SELECT * FROM (VALUES
    ('foo@example.com','User signed up'),
    ('foo@example.com','User logged in'),
    ('bar@example.com','User logged in'),
    ('baz@example.com','User updated profile'),
    ('bar@example.com','User signed up')
  ) AS t(email,event)
)
INSERT INTO user_events (account_id, event_id)
SELECT a.id, e.id
FROM pairs p
JOIN account a ON a.email = p.email
JOIN events e  ON e.event = p.event
ON CONFLICT (account_id, event_id) DO NOTHING;