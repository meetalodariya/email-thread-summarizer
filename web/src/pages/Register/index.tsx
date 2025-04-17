import { useState, useEffect, useRef, useCallback } from "react";
import {
  Box,
  Button,
  Flex,
  Heading,
  Text,
  VStack,
  Card,
} from "@chakra-ui/react";
import { FcGoogle } from "react-icons/fc";
import { useNavigate, useSearchParams } from "react-router";
import { toaster } from "@/components/ui/toaster";
import { authService } from "@/services/auth";
import { useAuth } from "@/providers/auth";

export const Register: React.FC<{
  isCallback?: boolean;
}> = ({ isCallback = false }) => {
  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate();
  const isFirstLoad = useRef(true);
  const [isLoading, setIsLoading] = useState(isCallback);
  const { signin, user } = useAuth();

  const handleRegister = useCallback(() => {
    setIsLoading(true);

    const code = searchParams.get("code");

    if (!code) {
      navigate("/auth");

      return;
    }

    const params = new URLSearchParams();
    params.delete("code");
    setSearchParams(
      (prevParams) => {
        prevParams.delete("code");
        return prevParams;
      },
      {
        preventScrollReset: true,
      }
    );

    authService
      .authenticate({ code })
      .then((response) => {
        const { token, name } = response;

        signin({ name, token }, () => {
          navigate("/inbox");
        });
      })
      .catch((err) => {
        console.log(err);

        toaster.create({
          title: "Registration failed. Try again.",
          type: "error",
          onStatusChange() {
            navigate("/auth");
          },
        });
      });
  }, [setIsLoading, searchParams, navigate, setSearchParams, signin]);

  useEffect(() => {
    if (user) {
      navigate("/inbox");
      return;
    }

    if (!isCallback || !isFirstLoad.current) {
      return;
    }

    handleRegister();

    return () => {
      isFirstLoad.current = false;
    };
  }, [isCallback, handleRegister, navigate, user]);

  return (
    <Flex minH="100vh" align="center" justify="center" minW="100vw">
      <Card.Root maxW="sm" overflow="hidden" shadow={"2xl"}>
        <Card.Body gap="2">
          <Box
            py="8"
            px={{ base: "4", md: "10" }}
            shadow="base"
            rounded={{ sm: "lg" }}
          >
            <VStack>
              <VStack textAlign="center">
                <Heading
                  size={{ base: "xl" }}
                  style={{
                    fontFamily: '"Winky Sans", sans-serif',
                    fontSize: "2rem",
                  }}
                >
                  ThreadSage
                </Heading>
                <Text fontSize={{ base: "sm" }} marginBottom={"18px"}>
                  Simplify your email threads with AI
                </Text>
              </VStack>
              <a
                href="http://localhost:8080/api/auth/google"
                style={{ textDecoration: "none", color: "inherit" }}
              >
                <Button
                  disabled={isLoading}
                  loading={isLoading}
                  loadingText="Signing in..."
                  w="full"
                  size="md"
                  variant="outline"
                  shadow={"md"}
                >
                  <FcGoogle />
                  Continue with Google
                </Button>
              </a>
            </VStack>
          </Box>
        </Card.Body>
        <Card.Footer gap="2">
          <Text fontSize="xs" textAlign="center">
            By continuing, you agree to our Terms of Service and Privacy Policy
          </Text>
        </Card.Footer>
      </Card.Root>
    </Flex>
  );
};
