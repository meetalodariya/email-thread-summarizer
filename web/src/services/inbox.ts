import { ApiResponse, GetInboxParams, InboxResponse } from "@/types/api";
import { api } from "./api";

export const inboxService = {
  getInbox: async (
    params: GetInboxParams,
    token?: string
  ): Promise<InboxResponse> => {
    const { data: response } = await api.get<InboxResponse>("/api/inbox", {
      params: {
        nextCursor: params.nextCursor || "",
        q: params.q || "",
      },
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    return response;
  },
};
