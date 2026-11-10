import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { StageFilters } from "../types";

const defaultFilters: StageFilters = { minBP: 100, maxBP: 200, title: "" };

const filterSlice = createSlice({
  name: "hypertensionFilters",
  initialState: defaultFilters,
  reducers: {
    setPressureRange: (state, action: PayloadAction<{ minBP: number; maxBP: number }>) => {
      state.minBP = action.payload.minBP;
      state.maxBP = action.payload.maxBP;
    },
    setTitleFilter: (state, action: PayloadAction<string>) => { state.title = action.payload; },
    resetFilters: () => ({ ...defaultFilters }),
  },
});

export const { setPressureRange, setTitleFilter, resetFilters } = filterSlice.actions;
export default filterSlice.reducer;
