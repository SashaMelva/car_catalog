-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.car_catalog
(
    reg_num character varying COLLATE pg_catalog."default" NOT NULL,
    mark character varying COLLATE pg_catalog."default" NOT NULL,
    model character varying COLLATE pg_catalog."default" NOT NULL,
    year integer NOT NULL,
    owner json,
    CONSTRAINT car_catalog_pkey PRIMARY KEY (reg_num)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS car_catalog;
-- +goose StatementEnd
