import { NavLink } from "react-router-dom";
export function BottomNav() {
  return <nav className="bottom-nav" aria-label="Основная навигация">
    <NavLink to="/feed" className={({isActive}) => isActive ? "active" : ""}>лента</NavLink>
    <span className="disabled-tab" aria-disabled="true">добавить</span>
    <NavLink to="/stages" className={({isActive}) => isActive ? "active" : ""}>плитка</NavLink>
    <span className="disabled-tab" aria-disabled="true">вход</span>
  </nav>;
}
