import { theme as baseTheme } from "@chakra-ui/theme";
import { type ThemeConfig } from "@chakra-ui/theme";

const config: ThemeConfig = {
  initialColorMode: "light",
  useSystemColorMode: false,
};

export const theme = {
  ...baseTheme,
  config,
};
