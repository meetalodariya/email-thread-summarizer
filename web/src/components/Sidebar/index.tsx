import { Box, Button, Flex } from "@chakra-ui/react";
import { IconType } from "react-icons/lib";

import { LuListTodo } from "react-icons/lu";
import { PiBagSimpleFill } from "react-icons/pi";
import { IoChatbubbles } from "react-icons/io5";
import { GrTransaction } from "react-icons/gr";
import { MdLocalOffer } from "react-icons/md";
import { FC } from "react";
import { Tab } from "@/types/api";

interface LinkItemProps {
  name: string;
  icon: IconType;
  code: Tab;
}

const LinkItems: Array<LinkItemProps> = [
  { name: "All", icon: LuListTodo, code: "" },
  { name: "Work", icon: PiBagSimpleFill, code: "work" },
  { name: "Personal", icon: IoChatbubbles, code: "personal" },
  { name: "Transactional", icon: GrTransaction, code: "transactional" },
  { name: "Promotional", icon: MdLocalOffer, code: "promotional" },
];

interface SidebarContentProps {
  tab: Tab;
  setTab: (tab: Tab) => void;
}

export const SidebarContent: FC<SidebarContentProps> = ({ tab, setTab }) => {
  return (
    <>
      {LinkItems.map(({ name, icon: Icon, code }) => (
        <Box key={name}>
          <Button
            variant={tab === code ? "subtle" : "ghost"}
            // colorPalette={""}
            w={"full"}
            textAlign={"start"}
            onClick={() => {
              setTab(code);
            }}
          >
            <Flex align={"center"} gap={"2"} mb={"2"} ml={"2"}>
              <Icon /> {name}
            </Flex>
          </Button>
        </Box>
      ))}
    </>
  );
};
