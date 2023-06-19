import { MantineProvider} from '@mantine/core';
import { Routes, Route, BrowserRouter } from "react-router-dom";
import HomePage from './pages/HomePage';
import MoneyRecordsPage from './pages/MoneyRecordsPage';


export default function App() {
  return (
    <MantineProvider withGlobalStyles withNormalizeCSS>
      <BrowserRouter>
        <div className="MerComponent">
          <Routes>
            <Route index element={<HomePage />} />
            <Route path="/money_records" element={<MoneyRecordsPage />} />
          </Routes>
        </div>
      </BrowserRouter>

    </MantineProvider>
  );
}