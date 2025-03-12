-- +goose Up
alter table entry_mappings
    add column id text not null default '<unset>';

-- +goose Down
alter table entry_mappings
    drop column id;
