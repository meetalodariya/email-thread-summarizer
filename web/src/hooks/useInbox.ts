import { AuthContext, useAuth } from "@/providers/auth";
import { UnauthorizedError } from "@/services/api";
import { inboxService } from "@/services/inbox";
import { ApiResponse, GetInboxParams, InboxResponse } from "@/types/api";
import { useInfiniteQuery } from "@tanstack/react-query";
import { useNavigate } from "react-router";

export const useGetInbox = (params: GetInboxParams) => {
  const { user, signout } = useAuth();
  const navigate = useNavigate();

  const {
    data,
    error,
    fetchNextPage,
    hasNextPage,
    isFetching,
    isFetchingNextPage,
    isError,
  } = useInfiniteQuery({
    queryKey: ["inbox", params.q],
    queryFn: async ({ pageParam }): Promise<InboxResponse> => {
      const response = await inboxService.getInbox(
        {
          ...params,
          nextCursor: pageParam,
        },
        user?.token
      );

      return response;
    },
    getNextPageParam: (lastPage) => lastPage.pagination.nextCursor || null,
    initialPageParam: "",
  });

  if (isError) {
    if (error instanceof UnauthorizedError) {
      signout(() => {
        navigate("/auth");
      });
    }
  }

  return {
    data,
    error,
    fetchNextPage,
    hasNextPage,
    isFetching,
    isFetchingNextPage,
    // status,
    // refetch,
  };
};
