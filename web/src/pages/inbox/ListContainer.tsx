import type { InboxResponse } from "@/types/api";
import { formatDate } from "@/utils/date";
import {
  Avatar,
  Button,
  CloseButton,
  Dialog,
  Flex,
  Portal,
  Spinner,
  Text,
} from "@chakra-ui/react";
import type { InfiniteData } from "@tanstack/react-query";
import { FC, useEffect, useRef, useState } from "react";
import { Fragment } from "react/jsx-runtime";
import Markdown from "react-markdown";
import { css } from "@emotion/react";

interface Props {
  fetchNextPage: () => void;
  hasNextPage: boolean;
  data: InfiniteData<InboxResponse, unknown> | undefined;
  isFetching: boolean;
  isFetchingNextPage: boolean;
}

const markdownCss = css`
  > h2:not(:first-child) {
    margin-top: 10px;
  }

  > h2 {
    margin-bottom: 10px;
    font-weight: bold;
    font-size: 1rem;
  }

  ul {
    padding-left: 12px;
  }

  li {
    list-style-type: circle;
  }
`;

export const ListContainer: FC<Props> = ({
  fetchNextPage,
  hasNextPage,
  data,
  isFetching,
  isFetchingNextPage,
}) => {
  const [dialogContent, setDialogContent] = useState({
    title: "",
    summary: "",
  });
  const observerRef = useRef<HTMLDivElement | null>(null);
  const [isObserverVisible, setIsObserverVisible] = useState(false);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        const entry = entries[0];
        setIsObserverVisible(entry.isIntersecting);
      },
      { threshold: 1.0 }
    );

    if (observerRef.current) {
      observer.observe(observerRef.current);
    }

    return () => observer.disconnect();
  }, []);

  // Fetch next page while observer is visible
  useEffect(() => {
    if (isObserverVisible && hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  }, [isObserverVisible, hasNextPage, isFetchingNextPage, fetchNextPage]);

  return (
    <>
      <Dialog.Root placement={"center"} motionPreset="slide-in-bottom">
        {data?.pages.map((group, i) => (
          <Fragment key={i}>
            {group.data.map((threadSummary) => (
              <Dialog.Trigger
                key={threadSummary.id}
                asChild
                onClick={() => {
                  setDialogContent({
                    title: threadSummary.threadSubject,
                    summary: threadSummary.summary,
                  });
                }}
              >
                <Flex
                  key={threadSummary.id}
                  p={2}
                  w="full"
                  align="center"
                  borderBottom="0.5px solid"
                  borderColor="gray.200"
                  bg={true ? "gray.100" : "white"}
                  cursor="pointer"
                  _hover={{ shadow: "lg" }}
                  mb={"0.5"}
                >
                  <Flex
                    flex="content"
                    w="full"
                    direction="row"
                    justify={"space-around"}
                  >
                    <Text fontWeight={true ? "bold" : "normal"} flex="1">
                      {...(threadSummary.recipients ?? []).map(
                        (recipient, index) => (
                          <span>
                            {recipient}{" "}
                            {index === threadSummary.recipients.length - 1
                              ? ""
                              : ","}
                          </span>
                        )
                      )}
                    </Text>
                    <Text fontWeight="medium" flex="2">
                      {threadSummary.threadSubject}
                    </Text>
                    <Text
                      fontSize="sm"
                      color="gray.500"
                      flex="1"
                      textAlign="right"
                    >
                      {formatDate(threadSummary.mostRecentEmailTimestamp)}
                    </Text>
                    {/* <Text fontSize="sm" color="gray.600">
                  {email.preview}
                </Text> */}
                  </Flex>
                </Flex>
              </Dialog.Trigger>
            ))}
          </Fragment>
        ))}
        <Portal>
          <Dialog.Backdrop />
          <Dialog.Positioner>
            <Dialog.Content>
              <Dialog.Header>
                <Dialog.Title>
                  Thread Summary: {dialogContent.title}
                </Dialog.Title>
              </Dialog.Header>
              <Dialog.Body>
                <div css={markdownCss}>
                  <Markdown>{dialogContent.summary}</Markdown>
                </div>
              </Dialog.Body>
              <Dialog.Footer>
                <Dialog.ActionTrigger asChild>
                  <Button variant="outline">Close</Button>
                </Dialog.ActionTrigger>
              </Dialog.Footer>
              <Dialog.CloseTrigger asChild>
                <CloseButton size="sm" />
              </Dialog.CloseTrigger>
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
        <div ref={observerRef} />
        {/* {!hasNextPage && !isFetching && (
        <span>Looks like you're all caught up!!</span>
      )} */}
        {isFetchingNextPage && <Spinner />}
      </Dialog.Root>
    </>
  );
};
