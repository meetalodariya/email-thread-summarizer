import type { InboxResponse, Tab } from "@/types/api";
import { formatDate } from "@/utils/date";
import {
  Avatar,
  Badge,
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
    padding-left: 14px;
  }

  li {
    list-style-type: circle;
  }
`;

const getCategoryBadge = (category: Tab) => {
  switch (category) {
    case "personal":
      return "blue";
    case "work":
      return "green";
    case "promotional":
      return "purple";
    case "transactional":
      return "orange";
  }
};

const getUrgencyColor = (urgencyScore) => {
  switch (urgencyScore) {
    case 1:
      return "#494CA2";
    case 2:
      return "#494CA2";
    case 3:
      return "#66b2b2";
    case 4:
      return "#006666";
    case 5:
      return "#FFA500";
    case 6:
      return "#FFD700";
    case 7:
      return "#FF8C00";
    case 8:
      return "#FF6347";
    case 9:
      return "#FF4500";
    case 10:
      return "#FF0000";
  }
};

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
    actionItems: "",
    urgencyScore: 0,
    category: "",
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
                    actionItems: threadSummary.actionItems,
                    category: threadSummary.category,
                    urgencyScore: threadSummary.urgencyScore,
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
                  style={{
                    backgroundColor:
                      getUrgencyColor(threadSummary.urgencyScore) + "30",
                  }}
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
                  <h2>
                    Urgency Score:{" "}
                    <span
                      style={{
                        color: getUrgencyColor(dialogContent.urgencyScore),
                        textShadow: "black 2px 1px 20px",
                      }}
                    >
                      {dialogContent.urgencyScore}
                    </span>
                  </h2>
                  <h2>
                    Category:{" "}
                    <Badge
                      style={{ textTransform: "capitalize" }}
                      colorPalette={getCategoryBadge(
                        dialogContent.category as Tab
                      )}
                    >
                      {dialogContent.category}
                    </Badge>
                  </h2>

                  <h2>Summary: </h2>
                  <Markdown>{dialogContent.summary}</Markdown>
                  <h2>Action Items: </h2>
                  <Markdown>{dialogContent.actionItems}</Markdown>
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
