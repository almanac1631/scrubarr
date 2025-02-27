-- +goose Up
alter table entry_mappings
    add column date_added date;

-- +goose Down
alter table entry_mappings
    drop column date_added;
