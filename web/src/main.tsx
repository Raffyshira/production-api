import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router";
import App from "./App";
import { ConfirmationPage } from "./ConfirmationPage";
import "./index.css";

const root = document.getElementById("root");

ReactDOM.createRoot(root).render(
  <BrowserRouter>
    <Routes>
      <Route path="/" element={<App />} />
      <Route path="/confirm/:token" element={<ConfirmationPage />} />
    </Routes>
  </BrowserRouter>,
);
