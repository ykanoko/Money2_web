import { MantineProvider } from "@mantine/core";
import { Routes, Route, BrowserRouter } from "react-router-dom";
import AddMoneyRecordPage from "./pages/AddMoneyRecordPage";
import RegisterPage from "./pages/RegisterPage";
import LoginPage from "./pages/LoginPage";
import MoneyRecordPage from "./pages/MoneyRecordsPage";

export default function App() {
  return (
    <MantineProvider withGlobalStyles withNormalizeCSS>
      <BrowserRouter>
        <Routes>
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route index element={<AddMoneyRecordPage />} />
          <Route path="/money_record" element={<MoneyRecordPage />} />
        </Routes>
      </BrowserRouter>
    </MantineProvider>
  );
}
