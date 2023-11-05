import { VStack } from "@chakra-ui/react";
import { LayoutEightPadding } from "../config";
import { Footer } from "../organisms/Footer";
import { Header } from "../organisms/Header";

interface LayoutProps {
  children: React.ReactNode;
}

export function Layout(props: LayoutProps) {
  return (
    <VStack minH="100vh" align="flex-start" justify="flex-start" bg="inherit">
      <Header />
      <VStack spacing={LayoutEightPadding} pl={LayoutEightPadding} pr={LayoutEightPadding} pb={LayoutEightPadding} width="100%" bg="inherit">
        {props.children}
      </VStack>
      <Footer />
    </VStack>
  )
}