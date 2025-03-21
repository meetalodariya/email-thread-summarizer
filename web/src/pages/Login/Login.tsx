import {
  Box,
  Button,
  Container,
  Flex,
  Heading,
  Text,
  VStack,

} from "@chakra-ui/react";
import { FcGoogle } from "react-icons/fc";

export const Login = () => {
  const isDark = colorMode === "dark";

  const handleGoogleLogin = async () => {
    // TODO: Implement Google login
    console.log("Google login clicked");
  };

  return (
    <Flex minH="100vh" align="center" justify="center" bg={isDark ? "gray.800" : "gray.50"}>
      <Container maxW="lg" py={{ base: "12", md: "24" }} px={{ base: "0", sm: "8" }}>
        <Box
          bg={isDark ? "gray.700" : "white"}
          py="8"
          px={{ base: "4", md: "10" }}
          shadow="base"
          rounded={{ sm: "lg" }}
        >
          <VStack spacing={6}>
            <VStack spacing={{ base: 2, md: 3 }} textAlign="center">
              <Heading size={{ base: "xs", md: "sm" }}>ThreadSage</Heading>
              <Text fontSize={{ base: "sm", md: "md" }} color={isDark ? "gray.400" : "gray.600"}>
                Simplify your email threads with AI
              </Text>
            </VStack>

            <Button
              w="full"
              size="lg"
              variant="outline"
              onClick={handleGoogleLogin}
              leftIcon={<FcGoogle />}
              _hover={{
                bg: isDark ? "gray.600" : "gray.50",
              }}
            > <a href="/auth/register/google"></a>
              Continue with Google
            </Button>

            <Text fontSize="xs" color={isDark ? "gray.400" : "gray.600"} textAlign="center">
              By continuing, you agree to our Terms of Service and Privacy Policy
            </Text>
          </VStack>
        </Box>
      </Container>
    </Flex>
  );
};
