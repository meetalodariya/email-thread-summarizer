-- Thread Summaries

CREATE TABLE public.thread_summaries (
  id bigserial NOT NULL,
  gmail_thread_id text NOT NULL,
  processed_email_ids text[] NULL,
  summary text NULL,
  thread_subject text NULL,
  user_id bigint NOT NULL,
  most_recent_email_timestamp timestamp with time zone NULL,
  created_at timestamp with time zone NULL,
  updated_at timestamp with time zone NULL,
  deleted_at timestamp with time zone NULL,
  search_vector tsvector NULL
);

ALTER TABLE public.thread_summaries ADD CONSTRAINT thread_summaries_pkey PRIMARY KEY (id);

ALTER TABLE thread_summaries ADD COLUMN search_vector tsvector;
UPDATE thread_summaries SET search_vector = to_tsvector('english', thread_summaries.thread_subject || ' ' || thread_summaries.summary);
CREATE INDEX idx_summary_search ON thread_summaries USING GIN(search_vector);

CREATE OR REPLACE FUNCTION search_vector_update() RETURNS trigger AS $$
BEGIN
  NEW.search_vector := to_tsvector('english', coalesce(NEW.summary, '') || ' ' || coalesce(NEW.thread_subject, ''));
  RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER tsvectorupdate
BEFORE INSERT OR UPDATE ON thread_summaries
FOR EACH ROW EXECUTE FUNCTION search_vector_update();

-- User

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

