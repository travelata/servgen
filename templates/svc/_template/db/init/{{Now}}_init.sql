-- +goose Up
CREATE ROLE {{.service.name}} LOGIN PASSWORD '{{.service.name}}' NOINHERIT CREATEDB;
CREATE SCHEMA {{.service.name}} AUTHORIZATION {{.service.name}};
GRANT USAGE ON SCHEMA {{.service.name}} TO PUBLIC;

-- +goose Down
DROP SCHEMA {{.service.name}};
DROP ROLE {{.service.name}};
