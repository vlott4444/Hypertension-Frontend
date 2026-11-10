import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { hypertensionStagesApi } from "../api/generated/HypertensionStagesApi";
import { apiError, isTransportError } from "../api/http";
import { unpackStage } from "../api/normalize";
import { MOCK_STAGES } from "../data/mock";
import type { HypertensionStage } from "../types";
import type { RootState } from "./store";
import { toggleHypertensionLike } from "./stagesSlice";

interface FeedState {
  current: HypertensionStage | null;
  source: "backend" | "mock" | "none";
  loading: boolean;
  error: string | null;
}
const initialState: FeedState = { current: null, source: "none", loading: false, error: null };
const mockItems = MOCK_STAGES.filter(item => item.status === "Published");

export const loadHypertensionFeed = createAsyncThunk<
  { stage: HypertensionStage; source: FeedState["source"] },
  { next?: boolean } | undefined, { state: RootState; rejectValue: string }
>("hypertensionFeed/load", async (arg, { getState, rejectWithValue }) => {
  const state = getState();
  const current = state.hypertensionFeed.current;
  const next = arg?.next === true;
  try {
    const response = next && current
      ? await hypertensionStagesApi.getHypertensionFeedById(current.id, true)
      : await hypertensionStagesApi.getHypertensionFeed();
    return { stage: unpackStage(response.data), source: "backend" };
  } catch (error) {
    if (!state.hypertensionAuth.isAuthenticated && isTransportError(error)) {
      const index = current ? mockItems.findIndex(s => s.id === current.id) : -1;
      return { stage: mockItems[next ? (index + 1) % mockItems.length : 0], source: "mock" };
    }
    return rejectWithValue(apiError(error));
  }
});

const feedSlice = createSlice({
  name: "hypertensionFeed",
  initialState,
  reducers: {},
  extraReducers: builder => {
    builder
      .addCase(loadHypertensionFeed.pending, state => { state.loading = true; state.error = null; })
      .addCase(loadHypertensionFeed.fulfilled, (state, action) => {
        state.loading = false;
        state.current = action.payload.stage;
        state.source = action.payload.source;
      })
      .addCase(toggleHypertensionLike.fulfilled, (state, action) => {
        if (state.current?.id !== action.payload.id) return;
        const wasLiked = action.payload.wasLiked;
        state.current.is_liked = !wasLiked;
        state.current.likes_count = Math.max(0, (state.current.likes_count || 0) + (wasLiked ? -1 : 1));
      })
      .addCase(loadHypertensionFeed.rejected, (state, action) => {
        state.loading = false;
        state.current = null;
        state.source = "none";
        state.error = action.payload || "Не удалось получить ленту";
      });
  },
});
export default feedSlice.reducer;
