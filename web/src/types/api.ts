export type Tab = "work" | "personal" | "promotional" | "transactional" | "";

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
  threadSubject: string;
  createdAt: string;
  updatedAt: string;
  mostRecentEmailTimestamp: string;
  recipients: string[];
  summary: string;
  actionItems: string;
  urgencyScore: number;
  category: string;
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
  category: string;
}

export interface ApiError {
  message: string;
  code: string;
}
