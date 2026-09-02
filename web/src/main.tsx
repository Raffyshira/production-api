import LoginPage from "@/LoginPage";
import SignupPage from "@/SignupPage";
import ExplorePage from "@/ExplorePage";
import { CookiesProvider } from "react-cookie";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router";
import App from "./App";
import { ThemeProvider } from "./components/Theme-Provider";
import { ConfirmationPage } from "./ConfirmationPage";
import "./index.css";
import { SinglePost } from "./SinglePost";

const root = document.getElementById("root");

ReactDOM.createRoot(root).render(
  <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
    <CookiesProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<App />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/confirm/:token" element={<ConfirmationPage />} />
          <Route path="/post/:postID" element={<SinglePost />} />
          <Route path="/signup" element={<SignupPage />} />
          <Route path="/explore" element={<ExplorePage />} />
        </Routes>
      </BrowserRouter>
    </CookiesProvider>
  </ThemeProvider>,
);
