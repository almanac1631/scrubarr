-- +goose Up
delete from entry_mappings;
alter table entry_mappings add column file_node varchar not null;

-- +goose Down
alter table entry_mappings drop column file_node;
