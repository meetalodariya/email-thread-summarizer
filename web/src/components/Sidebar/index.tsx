import { Box, Button, Flex } from "@chakra-ui/react";
import { IconType } from "react-icons/lib";
import { IoMdAlert } from "react-icons/io";
import { LuListTodo } from "react-icons/lu";
import { FC } from "react";
import { Tab } from "@/types/api";
import { GiThunderSkull } from "react-icons/gi";

interface LinkItemProps {
  name: string;
  icon: IconType;
  code: Tab;
}

const LinkItems: Array<LinkItemProps> = [
  { name: "Action Required", icon: LuListTodo, code: "action" },
  { name: "Important", icon: IoMdAlert, code: "important" },
  { name: "Malicious / Junk", icon: GiThunderSkull, code: "junk" },
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
            <Flex align={"center"} gap={"1"} mb={"2"} ml={"2"}>
              <Icon /> {name}
            </Flex>
          </Button>
        </Box>
      ))}
    </>
  );
};
