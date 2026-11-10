import { useState, type FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { TOKEN_STORAGE_KEY, apiError, http } from "../api/http";
import { authFailed, authStarted, loginSucceeded } from "../store/authSlice";
import { useAppDispatch, useAppSelector } from "../store/store";

function readToken(data: unknown): string {
  if (!data || typeof data !== "object") return "";
  const result = data as Record<string, unknown>;
  const nested = result.data && typeof result.data === "object" ? result.data as Record<string, unknown> : {};
  for (const value of [result.token, result.access_token, result.jwt, nested.token, nested.access_token]) {
    if (typeof value === "string" && value.trim()) return value;
  }
  return "";
}
export function LoginPage() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { busy, error, isAuthenticated } = useAppSelector(s => s.hypertensionAuth);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [register, setRegister] = useState(false);
  const [notice, setNotice] = useState("");
  const submit = async (event: FormEvent) => {
    event.preventDefault();
    if (!username.trim() || !password) return;
    dispatch(authStarted());
    try {
      if (register) {
        // Auth uses handwritten Axios calls and synchronous Redux reducers (no thunk/codegen).
        await http.post("/auth/register", { username: username.trim(), password });
        setNotice("Аккаунт создан. Теперь войдите.");
        setRegister(false);
        dispatch(authFailed(""));
        return;
      }
      const response = await http.post("/auth/login", { username: username.trim(), password });
      const token = readToken(response.data);
      if (!token) throw new Error("Сервер не вернул JWT. Проверьте реализацию /api/auth/login в Go.");
      localStorage.setItem(TOKEN_STORAGE_KEY, token);
      dispatch(loginSucceeded({ username: username.trim() }));
      navigate("/stages");
    } catch (error) { dispatch(authFailed(apiError(error))); }
  };
  return <main className="form-page"><h1>{register ? "Регистрация" : "Вход"}</h1>
    {isAuthenticated && <p>Вы уже вошли в аккаунт.</p>}
    <form onSubmit={event => { void submit(event); }} className="auth-form">
      <label htmlFor="login_username">Имя пользователя</label>
      <input id="login_username" autoComplete="username" value={username} disabled={busy} onChange={e => setUsername(e.target.value)} required/>
      <label htmlFor="login_password">Пароль</label>
      <input id="login_password" type="password" autoComplete={register ? "new-password" : "current-password"}
        value={password} disabled={busy} onChange={e => setPassword(e.target.value)} required/>
      {notice && <p className="api-connected">{notice}</p>}
      {error && <p className="api-warning" role="alert">{error}</p>}
      <button className="publish-button" type="submit" disabled={busy}>{register ? "Зарегистрироваться" : "Войти"}</button>
    </form>
    <button type="button" className="text-link" disabled={busy} onClick={() => {
      setRegister(v => !v); setNotice("");
    }}>{register ? "Уже есть аккаунт? Войти" : "Нет аккаунта? Зарегистрироваться"}</button>

  </main>;
}
