-- +goose Up
alter table entry_mappings
    add column size integer;

-- +goose Down
alter table entry_mappings
    drop column size;
