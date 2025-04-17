import { useState } from "react";
import { Navbar } from "@/components/Navbar";
import { useGetInbox } from "@/hooks/useInbox";
import { GetInboxParams, Tab } from "@/types/api";
import { useSearchParams } from "react-router";
import { Box, Flex } from "@chakra-ui/react";
import { SidebarContent } from "@/components/Sidebar";
import { ListContainer } from "./ListContainer";

const initialFilters = {
  q: "",
  nextCursor: "",
};

export const Inbox: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [filters, setFilters] = useState<GetInboxParams>({
    ...initialFilters,
    q: searchParams.get("q") || "",
  });
  const [tab, setTab] = useState<Tab>("action");

  const {
    data,
    error,
    fetchNextPage,
    hasNextPage,
    isFetching,
    isFetchingNextPage,
    // status,
    // refetch,
  } = useGetInbox(filters);

  return (
    <Box h="100vh" w="100vw">
      <Flex align="center" pos="sticky" top="0" left="0" right="0">
        <Navbar
          searchQuery={filters.q}
          onSearchChange={(q) => {
            setFilters({
              q,
              nextCursor: "",
            });
            setSearchParams(
              (prevParams) => {
                if (!q) {
                  prevParams.delete("q");
                } else {
                  prevParams.set("q", q);
                }
                return prevParams;
              },
              {
                preventScrollReset: true,
              }
            );
          }}
        />
      </Flex>
      <Flex>
        <Box
          transition="3s ease"
          borderRight="2px"
          w={{ base: "full", md: 60 }}
          // pos="fixed"

          pt={"2"}
          overflowY={"hidden"}
          flexShrink={"1"}
          shadow={"md"}
          minH={"calc(100vh - 64px)"}
        >
          <SidebarContent
            tab={tab}
            setTab={(tab) => {
              setTab(tab);
            }}
          />
        </Box>
        <Box
          p={1}
          w="full"
          overflowY={"auto"}
          flexGrow={"1"}
          maxH={"calc(100vh - 64px)"}
          id="scrollableDiv"
        >
          <ListContainer
            fetchNextPage={fetchNextPage}
            hasNextPage={hasNextPage}
            data={data}
            isFetching={isFetching}
            isFetchingNextPage={isFetchingNextPage}
          />
        </Box>
      </Flex>
    </Box>
  );
};
