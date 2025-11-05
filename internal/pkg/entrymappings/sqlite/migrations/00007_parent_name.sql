-- +goose Up
delete from entry_mappings;
alter table entry_mappings add column parent_name varchar not null;

-- +goose Down
alter table entry_mappings drop column parent_name;