import {
  Box,
  Flex,
  Avatar,
  HStack,
  Menu,
  Portal,
  Heading,
  InputGroup,
  Input,
} from "@chakra-ui/react";
import { HiMagnifyingGlass } from "react-icons/hi2";
import { useAuth } from "@/providers/auth";
import { useNavigate } from "react-router";
import React, { useState } from "react";

interface NavbarProps {
  onSearchChange: (q: string) => void;
  searchQuery: string;
}

export const Navbar: React.FC<NavbarProps> = ({
  onSearchChange,
  searchQuery,
}) => {
  const [search, setSearch] = useState("");
  const { signout, user } = useAuth();
  const navigate = useNavigate();

  return (
    <Box width="100%" px={4} backgroundColor={"whitesmoke"} shadow={"md"}>
      <Flex h={16} alignItems={"center"} justifyContent={"space-between"}>
        <HStack alignItems={"center"}>
          <Heading
            style={{
              fontFamily: '"Winky Sans", sans-serif',
              fontSize: "1.8rem",
            }}
          >
            ThreadSage
          </Heading>
        </HStack>
        <HStack as={"nav"} w="30vw">
          <InputGroup
            flex="1"
            startElement={<HiMagnifyingGlass size={"1rem"} />}
            backgroundColor={"white"}
          >
            <Input
              placeholder="Search email summaries"
              onChange={(e) => {
                setSearch(e.target.value);
              }}
              onKeyDown={(event) => {
                if (event.key === "Enter") {
                  onSearchChange(search);
                }
              }}
              defaultValue={searchQuery}
              value={search}
            />
          </InputGroup>
        </HStack>
        <Flex alignItems={"center"}>
          <Menu.Root>
            <Menu.Trigger border={"none"} padding={0} rounded={"4xl"}>
              <Avatar.Root size="sm">
                <Avatar.Fallback name={user?.name} />
              </Avatar.Root>
            </Menu.Trigger>
            <Portal>
              <Menu.Positioner>
                <Menu.Content>
                  <Menu.ItemGroup>
                    <Menu.Item
                      value="signout"
                      onClick={() => {
                        signout(() => {
                          navigate("/auth");
                        });
                      }}
                    >
                      Sign Out
                    </Menu.Item>
                    <Menu.Item value="signed_in" disabled cursor={"auto"}>
                      Signed in as: {user?.name}
                    </Menu.Item>
                  </Menu.ItemGroup>
                </Menu.Content>
              </Menu.Positioner>
            </Portal>
          </Menu.Root>
        </Flex>
      </Flex>
    </Box>
  );
};
