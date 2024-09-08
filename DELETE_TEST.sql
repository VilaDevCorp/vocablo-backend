CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO languages
(id, creation_date, code)
VALUES(uuid_generate_v4(), current_timestamp, 'en');

INSERT INTO users
(id, creation_date, username, email, "password", validated)
VALUES(uuid_generate_v4(), current_timestamp, 'johntest', 'johntest@gmail.com', '$2a$14$8nS81fSGzYIMiL3o3u4CUuy7knPG6dPElPALPOVYQysTbADsXeaZS', true);