export type Tab = "action" | "important" | "junk";

export interface RegisterRequest {
  code: string;
}

export interface RegisterResponse {
  token: string;
  name: string;
}

export interface ThreadSummary {
  id: string;
  gmailThreadID: string;
  summary: string;
  threadSubject: string;
  createdAt: string;
  updatedAt: string;
  mostRecentEmailTimestamp: string;
  recipients: string[];
}

export type ApiResponse<T> = { data: T | null; status?: number };

export type InboxResponse = {
  data: Array<ThreadSummary>;
  pagination: {
    nextCursor: string;
  };
};

export interface GetInboxParams {
  q: string;
  nextCursor: string;
}

export interface ApiError {
  message: string;
  code: string;
}
