import { NavLink, useNavigate } from "react-router-dom";
import { TOKEN_STORAGE_KEY, apiError, http } from "../api/http";
import { loggedOut, authStarted } from "../store/authSlice";
import { resetFilters } from "../store/filterSlice";
import { useAppDispatch, useAppSelector } from "../store/store";

export function BottomNav() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const isAuthenticated = useAppSelector(s => s.hypertensionAuth.isAuthenticated);
  const busy = useAppSelector(s => s.hypertensionAuth.busy || s.hypertensionStages.loading || s.hypertensionFeed.loading);
  const logout = async () => {
    dispatch(authStarted());
    let remoteError = "";
    try {
      // Auth is deliberately NOT a thunk or a generated API request.
      await http.post("/auth/logout");
    } catch (error) { remoteError = apiError(error); }
    // Always sign out locally; remote JWT revocation may have failed.
    localStorage.removeItem(TOKEN_STORAGE_KEY);
    dispatch(loggedOut());
    dispatch(resetFilters());
    navigate("/feed");
    if (remoteError) window.alert(`Вы вышли из интерфейса, но сервер не подтвердил отзыв JWT: ${remoteError}`);
  };
  return <nav className="bottom-nav" aria-label="Основная навигация">
    <NavLink to="/feed" className={({ isActive }) => isActive ? "active" : ""}>лента</NavLink>
    <NavLink to="/add" className={({ isActive }) => isActive ? "active" : ""}>добавить</NavLink>
    <NavLink to="/stages" className={({ isActive }) => isActive ? "active" : ""}>плитка</NavLink>
    {isAuthenticated ? <button type="button" className="nav-exit" disabled={busy} onClick={() => { void logout(); }}>выход</button>
      : <NavLink to="/login" className={({ isActive }) => isActive ? "active" : ""}>вход</NavLink>}
  </nav>;
}
