CREATE TABLE public.users (
  id bigserial NOT NULL,
  created_at timestamp with time zone NULL,
  updated_at timestamp with time zone NULL,
  deleted_at timestamp with time zone NULL,
  first_name text NOT NULL,
  last_name text NOT NULL,
  email text NOT NULL,
  gmail_access_token text NOT NULL,
  gmail_refresh_token text NOT NULL,
  gmail_token_expiry timestamp with time zone NOT NULL,
  is_gmail_token_valid boolean NULL,
  picture text NULL,
  last_scanned_timestamp timestamp with time zone NULL,
  last_processed_mail text NULL
);

ALTER TABLE public.users ADD CONSTRAINT users_pkey PRIMARY KEY (id);

