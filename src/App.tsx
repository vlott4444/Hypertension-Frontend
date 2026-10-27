import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { BottomNav } from "./components/BottomNav";
import { FeedPage } from "./pages/FeedPage";
import { GridPage } from "./pages/GridPage";
import { DetailPage } from "./pages/DetailPage";
export default function App(){return <BrowserRouter><div className="app-shell"><Routes><Route path="/" element={<Navigate to="/feed" replace/>}/><Route path="/feed" element={<FeedPage/>}/><Route path="/stages" element={<GridPage/>}/><Route path="/stages/:id" element={<DetailPage/>}/><Route path="*" element={<Navigate to="/feed" replace/>}/></Routes><BottomNav/></div></BrowserRouter>}
