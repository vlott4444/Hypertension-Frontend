import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

/** Intentionally not hydrated from localStorage; F5 demonstrably resets UI to guest. */
interface AuthState {
  isAuthenticated: boolean;
  username: string;
  busy: boolean;
  error: string | null;
}
const initialState: AuthState = { isAuthenticated: false, username: "", busy: false, error: null };
const authSlice = createSlice({
  name: "hypertensionAuth",
  initialState,
  reducers: {
    authStarted: state => { state.busy = true; state.error = null; },
    loginSucceeded: (state, action: PayloadAction<{ username: string }>) => {
      state.isAuthenticated = true;
      state.username = action.payload.username;
      state.busy = false;
      state.error = null;
    },
    authFailed: (state, action: PayloadAction<string>) => {
      state.busy = false;
      state.error = action.payload;
    },
    loggedOut: () => ({ ...initialState }),
  },
});
export const { authStarted, loginSucceeded, authFailed, loggedOut } = authSlice.actions;
export default authSlice.reducer;
