import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import axios from "axios";
import { hypertensionStagesApi } from "../api/generated/HypertensionStagesApi";
import { apiError, isTransportError } from "../api/http";
import { applyFilters, normalizeStage, unpackList, unpackStage } from "../api/normalize";
import { MOCK_STAGES } from "../data/mock";
import type { HypertensionStage } from "../types";
import type { RootState } from "./store";

interface StagesState {
  items: HypertensionStage[];
  selected: HypertensionStage | null;
  draft: HypertensionStage | null;
  source: "backend" | "mock" | "none";
  loading: boolean;
  error: string | null;
}
const initialState: StagesState = {
  items: [], selected: null, draft: null, source: "none", loading: false, error: null,
};

export const fetchHypertensionStages = createAsyncThunk<
  { items: HypertensionStage[]; source: StagesState["source"] }, void, { state: RootState; rejectValue: string }
>("hypertensionStages/fetchList", async (_, { getState, rejectWithValue }) => {
  const state = getState();
  const f = state.hypertensionFilters;
  try {
    // The filtering values are actually sent to Go, not just applied in React.
    const response = await hypertensionStagesApi.listHypertensionStages({
      min_sbp: f.minBP, max_sbp: f.maxBP, title: f.title.trim() || undefined,
    });
    return { items: applyFilters(unpackList(response.data), f), source: "backend" };
  } catch (error) {
    if (!state.hypertensionAuth.isAuthenticated && isTransportError(error)) {
      return { items: applyFilters(MOCK_STAGES, f), source: "mock" };
    }
    return rejectWithValue(apiError(error));
  }
});

export const fetchHypertensionStage = createAsyncThunk<
  HypertensionStage, number, { state: RootState; rejectValue: string }
>("hypertensionStages/fetchDetail", async (id, { getState, rejectWithValue }) => {
  try {
    const response = await hypertensionStagesApi.getHypertensionFeedById(id);
    return unpackStage(response.data);
  } catch (error) {
    if (!getState().hypertensionAuth.isAuthenticated && isTransportError(error)) {
      const stage = MOCK_STAGES.find(s => s.id === id && s.status === "Published");
      if (stage) return stage;
    }
    return rejectWithValue(apiError(error));
  }
});

export const fetchHypertensionDraft = createAsyncThunk<
  HypertensionStage | null, void, { rejectValue: string }
>("hypertensionStages/fetchDraft", async (_, { rejectWithValue }) => {
  try {
    const response = await hypertensionStagesApi.getHypertensionDraft();
    if (response.data == null || (typeof response.data === "object" &&
      "draft" in response.data && !response.data.draft)) return null;
    return unpackStage(response.data);
  } catch (error) {
    if (axios.isAxiosError(error) && error.response?.status === 404) return null;
    return rejectWithValue(apiError(error));
  }
});

export const createHypertensionStage = createAsyncThunk<
  HypertensionStage, FormData, { rejectValue: string }
>("hypertensionStages/createDraft", async (form, { rejectWithValue }) => {
  try {
    const response = await hypertensionStagesApi.createHypertensionStage(form);
    return { ...unpackStage(response.data), status: "Draft" };
  } catch (error) { return rejectWithValue(apiError(error)); }
});

export const publishHypertensionStage = createAsyncThunk<
  number, number, { rejectValue: string }
>("hypertensionStages/publish", async (id, { rejectWithValue }) => {
  try {
    await hypertensionStagesApi.publishHypertensionStage(id);
    return id;
  } catch (error) { return rejectWithValue(apiError(error)); }
});

export const deleteHypertensionStage = createAsyncThunk<
  number, number, { rejectValue: string }
>("hypertensionStages/delete", async (id, { rejectWithValue }) => {
  try {
    await hypertensionStagesApi.deleteHypertensionStage(id);
    return id;
  } catch (error) { return rejectWithValue(apiError(error)); }
});

export const toggleHypertensionLike = createAsyncThunk<
  { id: number; wasLiked: boolean }, { id: number; wasLiked: boolean }, { rejectValue: string }
>("hypertensionStages/toggleLike", async (payload, { rejectWithValue }) => {
  try {
    await hypertensionStagesApi.likeHypertensionStage(payload.id, { value: payload.wasLiked ? 0 : 1 });
    return payload;
  } catch (error) { return rejectWithValue(apiError(error)); }
});

const stagesSlice = createSlice({
  name: "hypertensionStages",
  initialState,
  reducers: {
    clearStagesError: state => { state.error = null; },
  },
  extraReducers: builder => {
    const pending = (state: StagesState) => { state.loading = true; state.error = null; };
    const rejected = (state: StagesState, action: { payload?: string }) => {
      state.loading = false;
      state.error = action.payload || "Ошибка обращения к API";
    };
    builder
      .addCase(fetchHypertensionStages.pending, pending)
      .addCase(fetchHypertensionStages.fulfilled, (state, action) => {
        state.loading = false;
        state.items = action.payload.items;
        state.source = action.payload.source;
      })
      .addCase(fetchHypertensionStages.rejected, (state, action) => {
        rejected(state, action);
        state.items = [];
        state.source = "none";
      })
      .addCase(fetchHypertensionStage.pending, pending)
      .addCase(fetchHypertensionStage.fulfilled, (state, action) => {
        state.loading = false;
        state.selected = action.payload;
      })
      .addCase(fetchHypertensionStage.rejected, (state, action) => {
        rejected(state, action);
        state.selected = null;
      })
      .addCase(fetchHypertensionDraft.pending, pending)
      .addCase(fetchHypertensionDraft.fulfilled, (state, action) => {
        state.loading = false;
        state.draft = action.payload;
      })
      .addCase(fetchHypertensionDraft.rejected, rejected)
      .addCase(createHypertensionStage.pending, pending)
      .addCase(createHypertensionStage.fulfilled, (state, action) => {
        state.loading = false;
        state.draft = action.payload;
      })
      .addCase(createHypertensionStage.rejected, rejected)
      .addCase(publishHypertensionStage.pending, pending)
      .addCase(publishHypertensionStage.fulfilled, (state, action) => {
        state.loading = false;
        if (state.draft?.id === action.payload) state.draft.status = "Published";
      })
      .addCase(publishHypertensionStage.rejected, rejected)
      .addCase(deleteHypertensionStage.pending, pending)
      .addCase(deleteHypertensionStage.fulfilled, (state, action) => {
        state.loading = false;
        state.items = state.items.filter(s => s.id !== action.payload);
        if (state.selected?.id === action.payload) state.selected = null;
      })
      .addCase(deleteHypertensionStage.rejected, rejected)
      .addCase(toggleHypertensionLike.pending, pending)
      .addCase(toggleHypertensionLike.fulfilled, (state, action) => {
        state.loading = false;
        const { id, wasLiked } = action.payload;
        for (const stage of [...state.items, ...(state.selected ? [state.selected] : [])]) {
          if (stage.id !== id) continue;
          stage.is_liked = !wasLiked;
          stage.likes_count = Math.max(0, (stage.likes_count || 0) + (wasLiked ? -1 : 1));
        }
      })
      .addCase(toggleHypertensionLike.rejected, rejected);
  },
});

export const { clearStagesError } = stagesSlice.actions;
export default stagesSlice.reducer;
