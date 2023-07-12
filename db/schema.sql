CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  login VARCHAR(64) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  last_auth TIMESTAMP,
  last_unauth TIMESTAMP,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  is_auth BOOLEAN NOT NULL DEFAULT(false)
);

CREATE TABLE test_variant (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE task (
  id BIGSERIAL PRIMARY KEY,
  task VARCHAR(256) NOT NULL,
  answer TEXT NOT NULL,
  options jsonb NOT NULL,
  variant_id BIGINT NOT NULL REFERENCES test_variant(id)
);

CREATE TABLE user_test (
  id BIGSERIAL PRIMARY KEY,
  start_at TIMESTAMP NOT NULL,
  user_id BIGINT NOT NULL REFERENCES users(id),
  variant_id BIGINT NOT NULL REFERENCES test_variant(id)
);

CREATE TABLE test_answer (
  id BIGSERIAL PRIMARY KEY,
  answer TEXT NOT NULL,
  test_id BIGINT NOT NULL REFERENCES user_test(id)
);

CREATE TABLE test_result (
  id BIGSERIAL PRIMARY KEY,
  percent SMALLINT NOT NULL DEFAULT(0),
  test_id BIGINT NOT NULL REFERENCES user_test(id),
  CHECK(percent BETWEEN 0 AND 100)
);

-- TEST DATA

INSERT INTO users 
(login, password, created_at, updated_at) 
VALUES ('admin', 'admin', NOW(), NOW());

-- CREATING 2 VARIANTS 

INSERT INTO test_variant
(name)
VALUES ('Variant 1');

INSERT INTO test_variant
(name)
VALUES ('Variant 2');

-- ADDING 3 TASKS TO VARIANT 1

INSERT INTO task
(task, answer, variant_id, options)
VALUES ('What is 2 * 2?', '4', 1, '["4", "2", "1", "6"]');

INSERT INTO task
(task, answer, variant_id, options)
VALUES ('What is 2 - 2?', '0', 1, '["3", "2", "0", "1"]');

INSERT INTO task
(task, answer, variant_id, options)
VALUES (
  '0 * 10000000 ?',
  '0',
  1,
  '["...", "0", "10000000", "I dont know"]'
);

-- ADDING 4 TASKS TO VARIANT 2

INSERT INTO task
(task, answer, variant_id, options)
VALUES ('What is 2 * 2?', '4', 2, '["4", "2", "1", "6"]');

INSERT INTO task
(task, answer, variant_id, options)
VALUES ('What is 2 + 2?', '4', 2, '["3", "2", "4", "1"]');

INSERT INTO task
(task, answer, variant_id, options)
VALUES ('What is 2 % 2?', '0', 2, '["60", "2", "0", "1"]');

INSERT INTO task
(task, answer, variant_id, options)
VALUES ('What is 2 / 2?', '1', 2, '["3", "2", "0", "1"]');