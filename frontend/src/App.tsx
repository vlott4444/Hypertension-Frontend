import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { BottomNav } from "./components/BottomNav";
import { BusyOverlay } from "./components/BusyOverlay";
import { FeedPage } from "./pages/FeedPage";
import { GridPage } from "./pages/GridPage";
import { DetailPage } from "./pages/DetailPage";
import { AddPage } from "./pages/AddPage";
import { LoginPage } from "./pages/LoginPage";

export default function App() {
  return <BrowserRouter><div className="app-shell">
    <Routes>
      <Route path="/" element={<Navigate to="/feed" replace/>}/>
      <Route path="/feed" element={<FeedPage/>}/>
      <Route path="/stages" element={<GridPage/>}/>
      <Route path="/stages/:id" element={<DetailPage/>}/>
      <Route path="/add" element={<AddPage/>}/>
      <Route path="/login" element={<LoginPage/>}/>
      <Route path="*" element={<Navigate to="/feed" replace/>}/>
    </Routes>
    <BottomNav/>
    <BusyOverlay/>
  </div></BrowserRouter>;
}
