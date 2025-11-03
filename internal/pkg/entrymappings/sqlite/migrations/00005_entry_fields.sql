-- +goose Up
delete from entry_mappings;

alter table entry_mappings add column file_path varchar not null;
alter table entry_mappings add column parent_id varchar not null;

-- +goose Down
alter table entry_mappings drop column file_path;
alter table entry_mappings drop column parent_id;
