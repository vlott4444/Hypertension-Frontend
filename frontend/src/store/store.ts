import { configureStore } from "@reduxjs/toolkit";
import { useDispatch, useSelector, type TypedUseSelectorHook } from "react-redux";
import authReducer from "./authSlice";
import filtersReducer from "./filterSlice";
import stagesReducer from "./stagesSlice";
import feedReducer from "./feedSlice";

export const store = configureStore({
  reducer: {
    hypertensionAuth: authReducer,
    hypertensionFilters: filtersReducer,
    hypertensionStages: stagesReducer,
    hypertensionFeed: feedReducer,
  },
  // Redux Toolkit enables redux-thunk and Redux DevTools by default.
});

if (import.meta.env.DEV) {
  (window as any).__REDUX_STORE__ = store;
}

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
export const useAppDispatch = useDispatch.withTypes<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
