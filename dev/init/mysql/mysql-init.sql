-- MySQL initialization (runs once on first container start)

USE pipeline_forge_dev;

CREATE TABLE IF NOT EXISTS account (
  id INT AUTO_INCREMENT PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS events (
  id INT AUTO_INCREMENT PRIMARY KEY,
  event TEXT NOT NULL,
  CONSTRAINT uq_event UNIQUE (event(255))
);

CREATE TABLE IF NOT EXISTS user_events (
  account_id INT NOT NULL,
  event_id INT NOT NULL,
  PRIMARY KEY (account_id, event_id),
  CONSTRAINT fk_user_events_account FOREIGN KEY (account_id) REFERENCES account(id),
  CONSTRAINT fk_user_events_events FOREIGN KEY (event_id) REFERENCES events(id)
);

INSERT IGNORE INTO account (email) VALUES
  ('foo@example.com'),
  ('bar@example.com'),
  ('baz@example.com');

INSERT IGNORE INTO events (event) VALUES
  ('User signed up'),
  ('User logged in'),
  ('User updated profile');

-- Map helper ids using subqueries so it works on fresh and reruns.
INSERT IGNORE INTO user_events (account_id, event_id)
SELECT a.id, e.id
FROM (SELECT 'foo@example.com' AS email, 'User signed up' AS event UNION ALL
      SELECT 'foo@example.com', 'User logged in' UNION ALL
      SELECT 'bar@example.com', 'User logged in' UNION ALL
      SELECT 'baz@example.com', 'User updated profile' UNION ALL
      SELECT 'bar@example.com', 'User signed up') x
JOIN account a ON a.email = x.email
JOIN events e ON e.event = x.event;