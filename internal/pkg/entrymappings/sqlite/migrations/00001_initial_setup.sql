-- +goose Up
create table entry_mappings
(
    retriever_id        text not null,
    name                text not null,
    api_resp            text not null,
    metadata_updated_at date default (datetime('now', 'utc')),
    primary key (retriever_id, name)
);

create table retrievers
(
    retriever_id        text not null,
    category            text not null,
    software_name       text not null,
    name                text not null,
    metadata_updated_at date default (datetime('now', 'utc')),
    primary key (retriever_id)
);

-- +goose Down
drop table entry_mappings;
drop table retrievers;
