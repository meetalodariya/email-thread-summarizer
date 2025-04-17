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


insert into "public"."thread_summaries" ("created_at", "deleted_at", "gmail_thread_id", "id", "most_recent_email_timestamp", "processed_email_ids", "search_vector", "summary", "thread_subject", "updated_at", "user_id") values ('2025-02-10 04:33:30+00', '0001-01-01 00:00:00+00', '194ee2249930a7db', '16', '2025-02-10 04:33:30+00', '{"194ee2249930a7db"}', '''/notifications.'':48 ''account'':14,30,33,55,74 ''action'':49 ''activ'':37,44,70 ''ad'':11 ''add'':22,60 ''alert'':2,5 ''check'':35,51 ''ensur'':72 ''googl'':7,16 ''item'':50 ''link'':41,65 ''may'':26 ''myaccount.google.com'':47 ''myaccount.google.com/notifications.'':46 ''new'':9 ''no-reply@accounts.google.com'':17 ''passkey'':10,24,62 ''provid'':40,66 ''regard'':8 ''review'':68 ''secur'':1,4,31,43,53,76 ''sender'':15 ''someon'':25 ''subject'':3 ''use'':28,63 ''view'':42', '
Subject: Security alert from Google regarding new passkey added to your account
Sender: Google <no-reply@accounts.google.com>
If you did not add the passkey, someone may be using your account. Secure your account by checking the activity at the provided link. View security activity at https://myaccount.google.com/notifications.
Action Items:
- Check and secure your account if you did not add the passkey. Use the link provided to review the activity and ensure your account''s security.', 'Security alert', '2025-03-22 01:22:38.562191+00', '5');
insert into "public"."thread_summaries" ("created_at", "deleted_at", "gmail_thread_id", "id", "most_recent_email_timestamp", "processed_email_ids", "search_vector", "summary", "thread_subject", "updated_at", "user_id") values ('2025-03-18 02:56:24+00', '0001-01-01 00:00:00+00', '195a72e5f8f05cee', '14', '2025-03-18 02:56:24+00', '{"195a72e5f8f05cee"}', '''account'':18,35,53 ''action'':36 ''activ'':27,40 ''advis'':23 ''alert'':2,8 ''associ'':41 ''chang'':20,46 ''check'':25,38 ''email'':7 ''googl'':4,34,52 ''gsummarizer@gmail.com'':11 ''item'':37 ''notif'':31,49 ''notifi'':10 ''phone'':15,45 ''recoveri'':14,44 ''review'':29,47 ''secur'':1,30,48 ''sent'':5 ''summari'':3', 'Summary:
Google sent an email alert to notify gsummarizer@gmail.com that the recovery phone for their account was changed. They are advised to check the activity and review security notifications on their Google Account.

Action Items:
- Check the activity associated with the recovery phone change.
- Review security notifications on the Google Account.', 'Security alert', '2025-03-22 01:22:34.300165+00', '5');
insert into "public"."thread_summaries" ("created_at", "deleted_at", "gmail_thread_id", "id", "most_recent_email_timestamp", "processed_email_ids", "search_vector", "summary", "thread_subject", "updated_at", "user_id") values ('2025-03-06 20:24:52+00', '0001-01-01 00:00:00+00', '1956d22136e08eab', '15', '2025-03-06 20:24:52+00', '{"1956d3b1ecd6d0a1","1956d22136e08eab"}', '''1'':80 ''2'':98 ''action'':78 ''alodariya'':84 ''ask'':87 ''brief'':12 ''certain'':20 ''clarif'':89 ''contribut'':58 ''detail'':25 ''eager'':56 ''email'':2,6,9,67,86 ''ensur'':60 ''express'':55 ''follow'':36 ''follow-up'':35 ''form'':49 ''guidanc'':43 ''howev'':64 ''inform'':77,106 ''inquir'':14 ''involv'':33 ''item'':79 ''limit'':76 ''matter'':21 ''meet'':50,83 ''messag'':13 ''necessari'':48,100 ''new'':5,8,66 ''next'':46 ''overal'':26 ''prepar'':52 ''process'':41,109 ''provid'':101 ''queri'':71 ''question'':111 ''refer'':96 ''relev'':103 ''respond'':81 ''seek'':42 ''sender'':54 ''simpl'':70 ''smooth'':62 ''specif'':40,92 ''specifi'':23 ''step'':47 ''summari'':3,27 ''test'':1 ''thread'':30,32 ''transit'':63 ''updat'':17,74,93,104 ''without'':22', '**Summary of New Email:**
The new email is a brief message inquiring about the update on a certain matter without specifying further details.

**Overall Summary of the Thread:**
The thread involves a follow-up on a specific process, seeking guidance on the next steps, necessary forms, meetings, and preparations. The sender expresses eagerness to contribute and ensure a smooth transition. However, the new email is a simple query about an update with limited information.

**Action Items:**
1. Respond to Meet Alodariya''s email asking for clarification on the specific update they are referring to.
2. If necessary, provide any relevant updates or information on the process in question.', 'test email', '2025-03-22 01:22:37.254442+00', '5');
